package module

import (
	"fmt"
	"regexp"
)

// ValidCategories lists all accepted module category values.
var ValidCategories = []string{"system", "apps", "ops", "setup"}

// ValidActionTypes lists all accepted action type values.
var ValidActionTypes = []string{"toggle", "select", "text", "list", "confirm", "execute"}

var nameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ActionConfig defines a single UI action exposed by a module.
type ActionConfig struct {
	Name         string   `json:"name"`
	Label        string   `json:"label"`
	Type         string   `json:"type"`
	Field        string   `json:"field"`
	Options      []string `json:"options,omitempty"`
	RequiresRoot bool     `json:"requires_root,omitempty"`
}

// Manifest represents a module's manifest.json metadata.
type Manifest struct {
	Name          string         `json:"name"`
	Version       string         `json:"version"`
	Description   string         `json:"description"`
	DescriptionES string         `json:"description_es"`
	Category      string         `json:"category"`
	Icon          string         `json:"icon,omitempty"`
	RequiresRoot  bool           `json:"requires_root"`
	Dependencies  []string       `json:"dependencies"`
	MinHubVersion string         `json:"min_hub_version"`
	Timeout       int            `json:"timeout"`
	Tags          []string       `json:"tags,omitempty"`
	Author        string         `json:"author,omitempty"`
	Actions       []ActionConfig `json:"actions,omitempty"`
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

	for _, a := range m.Actions {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("action %q: %w", a.Name, err)
		}
	}

	return nil
}

// Validate checks that the action config has valid fields.
func (a *ActionConfig) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("action name is required")
	}
	if a.Label == "" {
		return fmt.Errorf("action label is required")
	}
	if a.Type == "" {
		return fmt.Errorf("action type is required")
	}
	if !isValidActionType(a.Type) {
		return fmt.Errorf("invalid action type %q; must be one of: %v", a.Type, ValidActionTypes)
	}
	if (a.Type == "select" || a.Type == "list") && len(a.Options) == 0 {
		return fmt.Errorf("%s action must have options", a.Type)
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

func isValidActionType(t string) bool {
	for _, v := range ValidActionTypes {
		if v == t {
			return true
		}
	}
	return false
}
