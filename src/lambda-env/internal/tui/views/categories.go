package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lambdaos.dev/lambda-env/internal/hub"
)

var (
	categoryTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				MarginBottom(1)

	categoryItemStyle = lipgloss.NewStyle().
				PaddingLeft(2)

	categorySelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7D56F4"))
)

// CategorySelectedMsg is emitted when the user selects a category.
type CategorySelectedMsg struct {
	Category string
	Index    int
}

// CategoriesView is a sub-model for the category list screen.
type CategoriesView struct {
	categories       []string
	menu             []hub.MenuCategory
	cursor           int
	selectedCategory string
}

// NewCategoriesView creates a new CategoriesView.
func NewCategoriesView(cats []string, menu []hub.MenuCategory) *CategoriesView {
	return &CategoriesView{
		categories: cats,
		menu:       menu,
		cursor:     0,
	}
}

// Init implements tea.Model.
func (c *CategoriesView) Init() tea.Cmd {
	return nil
}

// Update handles user input for the categories view.
func (c *CategoriesView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if c.cursor > 0 {
				c.cursor--
			} else {
				c.cursor = len(c.categories) - 1
			}
		case tea.KeyDown:
			if c.cursor < len(c.categories)-1 {
				c.cursor++
			} else {
				c.cursor = 0
			}
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				switch msg.Runes[0] {
				case 'k':
					if c.cursor > 0 {
						c.cursor--
					} else {
						c.cursor = len(c.categories) - 1
					}
				case 'j':
					if c.cursor < len(c.categories)-1 {
						c.cursor++
					} else {
						c.cursor = 0
					}
				}
			}
		case tea.KeyEnter:
			if len(c.categories) > 0 && c.cursor < len(c.categories) {
				c.selectedCategory = c.categories[c.cursor]
				return c, func() tea.Msg {
					return CategorySelectedMsg{
						Category: c.selectedCategory,
						Index:    c.cursor,
					}
				}
			}
		}
	}
	return c, nil
}

// View renders the categories list.
func (c *CategoriesView) View() string {
	var b strings.Builder

	b.WriteString(categoryTitleStyle.Render("LambdaOS Settings"))
	b.WriteString("\n\n")

	if len(c.categories) == 0 {
		b.WriteString("No modules found.\n")
		return b.String()
	}

	catCount := make(map[string]int)
	for _, m := range c.menu {
		catCount[m.Name] = m.Count
	}

	for i, cat := range c.categories {
		cursor := "  "
		if c.cursor == i {
			cursor = "> "
		}
		count := catCount[cat]
		line := fmt.Sprintf("%s%s (%d)", cursor, cat, count)
		if c.cursor == i {
			b.WriteString(categorySelectedStyle.Render(line))
		} else {
			b.WriteString(categoryItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	return b.String()
}
