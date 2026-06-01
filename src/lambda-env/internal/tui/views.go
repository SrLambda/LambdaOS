package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/pkg/module"
)

// categoriesView is a placeholder for the categories sub-model.
// Full implementation in task 3a.2.
type categoriesView struct {
	categories []string
	menu       []hub.MenuCategory
	cursor     int
}

func newCategoriesView(cats []string, menu []hub.MenuCategory) *categoriesView {
	return &categoriesView{
		categories: cats,
		menu:       menu,
		cursor:     0,
	}
}

func (c *categoriesView) Init() tea.Cmd {
	return nil
}

func (c *categoriesView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *categoriesView) View() string {
	return "categories view placeholder"
}

// modulesView is a placeholder for the modules sub-model.
// Full implementation in task 3a.3.
type modulesView struct {
	modules  []module.Manifest
	category string
	cursor   int
}

func newModulesView(mods []module.Manifest, category string) *modulesView {
	return &modulesView{
		modules:  mods,
		category: category,
		cursor:   0,
	}
}

func (m *modulesView) Init() tea.Cmd {
	return nil
}

func (m *modulesView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *modulesView) View() string {
	return "modules view placeholder"
}
