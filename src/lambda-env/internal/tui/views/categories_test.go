package views

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lambdaos.dev/lambda-env/internal/hub"
	"lambdaos.dev/lambda-env/internal/tui/icons"
)

func TestCategoriesViewInitialState(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))

	if v.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", v.cursor)
	}
	if len(v.categories) != 2 {
		t.Errorf("categories count = %d, want 2", len(v.categories))
	}
}

func TestCategoriesViewEmptyList(t *testing.T) {
	v := NewCategoriesView([]string{}, []hub.MenuCategory{}, icons.NewProvider(false))

	view := v.View()
	if !strings.Contains(view, "No modules found") {
		t.Errorf("empty view = %q, want to contain 'No modules found'", view)
	}
}

func TestCategoriesViewDownNavigation(t *testing.T) {
	cats := []string{"system", "apps", "ops"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
		{Name: "ops", Count: 1},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))

	// Press down
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	cv := updated.(*CategoriesView)
	if cv.cursor != 1 {
		t.Errorf("cursor after down = %d, want 1", cv.cursor)
	}

	// Press down again
	updated, _ = cv.Update(tea.KeyMsg{Type: tea.KeyDown})
	cv = updated.(*CategoriesView)
	if cv.cursor != 2 {
		t.Errorf("cursor after second down = %d, want 2", cv.cursor)
	}
}

func TestCategoriesViewUpNavigation(t *testing.T) {
	cats := []string{"system", "apps", "ops"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
		{Name: "ops", Count: 1},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	v.cursor = 2 // Start at last item

	// Press up
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	cv := updated.(*CategoriesView)
	if cv.cursor != 1 {
		t.Errorf("cursor after up = %d, want 1", cv.cursor)
	}
}

func TestCategoriesViewWrapAroundDown(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	v.cursor = 1 // At last item

	// Press down should wrap to 0
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyDown})
	cv := updated.(*CategoriesView)
	if cv.cursor != 0 {
		t.Errorf("cursor after wrap down = %d, want 0", cv.cursor)
	}
}

func TestCategoriesViewWrapAroundUp(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	v.cursor = 0 // At first item

	// Press up should wrap to last item
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyUp})
	cv := updated.(*CategoriesView)
	if cv.cursor != 1 {
		t.Errorf("cursor after wrap up = %d, want 1", cv.cursor)
	}
}

func TestCategoriesViewSelectEmitsMessage(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))

	updated, cmd := v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd after enter, got nil")
	}

	msg := cmd()
	selected, ok := msg.(CategorySelectedMsg)
	if !ok {
		t.Fatalf("expected CategorySelectedMsg, got %T", msg)
	}
	if selected.Category != "system" {
		t.Errorf("selected category = %q, want %q", selected.Category, "system")
	}
	if selected.Index != 0 {
		t.Errorf("selected index = %d, want 0", selected.Index)
	}

	// Verify the sub-model was updated with selected state
	cv := updated.(*CategoriesView)
	if cv.selectedCategory != "system" {
		t.Errorf("sub-model selected category = %q, want %q", cv.selectedCategory, "system")
	}
}

func TestCategoriesViewJKeyNavigation(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))

	// Press 'j' should act like down
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	cv := updated.(*CategoriesView)
	if cv.cursor != 1 {
		t.Errorf("cursor after 'j' = %d, want 1", cv.cursor)
	}
}

func TestCategoriesViewKKeyNavigation(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	v.cursor = 1

	// Press 'k' should act like up
	updated, _ := v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	cv := updated.(*CategoriesView)
	if cv.cursor != 0 {
		t.Errorf("cursor after 'k' = %d, want 0", cv.cursor)
	}
}

func TestCategoriesViewShowsCategoryCount(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	view := v.View()

	if !strings.Contains(view, "system") {
		t.Errorf("view = %q, want to contain 'system'", view)
	}
	if !strings.Contains(view, "(3)") {
		t.Errorf("view = %q, want to contain '(3)'", view)
	}
}

func TestCategoriesViewSelectedCategory(t *testing.T) {
	cats := []string{"system", "apps", "ops"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
		{Name: "ops", Count: 1},
	}

	v := NewCategoriesView(cats, menu, icons.NewProvider(false))
	v.cursor = 1

	view := v.View()
	if !strings.Contains(view, "apps") {
		t.Errorf("view = %q, want to contain 'apps'", view)
	}
}

func TestCategoriesViewRendersWithIcons(t *testing.T) {
	cats := []string{"system", "apps"}
	menu := []hub.MenuCategory{
		{Name: "system", Count: 3},
		{Name: "apps", Count: 2},
	}

	tests := []struct {
		name      string
		nerdFonts bool
		wantIcon  string
	}{
		{"nerd mode", true, "\uf108"},
		{"fallback mode", false, "\u2699"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := icons.NewProvider(tt.nerdFonts)
			v := NewCategoriesView(cats, menu, p)
			view := v.View()
			if !strings.Contains(view, tt.wantIcon) {
				t.Errorf("view = %q, want to contain icon %q", view, tt.wantIcon)
			}
			if !strings.Contains(view, "system") {
				t.Errorf("view = %q, want to contain 'system'", view)
			}
		})
	}
}
