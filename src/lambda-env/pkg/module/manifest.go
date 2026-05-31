package module

import (
	"fmt"
	"regexp"
)

// ValidCategories lists all accepted module category values.
var ValidCategories = []string{"system", "apps", "ops", "setup"}

var nameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Manifest represents a module's manifest.json metadata.
type Manifest struct {
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Description   string   `json:"description"`
	DescriptionES string   `json:"description_es"`
	Category      string   `json:"category"`
	Icon          string   `json:"icon,omitempty"`
	RequiresRoot  bool     `json:"requires_root"`
	Dependencies  []string `json:"dependencies"`
	MinHubVersion string   `json:"min_hub_version"`
	Timeout       int      `json:"timeout"`
	Tags          []string `json:"tags,omitempty"`
	Author        string   `json:"author,omitempty"`
	// Path is the absolute directory containing this module's manifest and executable.
	// It is set by the discovery scanner and is not part of manifest.json.
	Path string `json:"-"`
}

// Response represents a module's JSON output on stdout.
type Response struct {
	Status        string                 `json:"status"`
	Action        string                 `json:"action"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Code          string                 `json:"code,omitempty"`
	Message       string                 `json:"message,omitempty"`
	MessageES     string                 `json:"message_es,omitempty"`
	Suggestion    string                 `json:"suggestion,omitempty"`
	SettingsDelta map[string]interface{} `json:"settings_delta,omitempty"`
}

// Validate checks that the manifest has all required fields and valid values.
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest name is required")
	}
	if !nameRegex.MatchString(m.Name) {
		return fmt.Errorf("manifest name %q must be lowercase with hyphens only", m.Name)
	}

	if m.Version == "" {
		return fmt.Errorf("manifest version is required")
	}

	if m.Description == "" {
		return fmt.Errorf("manifest description is required")
	}

	if m.DescriptionES == "" {
		return fmt.Errorf("manifest description_es is required")
	}

	if m.Category == "" {
		return fmt.Errorf("manifest category is required")
	}
	if !isValidCategory(m.Category) {
		return fmt.Errorf("manifest category %q is invalid; must be one of: %v", m.Category, ValidCategories)
	}

	if m.MinHubVersion == "" {
		return fmt.Errorf("manifest min_hub_version is required")
	}

	return nil
}

// Helper returns a human-readable summary of the manifest for display.
func (m *Manifest) Helper() string {
	rootHint := ""
	if m.RequiresRoot {
		rootHint = " [root]"
	}
	return fmt.Sprintf("%s%s — %s (%s)", m.Name, rootHint, m.Description, m.Category)
}

func isValidCategory(c string) bool {
	for _, v := range ValidCategories {
		if v == c {
			return true
		}
	}
	return false
}
