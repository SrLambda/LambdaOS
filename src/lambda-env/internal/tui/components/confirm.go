package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"lambdaos.dev/lambda-env/internal/tui/theme"
)

var (
	confirmOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(theme.Accent)).
				Padding(2, 4).
				Width(40)

	confirmMessageStyle = lipgloss.NewStyle().
				Bold(true).
				MarginBottom(1)

	confirmSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(theme.Accent)).
				Foreground(lipgloss.Color(theme.TextPrimary)).
				Padding(0, 1)

	confirmUnselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Dimmed)).
				Padding(0, 1)
)

// ConfirmResultMsg is emitted when the user makes a choice.
type ConfirmResultMsg struct {
	Confirmed bool
}

// Confirm is a yes/no confirmation dialog.
type Confirm struct {
	Message   string
	Selected  int // 0 = Yes, 1 = No
	Confirmed bool
}

// NewConfirm creates a new Confirm dialog with the given message.
func NewConfirm(message string) *Confirm {
	return &Confirm{
		Message:  message,
		Selected: 0,
	}
}

// Init implements tea.Model.
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// Update handles keyboard navigation and selection.
func (c *Confirm) Update(msg tea.Msg) (*Confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if c.Selected > 0 {
				c.Selected--
			} else {
				c.Selected = 1
			}
			return c, nil
		case tea.KeyRight:
			if c.Selected < 1 {
				c.Selected++
			} else {
				c.Selected = 0
			}
			return c, nil
		case tea.KeySpace:
			if c.Selected == 0 {
				c.Selected = 1
			} else {
				c.Selected = 0
			}
			return c, nil
		case tea.KeyEnter:
			c.Confirmed = c.Selected == 0
			return c, func() tea.Msg {
				return ConfirmResultMsg{Confirmed: c.Confirmed}
			}
		case tea.KeyEsc:
			c.Confirmed = false
			return c, func() tea.Msg {
				return ConfirmResultMsg{Confirmed: false}
			}
		}
	}
	return c, nil
}

// View renders the confirmation dialog overlay.
func (c *Confirm) View() string {
	var b strings.Builder
	b.WriteString(confirmMessageStyle.Render(c.Message))
	b.WriteString("\n\n")

	var yesStyle, noStyle lipgloss.Style
	if c.Selected == 0 {
		yesStyle = confirmSelectedStyle
		noStyle = confirmUnselectedStyle
	} else {
		yesStyle = confirmUnselectedStyle
		noStyle = confirmSelectedStyle
	}

	b.WriteString(yesStyle.Render("Yes"))
	b.WriteString("    ")
	b.WriteString(noStyle.Render("No"))

	return confirmOverlayStyle.Render(b.String())
}
