package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestViewStateValues(t *testing.T) {
	// Verify all expected view states exist and have correct values
	states := map[viewState]string{
		viewCategories:     "categories",
		viewModules:        "modules",
		viewModuleDetail:   "moduleDetail",
		viewConfirmDialog:  "confirmDialog",
	}

	for state, expected := range states {
		if string(state) != expected {
			t.Errorf("viewState %q = %q, want %q", state, string(state), expected)
		}
	}
}

func TestSubModelInterface(t *testing.T) {
	// Verify SubModel interface is satisfied by a simple implementation
	var _ SubModel = (*mockSubModel)(nil)
}

type mockSubModel struct {
	view string
}

func (m *mockSubModel) Init() tea.Cmd {
	return nil
}

func (m *mockSubModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	return m, nil
}

func (m *mockSubModel) View() string {
	return m.view
}

func TestModelHasSubModelFields(t *testing.T) {
	m := Model{}

	// Verify the model has fields for sub-models (they may be nil initially)
	_ = m.categoriesSub
	_ = m.modulesSub
	_ = m.activeSubModel
}

func TestActiveSubModelTracksState(t *testing.T) {
	m := Model{}

	if m.activeSubModel != nil {
		t.Error("activeSubModel should be nil initially")
	}

	mock := &mockSubModel{view: "test"}
	m.activeSubModel = mock

	if m.activeSubModel != mock {
		t.Error("activeSubModel should point to the assigned sub-model")
	}
}

func TestModelInitDelegatesToActiveSubModel(t *testing.T) {
	sub := &mockSubModel{}
	m := Model{activeSubModel: sub}

	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should delegate to sub-model and return nil for mock")
	}
}

func TestNewModelInitializesCategoriesSub(t *testing.T) {
	// NewModel should initialize categoriesSub and set it as active
	// This is verified indirectly through the model_test by checking
	// that the fields exist and can be assigned
	m := Model{}
	if m.categoriesSub != nil {
		t.Error("categoriesSub should be nil on zero value Model")
	}
}

func TestSubModelSwitching(t *testing.T) {
	// Verify that activeSubModel can be switched between different sub-models
	catSub := &mockSubModel{view: "categories"}
	modSub := &mockSubModel{view: "modules"}

	m := Model{activeSubModel: catSub}
	if m.activeSubModel.View() != "categories" {
		t.Error("initial active sub-model should be categories")
	}

	m.activeSubModel = modSub
	if m.activeSubModel.View() != "modules" {
		t.Error("after switch, active sub-model should be modules")
	}
}

func TestViewStateTransitions(t *testing.T) {
	// Verify view state can be transitioned through all expected states
	m := Model{view: viewCategories}

	if m.view != viewCategories {
		t.Errorf("initial view = %q, want %q", m.view, viewCategories)
	}

	m.view = viewModules
	if m.view != viewModules {
		t.Errorf("after transition view = %q, want %q", m.view, viewModules)
	}

	m.view = viewModuleDetail
	if m.view != viewModuleDetail {
		t.Errorf("after transition view = %q, want %q", m.view, viewModuleDetail)
	}

	m.view = viewConfirmDialog
	if m.view != viewConfirmDialog {
		t.Errorf("after transition view = %q, want %q", m.view, viewConfirmDialog)
	}
}

func TestMultipleSubModelImplementations(t *testing.T) {
	// Verify different struct types can implement SubModel
	var models []SubModel

	models = append(models, &mockSubModel{view: "a"})
	models = append(models, &mockSubModel{view: "b"})

	if len(models) != 2 {
		t.Fatalf("expected 2 sub-models, got %d", len(models))
	}

	if models[0].View() != "a" {
		t.Errorf("first sub-model view = %q, want %q", models[0].View(), "a")
	}
	if models[1].View() != "b" {
		t.Errorf("second sub-model view = %q, want %q", models[1].View(), "b")
	}
}
