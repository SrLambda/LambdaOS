package hub

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanPopulatesManifestIcon(t *testing.T) {
	// Create a temporary module directory with a valid manifest.
	tmpDir := t.TempDir()
	moduleDir := filepath.Join(tmpDir, "test-module")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	manifest := `{
		"name": "test-module",
		"version": "1.0.0",
		"description": "A test module",
		"description_es": "Un módulo de prueba",
		"category": "system",
		"min_hub_version": "0.1.0"
	}`
	manifestPath := filepath.Join(moduleDir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	// Use a mock icon provider that returns a known icon.
	mockProvider := &mockIconProvider{icon: "mock-icon"}

	// Replace module paths temporarily so only our test module is discovered.
	oldSystemPath := SystemModulesPath
	oldUserPath := UserModulesPath
	defer func() {
		SystemModulesPath = oldSystemPath
		UserModulesPath = oldUserPath
	}()
	SystemModulesPath = tmpDir
	UserModulesPath = "/nonexistent/user/path"

	manifests, err := Scan(mockProvider)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(manifests) != 1 {
		t.Fatalf("Scan() returned %d manifests, want 1", len(manifests))
	}

	m := manifests[0]
	if m.Icon == "" {
		t.Errorf("manifest.Icon = %q, want non-empty", m.Icon)
	}
	if m.Icon != "mock-icon" {
		t.Errorf("manifest.Icon = %q, want %q", m.Icon, "mock-icon")
	}
}

type mockIconProvider struct {
	icon string
}

func (m *mockIconProvider) Get(key string) string   { return m.icon }
func (m *mockIconProvider) ForCategory(catID string) string { return m.icon }
func (m *mockIconProvider) ForModule(modID string) string  { return m.icon }
func (m *mockIconProvider) ForWidget(widgetType string) string { return m.icon }
func (m *mockIconProvider) Width() int { return 1 }
