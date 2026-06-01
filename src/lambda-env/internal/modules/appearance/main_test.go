package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

func captureAppearanceResponse(t *testing.T, action, params, settingsPath string) module.Response {
	t.Helper()
	t.Setenv("LAMBDA_ENV_ACTION", action)
	if params != "" {
		t.Setenv("LAMBDA_ENV_PARAMS", params)
	} else {
		os.Unsetenv("LAMBDA_ENV_PARAMS")
	}
	t.Setenv("LAMBDA_ENV_SETTINGS", settingsPath)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var resp module.Response
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func TestHandleRunReturnsCurrentThemeAndAvailableThemes(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Appearance.Theme = "dark"
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "run", "", settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	if resp.Data["theme"] != "dark" {
		t.Errorf("theme = %v, want dark", resp.Data["theme"])
	}
	opts, ok := resp.Data["available_options"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected available_options map")
	}
	themes, ok := opts["set-theme"].([]interface{})
	if !ok || len(themes) != 4 {
		t.Fatalf("expected 4 themes, got %v", opts["set-theme"])
	}
}

func TestSetThemeValid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":"nord"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["appearance"].(map[string]interface{})
	if !ok || delta["theme"] != "nord" {
		t.Errorf("appearance delta theme = %v, want nord", delta["theme"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Appearance.Theme != "nord" {
		t.Errorf("Appearance.Theme = %q, want nord", loaded.Appearance.Theme)
	}
}

func TestSetThemeInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":"gruvbox"}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetWallpaperValid(t *testing.T) {
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
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-wallpaper", `{"value":"/home/user/wallpaper.png"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["appearance"].(map[string]interface{})
	if !ok || delta["wallpaper"] != "/home/user/wallpaper.png" {
		t.Errorf("wallpaper delta = %v", delta["wallpaper"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Appearance.Wallpaper != "/home/user/wallpaper.png" {
		t.Errorf("Wallpaper = %q, want /home/user/wallpaper.png", loaded.Appearance.Wallpaper)
	}
}

func TestSetWallpaperEmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-wallpaper", `{"value":""}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}

func TestSetFontSizeValid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-font-size", `{"value":18}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}
	delta, ok := resp.SettingsDelta["appearance"].(map[string]interface{})
	if !ok || delta["font_size"] != float64(18) {
		t.Errorf("font_size delta = %v (%T), want 18", delta["font_size"], delta["font_size"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Appearance.FontSize != 18 {
		t.Errorf("FontSize = %d, want 18", loaded.Appearance.FontSize)
	}
}

func TestSetFontSizeInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-font-size", `{"value":0}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error, got %s", resp.Status)
	}
}
