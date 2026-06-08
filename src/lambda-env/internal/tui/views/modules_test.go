package views

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestModulesViewInitialState(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))

	if v.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", v.cursor)
	}
	if len(v.modules) != 2 {
		t.Errorf("modules count = %d, want 2", len(v.modules))
	}
	if v.category != "system" {
		t.Errorf("category = %q, want %q", v.category, "system")
	}
}

func TestModulesViewEmptyList(t *testing.T) {
	v := NewModulesView([]module.Manifest{}, "apps", icons.NewProvider(false))

	view := v.View()
	if !strings.Contains(view, "No modules") {
		t.Errorf("empty view = %q, want to contain 'No modules'", view)
	}
	if !strings.Contains(view, "apps") {
		t.Errorf("empty view = %q, want to contain category name 'apps'", view)
	}
}

func TestModulesViewDownNavigation(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
		{Name: "appearance", Description: "Set theme"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	mv := updated.(*ModulesView)
	if mv.cursor != 1 {
		t.Errorf("cursor after down = %d, want 1", mv.cursor)
	}

	updated, _ = mv.Update(tea.KeyMsg{Type: tea.KeyDown})
	mv = updated.(*ModulesView)
	if mv.cursor != 2 {
		t.Errorf("cursor after second down = %d, want 2", mv.cursor)
	}
}

func TestModulesViewUpNavigation(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	v.cursor = 1

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	mv := updated.(*ModulesView)
	if mv.cursor != 0 {
		t.Errorf("cursor after up = %d, want 0", mv.cursor)
	}
}

func TestModulesViewWrapAroundDown(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	v.cursor = 1

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	mv := updated.(*ModulesView)
	if mv.cursor != 0 {
		t.Errorf("cursor after wrap down = %d, want 0", mv.cursor)
	}
}

func TestModulesViewWrapAroundUp(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	v.cursor = 0

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	mv := updated.(*ModulesView)
	if mv.cursor != 1 {
		t.Errorf("cursor after wrap up = %d, want 1", mv.cursor)
	}
}

func TestModulesViewSelectEmitsMessage(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))

	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd after enter, got nil")
	}

	msg := cmd()
	selected, ok := msg.(ModuleSelectedMsg)
	if !ok {
		t.Fatalf("expected ModuleSelectedMsg, got %T", msg)
	}
	if selected.Module.Name != "keyboard" {
		t.Errorf("selected module = %q, want %q", selected.Module.Name, "keyboard")
	}
	if selected.Index != 0 {
		t.Errorf("selected index = %d, want 0", selected.Index)
	}

	mv := updated.(*ModulesView)
	if mv.selectedModule != "keyboard" {
		t.Errorf("sub-model selected module = %q, want %q", mv.selectedModule, "keyboard")
	}
}

func TestModulesViewJKeyNavigation(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	mv := updated.(*ModulesView)
	if mv.cursor != 1 {
		t.Errorf("cursor after 'j' = %d, want 1", mv.cursor)
	}
}

func TestModulesViewKKeyNavigation(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	v.cursor = 1

	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	mv := updated.(*ModulesView)
	if mv.cursor != 0 {
		t.Errorf("cursor after 'k' = %d, want 0", mv.cursor)
	}
}

func TestModulesViewShowsModuleInfo(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	view := v.View()

	if !strings.Contains(view, "keyboard") {
		t.Errorf("view = %q, want to contain 'keyboard'", view)
	}
	if !strings.Contains(view, "Set keyboard layout") {
		t.Errorf("view = %q, want to contain description", view)
	}
	if !strings.Contains(view, "system") {
		t.Errorf("view = %q, want to contain category 'system'", view)
	}
}

func TestModulesViewSelectedModule(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
		{Name: "audio", Description: "Configure audio"},
		{Name: "appearance", Description: "Set theme"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))
	v.cursor = 1

	view := v.View()
	if !strings.Contains(view, "audio") {
		t.Errorf("view = %q, want to contain 'audio'", view)
	}
}

func TestModulesViewBackNavigation(t *testing.T) {
	mods := []module.Manifest{
		{Name: "keyboard", Description: "Set keyboard layout"},
	}

	v := NewModulesView(mods, "system", icons.NewProvider(false))

	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected cmd after esc, got nil")
	}

	msg := cmd()
	_, ok := msg.(BackMsg)
	if !ok {
		t.Fatalf("expected BackMsg, got %T", msg)
	}

	mv := updated.(*ModulesView)
	if mv.cursor != 0 {
		t.Errorf("cursor after back = %d, want 0", mv.cursor)
	}
}

func TestModulesViewRendersWithIcons(t *testing.T) {
	mods := []module.Manifest{
		{Name: "audio", Description: "Configure audio"},
		{Name: "display", Description: "Set display"},
	}

	tests := []struct {
		name      string
		nerdFonts bool
		wantIcon  string
	}{
		{"nerd mode", true, "\uf028"},
		{"fallback mode", false, "\u266a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := icons.NewProvider(tt.nerdFonts)
			v := NewModulesView(mods, "system", p)
			view := v.View()
			if !strings.Contains(view, tt.wantIcon) {
				t.Errorf("view = %q, want to contain icon %q", view, tt.wantIcon)
			}
			if !strings.Contains(view, "audio") {
				t.Errorf("view = %q, want to contain 'audio'", view)
			}
		})
	}
}
