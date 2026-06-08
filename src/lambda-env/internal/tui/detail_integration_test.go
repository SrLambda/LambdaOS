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

func TestFullFlowToggleActionAndBack(t *testing.T) {
	m := createIntegrationTestModel()

	// Step 1: categories → modules
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 1: view = %q, want %q", model.view, viewModules)
	}

	// Step 2: modules → detail view
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name:    "keyboard",
			Actions: []module.ActionConfig{{Name: "toggle-feature", Label: "Feature", Type: "toggle"}},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("step 2: view = %q, want %q", model.view, viewModuleDetail)
	}
	if model.detailSub == nil {
		t.Fatal("detailSub should be set")
	}

	// Step 3: toggle action in detail view
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("expected cmd after toggle action")
	}
	msg := cmd()
	actionMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if actionMsg.Action != "toggle-feature" {
		t.Errorf("action = %q, want %q", actionMsg.Action, "toggle-feature")
	}

	// Step 4: simulate execution response
	updated, _ = model.Update(execMsg{
		mod:      module.Manifest{Name: "keyboard"},
		response: &module.Response{Status: "ok", Message: "Feature toggled"},
		err:      nil,
	})
	model = updated.(Model)
	if !strings.Contains(model.statusBar.View(), "Feature toggled") {
		t.Errorf("statusBar = %q, want to contain 'Feature toggled'", model.statusBar.View())
	}

	// Step 5: back to modules
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 5: view = %q, want %q", model.view, viewModules)
	}

	// Step 6: back to categories
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewCategories {
		t.Fatalf("step 6: view = %q, want %q", model.view, viewCategories)
	}
}

func TestSelectActionWithDynamicOptions(t *testing.T) {
	m := createIntegrationTestModel()

	// Navigate to detail view with select action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "appearance",
			Actions: []module.ActionConfig{
				{Name: "theme", Label: "Theme", Type: "select", Options: []string{"dark"}},
			},
		},
		Index: 0,
	})
	model = updated.(Model)

	// Simulate dynamic options arriving
	updated, _ = model.Update(views.DynamicOptionsMsg{
		Options: map[string][]string{"theme": {"dark", "light", "nord"}},
		Values:  map[string]interface{}{"theme": "nord"},
	})
	model = updated.(Model)

	if model.detailSub == nil {
		t.Fatal("detailSub should exist")
	}
	if len(model.detailSub.Manifest().Actions[0].Options) != 3 {
		t.Errorf("expected 3 options after dynamic merge, got %d", len(model.detailSub.Manifest().Actions[0].Options))
	}

	// Change selection with right arrow
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = updated.(Model)

	// Execute the select action
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("expected cmd after select action")
	}
	msg := cmd()
	execMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Action != "theme" {
		t.Errorf("action = %q, want %q", execMsg.Action, "theme")
	}
}

func TestConfirmDialogBeforeDestructiveAction(t *testing.T) {
	m := createIntegrationTestModel()

	// Navigate to detail view with confirm action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "keyboard",
			Actions: []module.ActionConfig{
				{Name: "reset", Label: "Reset", Type: "confirm"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)

	// Press enter to show confirm dialog
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	// Confirm dialog should be rendered in the view
	view := model.View()
	if !strings.Contains(view, "Reset") {
		t.Errorf("view = %q, want to contain confirm message", view)
	}
	if !strings.Contains(view, "Yes") {
		t.Errorf("view = %q, want to contain 'Yes' button", view)
	}

	// Press Enter to confirm (Yes is selected by default)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	view = model.View()
	if strings.Contains(view, "Yes") && strings.Contains(view, "No") {
		t.Error("confirm dialog should be dismissed after confirming")
	}
	if cmd == nil {
		t.Fatal("expected cmd after confirm")
	}
	msg := cmd()
	execMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Action != "reset" {
		t.Errorf("action = %q, want %q", execMsg.Action, "reset")
	}
	if confirmed, _ := execMsg.Params["confirmed"].(bool); !confirmed {
		t.Errorf("confirmed = %v, want true", execMsg.Params["confirmed"])
	}
}

func TestErrorHandlingOnModuleExecution(t *testing.T) {
	m := createIntegrationTestModel()

	// Navigate to detail view with execute action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "keyboard",
			Actions: []module.ActionConfig{
				{Name: "apply", Label: "Apply", Type: "execute"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)

	// Simulate an execution error response
	updated, _ = model.Update(execMsg{
		mod:      module.Manifest{Name: "keyboard"},
		response: nil,
		err:      &testError{msg: "connection refused"},
	})
	model = updated.(Model)

	if model.statusBar == nil {
		t.Fatal("statusBar should be initialized")
	}
	view := model.statusBar.View()
	if !strings.Contains(view, "connection refused") {
		t.Errorf("statusBar = %q, want to contain error message", view)
	}
}

func createIntegrationTestModel() Model {
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
