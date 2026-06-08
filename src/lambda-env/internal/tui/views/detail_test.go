package views

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/pkg/module"
)

func TestDetailViewInitialState(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "toggle-feature", Label: "Enable Feature", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	if v.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", v.cursor)
	}
	if len(v.manifest.Actions) != 1 {
		t.Errorf("actions count = %d, want 1", len(v.manifest.Actions))
	}
}

func TestDetailViewEmptyActions(t *testing.T) {
	mod := module.Manifest{Name: "keyboard"}
	v := NewDetailView(mod, icons.NewProvider(false))

	view := v.View()
	if !strings.Contains(view, "No actions") {
		t.Errorf("empty view = %q, want to contain 'No actions'", view)
	}
}

func TestDetailViewRendersToggle(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "toggle-feature", Label: "Enable Feature", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	view := v.View()
	if !strings.Contains(view, "Enable Feature") {
		t.Errorf("view = %q, want to contain 'Enable Feature'", view)
	}
}

func TestDetailViewRendersSelect(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "select-theme", Label: "Theme", Type: "select", Options: []string{"dark", "light"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	view := v.View()
	if !strings.Contains(view, "Theme") {
		t.Errorf("view = %q, want to contain 'Theme'", view)
	}
	if !strings.Contains(view, "dark") {
		t.Errorf("view = %q, want to contain 'dark'", view)
	}
}

func TestDetailViewRendersText(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "set-variant", Label: "Variant", Type: "text"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	view := v.View()
	if !strings.Contains(view, "Variant") {
		t.Errorf("view = %q, want to contain 'Variant'", view)
	}
}

func TestDetailViewRendersConfirm(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "reset", Label: "Reset Settings", Type: "confirm"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	view := v.View()
	if !strings.Contains(view, "Reset Settings") {
		t.Errorf("view = %q, want to contain 'Reset Settings'", view)
	}
}

func TestDetailViewRendersExecute(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply Changes", Type: "execute"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	view := v.View()
	if !strings.Contains(view, "Apply Changes") {
		t.Errorf("view = %q, want to contain 'Apply Changes'", view)
	}
}

func TestDetailViewDownNavigation(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "a", Label: "A", Type: "toggle"},
			{Name: "b", Label: "B", Type: "toggle"},
			{Name: "c", Label: "C", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	dv := updated.(*DetailView)
	if dv.cursor != 1 {
		t.Errorf("cursor after down = %d, want 1", dv.cursor)
	}
}

func TestDetailViewUpNavigation(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "a", Label: "A", Type: "toggle"},
			{Name: "b", Label: "B", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	v.cursor = 1
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	dv := updated.(*DetailView)
	if dv.cursor != 0 {
		t.Errorf("cursor after up = %d, want 0", dv.cursor)
	}
}

func TestDetailViewWrapAroundDown(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "a", Label: "A", Type: "toggle"},
			{Name: "b", Label: "B", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	v.cursor = 1
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	dv := updated.(*DetailView)
	if dv.cursor != 0 {
		t.Errorf("cursor after wrap down = %d, want 0", dv.cursor)
	}
}

func TestDetailViewWrapAroundUp(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "a", Label: "A", Type: "toggle"},
			{Name: "b", Label: "B", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	v.cursor = 0
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	dv := updated.(*DetailView)
	if dv.cursor != 1 {
		t.Errorf("cursor after wrap up = %d, want 1", dv.cursor)
	}
}

func TestDetailViewToggleStateTracking(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "toggle-feature", Label: "Feature", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	if v.states[0].toggleOn != false {
		t.Errorf("initial toggle = %v, want false", v.states[0].toggleOn)
	}

	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeySpace})
	dv := updated.(*DetailView)
	if !dv.states[0].toggleOn {
		t.Errorf("toggle after space = %v, want true", dv.states[0].toggleOn)
	}
	if cmd == nil {
		t.Fatal("expected cmd after toggle, got nil")
	}

	msg := cmd()
	execMsg, ok := msg.(ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Name != "keyboard" {
		t.Errorf("msg.Name = %q, want %q", execMsg.Name, "keyboard")
	}
	if execMsg.Action != "toggle-feature" {
		t.Errorf("msg.Action = %q, want %q", execMsg.Action, "toggle-feature")
	}
}

func TestDetailViewSelectStateTracking(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "select-theme", Label: "Theme", Type: "select", Options: []string{"dark", "light", "nord"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	if v.states[0].selectIndex != 0 {
		t.Errorf("initial selectIndex = %d, want 0", v.states[0].selectIndex)
	}

	// Right should advance selection
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyRight})
	dv := updated.(*DetailView)
	if dv.states[0].selectIndex != 1 {
		t.Errorf("selectIndex after right = %d, want 1", dv.states[0].selectIndex)
	}

	// Left should go back
	updated, _ = dv.Update(tea.KeyMsg{Type: tea.KeyLeft})
	dv = updated.(*DetailView)
	if dv.states[0].selectIndex != 0 {
		t.Errorf("selectIndex after left = %d, want 0", dv.states[0].selectIndex)
	}
}

func TestDetailViewTextStateTracking(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "set-variant", Label: "Variant", Type: "text"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	if v.states[0].textValue != "" {
		t.Errorf("initial text = %q, want empty", v.states[0].textValue)
	}

	// Focus text input first
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv := updated.(*DetailView)
	if !dv.states[0].textFocused {
		t.Errorf("textFocused after enter = %v, want true", dv.states[0].textFocused)
	}

	// Type some characters
	updated, _ = dv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	dv = updated.(*DetailView)
	updated, _ = dv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	dv = updated.(*DetailView)
	updated, _ = dv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	dv = updated.(*DetailView)

	if dv.states[0].textValue != "dvo" {
		t.Errorf("textValue after typing = %q, want 'dvo'", dv.states[0].textValue)
	}
}

func TestDetailViewExecuteEmitsMessage(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "apply", Label: "Apply", Type: "execute"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd after enter on execute, got nil")
	}

	msg := cmd()
	execMsg, ok := msg.(ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Name != "keyboard" {
		t.Errorf("msg.Name = %q, want %q", execMsg.Name, "keyboard")
	}
	if execMsg.Action != "apply" {
		t.Errorf("msg.Action = %q, want %q", execMsg.Action, "apply")
	}

	dv := updated.(*DetailView)
	if dv.lastExecutedAction != "apply" {
		t.Errorf("lastExecutedAction = %q, want %q", dv.lastExecutedAction, "apply")
	}
}

func TestDetailViewConfirmShowsDialog(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "reset", Label: "Reset", Type: "confirm"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv := updated.(*DetailView)
	if !dv.showingConfirm {
		t.Errorf("showingConfirm after enter = %v, want true", dv.showingConfirm)
	}
}

func TestDetailViewConfirmDialogEmitsExecuteOnYes(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "reset", Label: "Reset", Type: "confirm"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	// Show confirm dialog
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv := updated.(*DetailView)

	// Press Enter to confirm "Yes"
	updated, cmd := dv.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv = updated.(*DetailView)
	if dv.showingConfirm {
		t.Error("showingConfirm should be false after confirming")
	}
	if cmd == nil {
		t.Fatal("expected cmd after confirm yes, got nil")
	}

	msg := cmd()
	execMsg, ok := msg.(ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Action != "reset" {
		t.Errorf("msg.Action = %q, want %q", execMsg.Action, "reset")
	}
}

func TestDetailViewConfirmDialogDismissedOnNo(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "reset", Label: "Reset", Type: "confirm"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv := updated.(*DetailView)

	// Move to "No" and press Enter
	dv.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated, cmd := dv.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dv = updated.(*DetailView)
	if dv.showingConfirm {
		t.Error("showingConfirm should be false after dismissing")
	}
	if cmd == nil {
		t.Fatal("expected cmd after confirm no, got nil")
	}

	msg := cmd()
	execMsg, ok := msg.(ActionExecuteMsg)
	if !ok {
		t.Fatalf("expected ActionExecuteMsg, got %T", msg)
	}
	if execMsg.Params["confirmed"] != false {
		t.Errorf("confirmed = %v, want false", execMsg.Params["confirmed"])
	}
}

func TestDetailViewBackOnEsc(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "a", Label: "A", Type: "toggle"},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected cmd after esc, got nil")
	}

	msg := cmd()
	_, ok := msg.(BackMsg)
	if !ok {
		t.Fatalf("expected BackMsg, got %T", msg)
	}

	// If text is focused, esc should blur instead of back
	v2 := NewDetailView(module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "text", Label: "Text", Type: "text"},
		},
	}, icons.NewProvider(false))
	v2.Update(tea.KeyMsg{Type: tea.KeyEnter}) // focus text
	updated, cmd = v2.Update(tea.KeyMsg{Type: tea.KeyEsc})
	dv := updated.(*DetailView)
	if dv.states[0].textFocused {
		t.Error("textFocused should be false after esc on focused text")
	}
	if cmd != nil {
		t.Error("expected nil cmd when blurring text, got non-nil")
	}
}

func TestDetailViewMergeDynamicOptions(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us", "dvorak"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))

	// Merge dynamic options that expand the list
	v.MergeDynamicOptions(
		map[string][]string{"layout": {"us", "dvorak", "colemak", "workman"}},
		map[string]interface{}{"layout": "colemak"},
	)

	if len(v.Manifest().Actions[0].Options) != 4 {
		t.Errorf("expected 4 options after merge, got %d", len(v.Manifest().Actions[0].Options))
	}
	if v.states[0].selectIndex != 2 {
		t.Errorf("expected selectIndex=2 (colemak), got %d", v.states[0].selectIndex)
	}
}

func TestDetailViewMergeDynamicOptionsFallbackToStatic(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us", "dvorak"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))

	// Merge with no dynamic options for this action — should keep static
	v.MergeDynamicOptions(
		map[string][]string{},
		map[string]interface{}{},
	)

	if len(v.Manifest().Actions[0].Options) != 2 {
		t.Errorf("expected 2 static options, got %d", len(v.Manifest().Actions[0].Options))
	}
}

func TestDetailViewDynamicOptionsErrorShowsWarning(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us", "dvorak"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	v.SetWarning("Dynamic options unavailable — using static list")

	view := v.View()
	if !strings.Contains(view, "unavailable") {
		t.Errorf("view = %q, want to contain warning about unavailable options", view)
	}
}

func TestDetailViewHandlesDynamicOptionsMsg(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, _ := v.Update(DynamicOptionsMsg{
		Options: map[string][]string{"layout": {"us", "dvorak", "colemak"}},
		Values:  map[string]interface{}{"layout": "dvorak"},
		Err:     nil,
	})
	dv := updated.(*DetailView)

	if len(dv.Manifest().Actions[0].Options) != 3 {
		t.Errorf("expected 3 options after DynamicOptionsMsg, got %d", len(dv.Manifest().Actions[0].Options))
	}
	if dv.states[0].selectIndex != 1 {
		t.Errorf("expected selectIndex=1 (dvorak), got %d", dv.states[0].selectIndex)
	}
}

func TestDetailViewHandlesDynamicOptionsMsgError(t *testing.T) {
	mod := module.Manifest{
		Name: "keyboard",
		Actions: []module.ActionConfig{
			{Name: "layout", Label: "Layout", Type: "select", Options: []string{"us"}},
		},
	}

	v := NewDetailView(mod, icons.NewProvider(false))
	updated, _ := v.Update(DynamicOptionsMsg{
		Options: nil,
		Values:  nil,
		Err:     fmt.Errorf("module timed out"),
	})
	dv := updated.(*DetailView)

	if dv.warning == "" {
		t.Error("expected warning after error DynamicOptionsMsg")
	}
	if !strings.Contains(dv.warning, "timed out") {
		t.Errorf("warning = %q, want to contain 'timed out'", dv.warning)
	}
}

func TestRenderActionLineWithIcons(t *testing.T) {
	tests := []struct {
		name         string
		nerdFonts    bool
		actType      string
		toggleOn     bool
		wantContains string
	}{
		{"toggle off nerd", true, "toggle", false, "\uf204"},
		{"toggle on nerd", true, "toggle", true, "\uf205"},
		{"toggle off fallback", false, "toggle", false, "\u25cb"},
		{"toggle on fallback", false, "toggle", true, "\u25cf"},
		{"select nerd", true, "select", false, "\uf0da"},
		{"select fallback", false, "select", false, "\u25b6"},
		{"confirm nerd", true, "confirm", false, "\uf059"},
		{"confirm fallback", false, "confirm", false, "\u2753"},
		{"execute nerd", true, "execute", false, "\uf04b"},
		{"execute fallback", false, "execute", false, "\u25b6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := module.Manifest{
				Name: "test",
				Actions: []module.ActionConfig{
					{Name: "act", Label: "Action", Type: tt.actType, Options: []string{"a", "b"}},
				},
			}
			p := icons.NewProvider(tt.nerdFonts)
			v := NewDetailView(mod, p)
			if tt.toggleOn {
				v.states[0].toggleOn = true
			}
			line := v.renderActionLine(0, v.manifest.Actions[0])
			if !strings.Contains(line, tt.wantContains) {
				t.Errorf("line = %q, want to contain %q", line, tt.wantContains)
			}
			if !strings.Contains(line, "Action") {
				t.Errorf("line = %q, want to contain 'Action'", line)
			}
		})
	}
}

func TestRenderActionLineLoadingState(t *testing.T) {
	mod := module.Manifest{
		Name: "test",
		Actions: []module.ActionConfig{
			{Name: "run", Label: "Run", Type: "execute"},
		},
	}
	p := icons.NewProvider(true)
	v := NewDetailView(mod, p)
	v.executing = []bool{true}
	line := v.renderActionLine(0, v.manifest.Actions[0])
	if !strings.Contains(line, "\uf021") {
		t.Errorf("line = %q, want to contain loading icon", line)
	}
	if !strings.Contains(line, "Run") {
		t.Errorf("line = %q, want to contain 'Run'", line)
	}
}

func TestRenderActionLineSuccessState(t *testing.T) {
	mod := module.Manifest{
		Name: "test",
		Actions: []module.ActionConfig{
			{Name: "run", Label: "Run", Type: "execute"},
		},
	}
	p := icons.NewProvider(true)
	v := NewDetailView(mod, p)
	v.lastResult = []actionResult{{status: "success", timestamp: time.Now()}}
	line := v.renderActionLine(0, v.manifest.Actions[0])
	if !strings.Contains(line, "\uf00c") {
		t.Errorf("line = %q, want to contain success icon", line)
	}
	if !strings.Contains(line, "Run") {
		t.Errorf("line = %q, want to contain 'Run'", line)
	}
}

func TestRenderActionLineErrorState(t *testing.T) {
	mod := module.Manifest{
		Name: "test",
		Actions: []module.ActionConfig{
			{Name: "run", Label: "Run", Type: "execute"},
		},
	}
	p := icons.NewProvider(true)
	v := NewDetailView(mod, p)
	v.lastResult = []actionResult{{status: "error", timestamp: time.Now()}}
	line := v.renderActionLine(0, v.manifest.Actions[0])
	if !strings.Contains(line, "\uf06a") {
		t.Errorf("line = %q, want to contain error icon", line)
	}
	if !strings.Contains(line, "Run") {
		t.Errorf("line = %q, want to contain 'Run'", line)
	}
}
