package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
)

func TestIntegrationSetTerminal(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	initial := settings.Defaults()
	initial.Qtile.Terminal = "kitty"

	data, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("failed to write settings.json: %v", err)
	}

	// Simulate set_terminal: change to foot
	delta := map[string]interface{}{
		"qtile": map[string]interface{}{
			"terminal": "foot",
		},
	}
	if err := settings.SaveDelta(settingsPath, delta); err != nil {
		t.Fatalf("SaveDelta() error = %v", err)
	}

	// Verify settings.json was updated
	s, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if s.Qtile.Terminal != "foot" {
		t.Errorf("expected terminal to be 'foot', got %q", s.Qtile.Terminal)
	}

	// Verify config.py generation reflects the change
	out, err := GenerateConfigPy(s.Qtile)
	if err != nil {
		t.Fatalf("GenerateConfigPy() error = %v", err)
	}
	if !strings.Contains(out, `terminal = "foot"`) {
		t.Error("expected config.py to contain terminal = \"foot\"")
	}
}
