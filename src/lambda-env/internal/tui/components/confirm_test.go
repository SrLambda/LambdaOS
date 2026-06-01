package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConfirmInitialState(t *testing.T) {
	c := NewConfirm("Are you sure?")
	if c.Message != "Are you sure?" {
		t.Errorf("Message = %q, want %q", c.Message, "Are you sure?")
	}
	if c.Selected != 0 {
		t.Errorf("Selected = %d, want 0 (Yes)", c.Selected)
	}
	if c.Confirmed != false {
		t.Errorf("Confirmed = %v, want false", c.Confirmed)
	}
}

func TestConfirmLeftRightNavigation(t *testing.T) {
	c := NewConfirm("Delete file?")

	updated, _ := c.Update(tea.KeyMsg{Type: tea.KeyRight})
	if updated.Selected != 1 {
		t.Errorf("Selected after right = %d, want 1 (No)", updated.Selected)
	}

	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if updated.Selected != 0 {
		t.Errorf("Selected after left = %d, want 0 (Yes)", updated.Selected)
	}
}

func TestConfirmEnterSelectsYes(t *testing.T) {
	c := NewConfirm("Proceed?")

	updated, cmd := c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !updated.Confirmed {
		t.Errorf("Confirmed after enter on Yes = %v, want true", updated.Confirmed)
	}
	if cmd == nil {
		t.Fatal("expected cmd after confirm, got nil")
	}

	msg := cmd()
	result, ok := msg.(ConfirmResultMsg)
	if !ok {
		t.Fatalf("expected ConfirmResultMsg, got %T", msg)
	}
	if !result.Confirmed {
		t.Errorf("result Confirmed = %v, want true", result.Confirmed)
	}
}

func TestConfirmEnterSelectsNo(t *testing.T) {
	c := NewConfirm("Proceed?")
	c.Selected = 1

	updated, cmd := c.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updated.Confirmed {
		t.Errorf("Confirmed after enter on No = %v, want false", updated.Confirmed)
	}
	if cmd == nil {
		t.Fatal("expected cmd after confirm, got nil")
	}

	msg := cmd()
	result, ok := msg.(ConfirmResultMsg)
	if !ok {
		t.Fatalf("expected ConfirmResultMsg, got %T", msg)
	}
	if result.Confirmed {
		t.Errorf("result Confirmed = %v, want false", result.Confirmed)
	}
}

func TestConfirmSpaceTogglesSelection(t *testing.T) {
	c := NewConfirm("Toggle?")

	updated, _ := c.Update(tea.KeyMsg{Type: tea.KeySpace})
	if updated.Selected != 1 {
		t.Errorf("Selected after space = %d, want 1", updated.Selected)
	}

	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeySpace})
	if updated.Selected != 0 {
		t.Errorf("Selected after second space = %d, want 0", updated.Selected)
	}
}

func TestConfirmViewShowsMessageAndButtons(t *testing.T) {
	c := NewConfirm("Delete everything?")

	view := c.View()
	if !strings.Contains(view, "Delete everything?") {
		t.Errorf("view = %q, want to contain message", view)
	}
	if !strings.Contains(view, "Yes") {
		t.Errorf("view = %q, want to contain 'Yes'", view)
	}
	if !strings.Contains(view, "No") {
		t.Errorf("view = %q, want to contain 'No'", view)
	}
}

func TestConfirmViewShowsSelectedButton(t *testing.T) {
	c := NewConfirm("Sure?")
	c.Selected = 1

	view := c.View()
	if !strings.Contains(view, "No") {
		t.Errorf("view = %q, want to contain 'No'", view)
	}
}

func TestConfirmEscCancels(t *testing.T) {
	c := NewConfirm("Cancel me?")

	updated, cmd := c.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if updated.Confirmed {
		t.Errorf("Confirmed after esc = %v, want false", updated.Confirmed)
	}
	if cmd == nil {
		t.Fatal("expected cmd after esc, got nil")
	}

	msg := cmd()
	_, ok := msg.(ConfirmResultMsg)
	if !ok {
		t.Fatalf("expected ConfirmResultMsg, got %T", msg)
	}
}

func TestConfirmNavigationWraps(t *testing.T) {
	c := NewConfirm("Wrap?")
	c.Selected = 1

	updated, _ := c.Update(tea.KeyMsg{Type: tea.KeyRight})
	if updated.Selected != 0 {
		t.Errorf("Selected after wrap right = %d, want 0", updated.Selected)
	}

	c.Selected = 0
	updated, _ = c.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if updated.Selected != 1 {
		t.Errorf("Selected after wrap left = %d, want 1", updated.Selected)
	}
}
