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
	// Version preserved.
	if s.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", s.Version, "1.0.0")
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
