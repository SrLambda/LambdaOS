package components

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"lambdaos.dev/lambda-env/internal/tui/theme"
)

var (
	textInputLabelStyle = lipgloss.NewStyle().Bold(true)
	textInputErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Error))
)

// TextInputSubmitMsg is emitted when the user presses Enter with a valid value.
type TextInputSubmitMsg struct {
	Label string
	Value string
}

// TextInputCancelMsg is emitted when the user presses Esc.
type TextInputCancelMsg struct{}

// TextInput wraps bubbles/textinput with validation and error display.
type TextInput struct {
	Label     string
	model     textinput.Model
	allowlist string
	regex     *regexp.Regexp
	minNum    *int
	maxNum    *int
	Err       error
}

// NewTextInput creates a new TextInput with the given label.
func NewTextInput(label string) *TextInput {
	m := textinput.New()
	m.Focus()
	return &TextInput{
		Label: label,
		model: m,
	}
}

// Value returns the current text value.
func (t *TextInput) Value() string {
	return t.model.Value()
}

// SetValue sets the text value programmatically.
func (t *TextInput) SetValue(v string) {
	t.model.SetValue(v)
}

// SetPlaceholder sets the placeholder text.
func (t *TextInput) SetPlaceholder(p string) {
	t.model.Placeholder = p
}

// SetMaxLength sets the maximum allowed length.
func (t *TextInput) SetMaxLength(n int) {
	t.model.CharLimit = n
}

// SetAllowlist restricts input to characters in the given string.
func (t *TextInput) SetAllowlist(chars string) {
	t.allowlist = chars
}

// SetRegex sets a regex pattern that the final value must match.
func (t *TextInput) SetRegex(pattern string) {
	t.regex = regexp.MustCompile(pattern)
}

// SetNumericRange sets min and max bounds for numeric input.
func (t *TextInput) SetNumericRange(min, max int) {
	t.minNum = &min
	t.maxNum = &max
}

// Focus focuses the underlying textinput model.
func (t *TextInput) Focus() *TextInput {
	t.model.Focus()
	return t
}

// Blur removes focus from the underlying textinput model.
func (t *TextInput) Blur() *TextInput {
	t.model.Blur()
	return t
}

// Focused returns whether the underlying model is focused.
func (t *TextInput) Focused() bool {
	return t.model.Focused()
}

// Init implements tea.Model.
func (t *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input and returns the updated text input.
func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			t.Err = t.validate()
			if t.Err == nil {
				return t, func() tea.Msg {
					return TextInputSubmitMsg{Label: t.Label, Value: t.model.Value()}
				}
			}
			return t, nil
		case tea.KeyEsc:
			t.model.SetValue("")
			t.Err = nil
			return t, func() tea.Msg {
				return TextInputCancelMsg{}
			}
		default:
			if msg.Type == tea.KeyRunes && t.allowlist != "" {
				for _, r := range msg.Runes {
					if !strings.ContainsRune(t.allowlist, r) {
						t.Err = fmt.Errorf("character %q is not allowed", r)
						return t, nil
					}
				}
			}
			var cmd tea.Cmd
			t.model, cmd = t.model.Update(msg)
			t.Err = nil
			return t, cmd
		}
	}

	var cmd tea.Cmd
	t.model, cmd = t.model.Update(msg)
	return t, cmd
}

func (t *TextInput) validate() error {
	val := t.model.Value()

	if t.regex != nil && !t.regex.MatchString(val) {
		return fmt.Errorf("value does not match required pattern")
	}

	if t.minNum != nil || t.maxNum != nil {
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("value must be a number")
		}
		if t.minNum != nil && n < *t.minNum {
			return fmt.Errorf("value must be at least %d", *t.minNum)
		}
		if t.maxNum != nil && n > *t.maxNum {
			return fmt.Errorf("value must be at most %d", *t.maxNum)
		}
	}

	return nil
}

// View renders the text input with label and any validation error.
func (t *TextInput) View() string {
	var b strings.Builder
	b.WriteString(textInputLabelStyle.Render(t.Label) + "\n")
	b.WriteString(t.model.View())
	if t.Err != nil {
		b.WriteString("\n")
		b.WriteString(textInputErrorStyle.Render(t.Err.Error()))
	}
	return b.String()
}
