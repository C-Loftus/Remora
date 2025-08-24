package main

import "testing"

func TestHotkeyList(t *testing.T) {
	app := NewApp()
	if len(app.GetHotKeys()) == 0 {
		t.Errorf("Hotkey list is empty")
	}
	if app.GetHotKeys()[0].effect != "lower speed" {
		t.Errorf("Expected effect to be lower speed, got %s", app.GetHotKeys()[0].effect)
	}
}
