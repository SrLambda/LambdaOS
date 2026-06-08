package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/internal/tui/theme"
	"lambdaos.dev/lambda-env/pkg/module"
)

var (
	moduleCategoryStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(theme.Success))

	moduleItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	moduleSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(theme.Accent))
)

// ModuleSelectedMsg is emitted when the user selects a module.
type ModuleSelectedMsg struct {
	Module module.Manifest
	Index  int
}

// BackMsg is emitted when the user wants to go back to the previous view.
type BackMsg struct{}

// ModulesView is a sub-model for the module list screen.
type ModulesView struct {
	modules        []module.Manifest
	category       string
	cursor         int
	selectedModule string
	iconProvider   icons.IconProvider
}

// NewModulesView creates a new ModulesView.
func NewModulesView(mods []module.Manifest, category string, provider icons.IconProvider) *ModulesView {
	return &ModulesView{
		modules:      mods,
		category:     category,
		cursor:       0,
		iconProvider: provider,
	}
}

// Init implements tea.Model.
func (m *ModulesView) Init() tea.Cmd {
	return nil
}

// Update handles user input for the modules view.
func (m *ModulesView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.modules) - 1
			}
		case tea.KeyDown:
			if m.cursor < len(m.modules)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				switch msg.Runes[0] {
				case 'k':
					if m.cursor > 0 {
						m.cursor--
					} else {
						m.cursor = len(m.modules) - 1
					}
				case 'j':
					if m.cursor < len(m.modules)-1 {
						m.cursor++
					} else {
						m.cursor = 0
					}
				}
			}
		case tea.KeyEnter:
			if len(m.modules) > 0 && m.cursor < len(m.modules) {
				m.selectedModule = m.modules[m.cursor].Name
				return m, func() tea.Msg {
					return ModuleSelectedMsg{
						Module: m.modules[m.cursor],
						Index:  m.cursor,
					}
				}
			}
		case tea.KeyEsc:
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}
	}
	return m, nil
}

// View renders the modules list.
func (m *ModulesView) View() string {
	var b strings.Builder

	if len(m.modules) == 0 {
		b.WriteString(fmt.Sprintf("No modules in %s.\n", m.category))
		return b.String()
	}

	b.WriteString(moduleCategoryStyle.Render(fmt.Sprintf("%s (%d)", m.category, len(m.modules))))
	b.WriteString("\n")

	width := m.iconProvider.Width()
	for i, mod := range m.modules {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		icon := m.iconProvider.ForModule(mod.Name)
		iconStr := icon + strings.Repeat(" ", width-1)
		line := fmt.Sprintf("%s%s %s — %s", cursor, iconStr, mod.Name, mod.Description)
		if m.cursor == i {
			b.WriteString(moduleSelectedStyle.Render(line))
		} else {
			b.WriteString(moduleItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	return b.String()
}
