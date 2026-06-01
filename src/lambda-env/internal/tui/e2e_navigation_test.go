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

// createE2ETestModel builds a Model with a realistic module catalog for E2E flows.
func createE2ETestModel() Model {
	cats := []string{"system", "apps", "ops"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 4},
		{Name: "apps", Count: 2},
		{Name: "ops", Count: 1},
	}

	modules := []module.Manifest{
		{Name: "keyboard", Description: "Set layout", Category: "system"},
		{Name: "audio", Description: "Configure audio", Category: "system"},
		{Name: "appearance", Description: "Set theme", Category: "system"},
		{Name: "defaults", Description: "Default apps", Category: "system"},
		{Name: "neovim", Description: "Edit config", Category: "apps"},
		{Name: "qtile", Description: "Window manager", Category: "apps"},
		{Name: "dotfiles", Description: "Manage dotfiles", Category: "ops"},
	}

	h := &hub.Hub{Modules: modules}

	m := Model{
		hub:           h,
		categories:    cats,
		categoriesSub: views.NewCategoriesView(cats, menu),
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

func TestE2EFullFlowToggleActionAndBack(t *testing.T) {
	m := createE2ETestModel()

	// categories → modules
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 1: view = %q, want %q", model.view, viewModules)
	}

	// modules → detail view with toggle action
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name:    "audio",
			Actions: []module.ActionConfig{{Name: "set-mute", Label: "Mute", Type: "toggle", Field: "audio.muted"}},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("step 2: view = %q, want %q", model.view, viewModuleDetail)
	}

	// toggle action in detail view
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
	if actionMsg.Action != "set-mute" {
		t.Errorf("action = %q, want %q", actionMsg.Action, "set-mute")
	}
	if val, _ := actionMsg.Params["value"].(bool); val != true {
		t.Errorf("toggle value = %v, want true", val)
	}

	// simulate execution response
	updated, _ = model.Update(execMsg{
		mod:      module.Manifest{Name: "audio"},
		response: &module.Response{Status: "ok", Message: "Mute toggled"},
		err:      nil,
	})
	model = updated.(Model)
	if !strings.Contains(model.statusBar.View(), "Mute toggled") {
		t.Errorf("statusBar = %q, want to contain 'Mute toggled'", model.statusBar.View())
	}

	// back to modules
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewModules {
		t.Fatalf("step 5: view = %q, want %q", model.view, viewModules)
	}

	// back to categories
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewCategories {
		t.Fatalf("step 6: view = %q, want %q", model.view, viewCategories)
	}
}

func TestE2EFullFlowSelectActionChangeTheme(t *testing.T) {
	m := createE2ETestModel()

	// categories → modules → detail
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "appearance",
			Actions: []module.ActionConfig{
				{Name: "set-theme", Label: "Theme", Type: "select", Options: []string{"dark", "light", "nord"}},
			},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("expected detail view, got %q", model.view)
	}

	// Simulate dynamic options and current value
	updated, _ = model.Update(views.DynamicOptionsMsg{
		Options: map[string][]string{"set-theme": {"dark", "light", "nord", "catppuccin"}},
		Values:  map[string]interface{}{"set-theme": "nord"},
	})
	model = updated.(Model)

	// Change selection with right arrow once (nord → catppuccin)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = updated.(Model)

	// Execute the select action
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("expected cmd after select action")
	}
	msg := cmd()
	actionMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if actionMsg.Action != "set-theme" {
		t.Errorf("action = %q, want %q", actionMsg.Action, "set-theme")
	}
	if val, _ := actionMsg.Params["value"].(string); val != "catppuccin" {
		t.Errorf("selected value = %q, want catppuccin", val)
	}

	// back → back
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewModules {
		t.Fatalf("expected modules view, got %q", model.view)
	}
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewCategories {
		t.Fatalf("expected categories view, got %q", model.view)
	}
}

func TestE2EFullFlowTextInputChangeVariant(t *testing.T) {
	m := createE2ETestModel()

	// categories → modules → detail with text action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "keyboard",
			Actions: []module.ActionConfig{
				{Name: "set-variant", Label: "Variant", Type: "text"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("expected detail view, got %q", model.view)
	}

	// Press Enter to focus text input
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if cmd != nil {
		t.Fatal("expected no cmd after focusing text input")
	}
	if model.detailSub == nil {
		t.Fatal("detailSub should exist")
	}

	// Type "dvorak" (bubbletea textinput receives rune messages)
	for _, r := range "dvorak" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(Model)
	}

	// Press Enter to submit
	updated, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("expected cmd after submitting text input")
	}
	msg := cmd()
	actionMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if actionMsg.Action != "set-variant" {
		t.Errorf("action = %q, want %q", actionMsg.Action, "set-variant")
	}
	if val, _ := actionMsg.Params["value"].(string); val != "dvorak" {
		t.Errorf("text value = %q, want dvorak", val)
	}

	// back → back
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewModules {
		t.Fatalf("expected modules view, got %q", model.view)
	}
	updated, _ = model.Update(views.BackMsg{})
	model = updated.(Model)
	if model.view != viewCategories {
		t.Fatalf("expected categories view, got %q", model.view)
	}
}

func TestE2EFullFlowConfirmDialogUnstow(t *testing.T) {
	m := createE2ETestModel()

	// categories → modules → detail with confirm action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "ops", Index: 2})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "dotfiles",
			Actions: []module.ActionConfig{
				{Name: "unstow", Label: "Unstow Dotfiles", Type: "confirm"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("expected detail view, got %q", model.view)
	}

	// Press enter to show confirm dialog
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	view := model.View()
	if !strings.Contains(view, "Unstow Dotfiles") {
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
	actionMsg, ok := msg.(views.ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if actionMsg.Action != "unstow" {
		t.Errorf("action = %q, want %q", actionMsg.Action, "unstow")
	}
	if confirmed, _ := actionMsg.Params["confirmed"].(bool); !confirmed {
		t.Errorf("confirmed = %v, want true", confirmed)
	}
}

func TestE2EHelpOverlayToggleInCategoriesView(t *testing.T) {
	m := createE2ETestModel()
	if m.view != viewCategories {
		t.Fatalf("expected categories view, got %q", m.view)
	}

	// Press '?' to show help
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model := updated.(Model)
	if !model.helpOverlay.Visible {
		t.Error("help overlay should be visible after '?' in categories view")
	}
	view := model.View()
	if !strings.Contains(view, "Toggle help") {
		t.Errorf("view = %q, want to contain help content", view)
	}

	// Press '?' again to hide help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(Model)
	if model.helpOverlay.Visible {
		t.Error("help overlay should be hidden after second '?' in categories view")
	}
}

func TestE2EHelpOverlayToggleInModulesView(t *testing.T) {
	m := createE2ETestModel()

	// Navigate to modules view
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	if model.view != viewModules {
		t.Fatalf("expected modules view, got %q", model.view)
	}

	// Press '?' to show help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(Model)
	if !model.helpOverlay.Visible {
		t.Error("help overlay should be visible after '?' in modules view")
	}

	// Press Esc to dismiss help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(Model)
	if model.helpOverlay.Visible {
		t.Error("help overlay should be hidden after Esc in modules view")
	}
}

func TestE2EHelpOverlayToggleInDetailView(t *testing.T) {
	m := createE2ETestModel()

	// Navigate to detail view
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name:    "audio",
			Actions: []module.ActionConfig{{Name: "set-mute", Label: "Mute", Type: "toggle"}},
		},
		Index: 0,
	})
	model = updated.(Model)
	if model.view != viewModuleDetail {
		t.Fatalf("expected detail view, got %q", model.view)
	}

	// Press '?' to show help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(Model)
	if !model.helpOverlay.Visible {
		t.Error("help overlay should be visible after '?' in detail view")
	}

	// Press '?' again to hide help
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updated.(Model)
	if model.helpOverlay.Visible {
		t.Error("help overlay should be hidden after second '?' in detail view")
	}
}

func TestE2EErrorHandlingModuleExecutionFailure(t *testing.T) {
	m := createE2ETestModel()

	// Navigate to detail view with execute action
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "audio",
			Actions: []module.ActionConfig{
				{Name: "run", Label: "Refresh", Type: "execute"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)

	// Simulate an execution error response
	updated, _ = model.Update(execMsg{
		mod:      module.Manifest{Name: "audio"},
		response: nil,
		err:      &testError{msg: "pactl connection refused"},
	})
	model = updated.(Model)

	if model.statusBar == nil {
		t.Fatal("statusBar should be initialized")
	}
	view := model.statusBar.View()
	if !strings.Contains(view, "pactl connection refused") {
		t.Errorf("statusBar = %q, want to contain error message", view)
	}
	if model.statusType != "error" {
		t.Errorf("statusType = %q, want error", model.statusType)
	}
}

func TestE2EErrorHandlingModuleWarningResponse(t *testing.T) {
	m := createE2ETestModel()

	// Navigate to detail view
	updated, _ := m.Update(views.CategorySelectedMsg{Category: "system", Index: 0})
	model := updated.(Model)
	updated, _ = model.Update(views.ModuleSelectedMsg{
		Module: module.Manifest{
			Name: "audio",
			Actions: []module.ActionConfig{
				{Name: "run", Label: "Refresh", Type: "execute"},
			},
		},
		Index: 0,
	})
	model = updated.(Model)

	// Simulate a warning response
	updated, _ = model.Update(execMsg{
		mod:      module.Manifest{Name: "audio"},
		response: &module.Response{Status: "warning", Message: "Volume already at max"},
		err:      nil,
	})
	model = updated.(Model)

	view := model.statusBar.View()
	if !strings.Contains(view, "Volume already at max") {
		t.Errorf("statusBar = %q, want to contain warning message", view)
	}
	if model.statusType != "warning" {
		t.Errorf("statusType = %q, want warning", model.statusType)
	}
}

func TestE2EEmptyCategoryShowsGracefulMessage(t *testing.T) {
	m := createE2ETestModel()

	// The test model does not have a "setup" category, so selecting it
	// would show no modules. We can simulate this by creating a model
	// with an empty category.
	cats := []string{"system", "setup"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 4},
		{Name: "setup", Count: 0},
	}
	modules := []module.Manifest{
		{Name: "audio", Description: "Audio", Category: "system"},
	}

	h := &hub.Hub{Modules: modules}
	emptyM := Model{
		hub:           h,
		categories:    cats,
		categoriesSub: views.NewCategoriesView(cats, menu),
		view:          viewCategories,
		cursor:        0,
		statusBar:     components.NewStatusBar().SetContext("categories"),
		helpOverlay:   m.helpOverlay,
	}
	emptyM.activeSubModel = emptyM.categoriesSub
	emptyM.helpOverlay.Visible = false

	updated, _ := emptyM.Update(views.CategorySelectedMsg{Category: "setup", Index: 1})
	model := updated.(Model)
	if model.view != viewModules {
		t.Fatalf("expected modules view, got %q", model.view)
	}

	view := model.View()
	if !strings.Contains(view, "No modules") {
		t.Errorf("view = %q, want to contain 'No modules'", view)
	}
}

func TestE2EProgrammaticTeaNewProgramSmoke(t *testing.T) {
	// This test verifies that tea.NewProgram can be instantiated with
	// the TUI model without panic, confirming bubbletea compatibility.
	m := createE2ETestModel()
	p := tea.NewProgram(m, tea.WithoutRenderer())
	if p == nil {
		t.Fatal("tea.NewProgram returned nil")
	}

	// We do not call p.Run() here because the TUI is an interactive
	// program that blocks; our E2E coverage uses direct model.Update
	// which is the idiomatic headless testing pattern for bubbletea.
}
