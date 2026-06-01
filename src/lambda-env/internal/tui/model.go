package tui

import (
	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"

	tea "github.com/charmbracelet/bubbletea"
)

// viewState tracks which screen the TUI is showing.
type viewState string

const (
	viewCategories    viewState = "categories"
	viewModules       viewState = "modules"
	viewModuleDetail  viewState = "moduleDetail"
	viewConfirmDialog viewState = "confirmDialog"
)

// SubModel is a sub-model that the parent Model can delegate to.
type SubModel interface {
	tea.Model
}

// Model is the bubbletea model for the LambdaOS settings TUI.
type Model struct {
	hub             *hub.Hub
	categories      []string
	currentCategory string
	modules         []module.Manifest
	cursor          int
	view            viewState
	statusMsg       string
	statusType      string // ok | error | warning
	quitting        bool

	// Sub-models
	categoriesSub  *views.CategoriesView
	modulesSub     *views.ModulesView
	detailSub      *views.DetailView
	activeSubModel SubModel

	// Components
	statusBar   *components.StatusBar
	helpOverlay *components.Help
}

// NewModel creates a Model from a Hub instance.
func NewModel(h *hub.Hub) Model {
	menu := h.BuildMenu()
	cats := make([]string, 0, len(menu))
	for _, c := range menu {
		cats = append(cats, c.Name)
	}

	m := Model{
		hub:        h,
		categories: cats,
		view:       viewCategories,
		cursor:     0,
	}

	// Initialize sub-models
	m.categoriesSub = views.NewCategoriesView(cats, h.BuildMenu())
	m.activeSubModel = m.categoriesSub

	// Initialize components
	m.statusBar = components.NewStatusBar().SetContext("categories")
	m.helpOverlay = components.NewHelp([]components.KeyBinding{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "enter", Desc: "Select"},
		{Key: "esc", Desc: "Back"},
		{Key: "?", Desc: "Toggle help"},
		{Key: "q", Desc: "Quit"},
	})
	m.helpOverlay.Visible = false

	return m
}

// Init is the bubbletea initialization command.
func (m Model) Init() tea.Cmd {
	if m.activeSubModel != nil {
		return m.activeSubModel.Init()
	}
	return nil
}
