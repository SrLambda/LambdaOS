package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/settings"
	"lambdaos.dev/lambda-env/pkg/module"
)

// v1_0_0Settings is a minimal Wave 2 (v1.0.0) settings fixture without
// the sections added in v1.1.0 and without use_global_theme.
var v1_0_0Settings = map[string]interface{}{
	"version": "1.0.0",
	"appearance": map[string]interface{}{
		"theme":     "dark",
		"font_size": 14,
		"opacity":   100,
		"wallpaper": "",
	},
	"display": map[string]interface{}{
		"active_profile": "default",
		"profiles":       []interface{}{},
	},
	"audio": map[string]interface{}{
		"default_sink": "",
		"volume":       75,
		"muted":        false,
	},
	"network": map[string]interface{}{
		"wifi_enabled":   true,
		"known_networks": []interface{}{},
	},
	"bluetooth": map[string]interface{}{
		"enabled":       true,
		"paired_devices": []interface{}{},
	},
	"keyboard": map[string]interface{}{
		"layout":  "us",
		"variant": "",
		"options": "",
	},
	"neovim": map[string]interface{}{
		"theme":          "tokyonight",
		"font":           "JetBrainsMono",
		"lines":          40,
		"columns":        120,
		"enable_lsp":     true,
		"enable_copilot": true,
		"enable_neotree": true,
		"lsp_servers":    []interface{}{"gopls", "pyright"},
	},
	"qtile": map[string]interface{}{
		"bar_position":         "top",
		"bar_size":             24,
		"layouts":              []interface{}{},
		"terminal":             "kitty",
		"browser":              "firefox",
		"default_file_manager": "thunar",
		"groups": []interface{}{
			map[string]interface{}{"name": "1"},
			map[string]interface{}{"name": "2"},
			map[string]interface{}{"name": "3"},
			map[string]interface{}{"name": "4"},
			map[string]interface{}{"name": "5"},
			map[string]interface{}{"name": "6"},
			map[string]interface{}{"name": "7"},
			map[string]interface{}{"name": "8"},
			map[string]interface{}{"name": "9"},
		},
	},
	"services": map[string]interface{}{
		"enabled": []interface{}{},
	},
}

func writeV1_0_0Settings(t *testing.T, path string) {
	t.Helper()
	data, err := json.MarshalIndent(v1_0_0Settings, "", "  ")
	if err != nil {
		t.Fatalf("marshal v1.0.0 settings: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write v1.0.0 settings: %v", err)
	}
}

func TestSchemaMigrationV1_0_0ToV1_1_0PreservesValues(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load v1.0.0 settings: %v", err)
	}

	// Migration should bump version.
	if loaded.Version != "1.2.0" {
		t.Errorf("version = %q, want %q", loaded.Version, "1.2.0")
	}

	// Existing values should be preserved.
	if loaded.Appearance.Theme != "dark" {
		t.Errorf("appearance.theme = %q, want %q", loaded.Appearance.Theme, "dark")
	}
	if loaded.Audio.Volume != 75 {
		t.Errorf("audio.volume = %d, want %d", loaded.Audio.Volume, 75)
	}
	if loaded.Keyboard.Layout != "us" {
		t.Errorf("keyboard.layout = %q, want %q", loaded.Keyboard.Layout, "us")
	}
	if loaded.Neovim.Theme != "tokyonight" {
		t.Errorf("neovim.theme = %q, want %q", loaded.Neovim.Theme, "tokyonight")
	}
	if loaded.Qtile.Terminal != "kitty" {
		t.Errorf("qtile.terminal = %q, want %q", loaded.Qtile.Terminal, "kitty")
	}
}

func TestSchemaMigrationPopulatesNewSectionsWithDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load v1.0.0 settings: %v", err)
	}

	// New sections added in v1.1.0 should have defaults.
	if loaded.Power.ScreenTimeout != 300 {
		t.Errorf("power.screen_timeout = %d, want %d", loaded.Power.ScreenTimeout, 300)
	}
	if loaded.Power.SleepTimeout != 600 {
		t.Errorf("power.sleep_timeout = %d, want %d", loaded.Power.SleepTimeout, 600)
	}
	if loaded.Power.LidCloseAction != "suspend" {
		t.Errorf("power.lid_close_action = %q, want %q", loaded.Power.LidCloseAction, "suspend")
	}

	if loaded.Defaults.Browser != "" {
		t.Errorf("defaults.browser = %q, want empty", loaded.Defaults.Browser)
	}
	if loaded.Defaults.Terminal != "" {
		t.Errorf("defaults.terminal = %q, want empty", loaded.Defaults.Terminal)
	}

	if len(loaded.Autostart.Enabled) != 0 {
		t.Errorf("autostart.enabled should be empty, got %v", loaded.Autostart.Enabled)
	}

	if loaded.Updates.AutoUpdate != false {
		t.Errorf("updates.auto_update = %v, want false", loaded.Updates.AutoUpdate)
	}
	if loaded.Updates.CheckInterval != 86400 {
		t.Errorf("updates.check_interval = %d, want %d", loaded.Updates.CheckInterval, 86400)
	}

	if !loaded.Security.FirewallEnabled {
		t.Errorf("security.firewall_enabled = %v, want true", loaded.Security.FirewallEnabled)
	}
	if loaded.Security.SudoTimeout != 5 {
		t.Errorf("security.sudo_timeout = %d, want %d", loaded.Security.SudoTimeout, 5)
	}

	if loaded.Fonts.Monospace != "JetBrainsMono" {
		t.Errorf("fonts.monospace = %q, want %q", loaded.Fonts.Monospace, "JetBrainsMono")
	}
	if loaded.Fonts.FontSize != 14 {
		t.Errorf("fonts.font_size = %d, want %d", loaded.Fonts.FontSize, 14)
	}

	if !loaded.Notifications.Enabled {
		t.Errorf("notifications.enabled = %v, want true", loaded.Notifications.Enabled)
	}
	if loaded.Notifications.TimeoutSeconds != 5 {
		t.Errorf("notifications.timeout_seconds = %d, want %d", loaded.Notifications.TimeoutSeconds, 5)
	}
}

func TestSchemaMigrationAddsUseGlobalTheme(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load v1.0.0 settings: %v", err)
	}

	if !loaded.Neovim.UseGlobalTheme {
		t.Errorf("neovim.use_global_theme = %v, want true", loaded.Neovim.UseGlobalTheme)
	}
	if !loaded.Qtile.UseGlobalTheme {
		t.Errorf("qtile.use_global_theme = %v, want true", loaded.Qtile.UseGlobalTheme)
	}
}

func TestSettingsDeltaFlowAudioVolume(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	// Fake audio module that emits a volume delta.
	manifest := `{
  "name": "audio",
  "version": "0.1.0",
  "description": "Manage audio",
  "description_es": "Gestionar audio",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "1.0.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"echo '{\"status\":\"ok\",\"action\":\"set-volume\",\"settings_delta\":{\"audio\":{\"volume\":42}}}'\n"
	writeModule(t, modulesDir, "audio", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var audioMod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "audio" {
			audioMod = &h.Modules[i]
			break
		}
	}
	if audioMod == nil {
		t.Fatal("audio module not found")
	}

	_, err = h.ExecuteAction(*audioMod, "set-volume", nil)
	if err != nil {
		t.Fatalf("ExecuteAction: %v", err)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load after delta: %v", err)
	}
	if loaded.Audio.Volume != 42 {
		t.Errorf("audio.volume = %d, want 42", loaded.Audio.Volume)
	}
	// Existing values should be preserved.
	if loaded.Audio.Muted != false {
		t.Errorf("audio.muted = %v, want false", loaded.Audio.Muted)
	}
	// Migration should have happened too.
	if loaded.Version != "1.2.0" {
		t.Errorf("version = %q, want 1.2.0", loaded.Version)
	}
	if !loaded.Neovim.UseGlobalTheme {
		t.Errorf("neovim.use_global_theme should be true after migration")
	}
}

func TestSettingsDeltaFlowKeyboardLayout(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "keyboard",
  "version": "0.1.0",
  "description": "Manage keyboard",
  "description_es": "Gestionar teclado",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "1.0.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"echo '{\"status\":\"ok\",\"action\":\"set-layout\",\"settings_delta\":{\"keyboard\":{\"layout\":\"es\"}}}'\n"
	writeModule(t, modulesDir, "keyboard", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var keyboardMod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "keyboard" {
			keyboardMod = &h.Modules[i]
			break
		}
	}
	if keyboardMod == nil {
		t.Fatal("keyboard module not found")
	}

	_, err = h.ExecuteAction(*keyboardMod, "set-layout", nil)
	if err != nil {
		t.Fatalf("ExecuteAction: %v", err)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load after delta: %v", err)
	}
	if loaded.Keyboard.Layout != "es" {
		t.Errorf("keyboard.layout = %q, want es", loaded.Keyboard.Layout)
	}
	// Existing keyboard values preserved.
	if loaded.Keyboard.Variant != "" {
		t.Errorf("keyboard.variant = %q, want empty", loaded.Keyboard.Variant)
	}
}

func TestSettingsDeltaFlowAppearanceTheme(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "appearance",
  "version": "0.1.0",
  "description": "Manage appearance",
  "description_es": "Gestionar apariencia",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "1.0.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"echo '{\"status\":\"ok\",\"action\":\"set-theme\",\"settings_delta\":{\"appearance\":{\"theme\":\"nord\"}}}'\n"
	writeModule(t, modulesDir, "appearance", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var appearanceMod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "appearance" {
			appearanceMod = &h.Modules[i]
			break
		}
	}
	if appearanceMod == nil {
		t.Fatal("appearance module not found")
	}

	_, err = h.ExecuteAction(*appearanceMod, "set-theme", nil)
	if err != nil {
		t.Fatalf("ExecuteAction: %v", err)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load after delta: %v", err)
	}
	if loaded.Appearance.Theme != "nord" {
		t.Errorf("appearance.theme = %q, want nord", loaded.Appearance.Theme)
	}
	// Existing appearance values preserved.
	if loaded.Appearance.FontSize != 14 {
		t.Errorf("appearance.font_size = %d, want 14", loaded.Appearance.FontSize)
	}
}

func TestSettingsDeltaFlowMultipleDeltas(t *testing.T) {
	home := setupTestHome(t)
	modulesDir := filepath.Join(home, ".local", "share", "lambda-env", "modules")

	manifest := `{
  "name": "multi-delta",
  "version": "0.1.0",
  "description": "Multi delta",
  "description_es": "Multi delta",
  "category": "system",
  "requires_root": false,
  "dependencies": [],
  "min_hub_version": "1.0.0"
}`
	script := "#!/usr/bin/env bash\n" +
		"echo '{\"status\":\"ok\",\"action\":\"run\",\"settings_delta\":{\"audio\":{\"volume\":60,\"muted\":true},\"appearance\":{\"theme\":\"light\"}}}'\n"
	writeModule(t, modulesDir, "multi-delta", manifest, script)

	settingsPath := filepath.Join(home, ".config", "lambdaos", "settings.json")
	writeV1_0_0Settings(t, settingsPath)

	h, err := hub.New(settingsPath)
	if err != nil {
		t.Fatalf("hub.New: %v", err)
	}
	defer h.Logger.Close()

	var mod *module.Manifest
	for i := range h.Modules {
		if h.Modules[i].Name == "multi-delta" {
			mod = &h.Modules[i]
			break
		}
	}
	if mod == nil {
		t.Fatal("multi-delta module not found")
	}

	_, err = h.ExecuteModule(*mod)
	if err != nil {
		t.Fatalf("ExecuteModule: %v", err)
	}

	loaded, err := settings.Load(settingsPath)
	if err != nil {
		t.Fatalf("Load after delta: %v", err)
	}
	if loaded.Audio.Volume != 60 {
		t.Errorf("audio.volume = %d, want 60", loaded.Audio.Volume)
	}
	if !loaded.Audio.Muted {
		t.Errorf("audio.muted = %v, want true", loaded.Audio.Muted)
	}
	if loaded.Appearance.Theme != "light" {
		t.Errorf("appearance.theme = %q, want light", loaded.Appearance.Theme)
	}
}
