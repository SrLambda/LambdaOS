package tui

import (
	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/pkg/module"

	tea "github.com/charmbracelet/bubbletea"
)

// viewState tracks which screen the TUI is showing.
type viewState string

const (
	viewCategories viewState = "categories"
	viewModules    viewState = "modules"
)

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
}

// NewModel creates a Model from a Hub instance.
func NewModel(h *hub.Hub) Model {
	menu := h.BuildMenu()
	cats := make([]string, 0, len(menu))
	for _, c := range menu {
		cats = append(cats, c.Name)
	}
	return Model{
		hub:        h,
		categories: cats,
		view:       viewCategories,
		cursor:     0,
	}
}

// Init is the bubbletea initialization command.
func (m Model) Init() tea.Cmd {
	return nil
}
