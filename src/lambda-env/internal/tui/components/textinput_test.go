package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextInputInitialState(t *testing.T) {
	input := NewTextInput("Hostname")
	if input.Label != "Hostname" {
		t.Errorf("initial Label = %q, want %q", input.Label, "Hostname")
	}
	if input.Value() != "" {
		t.Errorf("initial Value = %q, want empty", input.Value())
	}
	if input.Err != nil {
		t.Errorf("initial Err = %v, want nil", input.Err)
	}
}

func TestTextInputTyping(t *testing.T) {
	input := NewTextInput("Name")

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	if updated.Value() != "ab" {
		t.Errorf("Value after typing = %q, want %q", updated.Value(), "ab")
	}
}

func TestTextInputSubmitOnEnter(t *testing.T) {
	input := NewTextInput("Name")
	input.SetValue("test")

	updated, cmd := input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updated.Value() != "test" {
		t.Errorf("Value after enter = %q, want %q", updated.Value(), "test")
	}
	if cmd == nil {
		t.Fatal("expected cmd on enter, got nil")
	}

	msg := cmd()
	submit, ok := msg.(TextInputSubmitMsg)
	if !ok {
		t.Fatalf("expected TextInputSubmitMsg, got %T", msg)
	}
	if submit.Value != "test" {
		t.Errorf("submit msg Value = %q, want %q", submit.Value, "test")
	}
	if submit.Label != "Name" {
		t.Errorf("submit msg Label = %q, want %q", submit.Label, "Name")
	}
}

func TestTextInputCancelOnEsc(t *testing.T) {
	input := NewTextInput("Name")
	input.SetValue("test")

	updated, cmd := input.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if updated.Value() != "" {
		t.Errorf("Value after esc = %q, want empty", updated.Value())
	}
	if cmd == nil {
		t.Fatal("expected cmd on esc, got nil")
	}

	msg := cmd()
	_, ok := msg.(TextInputCancelMsg)
	if !ok {
		t.Fatalf("expected TextInputCancelMsg, got %T", msg)
	}
}

func TestTextInputAllowlistValidation(t *testing.T) {
	input := NewTextInput("Hex")
	input.SetAllowlist("0123456789abcdef")

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

	if updated.Value() != "" {
		t.Errorf("Value after invalid char = %q, want empty", updated.Value())
	}
	if updated.Err == nil {
		t.Error("expected validation error after invalid char, got nil")
	}
}

func TestTextInputRegexValidation(t *testing.T) {
	input := NewTextInput("Email")
	input.SetRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	input.SetValue("invalid-email")

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updated.Err == nil {
		t.Error("expected regex validation error, got nil")
	}
}

func TestTextInputMinMaxValidation(t *testing.T) {
	input := NewTextInput("Volume")
	input.SetNumericRange(0, 100)
	input.SetValue("150")

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updated.Err == nil {
		t.Error("expected max validation error, got nil")
	}

	input.SetValue("50")
	updated, _ = input.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated.Err != nil {
		t.Errorf("unexpected error for valid value: %v", updated.Err)
	}
}

func TestTextInputViewShowsError(t *testing.T) {
	input := NewTextInput("Volume")
	input.SetNumericRange(0, 100)
	input.SetValue("150")
	input.Update(tea.KeyMsg{Type: tea.KeyEnter})

	view := input.View()
	if !strings.Contains(view, "Error") && !strings.Contains(view, "error") && !strings.Contains(view, "must be") {
		t.Errorf("view = %q, want to contain error indication", view)
	}
}

func TestTextInputMaxLength(t *testing.T) {
	input := NewTextInput("Code")
	input.SetMaxLength(3)

	updated, _ := input.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	if updated.Value() != "abc" {
		t.Errorf("Value after 4 chars with max 3 = %q, want %q", updated.Value(), "abc")
	}
}

func TestTextInputPlaceholder(t *testing.T) {
	input := NewTextInput("Search")
	input.SetPlaceholder("type here...")

	view := input.View()
	if !strings.Contains(view, "type here...") {
		t.Errorf("view = %q, want to contain placeholder", view)
	}
}

func TestTextInputFocusBlur(t *testing.T) {
	input := NewTextInput("Name")
	input.SetValue("test")

	focused := input.Focus()
	if !focused.Focused() {
		t.Error("expected model to be focused after Focus()")
	}

	blurred := focused.Blur()
	if blurred.Focused() {
		t.Error("expected model to be blurred after Blur()")
	}
}
