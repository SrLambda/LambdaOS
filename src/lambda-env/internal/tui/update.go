package tui

import (
	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"

	tea "github.com/charmbracelet/bubbletea"
)

// execMsg carries the result of a module execution.
type execMsg struct {
	mod      module.Manifest
	response *module.Response
	err      error
}

// Update handles incoming messages and user input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global quit keys first.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyRunes:
			if len(keyMsg.Runes) == 1 {
				switch keyMsg.Runes[0] {
				case 'q':
					m.quitting = true
					return m, tea.Quit
				case '?':
					if m.helpOverlay != nil {
						m.helpOverlay.Visible = !m.helpOverlay.Visible
					}
					return m, nil
				}
			}
		}
	}

	// Handle view transition messages.
	switch msg := msg.(type) {
	case views.CategorySelectedMsg:
		return m.handleCategorySelected(msg)

	case views.ModuleSelectedMsg:
		return m.handleModuleSelected(msg)

	case views.BackMsg:
		return m.handleBack()

	case execMsg:
		return m.handleExecMsg(msg)
	}

	// If help overlay is visible, delegate to it first.
	if m.helpOverlay != nil && m.helpOverlay.Visible {
		updated, cmd := m.helpOverlay.Update(msg)
		m.helpOverlay = updated
		if !m.helpOverlay.Visible {
			// Help was dismissed; return the dismiss message
			return m, cmd
		}
		return m, cmd
	}

	// Delegate to active sub-model for everything else.
	if m.activeSubModel != nil {
		updated, cmd := m.activeSubModel.Update(msg)
		m.activeSubModel = updated
		return m, cmd
	}

	return m, nil
}

func (m Model) handleCategorySelected(msg views.CategorySelectedMsg) (tea.Model, tea.Cmd) {
	m.currentCategory = msg.Category
	if m.hub != nil {
		m.modules = filterByCategory(m.hub.Modules, m.currentCategory)
	}
	m.view = viewModules
	m.cursor = 0
	m.statusMsg = ""

	m.modulesSub = views.NewModulesView(m.modules, m.currentCategory)
	m.activeSubModel = m.modulesSub

	if m.statusBar != nil {
		m.statusBar.SetContext("modules").SetModule(m.currentCategory)
	}

	return m, nil
}

func (m Model) handleModuleSelected(msg views.ModuleSelectedMsg) (tea.Model, tea.Cmd) {
	// Wave 2 behavior: execute the module directly.
	// Wave 3b will navigate to detail view instead.
	if m.hub != nil {
		return m, executeCmd(m.hub, msg.Module)
	}
	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	m.view = viewCategories
	m.cursor = 0
	m.modules = nil
	m.currentCategory = ""
	m.statusMsg = ""

	m.activeSubModel = m.categoriesSub

	if m.statusBar != nil {
		m.statusBar.SetContext("categories").SetModule("")
	}

	return m, nil
}

func (m Model) handleExecMsg(msg execMsg) (tea.Model, tea.Cmd) {
	if m.statusBar == nil {
		m.statusBar = components.NewStatusBar()
	}

	if msg.err != nil {
		m.statusBar.SetSettingsState("error: " + msg.err.Error())
		m.statusType = "error"
	} else if msg.response != nil {
		switch msg.response.Status {
		case "ok":
			msgStr := msg.response.Message
			if msgStr == "" {
				msgStr = "Module executed successfully"
			}
			m.statusBar.SetSettingsState(msgStr)
			m.statusType = "ok"
		case "error":
			msgStr := msg.response.Message
			if msgStr == "" {
				msgStr = "Module reported an error"
			}
			m.statusBar.SetSettingsState(msgStr)
			m.statusType = "error"
		case "warning":
			msgStr := msg.response.Message
			if msgStr == "" {
				msgStr = "Warning from module"
			}
			m.statusBar.SetSettingsState(msgStr)
			m.statusType = "warning"
		default:
			m.statusBar.SetSettingsState("Unknown response status")
			m.statusType = "error"
		}
	}

	return m, nil
}

// filterByCategory returns modules whose Category matches cat.
func filterByCategory(mods []module.Manifest, cat string) []module.Manifest {
	var out []module.Manifest
	for _, m := range mods {
		if m.Category == cat {
			out = append(out, m)
		}
	}
	return out
}

// executeCmd returns a bubbletea command that runs a module.
func executeCmd(h *hub.Hub, mod module.Manifest) tea.Cmd {
	return func() tea.Msg {
		resp, err := h.ExecuteModule(mod)
		return execMsg{mod: mod, response: resp, err: err}
	}
}
