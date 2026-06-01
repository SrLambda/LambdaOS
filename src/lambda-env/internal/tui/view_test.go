package tui

import (
	"strings"
	"testing"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestViewDelegatesToActiveSubModel(t *testing.T) {
	m := createViewTestModel()

	view := m.View()
	if !strings.Contains(view, "LambdaOS Settings") {
		t.Errorf("view = %q, want to contain 'LambdaOS Settings'", view)
	}
}

func TestViewRendersStatusBar(t *testing.T) {
	m := createViewTestModel()
	m.statusBar.SetContext("categories")

	view := m.View()
	if !strings.Contains(view, "categories") {
		t.Errorf("view = %q, want to contain status bar context 'categories'", view)
	}
}

func TestViewRendersHelpOverlay(t *testing.T) {
	m := createViewTestModel()
	m.helpOverlay.Visible = true

	view := m.View()
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("view = %q, want to contain 'Keyboard Shortcuts'", view)
	}
}

func TestViewDoesNotRenderHelpWhenHidden(t *testing.T) {
	m := createViewTestModel()
	m.helpOverlay.Visible = false

	view := m.View()
	if strings.Contains(view, "Keyboard Shortcuts") {
		t.Errorf("view = %q, should not contain 'Keyboard Shortcuts' when hidden", view)
	}
}

func TestViewRendersModulesSubModel(t *testing.T) {
	m := createViewTestModel()
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout"},
	}, "system")
	m.activeSubModel = m.modulesSub
	m.statusBar.SetContext("modules")

	view := m.View()
	if !strings.Contains(view, "keyboard") {
		t.Errorf("view = %q, want to contain 'keyboard'", view)
	}
	if !strings.Contains(view, "system") {
		t.Errorf("view = %q, want to contain category 'system'", view)
	}
}

func TestViewWithExecStatus(t *testing.T) {
	m := createViewTestModel()
	m.statusBar.SetSettingsState("Layout applied")

	view := m.View()
	if !strings.Contains(view, "Layout applied") {
		t.Errorf("view = %q, want to contain 'Layout applied'", view)
	}
}

func TestViewQuittingReturnsEmpty(t *testing.T) {
	m := createViewTestModel()
	m.quitting = true

	view := m.View()
	if view != "" {
		t.Errorf("view = %q, want empty string when quitting", view)
	}
}

func createViewTestModel() Model {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	m := Model{
		categories:    cats,
		categoriesSub: views.NewCategoriesView(cats, menu),
		view:          viewCategories,
		cursor:        0,
		statusBar:     components.NewStatusBar(),
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
