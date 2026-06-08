package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestNavigationCategoriesToModulesAndBack(t *testing.T) {
	m := createNavTestModel()

	// Start at categories view
	if m.view != viewCategories {
		t.Fatalf("initial view = %q, want %q", m.view, viewCategories)
	}

	// Select a category
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)

	if model.view != viewModules {
		t.Errorf("after category selection view = %q, want %q", model.view, viewModules)
	}

	// Navigate back
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)

	if model.view != viewCategories {
		t.Errorf("after back view = %q, want %q", model.view, viewCategories)
	}
}

func TestSelectionPreservationOnNavigateBack(t *testing.T) {
	m := createNavTestModel()

	// Navigate to system category
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)

	// Navigate to apps category
	updated, _ = model.Update(views.CategorySelectedMsg{Category: "apps", Index: 1})
	model = updated.(Model)

	if model.currentCategory != "apps" {
		t.Errorf("currentCategory = %q, want %q", model.currentCategory, "apps")
	}

	// Navigate back to categories
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)

	// The categories sub-model should still exist
	if model.categoriesSub == nil {
		t.Fatal("categoriesSub should not be nil after navigating back")
	}

	// The view should be categories
	if model.view != viewCategories {
		t.Errorf("view = %q, want %q", model.view, viewCategories)
	}
}

func TestEmptyCategoryShowsNoModules(t *testing.T) {
	m := createNavTestModel()

	// Simulate selecting a category with no modules
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "empty", Index: 0})
	model := updated.(Model)

	view := model.View()
	if !strings.Contains(view, "No modules") {
		t.Errorf("view = %q, want to contain 'No modules'", view)
	}
}

func TestCategoriesWrapAroundDown(t *testing.T) {
	m := createNavTestModel()

	// Navigate to last category
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	model := updated.(Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)

	view := model.View()
	if !strings.Contains(view, "apps") {
		t.Errorf("view = %q, want to contain 'apps' at last position", view)
	}

	// Wrap around to first
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)

	view = model.View()
	if !strings.Contains(view, "system") {
		t.Errorf("view = %q, want to contain 'system' after wrap", view)
	}
}

func TestCategoriesWrapAroundUp(t *testing.T) {
	m := createNavTestModel()

	// Wrap around from first to last
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	model := updated.(Model)

	view := model.View()
	if !strings.Contains(view, "apps") {
		t.Errorf("view = %q, want to contain 'apps' after wrap up", view)
	}
}

func TestModulesWrapAroundDown(t *testing.T) {
	m := createNavTestModel()

	// Navigate to modules view
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)

	// Navigate to last module
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)

	view := model.View()
	if !strings.Contains(view, "appearance") {
		t.Errorf("view = %q, want to contain 'appearance' at last position", view)
	}

	// Wrap around to first
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)

	view = model.View()
	if !strings.Contains(view, "keyboard") {
		t.Errorf("view = %q, want to contain 'keyboard' after wrap", view)
	}
}

func TestModulesWrapAroundUp(t *testing.T) {
	m := createNavTestModel()

	// Navigate to modules view
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)

	// Wrap around from first to last
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(Model)

	view := model.View()
	if !strings.Contains(view, "appearance") {
		t.Errorf("view = %q, want to contain 'appearance' after wrap up", view)
	}
}

func TestNavigationThroughAllViews(t *testing.T) {
	m := createNavTestModel()

	// categories → modules → back → categories
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 1: view = %q, want %q", model.view, viewModules)
	}

	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewCategories {
		t.Fatalf("step 2: view = %q, want %q", model.view, viewCategories)
	}

	// categories → modules again
	updated, _ = model.Update(views.CategorySelectedMsg{Category: "apps", Index: 1})
	model = updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 3: view = %q, want %q", model.view, viewModules)
	}
	if model.currentCategory != "apps" {
		t.Errorf("step 3: currentCategory = %q, want %q", model.currentCategory, "apps")
	}
}

func createNavTestModel() Model {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	modules := []module.Manifest{
		{Name: "keyboard", Description: "Set layout", Category: "system"},
		{Name: "audio", Description: "Configure audio", Category: "system"},
		{Name: "appearance", Description: "Set theme", Category: "system"},
		{Name: "neovim", Description: "Edit config", Category: "apps"},
		{Name: "qtile", Description: "Window manager", Category: "apps"},
	}

	h := &hub.Hub{Modules: modules}

	m := Model{
		hub:           h,
		categories:    cats,
		iconProvider:  icons.NewProvider(false),
		categoriesSub: views.NewCategoriesView(cats, menu, icons.NewProvider(false)),
		view:          viewCategories,
		cursor:        0,
		statusBar:     components.NewStatusBar().SetContext("categories"),
		helpOverlay: components.NewHelp([]components.KeyBinding{
			{Key: "↑/↓", Desc: "Navigate"},
			{Key: "enter", Desc: "Select"},
			{Key: "esc", Desc: "Back"},
			{Key: "?", Desc: "Toggle help"},
			{Key: "q", Desc: "Quit"},
		}),
	}
	m.activeSubModel = m.categoriesSub
	m.helpOverlay.Visible = false
	return m
}
