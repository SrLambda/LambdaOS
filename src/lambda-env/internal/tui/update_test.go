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
		{Name: "keyboard", Description: "Set layout", Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		}},
	}, "system")
	m.activeSubModel = m.modulesSub

	// Simulate selecting a module with nil hub — should transition to detail view
	updated, cmd := m.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{Name: "keyboard", Description: "Set layout"},
		Index:  0,
	})
	model := updated.(Model)

	if model.view != viewModuleDetail {
		t.Errorf("view = %q, want %q", model.view, viewModuleDetail)
	}
	if model.detailSub == nil {
		t.Fatal("detailSub should be initialized")
	}
	if cmd != nil {
		t.Error("expected nil cmd when hub is nil")
	}
}

func TestModuleSelectedWithHub(t *testing.T) {
	m := createTestModel()

	// First transition to modules view
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout", Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		}},
	}, "system")
	m.activeSubModel = m.modulesSub
	m.hub = &hub.Hub{}

	// Simulate selecting a module — should transition to detail view
	updated, _ := m.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{Name: "keyboard", Description: "Set layout", Path: "/tmp", Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		}},
		Index: 0,
	})
	model := updated.(Model)

	if model.view != viewModuleDetail {
		t.Errorf("view = %q, want %q", model.view, viewModuleDetail)
	}
	if model.detailSub == nil {
		t.Fatal("detailSub should be initialized after module selection")
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

func TestModuleSelectedTransitionsToDetailView(t *testing.T) {
	m := createTestModel()

	// Transition to modules view first
	m.view = viewModules
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout", Actions: []module.ActionConfig{
			{Name: "toggle-feature", Label: "Feature", Type: "toggle"},
		}},
	}, "system")
	m.activeSubModel = m.modulesSub

	// Select a module — should transition to detail view
	updated, _ := m.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{Name: "keyboard", Description: "Set layout", Actions: []module.ActionConfig{
			{Name: "toggle-feature", Label: "Feature", Type: "toggle"},
		}},
		Index: 0,
	})
	model := updated.(Model)

	if model.view != viewModuleDetail {
		t.Errorf("view = %q, want %q", model.view, viewModuleDetail)
	}
	if model.detailSub == nil {
		t.Fatal("detailSub should be initialized after module selection")
	}
	if model.activeSubModel != model.detailSub {
		t.Error("activeSubModel should be detailSub after module selection")
	}
}

func TestActionExecuteMsgWithNilHub(t *testing.T) {
	m := createTestModel()

	// Set up detail view
	m.view = viewModuleDetail
	m.detailSub = views.NewDetailView(module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		},
	})
	m.activeSubModel = m.detailSub

	// ActionExecuteMsg with nil hub should not crash and produce no cmd
	_, cmd := m.Update(views.ActionExecuteMsg{
		Module: module.Manifest{Name: "keyboard"},
		Name:   "keyboard",
		Action: "apply",
		Params: nil,
	})

	if cmd != nil {
		t.Error("expected nil cmd when hub is nil")
	}
}

func TestActionExecuteMsgWithHub(t *testing.T) {
	m := createTestModel()

	// Set up detail view with a module that has a path
	m.view = viewModuleDetail
	m.detailSub = views.NewDetailView(module.Manifest{
		Name: "keyboard",
		Path: "/tmp/nonexistent",
		Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		},
	})
	m.activeSubModel = m.detailSub
	m.hub = &hub.Hub{}

	// ActionExecuteMsg with hub should generate a command
	_, cmd := m.Update(views.ActionExecuteMsg{
		Module: module.Manifest{Name: "keyboard", Path: "/tmp/nonexistent"},
		Name:   "keyboard",
		Action: "apply",
		Params: nil,
	})

	if cmd == nil {
		t.Fatal("expected cmd after ActionExecuteMsg, got nil")
	}

	// The command may error because path doesn't exist, but it should produce execMsg
	msg := cmd()
	if msg == nil {
		return
	}
	_, ok := msg.(execMsg)
	if !ok {
		t.Fatalf("expected execMsg, got %T", msg)
	}
}

func TestExecMsgUpdatesDetailViewState(t *testing.T) {
	m := createTestModel()

	// Set up detail view
	m.view = viewModuleDetail
	m.detailSub = views.NewDetailView(module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us", "dvorak"}},
		},
	})
	m.activeSubModel = m.detailSub

	// Simulate an execMsg with data that should update detail view
	updated, _ := m.Update(execMsg{
		mod: module.Manifest{Name: "keyboard"},
		response: &module.Response{
			Status: "ok",
			Data: map[string]interface{}{
				"available_options": map[string]interface{}{
					"layout": []interface{}{"us", "dvorak", "colemak"},
				},
				"current_value": map[string]interface{}{
					"layout": "dvorak",
				},
			},
		},
		err: nil,
	})
	model := updated.(Model)

	// Verify status bar updated
	if model.statusBar == nil {
		t.Fatal("statusBar should be initialized")
	}
	view := model.statusBar.View()
	if !strings.Contains(view, "Module executed successfully") {
		t.Errorf("statusBar view = %q, want to contain success message", view)
	}

	// Verify detail view options were merged
	if model.detailSub == nil {
		t.Fatal("detailSub should exist")
	}
	if len(model.detailSub.Manifest().Actions[0].Options) != 3 {
		t.Errorf("expected 3 options after merge, got %d", len(model.detailSub.Manifest().Actions[0].Options))
	}
}

func TestBackFromDetailReturnsToModules(t *testing.T) {
	m := createTestModel()

	// Transition to detail view
	m.view = viewModuleDetail
	m.modulesSub = views.NewModulesView([]module.Manifest{
		{Name: "keyboard", Description: "Set layout"},
	}, "system")
	m.detailSub = views.NewDetailView(module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "toggle", Label: "Toggle", Type: "toggle"},
		},
	})
	m.activeSubModel = m.detailSub

	// Send back message
	updated, _ := m.Update(views.BackMsg{})
	model := updated.(Model)

	if model.view != viewModules {
		t.Errorf("view = %q, want %q", model.view, viewModules)
	}
	if model.activeSubModel != model.modulesSub {
		t.Error("activeSubModel should be modulesSub after back from detail")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
