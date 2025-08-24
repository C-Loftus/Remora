package main

import "testing"

func TestDisplayServerType(t *testing.T) {
	if DetectDisplayServer() != "x11" {
		t.Errorf("Expected x11, got %s", DetectDisplayServer())
	}
}
