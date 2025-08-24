package main

import (
	"fmt"

	oc "github.com/c-loftus/orca-controller"
	"golang.design/x/hotkey"
)

type HotkeyWithMetadata struct {
	effect        string
	keysAsString  string
	hotkey        *hotkey.Hotkey
	functionToRun func(*oc.OrcaClient) error
}

func (h *HotkeyWithMetadata) ToString() string {
	return fmt.Sprintf("%s: %s", h.effect, h.keysAsString)
}

func NewHotkeyWithMetadata(effect string, keysAsString string, modifiers []hotkey.Modifier, key hotkey.Key) *HotkeyWithMetadata {
	return &HotkeyWithMetadata{
		hotkey:       hotkey.New(modifiers, key),
		effect:       effect,
		keysAsString: keysAsString,
	}
}

type DisplayServerType string

const (
	Unknown DisplayServerType = "unknown"
	Wayland DisplayServerType = "wayland"
	X11     DisplayServerType = "x11"
)
