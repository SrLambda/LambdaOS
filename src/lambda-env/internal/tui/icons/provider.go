package icons

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed icon-map.json
var iconMapJSON []byte

// IconProvider resolves icon glyphs for the active font mode.
type IconProvider interface {
	Get(key string) string
	ForCategory(catID string) string
	ForModule(modID string) string
	ForWidget(widgetType string) string
	Width() int
}

// provider is the default implementation of IconProvider.
type provider struct {
	nerdFonts bool
	data      map[string]map[string]iconEntry
}

type iconEntry struct {
	Nerd     string `json:"nerd"`
	Fallback string `json:"fallback"`
}

// NewProvider creates an IconProvider that resolves icons based on nerdFonts.
func NewProvider(nerdFonts bool) IconProvider {
	var data map[string]map[string]iconEntry
	if err := json.Unmarshal(iconMapJSON, &data); err != nil {
		data = make(map[string]map[string]iconEntry)
	}
	return &provider{nerdFonts: nerdFonts, data: data}
}

// Get resolves a dot-separated key (e.g. "modules.display") to an icon glyph.
// If the key is missing, it returns a safe default (·).
func (p *provider) Get(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "\u00b7"
	}
	cat, ok := p.data[parts[0]]
	if !ok {
		return "\u00b7"
	}
	entry, ok := cat[parts[1]]
	if !ok {
		return "\u00b7"
	}
	if p.nerdFonts {
		return entry.Nerd
	}
	return entry.Fallback
}

// ForModule resolves an icon for the given module ID.
func (p *provider) ForModule(modID string) string {
	return p.Get("modules." + modID)
}

// ForWidget resolves an icon for the given widget type.
func (p *provider) ForWidget(widgetType string) string {
	return p.Get("widgets." + widgetType)
}

// ForCategory resolves an icon for the given category ID.
func (p *provider) ForCategory(catID string) string {
	return p.Get("categories." + catID)
}

// Width returns the display cell width for the active mode.
// Nerd Font glyphs are typically 2 cells wide; Unicode fallbacks are 1.
func (p *provider) Width() int {
	if p.nerdFonts {
		return 2
	}
	return 1
}

// Default is the package-level provider initialized at startup.
// It auto-detects Nerd Font support via the LAMBDA_NERD_FONTS env var.
var Default IconProvider = NewProvider(DetectNerdFonts())
