package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
)

func TestIntegrationToggleLsp(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	initial := settings.Defaults()
	initial.Neovim.EnableLSP = true

	data, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("failed to write settings.json: %v", err)
	}

	// Simulate toggle: flip enable_lsp
	delta := map[string]interface{}{
		"neovim": map[string]interface{}{
			"enable_lsp": false,
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
	if s.Neovim.EnableLSP {
		t.Error("expected enable_lsp to be false after toggle")
	}

	// Verify lazy.lua generation reflects the toggle
	out, err := GenerateLazyLua(s.Neovim)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}
	if strings.Contains(out, `{ import = "plugins.lsp" }`) {
		t.Error("expected lazy.lua NOT to contain plugins.lsp when enable_lsp=false")
	}
}

func TestIntegrationToggleCopilot(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	initial := settings.Defaults()
	initial.Neovim.EnableCopilot = true

	data, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatalf("failed to write settings.json: %v", err)
	}

	// Simulate toggle: flip enable_copilot
	delta := map[string]interface{}{
		"neovim": map[string]interface{}{
			"enable_copilot": false,
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
	if s.Neovim.EnableCopilot {
		t.Error("expected enable_copilot to be false after toggle")
	}

	// Verify lazy.lua generation reflects the toggle
	out, err := GenerateLazyLua(s.Neovim)
	if err != nil {
		t.Fatalf("GenerateLazyLua() error = %v", err)
	}
	if strings.Contains(out, `{ import = "plugins.ai" }`) {
		t.Error("expected lazy.lua NOT to contain plugins.ai when enable_copilot=false")
	}
}
