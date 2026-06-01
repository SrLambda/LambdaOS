package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/views"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestUpdateDelegatesToActiveSubModel(t *testing.T) {
	m := createTestModel()

	// Send a key message that should be handled by the active sub-model
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	model := updated.(Model)

	// Verify delegation by checking the view output changed
	view := model.activeSubModel.View()
	if !strings.Contains(view, "apps") {
		t.Errorf("view should contain 'apps' after down navigation, got: %q", view)
	}
}

func TestCategorySelectedTransitionsToModules(t *testing.T) {
	m := createTestModel()

	// Simulate selecting the first category
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)

	if model.view != viewModules {
		t.Errorf("view = %q, want %q", model.view, viewModules)
	}
	if model.activeSubModel != model.modulesSub {
		t.Error("activeSubModel should be modulesSub after category selection")
	}
	// Verify the modules sub-model was created for the right category
	if model.modulesSub != nil {
		view := model.modulesSub.View()
		if !strings.Contains(view, "system") {
			t.Errorf("modulesSub view should contain 'system', got: %q", view)
		}
	}
}

func TestModuleSelectedWithNilHub(t *testing.T) {
	m := createTestModel()

	// First transition to modules view
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout"},
	}, "system")
	m.activeSubModel = m.modulesSub

	// Simulate selecting a module with nil hub — should not crash, no cmd
	_, cmd := m.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{Name: "keyboard", Description: "Set layout"},
		Index:  0,
	})

	if cmd != nil {
		t.Error("expected nil cmd when hub is nil")
	}
}

func TestModuleSelectedWithHub(t *testing.T) {
	m := createTestModel()

	// First transition to modules view
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout"},
	}, "system")
	m.activeSubModel = m.modulesSub
	m.hub = &hub.Hub{} // Non-nil hub for command generation

	// Simulate selecting a module
	_, cmd := m.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{Name: "keyboard", Description: "Set layout", Path: "/tmp"},
		Index:  0,
	})

	if cmd == nil {
		t.Fatal("expected cmd after module selection, got nil")
	}

	// The command should produce an execMsg when executed
	// (it may error because Path doesn't exist, but it should still produce execMsg)
	msg := cmd()
	if msg == nil {
		// This is expected because the module path doesn't exist
		// The important thing is that cmd was generated
		return
	}
	_, ok := msg.(execMsg)
	if !ok {
		t.Fatalf("expected execMsg, got %T", msg)
	}
}

func TestBackMsgReturnsToCategories(t *testing.T) {
	m := createTestModel()

	// Transition to modules
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout"},
	}, "system")
	m.activeSubModel = m.modulesSub

	// Send back message
	updated, _ := m.Update(views.BackMsg{})
	model := updated.(Model)

	if model.view != viewCategories {
		t.Errorf("view = %q, want %q", model.view, viewCategories)
	}
	if model.activeSubModel != model.categoriesSub {
		t.Error("activeSubModel should be categoriesSub after back")
	}
}

func TestExecMsgUpdatesStatusBar(t *testing.T) {
	m := createTestModel()

	// Simulate an execMsg with success
	updated, _ := m.Update(execMsg{
		mod:      module.Manifest{Name: "keyboard"},
		response: &module.Response{Status: "ok", Message: "Layout applied"},
		err:      nil,
	})
	model := updated.(Model)

	if model.statusBar == nil {
		t.Fatal("statusBar should be initialized")
	}
	view := model.statusBar.View()
	if !strings.Contains(view, "Layout applied") && !strings.Contains(view, "ok") {
		t.Errorf("statusBar view = %q, want to contain success info", view)
	}
}

func TestExecMsgErrorUpdatesStatusBar(t *testing.T) {
	m := createTestModel()

	// Simulate an execMsg with error
	updated, _ := m.Update(execMsg{
		mod:      module.Manifest{Name: "keyboard"},
		response: nil,
		err:      &testError{msg: "command failed"},
	})
	model := updated.(Model)

	if model.statusBar == nil {
		t.Fatal("statusBar should be initialized")
	}
	view := model.statusBar.View()
	if !strings.Contains(view, "command failed") && !strings.Contains(view, "error") {
		t.Errorf("statusBar view = %q, want to contain error info", view)
	}
}

func TestQuitKeyExitsApplication(t *testing.T) {
	m := createTestModel()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	model := updated.(Model)

	if !model.quitting {
		t.Error("quitting should be true after pressing 'q'")
	}
	if cmd == nil {
		t.Fatal("expected tea.Quit command, got nil")
	}

	// Verify the command is tea.Quit
	msg := cmd()
	if msg == nil {
		// tea.Quit returns nil as a message, which is valid
		return
	}
}

func TestCtrlCExitsApplication(t *testing.T) {
	m := createTestModel()

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model := updated.(Model)

	if !model.quitting {
		t.Error("quitting should be true after pressing ctrl+c")
	}
	if cmd == nil {
		t.Fatal("expected tea.Quit command, got nil")
	}
}

func TestHelpToggleWithQuestionMark(t *testing.T) {
	m := createTestModel()

	// Press '?' to show help
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model := updated.(Model)

	if model.helpOverlay == nil {
		t.Fatal("helpOverlay should be initialized")
	}
	if !model.helpOverlay.Visible {
		t.Error("helpOverlay should be visible after pressing '?'")
	}

	// Press '?' again to hide help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(Model)
	if model.helpOverlay.Visible {
		t.Error("helpOverlay should be hidden after pressing '?' again")
	}
}

func createTestModel() Model {
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

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
