package main

import (
	"context"
	"time"

	"fmt"
	"image/png"
	"os"

	oc "github.com/c-loftus/orca-controller"
	"github.com/charmbracelet/log"
	"github.com/kbinani/screenshot"
	"github.com/ollama/ollama/api"
	"golang.design/x/hotkey"
)

func takeScreenshot() (string, error) {
	const activeDisplayIndex = 0
	bounds := screenshot.GetDisplayBounds(activeDisplayIndex)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%d_%dx%d.png", activeDisplayIndex, bounds.Dx(), bounds.Dy())
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return "", err
	}

	return fileName, nil
}

func createClient() *oc.OrcaClient {
	for {
		client, err := oc.NewOrcaClient()
		_ = client.SpeechAndVerbosityManager.InterruptSpeech(false)
		err2 := client.PresentMessage("Rotor connected")
		if err == nil && err2 == nil {
			log.Info("Orca client created")
			return client
		}
		log.Error("Failed to create Orca client, retrying...")
		time.Sleep(2 * time.Second)
	}
}

func handleKeys() error {

	client := createClient()

	lowerSpeed := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF11)
	err := lowerSpeed.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	raiseSpeed := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF12)
	err = raiseSpeed.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	changeVerbosity := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF10)
	err = changeVerbosity.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	toggleSpeech := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF8)
	err = toggleSpeech.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	processScreenshot := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF9)
	err = processScreenshot.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	for {
		select {

		case <-toggleSpeech.Keydown():
			err := client.SpeechAndVerbosityManager.ToggleSpeech(true)
			if err != nil {
				log.Error(err)
				continue
			}

		case <-processScreenshot.Keydown():
			name, err := takeScreenshot()
			if err != nil {
				log.Error(err)
				continue
			}
			err = client.PresentMessage("Screenshot taken")
			if err != nil {
				log.Error(err)
				continue
			}
			ollamaClient, err := api.ClientFromEnvironment()
			if err != nil {
				log.Fatal(err)
			}

			asBytes, err := os.ReadFile(name)
			if err != nil {
				log.Error(err)
				continue
			}

			messages := api.Message{
				Role:   "user",
				Images: []api.ImageData{asBytes},
			}

			chatReq := api.ChatRequest{
				Model: "qwen2.5vl",
				Messages: []api.Message{
					messages,
				},
			}
			var allContent string
			respFunc := func(resp api.ChatResponse) error {
				allContent += resp.Message.Content
				return nil
			}
			client.PresentMessage("Processing screenshot...")
			log.Info("Processing screenshot...")
			if err := ollamaClient.Chat(context.Background(), &chatReq, respFunc); err != nil {
				log.Error(err)
				continue
			}
			log.Info(allContent)

			err = client.PresentMessage(allContent)
			if err != nil {
				log.Error(err)
				continue
			}
			err = os.Remove(name)
			if err != nil {
				log.Error(err)
				continue
			}
		case <-raiseSpeed.Keydown():

			err := client.SpeechAndVerbosityManager.InterruptSpeech(true)
			if err != nil {
				log.Error(err)
				continue
			}
			err = client.SpeechAndVerbosityManager.SetRate(100)
			if err != nil {
				log.Error(err)
				continue
			}
			rate, err := client.SpeechAndVerbosityManager.Rate()
			log.Info("Increased rate to " + fmt.Sprint(rate))
			if err != nil {
				log.Error(err)
				continue
			}
			err = client.PresentMessage("Rate " + fmt.Sprint(rate))
			if err != nil {
				log.Error(err)
				continue
			}

		case <-changeVerbosity.Keydown():

			err := client.SpeechAndVerbosityManager.ToggleVerbosity(true)
			if err != nil {
				log.Error(err)
				continue
			}
		case <-lowerSpeed.Keydown():
			err := client.SpeechAndVerbosityManager.SetRate(25)
			if err != nil {
				log.Error(err)
				continue
			}
			rate, err := client.SpeechAndVerbosityManager.Rate()
			if err != nil {
				log.Error(err)
				continue
			}
			log.Info("Decreased rate to " + fmt.Sprint(rate))
			err = client.PresentMessage("Rate " + fmt.Sprint(rate))
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func main() {

	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	models, err := ollamaClient.List(context.Background())
	if err != nil {
		log.Fatal(err)
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

		err = ollamaClient.Pull(context.Background(), req, progressFunc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Info("qwen2.5vl vision model found; skipping pull")
	}

	err = handleKeys()
	if err != nil {
		log.Fatal(err)
	}
}
