package main

import (
	"fmt"
	"time"

	oc "github.com/c-loftus/orca-controller"
	"github.com/charmbracelet/log"
	"golang.design/x/hotkey"
)

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

	for {
		select {
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
	err := handleKeys()
	if err != nil {
		log.Fatal(err)
	}
}
