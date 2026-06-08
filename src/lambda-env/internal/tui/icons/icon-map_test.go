package icons

import (
	"encoding/json"
	"testing"
)

func TestIconMapJSONStructure(t *testing.T) {
	var data map[string]map[string]struct {
		Nerd     string `json:"nerd"`
		Fallback string `json:"fallback"`
	}

	if err := json.Unmarshal(iconMapJSON, &data); err != nil {
		t.Fatalf("failed to unmarshal icon-map.json: %v", err)
	}

	// Verify top-level categories exist.
	requiredCategories := []string{"categories", "modules", "widgets"}
	for _, cat := range requiredCategories {
		if _, ok := data[cat]; !ok {
			t.Errorf("missing top-level category %q", cat)
		}
	}

	// Verify modules.
	requiredModules := []string{
		"display", "audio", "network", "bluetooth",
		"security", "neovim", "qtile", "dotfiles", "logs", "storage",
	}
	for _, mod := range requiredModules {
		entry, ok := data["modules"][mod]
		if !ok {
			t.Errorf("missing module %q", mod)
			continue
		}
		if entry.Nerd == "" {
			t.Errorf("module %q has empty nerd glyph", mod)
		}
		if entry.Fallback == "" {
			t.Errorf("module %q has empty fallback glyph", mod)
		}
	}

	// Verify widgets.
	requiredWidgets := []string{
		"toggle_on", "toggle_off", "loading", "success",
		"error", "warning", "search", "confirm", "lock",
	}
	for _, w := range requiredWidgets {
		entry, ok := data["widgets"][w]
		if !ok {
			t.Errorf("missing widget %q", w)
			continue
		}
		if entry.Nerd == "" {
			t.Errorf("widget %q has empty nerd glyph", w)
		}
		if entry.Fallback == "" {
			t.Errorf("widget %q has empty fallback glyph", w)
		}
	}

	// Verify categories.
	requiredCats := []string{"system", "apps", "ops"}
	for _, c := range requiredCats {
		entry, ok := data["categories"][c]
		if !ok {
			t.Errorf("missing category %q", c)
			continue
		}
		if entry.Nerd == "" {
			t.Errorf("category %q has empty nerd glyph", c)
		}
		if entry.Fallback == "" {
			t.Errorf("category %q has empty fallback glyph", c)
		}
	}
}
