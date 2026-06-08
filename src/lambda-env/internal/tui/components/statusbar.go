package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"lambdaos.dev/lambda-env/internal/tui/theme"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Dimmed)).
			Background(lipgloss.Color(theme.StatusBg)).
			Padding(0, 1).
			Width(80)

	statusBarContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Accent)).
				Bold(true)

	statusBarModuleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Success))

	statusBarModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Error))

	statusBarStateStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Warn))
)

// StatusBar is a persistent bar showing TUI context and state.
type StatusBar struct {
	Context       string
	Module        string
	SettingsState string
	Modified      bool
	Width         int
}

// NewStatusBar creates a new empty StatusBar.
func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

// SetContext sets the current view context.
func (s *StatusBar) SetContext(c string) *StatusBar {
	s.Context = c
	return s
}

// SetModule sets the current module name.
func (s *StatusBar) SetModule(m string) *StatusBar {
	s.Module = m
	return s
}

// SetSettingsState sets the settings state text.
func (s *StatusBar) SetSettingsState(st string) *StatusBar {
	s.SettingsState = st
	return s
}

// SetModified sets the modified indicator.
func (s *StatusBar) SetModified(m bool) *StatusBar {
	s.Modified = m
	return s
}

// SetWidth sets the status bar width to match the terminal.
func (s *StatusBar) SetWidth(w int) *StatusBar {
	s.Width = w
	return s
}

// View renders the status bar.
func (s *StatusBar) View() string {
	var parts []string

	if s.Context != "" {
		parts = append(parts, statusBarContextStyle.Render(s.Context))
	}
	if s.Module != "" {
		parts = append(parts, statusBarModuleStyle.Render(s.Module))
	}
	if s.SettingsState != "" {
		parts = append(parts, statusBarStateStyle.Render(s.SettingsState))
	}
	if s.Modified {
		parts = append(parts, statusBarModifiedStyle.Render("*"))
	}

	width := s.Width
	if width == 0 {
		width = 80
	}
	style := statusBarStyle.Width(width)

	if len(parts) == 0 {
		return style.Render(" ")
	}

	return style.Render(strings.Join(parts, "  |  "))
}
