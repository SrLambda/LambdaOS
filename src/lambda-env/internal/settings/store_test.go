package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// File does not exist — should return defaults.
	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load missing file: unexpected error: %v", err)
	}

	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}
	if s.Appearance.Theme != "dark" {
		t.Errorf("Appearance.Theme = %q, want %q", s.Appearance.Theme, "dark")
	}
	if s.Audio.Volume != 75 {
		t.Errorf("Audio.Volume = %d, want %d", s.Audio.Volume, 75)
	}
	if s.Display.ActiveProfile != "default" {
		t.Errorf("Display.ActiveProfile = %q, want %q", s.Display.ActiveProfile, "default")
	}
	if len(s.Display.Profiles) != 0 {
		t.Errorf("Display.Profiles len = %d, want 0", len(s.Display.Profiles))
	}
	if !s.Network.WifiEnabled {
		t.Error("Network.WifiEnabled = false, want true")
	}
	if s.Keyboard.Layout != "us" {
		t.Errorf("Keyboard.Layout = %q, want %q", s.Keyboard.Layout, "us")
	}
	if s.Neovim.Theme != "tokyonight" {
		t.Errorf("Neovim.Theme = %q, want %q", s.Neovim.Theme, "tokyonight")
	}
	if s.Qtile.BarPosition != "top" {
		t.Errorf("Qtile.BarPosition = %q, want %q", s.Qtile.BarPosition, "top")
	}
}

func TestLoadPartial(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	partial := `{"version":"1.0.0","display":{"active_profile":"home"}}`
	if err := os.WriteFile(path, []byte(partial), 0644); err != nil {
		t.Fatalf("write partial settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load partial file: unexpected error: %v", err)
	}

	// Preserved value.
	if s.Display.ActiveProfile != "home" {
		t.Errorf("Display.ActiveProfile = %q, want %q", s.Display.ActiveProfile, "home")
	}
	// Default filled in.
	if s.Audio.Volume != 75 {
		t.Errorf("Audio.Volume = %d, want %d", s.Audio.Volume, 75)
	}
	if s.Appearance.Theme != "dark" {
		t.Errorf("Appearance.Theme = %q, want %q", s.Appearance.Theme, "dark")
	}
	// Version migrated to current.
	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}
}

func TestSaveAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	s := Defaults()
	s.Appearance.Theme = "gruvbox"

	if err := Save(path, &s); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	// File must exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("settings.json was not created")
	}

	// Content must be valid JSON and match.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var loaded Settings
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal saved settings: %v", err)
	}
	if loaded.Appearance.Theme != "gruvbox" {
		t.Errorf("saved theme = %q, want %q", loaded.Appearance.Theme, "gruvbox")
	}

	// No temp files should remain.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp") {
			t.Errorf("temp file left behind: %s", e.Name())
		}
	}
}

func TestSaveDeltaMerge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// Seed with some known values.
	s := Defaults()
	s.Audio.Volume = 50
	if err := Save(path, &s); err != nil {
		t.Fatalf("Save initial: %v", err)
	}

	// Apply delta that changes only display.active_profile.
	delta := map[string]interface{}{
		"display": map[string]interface{}{
			"active_profile": "office",
		},
	}
	if err := SaveDelta(path, delta); err != nil {
		t.Fatalf("SaveDelta: %v", err)
	}

	// Load and verify.
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load after delta: %v", err)
	}

	if loaded.Display.ActiveProfile != "office" {
		t.Errorf("Display.ActiveProfile = %q, want %q", loaded.Display.ActiveProfile, "office")
	}
	// Audio.Volume must be preserved.
	if loaded.Audio.Volume != 50 {
		t.Errorf("Audio.Volume = %d, want %d", loaded.Audio.Volume, 50)
	}
	// Appearance defaults must remain.
	if loaded.Appearance.Theme != "dark" {
		t.Errorf("Appearance.Theme = %q, want %q", loaded.Appearance.Theme, "dark")
	}
}

func TestDowngradeRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// Write a future version.
	future := `{"version":"2.0.0","appearance":{"theme":"dark","font_size":14,"opacity":100,"wallpaper":""},"display":{"active_profile":"default","profiles":[]},"audio":{"default_sink":"","volume":75,"muted":false},"network":{"wifi_enabled":true,"known_networks":[]},"bluetooth":{"enabled":true,"paired_devices":[]},"keyboard":{"layout":"us","variant":"","options":""},"neovim":{"theme":"tokyonight","font":"JetBrainsMono","lines":40,"columns":120},"qtile":{"bar_position":"top","bar_size":24,"layouts":[]},"services":{"enabled":[]}}`
	if err := os.WriteFile(path, []byte(future), 0644); err != nil {
		t.Fatalf("write future settings: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load future version: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "downgrade") {
		t.Errorf("error = %q, want to contain 'downgrade'", err.Error())
	}
}

func TestMigrationAddsMissingFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// Old version with only appearance partially set.
	old := `{"version":"0.9.0","appearance":{"theme":"gruvbox"}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load old version: unexpected error: %v", err)
	}

	// Version must be bumped.
	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}
	// Existing user value preserved.
	if s.Appearance.Theme != "gruvbox" {
		t.Errorf("Appearance.Theme = %q, want %q", s.Appearance.Theme, "gruvbox")
	}
	// Missing fields filled with defaults.
	if s.Appearance.FontSize != 14 {
		t.Errorf("Appearance.FontSize = %d, want %d", s.Appearance.FontSize, 14)
	}
	if s.Audio.Volume != 75 {
		t.Errorf("Audio.Volume = %d, want %d", s.Audio.Volume, 75)
	}
	if s.Keyboard.Layout != "us" {
		t.Errorf("Keyboard.Layout = %q, want %q", s.Keyboard.Layout, "us")
	}
}

func TestSaveDeltaEmptyNoOp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	s := Defaults()
	s.Audio.Volume = 42
	if err := Save(path, &s); err != nil {
		t.Fatalf("Save initial: %v", err)
	}

	// Read file content before.
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile before: %v", err)
	}

	// Apply empty delta.
	if err := SaveDelta(path, map[string]interface{}{}); err != nil {
		t.Fatalf("SaveDelta empty: %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile after: %v", err)
	}

	if string(before) != string(after) {
		t.Error("empty delta modified the file")
	}
}

func TestValidateInvalidVolume(t *testing.T) {
	s := Defaults()
	s.Audio.Volume = 150
	if err := s.Validate(); err == nil {
		t.Fatal("Validate high volume: expected error, got nil")
	}

	s.Audio.Volume = -1
	if err := s.Validate(); err == nil {
		t.Fatal("Validate negative volume: expected error, got nil")
	}
}

func TestValidateMissingVersion(t *testing.T) {
	s := Defaults()
	s.Version = ""
	if err := s.Validate(); err == nil {
		t.Fatal("Validate missing version: expected error, got nil")
	}
}

func TestValidateActiveProfileReference(t *testing.T) {
	s := Defaults()
	s.Display.Profiles = []OutputProfile{
		{Name: "home", Outputs: []OutputConfig{}},
		{Name: "office", Outputs: []OutputConfig{}},
	}
	s.Display.ActiveProfile = "invalid"
	if err := s.Validate(); err == nil {
		t.Fatal("Validate invalid profile: expected error, got nil")
	}

	// Valid profile should pass.
	s.Display.ActiveProfile = "home"
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate valid profile: unexpected error: %v", err)
	}
}

func TestValidateNeovimDefaults(t *testing.T) {
	s := Defaults()
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate neovim defaults: unexpected error: %v", err)
	}
	if !s.Neovim.EnableLSP {
		t.Error("Neovim.EnableLSP = false, want true")
	}
	if !s.Neovim.EnableCopilot {
		t.Error("Neovim.EnableCopilot = false, want true")
	}
	if !s.Neovim.EnableNeotree {
		t.Error("Neovim.EnableNeotree = false, want true")
	}
	if len(s.Neovim.LspServers) != 2 {
		t.Fatalf("Neovim.LspServers len = %d, want 2", len(s.Neovim.LspServers))
	}
	if s.Neovim.LspServers[0] != "gopls" || s.Neovim.LspServers[1] != "pyright" {
		t.Errorf("Neovim.LspServers = %v, want [gopls, pyright]", s.Neovim.LspServers)
	}
}

func TestValidateQtileDefaults(t *testing.T) {
	s := Defaults()
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate qtile defaults: unexpected error: %v", err)
	}
	if s.Qtile.Terminal != "kitty" {
		t.Errorf("Qtile.Terminal = %q, want %q", s.Qtile.Terminal, "kitty")
	}
	if s.Qtile.Browser != "firefox" {
		t.Errorf("Qtile.Browser = %q, want %q", s.Qtile.Browser, "firefox")
	}
	if s.Qtile.DefaultFileManager != "thunar" {
		t.Errorf("Qtile.DefaultFileManager = %q, want %q", s.Qtile.DefaultFileManager, "thunar")
	}
	if len(s.Qtile.Groups) != 9 {
		t.Fatalf("Qtile.Groups len = %d, want 9", len(s.Qtile.Groups))
	}
	for i := 0; i < 9; i++ {
		if s.Qtile.Groups[i].Name != string(rune('1'+i)) {
			t.Errorf("Qtile.Groups[%d].Name = %q, want %q", i, s.Qtile.Groups[i].Name, string(rune('1'+i)))
		}
	}
}

func TestValidateInvalidTerminal(t *testing.T) {
	s := Defaults()
	s.Qtile.Terminal = "gnome-terminal"
	if err := s.Validate(); err == nil {
		t.Fatal("Validate invalid terminal: expected error, got nil")
	}
}

func TestValidateInvalidBrowser(t *testing.T) {
	s := Defaults()
	s.Qtile.Browser = "safari"
	if err := s.Validate(); err == nil {
		t.Fatal("Validate invalid browser: expected error, got nil")
	}
}

func TestValidateEmptyTerminalSkipped(t *testing.T) {
	s := Defaults()
	s.Qtile.Terminal = ""
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate empty terminal: unexpected error: %v", err)
	}
}

func TestLoadDefaultsIncludesNewSections(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load missing file: unexpected error: %v", err)
	}

	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}

	// Power defaults.
	if s.Power.ScreenTimeout != 300 {
		t.Errorf("Power.ScreenTimeout = %d, want %d", s.Power.ScreenTimeout, 300)
	}
	if s.Power.SleepTimeout != 600 {
		t.Errorf("Power.SleepTimeout = %d, want %d", s.Power.SleepTimeout, 600)
	}
	if s.Power.LidCloseAction != "suspend" {
		t.Errorf("Power.LidCloseAction = %q, want %q", s.Power.LidCloseAction, "suspend")
	}

	// Defaults defaults.
	if s.Defaults.Browser != "" {
		t.Errorf("Defaults.Browser = %q, want empty", s.Defaults.Browser)
	}
	if s.Defaults.Terminal != "" {
		t.Errorf("Defaults.Terminal = %q, want empty", s.Defaults.Terminal)
	}
	if s.Defaults.Editor != "" {
		t.Errorf("Defaults.Editor = %q, want empty", s.Defaults.Editor)
	}
	if s.Defaults.FileManager != "" {
		t.Errorf("Defaults.FileManager = %q, want empty", s.Defaults.FileManager)
	}

	// Autostart defaults.
	if len(s.Autostart.Enabled) != 0 {
		t.Errorf("Autostart.Enabled len = %d, want 0", len(s.Autostart.Enabled))
	}

	// Updates defaults.
	if s.Updates.AutoUpdate != false {
		t.Error("Updates.AutoUpdate = true, want false")
	}
	if s.Updates.CheckInterval != 86400 {
		t.Errorf("Updates.CheckInterval = %d, want %d", s.Updates.CheckInterval, 86400)
	}
	if len(s.Updates.ExcludePackages) != 0 {
		t.Errorf("Updates.ExcludePackages len = %d, want 0", len(s.Updates.ExcludePackages))
	}

	// Security defaults.
	if s.Security.FirewallEnabled != true {
		t.Error("Security.FirewallEnabled = false, want true")
	}
	if s.Security.SudoTimeout != 5 {
		t.Errorf("Security.SudoTimeout = %d, want %d", s.Security.SudoTimeout, 5)
	}
	if s.Security.ScreenLockTimeout != 300 {
		t.Errorf("Security.ScreenLockTimeout = %d, want %d", s.Security.ScreenLockTimeout, 300)
	}

	// Fonts defaults.
	if s.Fonts.Monospace != "JetBrainsMono" {
		t.Errorf("Fonts.Monospace = %q, want %q", s.Fonts.Monospace, "JetBrainsMono")
	}
	if s.Fonts.SansSerif != "Noto Sans" {
		t.Errorf("Fonts.SansSerif = %q, want %q", s.Fonts.SansSerif, "Noto Sans")
	}
	if s.Fonts.Serif != "Noto Serif" {
		t.Errorf("Fonts.Serif = %q, want %q", s.Fonts.Serif, "Noto Serif")
	}
	if s.Fonts.FontSize != 14 {
		t.Errorf("Fonts.FontSize = %d, want %d", s.Fonts.FontSize, 14)
	}

	// Notifications defaults.
	if s.Notifications.Enabled != true {
		t.Error("Notifications.Enabled = false, want true")
	}
	if s.Notifications.DoNotDisturb != false {
		t.Error("Notifications.DoNotDisturb = true, want false")
	}
	if s.Notifications.TimeoutSeconds != 5 {
		t.Errorf("Notifications.TimeoutSeconds = %d, want %d", s.Notifications.TimeoutSeconds, 5)
	}
}

func TestMigrationV100ToV110(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// A realistic v1.0.0 settings file with some user values.
	old := `{"version":"1.0.0","appearance":{"theme":"gruvbox","font_size":12,"opacity":100,"wallpaper":""},"display":{"active_profile":"home","profiles":[]},"audio":{"default_sink":"alsa_output","volume":80,"muted":false},"network":{"wifi_enabled":false,"known_networks":[]},"bluetooth":{"enabled":false,"paired_devices":[]},"keyboard":{"layout":"es","variant":"","options":""},"neovim":{"theme":"catppuccin","font":"FiraCode","lines":50,"columns":100,"enable_lsp":true,"enable_copilot":false,"enable_neotree":true,"lsp_servers":["gopls"]},"qtile":{"bar_position":"bottom","bar_size":30,"layouts":[],"terminal":"foot","browser":"chromium","default_file_manager":"pcmanfm","groups":[]},"services":{"enabled":[]}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load v1.0.0: unexpected error: %v", err)
	}

	// Version bumped.
	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}

	// Existing values preserved.
	if s.Appearance.Theme != "gruvbox" {
		t.Errorf("Appearance.Theme = %q, want %q", s.Appearance.Theme, "gruvbox")
	}
	if s.Audio.Volume != 80 {
		t.Errorf("Audio.Volume = %d, want %d", s.Audio.Volume, 80)
	}
	if s.Keyboard.Layout != "es" {
		t.Errorf("Keyboard.Layout = %q, want %q", s.Keyboard.Layout, "es")
	}
	if s.Neovim.Theme != "catppuccin" {
		t.Errorf("Neovim.Theme = %q, want %q", s.Neovim.Theme, "catppuccin")
	}
	if s.Qtile.BarPosition != "bottom" {
		t.Errorf("Qtile.BarPosition = %q, want %q", s.Qtile.BarPosition, "bottom")
	}
	if s.Network.WifiEnabled != false {
		t.Error("Network.WifiEnabled = true, want false")
	}

	// New sections filled with defaults.
	if s.Power.ScreenTimeout != 300 {
		t.Errorf("Power.ScreenTimeout = %d, want %d", s.Power.ScreenTimeout, 300)
	}
	if s.Security.FirewallEnabled != true {
		t.Error("Security.FirewallEnabled = false, want true")
	}
	if s.Fonts.Monospace != "JetBrainsMono" {
		t.Errorf("Fonts.Monospace = %q, want %q", s.Fonts.Monospace, "JetBrainsMono")
	}
	if s.Notifications.Enabled != true {
		t.Error("Notifications.Enabled = false, want true")
	}
	if len(s.Autostart.Enabled) != 0 {
		t.Errorf("Autostart.Enabled len = %d, want 0", len(s.Autostart.Enabled))
	}
	if s.Updates.AutoUpdate != false {
		t.Error("Updates.AutoUpdate = true, want false")
	}
	if s.Defaults.Browser != "" {
		t.Errorf("Defaults.Browser = %q, want empty", s.Defaults.Browser)
	}
}

func TestMigrationV100ToV110NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// v1.0.0 with audio.volume set to 80 (default is 75).
	old := `{"version":"1.0.0","audio":{"volume":80}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load v1.0.0: unexpected error: %v", err)
	}

	// User value preserved despite different default.
	if s.Audio.Volume != 80 {
		t.Errorf("Audio.Volume = %d, want 80 (user value preserved)", s.Audio.Volume)
	}
}

func TestAudioDefaults(t *testing.T) {
	s := Defaults()
	if s.Audio.DefaultSource != "" {
		t.Errorf("Audio.DefaultSource = %q, want empty", s.Audio.DefaultSource)
	}
	if s.Audio.Profile != "" {
		t.Errorf("Audio.Profile = %q, want empty", s.Audio.Profile)
	}
	if len(s.Audio.Profiles) != 0 {
		t.Errorf("Audio.Profiles len = %d, want 0", len(s.Audio.Profiles))
	}
}

func TestAudioProfileSerialization(t *testing.T) {
	s := Defaults()
	s.Audio.Profiles = []AudioProfile{
		{Name: "Headphones", DefaultSink: "alsa_output.usb", DefaultSource: "alsa_input.usb", Volume: 80},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var loaded Settings
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(loaded.Audio.Profiles) != 1 {
		t.Fatalf("Audio.Profiles len = %d, want 1", len(loaded.Audio.Profiles))
	}
	if loaded.Audio.Profiles[0].Name != "Headphones" {
		t.Errorf("Profiles[0].Name = %q, want %q", loaded.Audio.Profiles[0].Name, "Headphones")
	}
	if loaded.Audio.Profiles[0].Volume != 80 {
		t.Errorf("Profiles[0].Volume = %d, want 80", loaded.Audio.Profiles[0].Volume)
	}
}

func TestMigrationV110ToV120AddsAudioFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// v1.1.0 settings without new audio fields.
	old := `{"version":"1.1.0","audio":{"default_sink":"alsa_output","volume":60,"muted":false}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load v1.1.0: unexpected error: %v", err)
	}

	if s.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, CurrentVersion)
	}
	if s.Audio.Volume != 60 {
		t.Errorf("Audio.Volume = %d, want 60", s.Audio.Volume)
	}
	if s.Audio.DefaultSource != "" {
		t.Errorf("Audio.DefaultSource = %q, want empty", s.Audio.DefaultSource)
	}
	if s.Audio.Profile != "" {
		t.Errorf("Audio.Profile = %q, want empty", s.Audio.Profile)
	}
	if len(s.Audio.Profiles) != 0 {
		t.Errorf("Audio.Profiles len = %d, want 0", len(s.Audio.Profiles))
	}
}

func TestValidateNewSections(t *testing.T) {
	// Power validation.
	s := Defaults()
	s.Power.ScreenTimeout = -1
	if err := s.Validate(); err == nil {
		t.Fatal("Validate negative screen_timeout: expected error, got nil")
	}

	// Notifications validation.
	s = Defaults()
	s.Notifications.TimeoutSeconds = -1
	if err := s.Validate(); err == nil {
		t.Fatal("Validate negative timeout_seconds: expected error, got nil")
	}

	// Security validation.
	s = Defaults()
	s.Security.SudoTimeout = -1
	if err := s.Validate(); err == nil {
		t.Fatal("Validate negative sudo_timeout: expected error, got nil")
	}

	// Fonts validation.
	s = Defaults()
	s.Fonts.FontSize = 0
	if err := s.Validate(); err == nil {
		t.Fatal("Validate zero font_size: expected error, got nil")
	}

	// Updates validation.
	s = Defaults()
	s.Updates.CheckInterval = -1
	if err := s.Validate(); err == nil {
		t.Fatal("Validate negative check_interval: expected error, got nil")
	}
}

func TestUseGlobalThemeDefaultsToTrue(t *testing.T) {
	s := Defaults()
	if !s.Neovim.UseGlobalTheme {
		t.Error("Neovim.UseGlobalTheme = false, want true")
	}
	if !s.Qtile.UseGlobalTheme {
		t.Error("Qtile.UseGlobalTheme = false, want true")
	}
}

func TestUseGlobalThemeRequiresNonEmptyTheme(t *testing.T) {
	// Neovim use_global_theme=true with empty theme should fail.
	s := Defaults()
	s.Neovim.UseGlobalTheme = true
	s.Appearance.Theme = ""
	err := s.Validate()
	if err == nil {
		t.Fatal("Validate empty theme with use_global_theme=true: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "appearance.theme") {
		t.Errorf("error = %q, want to contain 'appearance.theme'", err.Error())
	}

	// Qtile use_global_theme=true with empty theme should fail.
	s = Defaults()
	s.Qtile.UseGlobalTheme = true
	s.Appearance.Theme = ""
	if err := s.Validate(); err == nil {
		t.Fatal("Validate empty theme with qtile use_global_theme=true: expected error, got nil")
	}

	// use_global_theme=false with empty theme should pass.
	s = Defaults()
	s.Neovim.UseGlobalTheme = false
	s.Qtile.UseGlobalTheme = false
	s.Appearance.Theme = ""
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate use_global_theme=false with empty theme: unexpected error: %v", err)
	}
}

func TestMigrationV100ToV110AddsUseGlobalTheme(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// Minimal v1.0.0 file.
	old := `{"version":"1.0.0","neovim":{"theme":"tokyonight"},"qtile":{"bar_position":"top"}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load v1.0.0: unexpected error: %v", err)
	}

	if !s.Neovim.UseGlobalTheme {
		t.Error("Neovim.UseGlobalTheme = false after migration, want true")
	}
	if !s.Qtile.UseGlobalTheme {
		t.Error("Qtile.UseGlobalTheme = false after migration, want true")
	}
	// Existing values preserved.
	if s.Neovim.Theme != "tokyonight" {
		t.Errorf("Neovim.Theme = %q, want %q", s.Neovim.Theme, "tokyonight")
	}
	if s.Qtile.BarPosition != "top" {
		t.Errorf("Qtile.BarPosition = %q, want %q", s.Qtile.BarPosition, "top")
	}
}
