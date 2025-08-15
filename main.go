// main.go
package main

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"time"

	oc "github.com/c-loftus/orca-controller"
	"github.com/charmbracelet/log"
	"github.com/kbinani/screenshot"
	"github.com/ollama/ollama/api"
	"golang.design/x/hotkey"

	// gotk4 packages (GTK4 + GLib)
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	modelName          = "qwen2.5vl"
	modelQualifiedName = "qwen2.5vl:latest"
)

type AppUI struct {
	application *gtk.Application

	// widgets
	connLabel    *gtk.Label
	modelLabel   *gtk.Label
	modelSpinner *gtk.Spinner
	ollamaLabel  *gtk.Label
	lastMsgView  *gtk.Label
	shortcutsBox *gtk.Box
}

func newAppUI(app *gtk.Application) *AppUI {
	// build main window and widgets
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("Orca Controller")
	win.SetDefaultSize(640, 480)

	root := gtk.NewBox(gtk.OrientationVertical, 8)
	root.SetMarginTop(10)
	root.SetMarginBottom(10)
	root.SetMarginStart(10)
	root.SetMarginEnd(10)

	// Connection row
	connRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
	connLabel := gtk.NewLabel("Disconnected ❌")
	connRow.Append(gtk.NewLabel("Connection:"))
	connRow.Append(connLabel)

	// Model status row
	modelRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
	modelLabel := gtk.NewLabel("Model: idle")
	modelSpinner := gtk.NewSpinner()
	modelSpinner.Stop()
	modelRow.Append(gtk.NewLabel("Model status:"))
	modelRow.Append(modelLabel)
	modelRow.Append(modelSpinner)

	// Ollama row (pull / presence)
	ollamaLabel := gtk.NewLabel("") // will show pull progress / presence

	// Shortcuts
	shortcutsCardLabel := gtk.NewLabel("Keyboard shortcuts:")
	shortcutsBox := gtk.NewBox(gtk.OrientationVertical, 4)
	for _, r := range shortcutRows() {
		row := gtk.NewBox(gtk.OrientationHorizontal, 6)
		row.Append(gtk.NewLabel(r[0]))
		row.Append(gtk.NewLabel(" → "))
		row.Append(gtk.NewLabel(r[1]))
		shortcutsBox.Append(row)
	}

	// Last model message
	lastMsgLabelTitle := gtk.NewLabel("Last model message:")
	lastMsgView := gtk.NewLabel("")
	lastMsgView.SetWrap(true)
	lastMsgView.SetHAlign(gtk.AlignStart)
	lastMsgView.SetVAlign(gtk.AlignStart)
	lastMsgView.SetMaxWidthChars(80)

	// pack everything
	root.Append(connRow)
	root.Append(modelRow)

	// separator
	sep1 := gtk.NewSeparator(gtk.OrientationHorizontal)
	root.Append(sep1)

	root.Append(ollamaLabel)
	root.Append(shortcutsCardLabel)
	root.Append(shortcutsBox)

	sep2 := gtk.NewSeparator(gtk.OrientationHorizontal)
	root.Append(sep2)

	root.Append(lastMsgLabelTitle)
	root.Append(lastMsgView)

	win.SetChild(root)
	win.Present()

	return &AppUI{
		application:  app,
		connLabel:    connLabel,
		modelLabel:   modelLabel,
		modelSpinner: modelSpinner,
		ollamaLabel:  ollamaLabel,
		lastMsgView:  lastMsgView,
		shortcutsBox: shortcutsBox,
	}
}

// helper to schedule UI work on main loop
func idleAdd(f func()) {
	// glib.IdleAdd takes func() bool; returning false means run once.
	glib.IdleAdd(func() bool {
		f()
		return false
	})
}

// ------------ your existing code adapted ---------------

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

func createClient(ui *AppUI) *oc.OrcaClient {
	for {
		client, err := oc.NewOrcaClient()
		// best-effort: avoid nil panic if NewOrcaClient returns nil
		if client != nil {
			_ = client.SpeechAndVerbosityManager.InterruptSpeech(false)
		}
		var err2 error
		if client != nil {
			err2 = client.PresentMessage("Rotor connected")
		}
		if err == nil && err2 == nil {
			// update UI on main loop
			idleAdd(func() {
				ui.connLabel.SetText("Connected ✅")
			})
			log.Info("Orca client created")
			return client
		}

		idleAdd(func() {
			ui.connLabel.SetText("Disconnected ❌ (retrying in 2s)")
		})
		log.Error("Failed to create Orca client, retrying...")
		time.Sleep(2 * time.Second)
	}
}

func handleKeys(ui *AppUI) error {
	client := createClient(ui)

	lowerSpeed := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF11)
	if err := lowerSpeed.Register(); err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}
	raiseSpeed := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF12)
	if err := raiseSpeed.Register(); err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}
	changeVerbosity := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF10)
	if err := changeVerbosity.Register(); err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}
	toggleSpeech := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF8)
	if err := toggleSpeech.Register(); err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}
	processScreenshot := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF9)
	if err := processScreenshot.Register(); err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return err
	}

	for {
		select {
		case <-toggleSpeech.Keydown():
			if err := client.SpeechAndVerbosityManager.ToggleSpeech(true); err != nil {
				log.Error(err)
			}

		case <-processScreenshot.Keydown():
			// mark model busy
			idleAdd(func() {
				ui.modelLabel.SetText("Model: processing…")
				ui.modelSpinner.Start()
			})

			name, err := takeScreenshot()
			if err != nil {
				log.Error(err)
				idleAdd(func() {
					ui.modelLabel.SetText("Model: error (screenshot)")
					ui.modelSpinner.Stop()
				})
				continue
			}
			if err := client.PresentMessage("Screenshot taken"); err != nil {
				log.Error(err)
			}

			ollamaClient, err := api.ClientFromEnvironment()
			if err != nil {
				log.Fatal(err)
			}

			asBytes, err := os.ReadFile(name)
			if err != nil {
				log.Error(err)
				idleAdd(func() {
					ui.modelLabel.SetText("Model: error (read image)")
					ui.modelSpinner.Stop()
				})
				continue
			}

			messages := api.Message{
				Role:   "user",
				Images: []api.ImageData{asBytes},
			}

			chatReq := api.ChatRequest{
				Model: modelName,
				Messages: []api.Message{
					messages,
				},
			}
			var allContent string
			respFunc := func(resp api.ChatResponse) error {
				allContent += resp.Message.Content
				return nil
			}

			_ = client.PresentMessage("Processing screenshot...")
			log.Info("Processing screenshot...")
			if err := ollamaClient.Chat(context.Background(), &chatReq, respFunc); err != nil {
				log.Error(err)
				idleAdd(func() {
					ui.modelLabel.SetText("Model: error (chat)")
					ui.modelSpinner.Stop()
				})
				continue
			}
			log.Info(allContent)

			if err := client.PresentMessage(allContent); err != nil {
				log.Error(err)
			}

			idleAdd(func() {
				ui.lastMsgView.SetText(allContent)
				ui.modelLabel.SetText("Model: idle")
				ui.modelSpinner.Stop()
			})

			_ = os.Remove(name)

		case <-raiseSpeed.Keydown():
			if err := client.SpeechAndVerbosityManager.InterruptSpeech(true); err != nil {
				log.Error(err)
				continue
			}
			if err := client.SpeechAndVerbosityManager.SetRate(100); err != nil {
				log.Error(err)
				continue
			}
			rate, err := client.SpeechAndVerbosityManager.Rate()
			if err != nil {
				log.Error(err)
				continue
			}
			log.Info("Increased rate to " + fmt.Sprint(rate))
			if err := client.PresentMessage("Rate " + fmt.Sprint(rate)); err != nil {
				log.Error(err)
			}

		case <-changeVerbosity.Keydown():
			if err := client.SpeechAndVerbosityManager.ToggleVerbosity(true); err != nil {
				log.Error(err)
			}

		case <-lowerSpeed.Keydown():
			if err := client.SpeechAndVerbosityManager.SetRate(25); err != nil {
				log.Error(err)
				continue
			}
			rate, err := client.SpeechAndVerbosityManager.Rate()
			if err != nil {
				log.Error(err)
				continue
			}
			log.Info("Decreased rate to " + fmt.Sprint(rate))
			if err := client.PresentMessage("Rate " + fmt.Sprint(rate)); err != nil {
				log.Error(err)
			}
		}
	}
}

func ensureModel(ui *AppUI, ollamaClient *api.Client) error {
	models, err := ollamaClient.List(context.Background())
	if err != nil {
		return err
	}
	var found bool
	for _, m := range models.Models {
		if m.Name == modelQualifiedName {
			found = true
			break
		}
	}
	if found {
		idleAdd(func() { ui.ollamaLabel.SetText(fmt.Sprintf("%s found; skipping pull", modelQualifiedName)) })
		return nil
	}

	idleAdd(func() { ui.ollamaLabel.SetText(fmt.Sprintf("%s not found; pulling…", modelName)) })
	req := &api.PullRequest{Model: modelName}
	progressFunc := func(resp api.ProgressResponse) error {
		idleAdd(func() {
			ui.ollamaLabel.SetText(fmt.Sprintf("Pulling %s — status=%s, total=%d, completed=%d",
				modelName, resp.Status, resp.Total, resp.Completed))
		})
		return nil
	}
	return ollamaClient.Pull(context.Background(), req, progressFunc)
}

func shortcutRows() [][]string {
	return [][]string{
		{"Ctrl + Shift + F8", "Toggle Speech"},
		{"Ctrl + Shift + F9", "Process Screenshot"},
		{"Ctrl + Shift + F10", "Toggle Verbosity"},
		{"Ctrl + Shift + F11", "Lower Speech Speed"},
		{"Ctrl + Shift + F12", "Raise Speech Speed"},
	}
}

func main() {
	// create GTK application
	const appID = "com.example.orca"
	application := gtk.NewApplication(appID, gio.ApplicationDefaultFlags)

	application.Connect("activate", func() {
		ui := newAppUI(application)

		// Ensure Ollama model in background and show progress
		go func() {
			ollamaClient, err := api.ClientFromEnvironment()
			if err != nil {
				log.Fatal(err)
			}
			if err := ensureModel(ui, ollamaClient); err != nil {
				log.Error("ollama error:", err)
				idleAdd(func() { ui.ollamaLabel.SetText("Ollama error: " + err.Error()) })
			}
		}()

		// Run hotkeys & core logic in background
		go func() {
			if err := handleKeys(ui); err != nil {
				log.Fatal(err)
			}
		}()

		// show window
		win := application.ActiveWindow()
		if win == nil {
			// alternative: the constructor already created a window in newAppUI
			// but ensure it's shown
			// The newAppUI used ApplicationWindowNew which is already added to app
		}
		// present application window
		// The ApplicationWindow created by newAppUI will be shown by the runtime when app.Run is called.
	})

	// run the GTK application
	application.Run(os.Args)
}
