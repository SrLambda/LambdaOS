package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHelpInitialState(t *testing.T) {
	bindings := []KeyBinding{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "enter", Desc: "Select"},
	}
	h := NewHelp(bindings)
	if len(h.Bindings) != 2 {
		t.Fatalf("Bindings len = %d, want 2", len(h.Bindings))
	}
	if !h.Visible {
		t.Error("Visible = false, want true")
	}
}

func TestHelpDismissOnEsc(t *testing.T) {
	bindings := []KeyBinding{{Key: "q", Desc: "Quit"}}
	h := NewHelp(bindings)

	updated, cmd := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updated.Visible {
		t.Error("Visible after esc = true, want false")
	}
	if cmd == nil {
		t.Fatal("expected cmd after dismiss, got nil")
	}

	msg := cmd()
	_, ok := msg.(HelpDismissedMsg)
	if !ok {
		t.Fatalf("expected HelpDismissedMsg, got %T", msg)
	}
}

func TestHelpDismissOnQuestionMark(t *testing.T) {
	bindings := []KeyBinding{{Key: "q", Desc: "Quit"}}
	h := NewHelp(bindings)

	updated, _ := h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if updated.Visible {
		t.Error("Visible after ? = true, want false")
	}
}

func TestHelpOtherKeysIgnored(t *testing.T) {
	bindings := []KeyBinding{{Key: "q", Desc: "Quit"}}
	h := NewHelp(bindings)

	updated, _ := h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if !updated.Visible {
		t.Error("Visible after 'a' = false, want true")
	}
}

func TestHelpViewShowsBindings(t *testing.T) {
	bindings := []KeyBinding{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "enter", Desc: "Select"},
		{Key: "esc", Desc: "Back"},
	}
	h := NewHelp(bindings)

	view := h.View()
	if !strings.Contains(view, "Navigate") {
		t.Errorf("view = %q, want to contain 'Navigate'", view)
	}
	if !strings.Contains(view, "Select") {
		t.Errorf("view = %q, want to contain 'Select'", view)
	}
	if !strings.Contains(view, "Back") {
		t.Errorf("view = %q, want to contain 'Back'", view)
	}
}

func TestHelpViewShowsDismissHint(t *testing.T) {
	bindings := []KeyBinding{{Key: "q", Desc: "Quit"}}
	h := NewHelp(bindings)

	view := h.View()
	if !strings.Contains(view, "esc") && !strings.Contains(view, "?") {
		t.Errorf("view = %q, want to contain dismiss hint", view)
	}
}

func TestHelpEmptyBindings(t *testing.T) {
	h := NewHelp([]KeyBinding{})
	view := h.View()
	if view == "" {
		t.Error("view should not be empty even with no bindings")
	}
}

func TestHelpToggleVisibility(t *testing.T) {
	bindings := []KeyBinding{{Key: "q", Desc: "Quit"}}
	h := NewHelp(bindings)

	// Dismiss
	updated, _ := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updated.Visible {
		t.Error("Visible after first esc = true, want false")
	}

	// Toggle back on with ?
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !updated.Visible {
		t.Error("Visible after ? = false, want true")
	}
}
