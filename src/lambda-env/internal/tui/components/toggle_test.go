package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestToggleInitialState(t *testing.T) {
	toggle := NewToggle("Test Feature", false)
	if toggle.Value != false {
		t.Errorf("initial Value = %v, want false", toggle.Value)
	}
	if toggle.Label != "Test Feature" {
		t.Errorf("initial Label = %q, want %q", toggle.Label, "Test Feature")
	}
	if toggle.Focused != false {
		t.Errorf("initial Focused = %v, want false", toggle.Focused)
	}
}

func TestToggleSpaceFlipsState(t *testing.T) {
	toggle := NewToggle("Feature", false)

	updated, cmd := toggle.Update(tea.KeyMsg{Type: tea.KeySpace})

	if updated.Value != true {
		t.Errorf("Value after space = %v, want true", updated.Value)
	}
	if cmd == nil {
		t.Error("expected a cmd after toggle change, got nil")
	}
}

func TestToggleEnterFlipsState(t *testing.T) {
	toggle := NewToggle("Feature", true)

	updated, cmd := toggle.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updated.Value != false {
		t.Errorf("Value after enter = %v, want false", updated.Value)
	}
	if cmd == nil {
		t.Error("expected a cmd after toggle change, got nil")
	}
}

func TestToggleOtherKeysIgnored(t *testing.T) {
	toggle := NewToggle("Feature", false)

	updated, _ := toggle.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	if updated.Value != false {
		t.Errorf("Value after 'a' = %v, want false", updated.Value)
	}
}

func TestToggleEmitsChangedMessage(t *testing.T) {
	toggle := NewToggle("Feature", false)

	_, cmd := toggle.Update(tea.KeyMsg{Type: tea.KeySpace})

	if cmd == nil {
		t.Fatal("expected cmd, got nil")
	}

	msg := cmd()
	changed, ok := msg.(ToggleChangedMsg)
	if !ok {
		t.Fatalf("expected ToggleChangedMsg, got %T", msg)
	}
	if changed.Value != true {
		t.Errorf("changed msg Value = %v, want true", changed.Value)
	}
	if changed.Label != "Feature" {
		t.Errorf("changed msg Label = %q, want %q", changed.Label, "Feature")
	}
}

func TestToggleViewShowsState(t *testing.T) {
	on := NewToggle("Wifi", true)
	off := NewToggle("Wifi", false)

	onView := on.View()
	offView := off.View()

	if !strings.Contains(onView, "On") {
		t.Errorf("on view = %q, want to contain 'On'", onView)
	}
	if !strings.Contains(offView, "Off") {
		t.Errorf("off view = %q, want to contain 'Off'", offView)
	}
}

func TestToggleFocusState(t *testing.T) {
	toggle := NewToggle("Feature", false)
	toggle.Focused = true

	view := toggle.View()
	if !strings.Contains(view, "Feature") {
		t.Errorf("focused view = %q, want to contain label", view)
	}
}

func TestToggleDoubleToggle(t *testing.T) {
	toggle := NewToggle("Feature", false)

	updated, _ := toggle.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeySpace})

	if updated.Value != false {
		t.Errorf("Value after double toggle = %v, want false", updated.Value)
	}
}
