package main

import (
	"strings"
	"testing"
)

func TestHotkeyList(t *testing.T) {
	app := NewApp()
	if len(app.GetHotKeys()) == 0 {
		t.Errorf("Hotkey list is empty")
	}
	if !strings.Contains(app.GetHotKeys()[0], "lower speed") {
		t.Errorf("Expected effect contain lower speed, got '%s'", app.GetHotKeys()[0])
	}
}
