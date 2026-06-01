package main

import (
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestIntegrationDefaultsFullFlow(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"xdg-settings set default-web-browser firefox.desktop": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default kitty.desktop x-scheme-handler/terminal": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default nvim.desktop text/plain": {
				Stdout:   "",
				ExitCode: 0,
			},
			"xdg-mime default thunar.desktop inode/directory": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	appsDir = tmpDir
	for _, app := range []string{"firefox", "chromium", "kitty", "thunar", "nvim"} {
		if err := os.WriteFile(filepath.Join(tmpDir, app+".desktop"), []byte("[Desktop Entry]"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	settingsDir := t.TempDir()
	settingsPath := filepath.Join(settingsDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Run action returns dynamic app options
	resp := captureDefaultsResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("run expected ok, got %s", resp.Status)
	}
	apps, ok := resp.Data["available_apps"].([]interface{})
	if !ok || len(apps) != 5 {
		t.Fatalf("expected 5 apps in run response, got %v", resp.Data["available_apps"])
	}

	// Set browser
	resp = captureDefaultsResponse(t, "set-browser", `{"value":"firefox"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-browser expected ok, got %s", resp.Status)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Defaults.Browser != "firefox" {
		t.Errorf("browser = %q, want firefox", loaded.Defaults.Browser)
	}

	// Set terminal
	resp = captureDefaultsResponse(t, "set-terminal", `{"value":"kitty"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-terminal expected ok, got %s", resp.Status)
	}
	loaded, err = settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Defaults.Terminal != "kitty" {
		t.Errorf("terminal = %q, want kitty", loaded.Defaults.Terminal)
	}

	// Apply all defaults via batch
	s = settings.Defaults()
	s.Defaults.Browser = "firefox"
	s.Defaults.Terminal = "kitty"
	s.Defaults.Editor = "nvim"
	s.Defaults.FileManager = "thunar"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp = captureDefaultsResponse(t, "apply", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("apply expected ok, got %s", resp.Status)
	}
	results, ok := resp.Data["results"].(map[string]interface{})
	if !ok || len(results) != 4 {
		t.Fatalf("expected 4 apply results, got %v", resp.Data["results"])
	}
}
