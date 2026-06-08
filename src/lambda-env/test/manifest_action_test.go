package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"
)

// moduleActionExpectations maps each module to the actions it should support.
// This is the contract between manifest.json and main.go switch cases.
var moduleActionExpectations = map[string][]struct {
	Name string
	Type string
}{
	"audio": {
		{Name: "run", Type: "execute"},
		{Name: "set-volume", Type: "text"},
		{Name: "set-mute", Type: "toggle"},
		{Name: "set-sink", Type: "select"},
		{Name: "set-source", Type: "select"},
		{Name: "set-profile", Type: "select"},
		{Name: "set-app-volume", Type: "select"},
	},
	"keyboard": {
		{Name: "run", Type: "execute"},
		{Name: "set-layout", Type: "select"},
		{Name: "set-variant", Type: "select"},
		{Name: "set-compose", Type: "select"},
		{Name: "set-options", Type: "select"},
	},
	"appearance": {
		{Name: "run", Type: "execute"},
		{Name: "set-theme", Type: "select"},
		{Name: "set-wallpaper", Type: "text"},
		{Name: "set-font-size", Type: "text"},
	},
	"defaults": {
		{Name: "run", Type: "execute"},
		{Name: "set-browser", Type: "select"},
		{Name: "set-terminal", Type: "select"},
		{Name: "set-editor", Type: "select"},
		{Name: "set-file-manager", Type: "select"},
		{Name: "apply", Type: "confirm"},
	},
	"dotfiles": {
		{Name: "stow", Type: "execute"},
		{Name: "unstow", Type: "confirm"},
		{Name: "backup", Type: "execute"},
	},
	"neovim": {
		{Name: "toggle-lsp", Type: "toggle"},
		{Name: "toggle-copilot", Type: "toggle"},
		{Name: "toggle-neotree", Type: "toggle"},
		{Name: "set-theme", Type: "select"},
		{Name: "apply", Type: "execute"},
	},
	"qtile": {
		{Name: "set-terminal", Type: "select"},
		{Name: "set-browser", Type: "select"},
		{Name: "set-file-manager", Type: "select"},
		{Name: "reload", Type: "execute"},
	},
	"power": {
		{Name: "run", Type: "execute"},
		{Name: "set-screen-timeout", Type: "text"},
		{Name: "set-sleep-timeout", Type: "text"},
		{Name: "set-lid-close-action", Type: "select"},
	},
	"display": {
		{Name: "run", Type: "execute"},
		{Name: "set-mode", Type: "select"},
		{Name: "set-position", Type: "text"},
		{Name: "set-primary", Type: "toggle"},
		{Name: "save-profile", Type: "text"},
		{Name: "load-profile", Type: "select"},
	},
	"hardware-dashboard": {
		{Name: "run", Type: "execute"},
	},
}

func loadModuleManifest(t *testing.T, name string) module.Manifest {
	t.Helper()
	modPath := filepath.Join("..", "internal", "modules", name, "manifest.json")
	data, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatalf("read manifest for %s: %v", name, err)
	}
	var m module.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse manifest for %s: %v", name, err)
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("validate manifest for %s: %v", name, err)
	}
	return m
}

func TestAllModuleManifestsAreValid(t *testing.T) {
	for name := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)
			if m.Name != name {
				t.Errorf("manifest name = %q, want %q", m.Name, name)
			}
			if m.Category == "" {
				t.Errorf("manifest category is empty")
			}
			if m.Version == "" {
				t.Errorf("manifest version is empty")
			}
		})
	}
}

func TestManifestActionsMatchExpected(t *testing.T) {
	for name, expected := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)

			if len(m.Actions) == 0 {
				t.Fatalf("module %s has no actions in manifest", name)
			}

			if len(m.Actions) != len(expected) {
				t.Errorf("action count = %d, want %d", len(m.Actions), len(expected))
			}

			actionMap := make(map[string]string, len(m.Actions))
			for _, a := range m.Actions {
				actionMap[a.Name] = a.Type
			}

			for _, exp := range expected {
				actualType, ok := actionMap[exp.Name]
				if !ok {
					t.Errorf("action %q missing from manifest", exp.Name)
					continue
				}
				if actualType != exp.Type {
					t.Errorf("action %q type = %q, want %q", exp.Name, actualType, exp.Type)
				}
			}
		})
	}
}

func TestManifestActionTypesAreValid(t *testing.T) {
	validTypes := map[string]bool{
		"toggle":  true,
		"select":  true,
		"text":    true,
		"list":    true,
		"confirm": true,
		"execute": true,
	}

	for name := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)
			for _, a := range m.Actions {
				if !validTypes[a.Type] {
					t.Errorf("action %q has invalid type %q", a.Name, a.Type)
				}
			}
		})
	}
}

func TestManifestSelectActionsHaveOptions(t *testing.T) {
	for name := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)
			for _, a := range m.Actions {
				if (a.Type == "select" || a.Type == "list") && len(a.Options) == 0 {
					t.Errorf("action %q (%s) has no options", a.Name, a.Type)
				}
			}
		})
	}
}

func TestDetailViewRendersAllModules(t *testing.T) {
	for name := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)
			dv := views.NewDetailView(m, icons.NewProvider(false))
			rendered := dv.View()

			if rendered == "" {
				t.Fatal("detail view rendered empty string")
			}

			if !strings.Contains(rendered, m.Name) {
				t.Errorf("view missing module name %q", m.Name)
			}

			for _, a := range m.Actions {
				if !strings.Contains(rendered, a.Label) {
					t.Errorf("view missing action label %q", a.Label)
				}
			}
		})
	}
}

func TestDetailViewHandlesEmptyActionsGracefully(t *testing.T) {
	m := module.Manifest{
		Name:    "empty-module",
		Actions: []module.ActionConfig{},
	}
	dv := views.NewDetailView(m, icons.NewProvider(false))
	rendered := dv.View()
	if !strings.Contains(rendered, "No actions available") {
		t.Errorf("view = %q, want to contain 'No actions available'", rendered)
	}
}

func TestManifestActionsMatchMainGoSwitchCases(t *testing.T) {
	// For each module, read main.go and verify the switch cases mention
	// every action name defined in the manifest. This catches the
	// hyphen/underscore mismatch class of bugs.
	for name := range moduleActionExpectations {
		t.Run(name, func(t *testing.T) {
			m := loadModuleManifest(t, name)
			mainPath := filepath.Join("..", "internal", "modules", name, "main.go")
			src, err := os.ReadFile(mainPath)
			if err != nil {
				t.Fatalf("read main.go for %s: %v", name, err)
			}
			srcStr := string(src)

			// Find the action switch block.
			if !strings.Contains(srcStr, "switch action") {
				t.Fatalf("main.go for %s missing 'switch action'", name)
			}

			for _, a := range m.Actions {
				// Each manifest action should appear in a case statement.
				casePattern := `case "` + a.Name + `"`
				if !strings.Contains(srcStr, casePattern) {
					t.Errorf("manifest action %q (%s) not found in main.go switch cases", a.Name, a.Type)
				}
			}
		})
	}
}
