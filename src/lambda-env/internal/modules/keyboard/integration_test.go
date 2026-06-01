package main

import (
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestIntegrationKeyboardFullFlow(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"setxkbmap -layout": {
				Stdout:   "us\nes\nde\n",
				ExitCode: 0,
			},
			"setxkbmap -layout es": {
				Stdout:   "",
				ExitCode: 0,
			},
			"setxkbmap -layout es -variant intl": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Run action returns dynamic layout options
	resp := captureKeyboardResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("run expected ok, got %s", resp.Status)
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok || len(opts["set-layout"].([]interface{})) != 3 {
		t.Fatalf("expected 3 layouts in run response")
	}

	// Set layout
	resp = captureKeyboardResponse(t, "set-layout", `{"value":"es"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-layout expected ok, got %s", resp.Status)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Keyboard.Layout != "es" {
		t.Errorf("layout = %q, want es", loaded.Keyboard.Layout)
	}

	// Set variant
	resp = captureKeyboardResponse(t, "set-variant", `{"value":"intl"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-variant expected ok, got %s", resp.Status)
	}
	loaded, err = settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Keyboard.Variant != "intl" {
		t.Errorf("variant = %q, want intl", loaded.Keyboard.Variant)
	}
}
