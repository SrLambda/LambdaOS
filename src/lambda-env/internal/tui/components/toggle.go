package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	toggleLabelStyle = lipgloss.NewStyle().
				Bold(true)

	toggleOnStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575"))

	toggleOffStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262"))

	toggleCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4"))
)

// ToggleChangedMsg is emitted when the toggle value changes.
type ToggleChangedMsg struct {
	Label string
	Value bool
}

// Toggle is a boolean toggle component.
type Toggle struct {
	Label   string
	Value   bool
	Focused bool
}

// NewToggle creates a new Toggle with the given label and initial value.
func NewToggle(label string, value bool) Toggle {
	return Toggle{
		Label: label,
		Value: value,
	}
}

// Init implements tea.Model.
func (t Toggle) Init() tea.Cmd {
	return nil
}

// Update handles user input and returns the updated toggle.
func (t Toggle) Update(msg tea.Msg) (Toggle, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace, tea.KeyEnter:
			t.Value = !t.Value
			return t, func() tea.Msg {
				return ToggleChangedMsg{Label: t.Label, Value: t.Value}
			}
		}
	}
	return t, nil
}

// View renders the toggle as a string.
func (t Toggle) View() string {
	var indicator string
	if t.Value {
		indicator = toggleOnStyle.Render("● On")
	} else {
		indicator = toggleOffStyle.Render("○ Off")
	}

	cursor := "  "
	if t.Focused {
		cursor = toggleCursorStyle.Render("> ")
	}

	return cursor + toggleLabelStyle.Render(t.Label) + "  " + indicator
}
