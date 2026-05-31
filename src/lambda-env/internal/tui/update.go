package tui

import (
	"lambdaos.dev/lambda-env/internal/hub"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				// wrap around
				if m.view == viewCategories {
					m.cursor = len(m.categories) - 1
				} else {
					m.cursor = len(m.modules) - 1
				}
			}
		case "down", "j":
			if m.view == viewCategories {
				if m.cursor < len(m.categories)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			} else {
				if m.cursor < len(m.modules)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}
		case "enter":
			if m.view == viewCategories && len(m.categories) > 0 {
				m.currentCategory = m.categories[m.cursor]
				m.modules = filterByCategory(m.hub.Modules, m.currentCategory)
				m.view = viewModules
				m.cursor = 0
				m.statusMsg = ""
			} else if m.view == viewModules && len(m.modules) > 0 {
				mod := m.modules[m.cursor]
				return m, executeCmd(m.hub, mod)
			}
		case "esc":
			if m.view == viewModules {
				m.view = viewCategories
				m.cursor = 0
				m.modules = nil
				m.currentCategory = ""
				m.statusMsg = ""
			}
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case execMsg:
		if msg.err != nil {
			m.statusMsg = "Error: " + msg.err.Error()
			m.statusType = "error"
		} else if msg.response != nil {
			switch msg.response.Status {
			case "ok":
				m.statusMsg = msg.response.Message
				if m.statusMsg == "" {
					m.statusMsg = "Module executed successfully"
				}
				m.statusType = "ok"
			case "error":
				m.statusMsg = msg.response.Message
				if m.statusMsg == "" {
					m.statusMsg = "Module reported an error"
				}
				m.statusType = "error"
			case "warning":
				m.statusMsg = msg.response.Message
				if m.statusMsg == "" {
					m.statusMsg = "Warning from module"
				}
				m.statusType = "warning"
			default:
				m.statusMsg = "Unknown response status"
				m.statusType = "error"
			}
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
