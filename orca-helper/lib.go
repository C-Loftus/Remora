package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	oc "github.com/c-loftus/orca-controller"

	"github.com/charmbracelet/log"
	"github.com/kbinani/screenshot"
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

func takeScreenshot() (string, error) {
	const activeDisplayIndex = 0
	bounds := screenshot.GetDisplayBounds(activeDisplayIndex)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%d_%dx%d.png", activeDisplayIndex, bounds.Dx(), bounds.Dy())

	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return "", err
	}

	return fileName, nil
}
