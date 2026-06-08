package theme

import "github.com/charmbracelet/lipgloss"

// Color constants for the LambdaOS TUI palette.
const (
	Bg          = "#0D0D0D"
	StatusBg    = "#1A1A1A"
	Accent    = "#8B6AF4"
	Success     = "#04B575"
	Error       = "#FF4672"
	Warn        = "#F4D03F"
	Dimmed      = "#909090"
	TextPrimary = "#FFFFFF"
)

// Common lipgloss style variables that consume the color constants above.
var (
	AccentStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(Accent))
	SuccessStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(Success))
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(Error))
	WarnStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(Warn))
	DimmedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(Dimmed))
	StatusBarStyle = lipgloss.NewStyle().Background(lipgloss.Color(StatusBg)).Foreground(lipgloss.Color(Dimmed))
)
