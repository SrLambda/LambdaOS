package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStatusBarInitialState(t *testing.T) {
	sb := NewStatusBar()
	if sb.Context != "" {
		t.Errorf("Context = %q, want empty", sb.Context)
	}
	if sb.Module != "" {
		t.Errorf("Module = %q, want empty", sb.Module)
	}
	if sb.SettingsState != "" {
		t.Errorf("SettingsState = %q, want empty", sb.SettingsState)
	}
	if sb.Modified {
		t.Error("Modified = true, want false")
	}
}

func TestStatusBarSetters(t *testing.T) {
	sb := NewStatusBar()
	sb.SetContext("categories").
		SetModule("keyboard").
		SetSettingsState("unsaved").
		SetModified(true)

	if sb.Context != "categories" {
		t.Errorf("Context = %q, want %q", sb.Context, "categories")
	}
	if sb.Module != "keyboard" {
		t.Errorf("Module = %q, want %q", sb.Module, "keyboard")
	}
	if sb.SettingsState != "unsaved" {
		t.Errorf("SettingsState = %q, want %q", sb.SettingsState, "unsaved")
	}
	if !sb.Modified {
		t.Error("Modified = false, want true")
	}
}

func TestStatusBarViewShowsContext(t *testing.T) {
	sb := NewStatusBar()
	sb.SetContext("modules")

	view := sb.View()
	if !strings.Contains(view, "modules") {
		t.Errorf("view = %q, want to contain 'modules'", view)
	}
}

func TestStatusBarViewShowsModule(t *testing.T) {
	sb := NewStatusBar()
	sb.SetModule("audio")

	view := sb.View()
	if !strings.Contains(view, "audio") {
		t.Errorf("view = %q, want to contain 'audio'", view)
	}
}

func TestStatusBarViewShowsSettingsState(t *testing.T) {
	sb := NewStatusBar()
	sb.SetSettingsState("saved")

	view := sb.View()
	if !strings.Contains(view, "saved") {
		t.Errorf("view = %q, want to contain 'saved'", view)
	}
}

func TestStatusBarViewShowsModified(t *testing.T) {
	sb := NewStatusBar()
	sb.SetModified(true)

	view := sb.View()
	if !strings.Contains(view, "*") && !strings.Contains(view, "modified") {
		t.Errorf("view = %q, want to contain modified indicator", view)
	}
}

func TestStatusBarViewWithAllFields(t *testing.T) {
	sb := NewStatusBar()
	sb.SetContext("detail").
		SetModule("appearance").
		SetSettingsState("unsaved").
		SetModified(true)

	view := sb.View()
	if !strings.Contains(view, "detail") {
		t.Errorf("view = %q, want to contain 'detail'", view)
	}
	if !strings.Contains(view, "appearance") {
		t.Errorf("view = %q, want to contain 'appearance'", view)
	}
	if !strings.Contains(view, "unsaved") {
		t.Errorf("view = %q, want to contain 'unsaved'", view)
	}
}

func TestStatusBarEmptyView(t *testing.T) {
	sb := NewStatusBar()
	view := sb.View()
	if view == "" {
		t.Error("view should not be empty even with no fields set")
	}
}

func TestStatusBarChaining(t *testing.T) {
	sb := NewStatusBar().
		SetContext("categories").
		SetModule("").
		SetSettingsState("default").
		SetModified(false)

	if sb.Context != "categories" {
		t.Errorf("Context = %q, want %q", sb.Context, "categories")
	}
	if sb.Module != "" {
		t.Errorf("Module = %q, want empty", sb.Module)
	}
	if sb.SettingsState != "default" {
		t.Errorf("SettingsState = %q, want %q", sb.SettingsState, "default")
	}
	if sb.Modified {
		t.Error("Modified = true, want false")
	}
}

func TestStatusBarDynamicWidth(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		wantWidth  int
		setContext bool
	}{
		{
			name:       "width 120 spans full terminal",
			width:      120,
			wantWidth:  120,
			setContext: true,
		},
		{
			name:       "width 40 adapts to narrow terminal",
			width:      40,
			wantWidth:  40,
			setContext: true,
		},
		{
			name:       "zero width falls back to default 80",
			width:      0,
			wantWidth:  80,
			setContext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatusBar()
			if tt.setContext {
				sb.SetContext("test")
			}
			sb.SetWidth(tt.width)

			view := sb.View()
			got := lipgloss.Width(view)
			if got != tt.wantWidth {
				t.Errorf("lipgloss.Width(view) = %d, want %d", got, tt.wantWidth)
			}
		})
	}
}
