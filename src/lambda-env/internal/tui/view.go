package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	categoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7D56F4"))

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4672"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F4D03F"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

// View renders the TUI as a string.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("LambdaOS Settings"))
	b.WriteString("\n\n")

	if m.view == viewCategories {
		m.renderCategories(&b)
	} else {
		m.renderModules(&b)
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		switch m.statusType {
		case "ok":
			b.WriteString(okStyle.Render(m.statusMsg))
		case "error":
			b.WriteString(errorStyle.Render(m.statusMsg))
		case "warning":
			b.WriteString(warningStyle.Render(m.statusMsg))
		default:
			b.WriteString(m.statusMsg)
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("↑/↓ or k/j navigate • enter select • esc back • q quit"))

	return b.String()
}

func (m Model) renderCategories(b *strings.Builder) {
	if len(m.categories) == 0 {
		b.WriteString("No modules found.\n")
		return
	}

	menu := m.hub.BuildMenu()
	catCount := make(map[string]int)
	for _, c := range menu {
		catCount[c.Name] = c.Count
	}

	for i, cat := range m.categories {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		count := catCount[cat]
		line := fmt.Sprintf("%s%s (%d)", cursor, cat, count)
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(line))
		} else {
			b.WriteString(itemStyle.Render(line))
		}
		b.WriteString("\n")
	}
}

func (m Model) renderModules(b *strings.Builder) {
	if len(m.modules) == 0 {
		b.WriteString(fmt.Sprintf("No modules in %s.\n", m.currentCategory))
		return
	}

	b.WriteString(categoryStyle.Render(fmt.Sprintf("%s (%d)", m.currentCategory, len(m.modules))))
	b.WriteString("\n")

	for i, mod := range m.modules {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		line := fmt.Sprintf("%s%s — %s", cursor, mod.Name, mod.Description)
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(line))
		} else {
			b.WriteString(itemStyle.Render(line))
		}
		b.WriteString("\n")
	}
}
