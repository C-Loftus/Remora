package main

import (
	"fmt"

	oc "github.com/c-loftus/orca-controller"
	"github.com/charmbracelet/log"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"golang.design/x/hotkey"
)

var (
	screenCurtainEnabled     bool
	savedBrightness          int
	overlayWindow            xproto.Window
	xConn                    *xgb.Conn
	mostRecentOllamaResponse string
	mostRecentOcrResponse    string
)

var hotkeyList = []HotkeyWithMetadata{
	{
		effect:       "ocr screen",
		keysAsString: "Ctrl+Shift+F5",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF5),
		functionToRun: func(client *oc.OrcaClient) error {

			defer func() {
				// the tesseract go bindings are a bit flakey at times
				// it is best to register this panic handler here so
				// if something goes wrong, we can present a message
				// instead of having it crash
				if r := recover(); r != nil {
					err := fmt.Errorf("panic in ocr function: %v", r)
					client.PresentMessage("Unexpected error (panic) while running OCR")
					mostRecentOcrResponse = err.Error()
					log.Error(err)
				}
			}()

			screenshot, err := takeScreenshot()
			if err != nil {
				return err
			}
			ocrResult, err := ocr(screenshot)
			client.SpeechAndVerbosityManager.InterruptSpeech(false)
			if err != nil {
				client.PresentMessage("Error running OCR")
				mostRecentOcrResponse = err.Error()
				return err
			}
			client.PresentMessage("Finished running OCR")
			mostRecentOcrResponse = ocrResult
			return nil
		},
	},
	{
		effect:       "toggle verbosity",
		keysAsString: "Ctrl+Shift+F6",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF6),
		functionToRun: func(client *oc.OrcaClient) error {
			return client.SpeechAndVerbosityManager.ToggleVerbosity(true)
		},
	},
	{
		effect:        "describe screen with local AI model",
		keysAsString:  "Ctrl+Shift+F7",
		hotkey:        hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF7),
		functionToRun: takeScreenshotAndSendToLlm,
	},
	{
		effect:       "toggle speech",
		keysAsString: "Ctrl+Shift+F8",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF8),
		functionToRun: func(client *oc.OrcaClient) error {
			err := client.SpeechAndVerbosityManager.InterruptSpeech(true)
			if err != nil {
				return err
			}
			return client.SpeechAndVerbosityManager.ToggleSpeech(true)
		},
	},
	{
		effect:       "screen curtain",
		keysAsString: "Ctrl+Shift+F9",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF9),
		functionToRun: func(client *oc.OrcaClient) error {

			if screenCurtainEnabled {
				_ = client.PresentMessage("Disabling screen curtain")
			} else {
				_ = client.PresentMessage("Enabling screen curtain")
			}

			return toggleScreenCurtain()
		},
	},
	{
		effect:       "change verbosity",
		keysAsString: "Ctrl+Shift+F10",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF10),
	},
	{
		effect:       "slow speed",
		keysAsString: "Ctrl+Shift+F11",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF11),
		functionToRun: func(client *oc.OrcaClient) error {
			err := client.PresentMessage("Slow speed")
			if err != nil {
				return err
			}
			return client.SpeechAndVerbosityManager.SetRate(20)
		},
	},
	{
		effect:       "fast speed",
		keysAsString: "Ctrl+Shift+F12",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF12),
		functionToRun: func(client *oc.OrcaClient) error {
			err := client.PresentMessage("Fast speed")
			if err != nil {
				return err
			}
			return client.SpeechAndVerbosityManager.SetRate(100)
		},
	},
}
