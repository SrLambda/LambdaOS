package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestToggleStateTransitions verifies rapid state transitions.
func TestToggleStateTransitions(t *testing.T) {
	toggle := NewToggle("Feature", false)

	states := []bool{}
	for i := 0; i < 5; i++ {
		updated, _ := toggle.Update(tea.KeyMsg{Type: tea.KeySpace})
		states = append(states, updated.Value)
		toggle = updated
	}

	expected := []bool{true, false, true, false, true}
	for i, exp := range expected {
		if states[i] != exp {
			t.Errorf("state[%d] = %v, want %v", i, states[i], exp)
		}
	}
}

func TestToggleNoMessageOnOtherKeys(t *testing.T) {
	toggle := NewToggle("Feature", false)

	_, cmd := toggle.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if cmd != nil {
		msg := cmd()
		if msg != nil {
			t.Errorf("unexpected message on 'x': %T", msg)
		}
	}
}

func TestToggleEmptyLabel(t *testing.T) {
	toggle := NewToggle("", false)
	view := toggle.View()
	if view == "" {
		t.Error("view should not be empty even with empty label")
	}
}

// TestTextInputStateTransitions verifies typing, submit, and cancel flow.
func TestTextInputStateTransitions(t *testing.T) {
	input := NewTextInput("Name")

	// Type "abc"
	for _, r := range []rune{'a', 'b', 'c'} {
		updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		input = updated
	}
	if input.Value() != "abc" {
		t.Errorf("Value after typing = %q, want %q", input.Value(), "abc")
	}

	// Submit
	updated, cmd := input.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.Err != nil {
		t.Errorf("unexpected error on submit: %v", updated.Err)
	}
	msg := cmd()
	submit, ok := msg.(TextInputSubmitMsg)
	if !ok {
		t.Fatalf("expected TextInputSubmitMsg, got %T", msg)
	}
	if submit.Value != "abc" {
		t.Errorf("submit value = %q, want %q", submit.Value, "abc")
	}

	// Cancel should clear value
	updated, _ = input.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updated.Value() != "" {
		t.Errorf("Value after cancel = %q, want empty", updated.Value())
	}
}

func TestTextInputBackspaceRemovesChar(t *testing.T) {
	input := NewTextInput("Code")
	input.SetValue("ab")

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if updated.Value() != "a" {
		t.Errorf("Value after backspace = %q, want %q", updated.Value(), "a")
	}
}

func TestTextInputEmptySubmit(t *testing.T) {
	input := NewTextInput("Required")
	input.SetRegex(`.+`) // non-empty

	updated, cmd := input.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.Err == nil {
		t.Error("expected error on empty submit, got nil")
	}
	if cmd != nil {
		msg := cmd()
		if msg != nil {
			t.Errorf("unexpected message on failed submit: %T", msg)
		}
	}
}

// TestConfirmStateTransitions verifies navigation and selection flow.
func TestConfirmStateTransitions(t *testing.T) {
	c := NewConfirm("Proceed?")

	// Navigate to No
	updated, _ := c.Update(tea.KeyMsg{Type: tea.KeyRight})
	if updated.Selected != 1 {
		t.Fatalf("Selected = %d, want 1", updated.Selected)
	}

	// Navigate back to Yes
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if updated.Selected != 0 {
		t.Fatalf("Selected = %d, want 0", updated.Selected)
	}

	// Select
	updated, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !updated.Confirmed {
		t.Error("Confirmed = false, want true")
	}
	msg := cmd()
	result := msg.(ConfirmResultMsg)
	if !result.Confirmed {
		t.Error("result.Confirmed = false, want true")
	}
}

func TestConfirmRapidNavigation(t *testing.T) {
	c := NewConfirm("Fast?")

	// Rapid right presses
	for i := 0; i < 10; i++ {
		updated, _ := c.Update(tea.KeyMsg{Type: tea.KeyRight})
		c = updated
	}
	if c.Selected != 0 {
		t.Errorf("Selected after 10 right presses = %d, want 0", c.Selected)
	}
}

// TestHelpStateTransitions verifies visibility toggling.
func TestHelpStateTransitions(t *testing.T) {
	h := NewHelp([]KeyBinding{
		{Key: "a", Desc: "Action A"},
		{Key: "b", Desc: "Action B"},
	})

	if !h.Visible {
		t.Fatal("initial Visible = false, want true")
	}

	// Dismiss with esc
	updated, _ := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updated.Visible {
		t.Error("Visible after esc = true, want false")
	}

	// Toggle back with ?
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !updated.Visible {
		t.Error("Visible after ? = false, want true")
	}

	// Toggle off again with ?
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if updated.Visible {
		t.Error("Visible after second ? = true, want false")
	}
}

func TestHelpViewWhenInvisible(t *testing.T) {
	h := NewHelp([]KeyBinding{{Key: "q", Desc: "Quit"}})
	h.Visible = false
	view := h.View()
	if view != "" {
		t.Errorf("view when invisible = %q, want empty", view)
	}
}

// TestStatusBarNoKeyboardInteraction verifies status bar ignores keys.
func TestStatusBarNoKeyboardInteraction(t *testing.T) {
	sb := NewStatusBar()
	sb.SetContext("test")

	// StatusBar has no Update method, so nothing to test for keyboard interaction.
	// This test documents that status bar is purely presentational.
	view := sb.View()
	if view == "" {
		t.Error("view should not be empty")
	}
}
