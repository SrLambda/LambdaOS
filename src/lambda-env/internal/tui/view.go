package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4672"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F4D03F"))
)

// View renders the TUI as a string.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Render active sub-model view
	if m.activeSubModel != nil {
		b.WriteString(m.activeSubModel.View())
	}

	// Render status bar (persistent across all views)
	if m.statusBar != nil {
		b.WriteString("\n")
		b.WriteString(m.statusBar.View())
		b.WriteString("\n")
	}

	// Render help overlay if visible
	if m.helpOverlay != nil && m.helpOverlay.Visible {
		// Overlay is rendered on top - in a real terminal this would be
		// a modal overlay, but for simplicity we render it inline
		b.WriteString(m.helpOverlay.View())
	}

	return b.String()
}
