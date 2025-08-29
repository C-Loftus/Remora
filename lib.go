package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strconv"
	"strings"

	oc "github.com/c-loftus/orca-controller"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/kbinani/screenshot"
	"github.com/ollama/ollama/api"
	"github.com/otiai10/gosseract/v2"

	"github.com/charmbracelet/log"
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

func SpeakAndLog(client *oc.OrcaClient, text string) error {
	log.Info(text)

	err := client.SpeechAndVerbosityManager.InterruptSpeech(false)
	if err != nil {
		return err
	}
	err = client.PresentMessage(text)
	if err != nil {
		return err
	}
	return nil
}

func takeScreenshotAndSendToLlm(client *oc.OrcaClient) error {
	if ollamaProcessing.Load() {
		SpeakAndLog(client, "Currently processing a screenshot")
		return nil
	}
	log.Info("Taking screenshot")
	name, err := takeScreenshot()
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

	chatReq := api.ChatRequest{
		Model: "qwen2.5vl",
		Messages: []api.Message{
			{
				Role:   "user",
				Images: []api.ImageData{asBytes},
			},
		},
	}

	// if the user specified a prompt, use it
	if visionModelPrompt != "" {
		chatReq.Messages = append(chatReq.Messages, api.Message{
			Role:    "user",
			Content: visionModelPrompt,
		})
	}

	var allContent string
	respFunc := func(resp api.ChatResponse) error {
		allContent += resp.Message.Content
		return nil
	}
	SpeakAndLog(client, "Processing screenshot")
	ollamaProcessing.Store(true)
	if err := ollamaClient.Chat(context.Background(), &chatReq, respFunc); err != nil {
		return err
	}
	ollamaProcessing.Store(false)
	mostRecentOllamaResponse = allContent

	SpeakAndLog(client, allContent)
	return os.Remove(name)
}

func ocr(imageName string) (string, error) {
	ocr_client := gosseract.NewClient()
	defer ocr_client.Close()
	ocr_client.SetImage(imageName)
	text, err := ocr_client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
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

func getBrightness() (int, error) {
	out, err := exec.Command("ddcutil", "getvcp", "0x10").Output()
	if err != nil {
		return 0, err
	}

	text := string(bytes.TrimSpace(out))
	idx := strings.Index(text, "current value =")
	if idx < 0 {
		return 0, fmt.Errorf("could not find current value in: %s", text)
	}

	rest := text[idx+len("current value ="):]
	parts := strings.Split(rest, ",")
	val, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, fmt.Errorf("failed to parse brightness: %w", err)
	}

	return val, nil
}

func setBrightness(value int) error {
	cmd := exec.Command("ddcutil", "setvcp", "0x10", strconv.Itoa(value))
	return cmd.Run()
}

func createOverlay(X *xgb.Conn) (xproto.Window, error) {
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	win, err := xproto.NewWindowId(X)
	if err != nil {
		return 0, err
	}

	xproto.CreateWindow(
		X,
		screen.RootDepth,
		win,
		screen.Root,
		0, 0,
		screen.WidthInPixels,
		screen.HeightInPixels,
		0,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwBackPixel|xproto.CwOverrideRedirect,
		[]uint32{0, 1}, // black background, override redirect
	)
	xproto.MapWindow(X, win)

	return win, nil
}

func toggleScreenCurtain() error {
	if !screenCurtainEnabled.Load() {
		log.Info("Enabling screen curtain")

		brightness, err := getBrightness()
		if err != nil {
			return err
		}
		savedBrightness.Store(int64(brightness))
		if err := setBrightness(0); err != nil {
			return err
		}

		X, err := xgb.NewConn()
		if err != nil {
			return err
		}
		win, err := createOverlay(X)
		if err != nil {
			return err
		}

		overlayWindow = win
		xConn = X
		screenCurtainEnabled.Store(true)
	} else {
		// Disable screen curtain
		log.Info("Disabling screen curtain")

		if xConn != nil {
			xproto.DestroyWindow(xConn, overlayWindow)
			xConn.Close()
			xConn = nil
		}
		if err := setBrightness(int(savedBrightness.Load())); err != nil {
			return err
		}
		screenCurtainEnabled.Store(false)
	}
	return nil
}
