package main

import (
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestIntegrationAppearanceFullFlow(t *testing.T) {
	oldExecutor := executor
	executor = &module.MockExecutor{
		Responses: map[string]module.MockResponse{
			"feh --bg-scale /home/user/wallpaper.png": {
				Stdout:   "",
				ExitCode: 0,
			},
		},
	}
	defer func() { executor = oldExecutor }()

	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Neovim.UseGlobalTheme = true
	s.Qtile.UseGlobalTheme = true
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Run action returns dynamic theme options
	resp := captureAppearanceResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("run expected ok, got %s", resp.Status)
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok || len(opts["set-theme"].([]interface{})) != 4 {
		t.Fatalf("expected 4 themes in run response")
	}

	// Set theme with sync to neovim and qtile
	resp = captureAppearanceResponse(t, "set-theme", `{"value":"nord"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-theme expected ok, got %s", resp.Status)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Appearance.Theme != "nord" {
		t.Errorf("theme = %q, want nord", loaded.Appearance.Theme)
	}
	if loaded.Neovim.Theme != "nord" {
		t.Errorf("neovim theme = %q, want nord", loaded.Neovim.Theme)
	}
	if loaded.Qtile.ColorScheme != "nord" {
		t.Errorf("qtile color_scheme = %q, want nord", loaded.Qtile.ColorScheme)
	}

	// Set wallpaper
	resp = captureAppearanceResponse(t, "set-wallpaper", `{"value":"/home/user/wallpaper.png"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-wallpaper expected ok, got %s", resp.Status)
	}
	loaded, err = settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Appearance.Wallpaper != "/home/user/wallpaper.png" {
		t.Errorf("wallpaper = %q, want /home/user/wallpaper.png", loaded.Appearance.Wallpaper)
	}
}

func TestIntegrationAppearanceNoSyncWhenDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Neovim.UseGlobalTheme = false
	s.Neovim.Theme = "gruvbox"
	s.Qtile.UseGlobalTheme = false
	s.Qtile.ColorScheme = "custom"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":"dark"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("set-theme expected ok, got %s", resp.Status)
	}
	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Neovim.Theme != "gruvbox" {
		t.Errorf("neovim theme = %q, want gruvbox", loaded.Neovim.Theme)
	}
	if loaded.Qtile.ColorScheme != "custom" {
		t.Errorf("qtile color_scheme = %q, want custom", loaded.Qtile.ColorScheme)
	}
}
