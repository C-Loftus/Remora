package main

import (
	oc "github.com/c-loftus/orca-controller"
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
)

var hotkeyList = []HotkeyWithMetadata{
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
	{
		effect:       "change verbosity",
		keysAsString: "Ctrl+Shift+F10",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF10),
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
		effect:        "describe screen",
		keysAsString:  "Ctrl+Shift+F9",
		hotkey:        hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF9),
		functionToRun: takeScreenshotAndSendToLlm,
	},
	{
		effect:       "screen curtain",
		keysAsString: "Ctrl+Shift+F7",
		hotkey:       hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyF7),
		functionToRun: func(client *oc.OrcaClient) error {
			if screenCurtainEnabled {
				err := client.PresentMessage("Disabling screen curtain")
				if err != nil {
					return err
				}
			} else {
				err := client.PresentMessage("Enabling screen curtain")
				if err != nil {
					return err
				}
			}

			return toggleScreenCurtain()
		},
	},
}
