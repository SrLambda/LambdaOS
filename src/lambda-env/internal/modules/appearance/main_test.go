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

func TestThemeMap(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"dark", "tokyonight"},
		{"light", "tokyonight-light"},
		{"nord", "nord"},
		{"catppuccin", "catppuccin-mocha"},
	}
	for _, c := range cases {
		got, ok := themeMap[c.input]
		if !ok {
			t.Errorf("themeMap missing key %q", c.input)
			continue
		}
		if got != c.want {
			t.Errorf("themeMap[%q] = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestSetThemeSyncsNeovimAndQtile(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Neovim.UseGlobalTheme = true
	s.Qtile.UseGlobalTheme = true
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":"catppuccin"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	neovimDelta, ok := resp.SettingsDelta["neovim"].(map[string]interface{})
	if !ok || neovimDelta["theme"] != "catppuccin-mocha" {
		t.Errorf("neovim delta theme = %v, want catppuccin-mocha", neovimDelta["theme"])
	}
	qtileDelta, ok := resp.SettingsDelta["qtile"].(map[string]interface{})
	if !ok || qtileDelta["color_scheme"] != "catppuccin-mocha" {
		t.Errorf("qtile delta color_scheme = %v, want catppuccin-mocha", qtileDelta["color_scheme"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Neovim.Theme != "catppuccin-mocha" {
		t.Errorf("Neovim.Theme = %q, want catppuccin-mocha", loaded.Neovim.Theme)
	}
	if loaded.Qtile.ColorScheme != "catppuccin-mocha" {
		t.Errorf("Qtile.ColorScheme = %q, want catppuccin-mocha", loaded.Qtile.ColorScheme)
	}
}

func TestSetThemePartialSync(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	s.Neovim.UseGlobalTheme = false
	s.Neovim.Theme = "gruvbox"
	s.Qtile.UseGlobalTheme = true
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":"dark"}`, settingsPath)
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
	}

	_, hasNeovim := resp.SettingsDelta["neovim"]
	if hasNeovim {
		t.Error("expected no neovim delta when use_global_theme=false")
	}
	qtileDelta, ok := resp.SettingsDelta["qtile"].(map[string]interface{})
	if !ok || qtileDelta["color_scheme"] != "tokyonight" {
		t.Errorf("qtile delta = %v", resp.SettingsDelta["qtile"])
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Neovim.Theme != "gruvbox" {
		t.Errorf("Neovim.Theme = %q, want gruvbox", loaded.Neovim.Theme)
	}
	if loaded.Qtile.ColorScheme != "tokyonight" {
		t.Errorf("Qtile.ColorScheme = %q, want tokyonight", loaded.Qtile.ColorScheme)
	}
}

func TestSetThemeEmptyDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	s := settings.Defaults()
	if err := settings.Save(settingsPath, &s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	resp := captureAppearanceResponse(t, "set-theme", `{"value":""}`, settingsPath)
	if resp.Status != "error" {
		t.Fatalf("expected status error for empty theme, got %s", resp.Status)
	}
}

func TestSetThemeMapAllPresets(t *testing.T) {
	cases := []struct {
		theme          string
		expectedNeovim string
		expectedQtile  string
	}{
		{"dark", "tokyonight", "tokyonight"},
		{"light", "tokyonight-light", "tokyonight-light"},
		{"nord", "nord", "nord"},
		{"catppuccin", "catppuccin-mocha", "catppuccin-mocha"},
	}

	for _, c := range cases {
		t.Run(c.theme, func(t *testing.T) {
			tmpDir := t.TempDir()
			settingsPath := filepath.Join(tmpDir, "settings.json")
			s := settings.Defaults()
			s.Neovim.UseGlobalTheme = true
			s.Qtile.UseGlobalTheme = true
			if err := settings.Save(settingsPath, &s); err != nil {
				t.Fatalf("Save: %v", err)
			}

			resp := captureAppearanceResponse(t, "set-theme", `{"value":"`+c.theme+`"}`, settingsPath)
			if resp.Status != "ok" {
				t.Fatalf("expected status ok, got %s: %s", resp.Status, resp.Message)
			}

			loaded, err := settings.Load(settingsPath)
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if loaded.Appearance.Theme != c.theme {
				t.Errorf("Appearance.Theme = %q, want %q", loaded.Appearance.Theme, c.theme)
			}
			if loaded.Neovim.Theme != c.expectedNeovim {
				t.Errorf("Neovim.Theme = %q, want %q", loaded.Neovim.Theme, c.expectedNeovim)
			}
			if loaded.Qtile.ColorScheme != c.expectedQtile {
				t.Errorf("Qtile.ColorScheme = %q, want %q", loaded.Qtile.ColorScheme, c.expectedQtile)
			}
		})
	}
}
