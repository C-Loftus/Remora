package main

import (
	"context"
	"image/png"
	"os"
	"strings"
	"time"

	oc "github.com/c-loftus/orca-controller"
	"github.com/kbinani/screenshot"
	"github.com/ollama/ollama/api"

	"github.com/charmbracelet/log"
)

type DisplayServerType string

const (
	Unknown DisplayServerType = "unknown"
	Wayland DisplayServerType = "wayland"
	X11     DisplayServerType = "x11"
)

// DetectDisplayServer returns "wayland", "x11", or "unknown".
func DetectDisplayServer() DisplayServerType {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}
	if os.Getenv("DISPLAY") != "" {
		return "x11"
	}

	sessionType := strings.ToLower(os.Getenv("XDG_SESSION_TYPE"))
	switch sessionType {
	case "wayland":
		return Wayland
	case "x11":
		return X11
	default:
		return Unknown
	}
}

func takeScreenshotAndSendToLlm(client *oc.OrcaClient) error {
	name, err := takeScreenshot()
	if err != nil {
		return err
	}
	err = client.PresentMessage("Screenshot taken")
	if err != nil {
		return err
	}
	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	asBytes, err := os.ReadFile(name)
	if err != nil {
		return err
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
		return err
	}
	log.Info(allContent)

	err = client.PresentMessage(allContent)
	if err != nil {
		return err
	}
	return os.Remove(name)
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

func takeScreenshot() (string, error) {
	const activeDisplayIndex = 0
	bounds := screenshot.GetDisplayBounds(activeDisplayIndex)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}
	file, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return "", err
	}

	return file.Name(), nil
}
