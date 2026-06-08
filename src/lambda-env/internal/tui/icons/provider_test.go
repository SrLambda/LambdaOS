package icons

import (
	"testing"
)

func TestProviderGetNerdMode(t *testing.T) {
	p := NewProvider(true)

	tests := []struct {
		key  string
		want string
	}{
		{"modules.display", "\uf26c"},
		{"modules.audio", "\uf028"},
		{"widgets.toggle_on", "\uf205"},
		{"categories.system", "\uf108"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := p.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestProviderGetFallbackMode(t *testing.T) {
	p := NewProvider(false)

	tests := []struct {
		key  string
		want string
	}{
		{"modules.display", "\u25a3"},
		{"modules.audio", "\u266a"},
		{"widgets.toggle_on", "\u25cf"},
		{"categories.system", "\u2699"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := p.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestProviderGetMissingKey(t *testing.T) {
	p := NewProvider(true)

	got := p.Get("modules.nonexistent")
	if got != "\u00b7" {
		t.Errorf("Get(missing) = %q, want default \u00b7", got)
	}
}

func TestProviderGetMalformedKey(t *testing.T) {
	p := NewProvider(true)

	got := p.Get("badkey")
	if got != "\u00b7" {
		t.Errorf("Get(badkey) = %q, want default \u00b7", got)
	}
}

func TestProviderConvenienceMethods(t *testing.T) {
	p := NewProvider(true).(*provider)

	if got := p.ForModule("display"); got != "\uf26c" {
		t.Errorf("ForModule(display) = %q, want \uf26c", got)
	}
	if got := p.ForWidget("toggle_on"); got != "\uf205" {
		t.Errorf("ForWidget(toggle_on) = %q, want \uf205", got)
	}
	if got := p.ForCategory("system"); got != "\uf108" {
		t.Errorf("ForCategory(system) = %q, want \uf108", got)
	}
}

func TestProviderWidth(t *testing.T) {
	nerd := NewProvider(true).(*provider)
	fallback := NewProvider(false).(*provider)

	if nerd.Width() != 2 {
		t.Errorf("nerd Width = %d, want 2", nerd.Width())
	}
	if fallback.Width() != 1 {
		t.Errorf("fallback Width = %d, want 1", fallback.Width())
	}
}
