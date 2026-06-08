package module

import (
	"os"
	"path/filepath"
	"testing"

	"lambdaos.dev/lambda-env/internal/settings"
)

// TestFoundationSchemaMigration verifies that loading a v1.0.0 settings file
// automatically migrates to v1.1.0, preserving existing values and adding
// the 7 new sections with defaults.
func TestFoundationSchemaMigration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	old := `{"version":"1.0.0","appearance":{"theme":"nord","font_size":10},"audio":{"volume":60}}`
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatalf("write old settings: %v", err)
	}

	s, err := settings.Load(path)
	if err != nil {
		t.Fatalf("Load v1.0.0: unexpected error: %v", err)
	}

	if s.Version != settings.CurrentVersion {
		t.Errorf("Version = %q, want %q", s.Version, settings.CurrentVersion)
	}

	// Preserved.
	if s.Appearance.Theme != "nord" {
		t.Errorf("Appearance.Theme = %q, want %q", s.Appearance.Theme, "nord")
	}
	if s.Appearance.FontSize != 10 {
		t.Errorf("Appearance.FontSize = %d, want %d", s.Appearance.FontSize, 10)
	}
	if s.Audio.Volume != 60 {
		t.Errorf("Audio.Volume = %d, want %d", s.Audio.Volume, 60)
	}

	// New sections added with defaults.
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
	if s.Updates.AutoUpdate != false {
		t.Error("Updates.AutoUpdate = true, want false")
	}
	if len(s.Autostart.Enabled) != 0 {
		t.Errorf("Autostart.Enabled len = %d, want 0", len(s.Autostart.Enabled))
	}
	if s.Defaults.Browser != "" {
		t.Errorf("Defaults.Browser = %q, want empty", s.Defaults.Browser)
	}

	// use_global_theme added.
	if !s.Neovim.UseGlobalTheme {
		t.Error("Neovim.UseGlobalTheme = false, want true")
	}
	if !s.Qtile.UseGlobalTheme {
		t.Error("Qtile.UseGlobalTheme = false, want true")
	}
}

// TestFoundationManifestActionsParsing verifies that manifests with actions
// are parsed and validated correctly, including type checking.
func TestFoundationManifestActionsParsing(t *testing.T) {
	m := Manifest{
		Name:          "test-module",
		Version:       "1.0.0",
		Description:   "Test module",
		DescriptionES: "Módulo de prueba",
		Category:      "system",
		MinHubVersion: "1.0.0",
		Actions: []ActionConfig{
			{Name: "toggle-feature", Label: "Feature", Type: "toggle", Field: "test.feature"},
			{Name: "select-theme", Label: "Theme", Type: "select", Field: "test.theme", Options: []string{"a", "b"}},
			{Name: "confirm-delete", Label: "Delete", Type: "confirm", Field: ""},
			{Name: "execute-action", Label: "Run", Type: "execute", Field: ""},
			{Name: "text-input", Label: "Input", Type: "text", Field: "test.value"},
			{Name: "list-items", Label: "Items", Type: "list", Field: "test.items", Options: []string{"x", "y"}},
		},
	}

	if err := m.Validate(); err != nil {
		t.Fatalf("valid manifest with all action types: unexpected error: %v", err)
	}

	// Invalid type should fail.
	m2 := Manifest{
		Name:          "bad",
		Version:       "1.0.0",
		Description:   "Bad",
		DescriptionES: "Malo",
		Category:      "system",
		MinHubVersion: "1.0.0",
		Actions: []ActionConfig{
			{Name: "bad", Label: "Bad", Type: "slider", Field: ""},
		},
	}
	if err := m2.Validate(); err == nil {
		t.Fatal("invalid action type: expected error, got nil")
	}
}

// TestFoundationMockExecutorPattern verifies the MockExecutor pattern
// works for testing modules without real CLI dependencies.
func TestFoundationMockExecutorPattern(t *testing.T) {
	mock := &MockExecutor{
		Responses: map[string]MockResponse{
			"setxkbmap us": {
				Stdout:   "",
				Stderr:   "",
				ExitCode: 0,
				Err:      nil,
			},
			"gsettings set org.gnome.desktop.interface gtk-theme Dracula": {
				Stdout:   "",
				Stderr:   "",
				ExitCode: 0,
				Err:      nil,
			},
		},
	}

	stdout, stderr, exitCode, err := mock.Run("setxkbmap", "us")
	if err != nil {
		t.Fatalf("mock setxkbmap: unexpected error: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0", exitCode)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("unexpected output: stdout=%q stderr=%q", stdout, stderr)
	}

	stdout, stderr, exitCode, err = mock.Run("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "Dracula")
	if err != nil {
		t.Fatalf("mock gsettings: unexpected error: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0", exitCode)
	}
}
