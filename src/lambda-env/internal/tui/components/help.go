package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"lambdaos.dev/lambda-env/internal/tui/theme"
)

var (
	helpOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(theme.Accent)).
				Padding(2, 4).
				Width(50)

	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	helpDismissStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Dimmed)).
				MarginTop(1)
)

// KeyBinding describes a keyboard shortcut and its purpose.
type KeyBinding struct {
	Key  string
	Desc string
}

// HelpDismissedMsg is emitted when the help overlay is dismissed.
type HelpDismissedMsg struct{}

// Help is a dismissible overlay showing available key bindings.
type Help struct {
	Bindings []KeyBinding
	Visible  bool
}

// NewHelp creates a new Help overlay with the given bindings.
func NewHelp(bindings []KeyBinding) *Help {
	return &Help{
		Bindings: bindings,
		Visible:  true,
	}
}

// Init implements tea.Model.
func (h *Help) Init() tea.Cmd {
	return nil
}

// Update handles dismissal keys.
func (h *Help) Update(msg tea.Msg) (*Help, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			h.Visible = false
			return h, func() tea.Msg {
				return HelpDismissedMsg{}
			}
		case tea.KeyRunes:
			if len(msg.Runes) == 1 && msg.Runes[0] == '?' {
				h.Visible = !h.Visible
				return h, nil
			}
		}
	}
	return h, nil
}

// View renders the help overlay.
func (h *Help) View() string {
	if !h.Visible {
		return ""
	}

	var b strings.Builder
	b.WriteString(helpTitleStyle.Render("Keyboard Shortcuts") + "\n\n")

	for _, kb := range h.Bindings {
		b.WriteString(kb.Key + "  " + kb.Desc + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpDismissStyle.Render("Press esc or ? to dismiss"))

	return helpOverlayStyle.Render(b.String())
}
