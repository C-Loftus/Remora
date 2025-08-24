package main

import (
	"context"
	"embed"
	"sync/atomic"
	"time"

	oc "github.com/c-loftus/orca-controller"
	"github.com/charmbracelet/log"
	"github.com/ollama/ollama/api"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func handleKeys(app *App) error {

	for _, hotkey := range hotkeyList {
		if err := hotkey.hotkey.Register(); err != nil {
			return err
		}
		// spin off a goroutine for each hotkey listener
		go func(hk HotkeyWithMetadata, client *oc.OrcaClient) {
			log.Info("hotkey registered", "effect", hk.effect, "keys", hk.keysAsString)
			for range hk.hotkey.Keydown() {
				if hk.functionToRun != nil {
					err := hk.functionToRun(client)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}(hotkey, app.orcaConnection.OrcaClient)
	}

	// block forever so goroutines keep running
	select {}
}

func setupOllama() error {
	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	models, err := ollamaClient.List(context.Background())
	if err != nil {
		return err
	}
	var foundQwen bool
	for _, model := range models.Models {
		if model.Name == "qwen2.5vl:latest" {
			foundQwen = true
			break
		}
	}

	if !foundQwen {
		log.Info("qwen2.5vl not found; pulling")
		req := &api.PullRequest{
			Model: "qwen2.5vl",
		}
		progressFunc := func(resp api.ProgressResponse) error {
			log.Info("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
			return nil
		}

		return ollamaClient.Pull(context.Background(), req, progressFunc)
	}
	log.Info("qwen2.5vl vision model found; skipping pull")
	return nil
}

func main() {

	log.Info("Starting orca-helper. Please note that this program works best on X11 due to the lack of accessibility features in Wayland.")

	if err := setupOllama(); err != nil {
		log.Errorf("failed to setup ollama: %v", err)
	}

	// Create an instance of the app structure
	app := NewApp()

	var clientCreated atomic.Bool

	go func() {
		for {
			success := app.TryCreateClient()
			if success {
				clientCreated.Store(true)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			if clientCreated.Load() {
				break
			}
			time.Sleep(1 * time.Second)
		}
		err := handleKeys(app)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "orca-helper",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
