package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Background(lipgloss.Color("#1A1A1A")).
				Padding(0, 1).
				Width(80)

	statusBarContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	statusBarModuleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575"))

	statusBarModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4672"))

	statusBarStateStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F4D03F"))
)

// StatusBar is a persistent bar showing TUI context and state.
type StatusBar struct {
	Context       string
	Module        string
	SettingsState string
	Modified      bool
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

	if len(parts) == 0 {
		return statusBarStyle.Render(" ")
	}

	return statusBarStyle.Render(strings.Join(parts, "  |  "))
}
