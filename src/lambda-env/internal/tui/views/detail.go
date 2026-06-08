package views

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lambdaos.dev/lambda-env/internal/tui/components"
	"lambdaos.dev/lambda-env/internal/tui/icons"
	"lambdaos.dev/lambda-env/internal/tui/theme"
	"lambdaos.dev/lambda-env/pkg/module"
)

var (
	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(theme.Accent)).
				MarginBottom(1)

	detailItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	detailSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(theme.Accent))

	detailCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Accent))
)

// ActionExecuteMsg is emitted when the user triggers an action.
type ActionExecuteMsg struct {
	Module module.Manifest
	Name   string
	Action string
	Params map[string]interface{}
}

// DynamicOptionsMsg carries the result of a background module query.
type DynamicOptionsMsg struct {
	Options map[string][]string
	Values  map[string]interface{}
	Err     error
}

// widgetState holds the runtime state for a single action widget.
type widgetState struct {
	toggleOn    bool
	selectIndex int
	textValue   string
	textFocused bool
}

// actionResult holds the result of an executed action for feedback display.
type actionResult struct {
	status    string
	timestamp time.Time
}

// DetailView is a sub-model for the module detail screen.
type DetailView struct {
	manifest           module.Manifest
	cursor             int
	states             []widgetState
	textInputs         []*components.TextInput
	showingConfirm     bool
	confirmDialog      *components.Confirm
	lastExecutedAction string
	lastExecutedIndex  int
	warning            string
	iconProvider       icons.IconProvider
	executing          []bool
	lastResult         []actionResult
}

// NewDetailView creates a new DetailView for the given manifest.
func NewDetailView(mod module.Manifest, provider icons.IconProvider) *DetailView {
	states := make([]widgetState, len(mod.Actions))
	textInputs := make([]*components.TextInput, len(mod.Actions))
	executing := make([]bool, len(mod.Actions))
	lastResult := make([]actionResult, len(mod.Actions))

	for i, a := range mod.Actions {
		if a.Type == "text" {
			ti := components.NewTextInput(a.Label)
			ti.SetPlaceholder(a.Label)
			textInputs[i] = ti
		}
	}

	return &DetailView{
		manifest:     mod,
		cursor:       0,
		states:       states,
		textInputs:   textInputs,
		iconProvider: provider,
		executing:    executing,
		lastResult:   lastResult,
	}
}

// Init implements tea.Model.
func (d *DetailView) Init() tea.Cmd {
	return nil
}

// Update handles user input for the detail view.
func (d *DetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If confirm dialog is showing, delegate to it.
	if d.showingConfirm && d.confirmDialog != nil {
		updated, cmd := d.confirmDialog.Update(msg)
		d.confirmDialog = updated
		if result, ok := msg.(components.ConfirmResultMsg); ok {
			d.showingConfirm = false
			return d, d.emitConfirmAction(d.manifest.Actions[d.cursor].Name, result.Confirmed)
		}
		// Also handle if the update itself produced a result message
		if cmd != nil {
			msg := cmd()
			if result, ok := msg.(components.ConfirmResultMsg); ok {
				d.showingConfirm = false
				return d, d.emitConfirmAction(d.manifest.Actions[d.cursor].Name, result.Confirmed)
			}
			return d, cmd
		}
		return d, nil
	}

	switch msg := msg.(type) {
	case DynamicOptionsMsg:
		if msg.Err != nil {
			d.SetWarning(fmt.Sprintf("Dynamic options unavailable: %v — using static list", msg.Err))
			return d, nil
		}
		d.MergeDynamicOptions(msg.Options, msg.Values)
		return d, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			d.moveCursor(-1)
			return d, nil
		case tea.KeyDown:
			d.moveCursor(1)
			return d, nil
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				r := msg.Runes[0]
				// If text input is focused, delegate typing to it.
				if d.cursor < len(d.states) && d.states[d.cursor].textFocused {
					if d.textInputs[d.cursor] != nil {
						updated, cmd := d.textInputs[d.cursor].Update(msg)
						d.textInputs[d.cursor] = updated
						d.states[d.cursor].textValue = updated.Value()
						return d, cmd
					}
				}
				switch r {
				case 'k':
					d.moveCursor(-1)
					return d, nil
				case 'j':
					d.moveCursor(1)
					return d, nil
				}
			}
		case tea.KeyEsc:
			// If text is focused, blur it instead of going back.
			if d.cursor < len(d.states) && d.states[d.cursor].textFocused {
				d.states[d.cursor].textFocused = false
				if d.textInputs[d.cursor] != nil {
					d.textInputs[d.cursor].Blur()
				}
				return d, nil
			}
			return d, func() tea.Msg { return BackMsg{} }
		case tea.KeySpace, tea.KeyEnter:
			return d.handleAction()
		case tea.KeyLeft:
			if d.cursor < len(d.manifest.Actions) {
				act := d.manifest.Actions[d.cursor]
				if (act.Type == "select" || act.Type == "list") && d.cursor < len(d.states) {
					d.states[d.cursor].selectIndex--
					if d.states[d.cursor].selectIndex < 0 {
						d.states[d.cursor].selectIndex = len(act.Options) - 1
					}
					return d, nil
				}
			}
		case tea.KeyRight:
			if d.cursor < len(d.manifest.Actions) {
				act := d.manifest.Actions[d.cursor]
				if (act.Type == "select" || act.Type == "list") && d.cursor < len(d.states) {
					d.states[d.cursor].selectIndex++
					if d.states[d.cursor].selectIndex >= len(act.Options) {
						d.states[d.cursor].selectIndex = 0
					}
					return d, nil
				}
			}
		}

		// If text input is focused, delegate typing to it.
		if d.cursor < len(d.states) && d.states[d.cursor].textFocused {
			if d.textInputs[d.cursor] != nil {
				updated, cmd := d.textInputs[d.cursor].Update(msg)
				d.textInputs[d.cursor] = updated
				d.states[d.cursor].textValue = updated.Value()
				return d, cmd
			}
		}
	}

	return d, nil
}

func (d *DetailView) moveCursor(delta int) {
	if len(d.manifest.Actions) == 0 {
		return
	}
	d.cursor += delta
	if d.cursor < 0 {
		d.cursor = len(d.manifest.Actions) - 1
	} else if d.cursor >= len(d.manifest.Actions) {
		d.cursor = 0
	}
}

func (d *DetailView) handleAction() (tea.Model, tea.Cmd) {
	if d.cursor >= len(d.manifest.Actions) {
		return d, nil
	}
	act := d.manifest.Actions[d.cursor]

	switch act.Type {
	case "toggle":
		d.states[d.cursor].toggleOn = !d.states[d.cursor].toggleOn
		return d, d.emitAction(act.Name, d.states[d.cursor].toggleOn)

	case "select", "list":
		selected := ""
		if d.states[d.cursor].selectIndex < len(act.Options) {
			selected = act.Options[d.states[d.cursor].selectIndex]
		}
		return d, d.emitAction(act.Name, selected)

	case "text":
		if d.states[d.cursor].textFocused {
			// Submit text value
			d.states[d.cursor].textFocused = false
			if d.textInputs[d.cursor] != nil {
				d.textInputs[d.cursor].Blur()
			}
			return d, d.emitAction(act.Name, d.states[d.cursor].textValue)
		}
		// Focus text input
		d.states[d.cursor].textFocused = true
		if d.textInputs[d.cursor] != nil {
			d.textInputs[d.cursor].Focus()
		}
		return d, nil

	case "confirm":
		d.showingConfirm = true
		d.confirmDialog = components.NewConfirm(act.Label)
		return d, nil

	case "execute":
		return d, d.emitAction(act.Name, nil)
	}

	return d, nil
}

func (d *DetailView) emitAction(actionName string, value interface{}) tea.Cmd {
	d.lastExecutedAction = actionName
	d.lastExecutedIndex = d.cursor
	params := make(map[string]interface{})
	if value != nil {
		params["value"] = value
	}
	return func() tea.Msg {
		return ActionExecuteMsg{
			Module: d.manifest,
			Name:   d.manifest.Name,
			Action: actionName,
			Params: params,
		}
	}
}

func (d *DetailView) emitConfirmAction(actionName string, confirmed bool) tea.Cmd {
	d.lastExecutedAction = actionName
	d.lastExecutedIndex = d.cursor
	params := map[string]interface{}{
		"confirmed": confirmed,
	}
	return func() tea.Msg {
		return ActionExecuteMsg{
			Module: d.manifest,
			Name:   d.manifest.Name,
			Action: actionName,
			Params: params,
		}
	}
}

// SetExecuting sets the executing state for the given action index.
func (d *DetailView) SetExecuting(idx int, executing bool) {
	if idx >= 0 && idx < len(d.executing) {
		d.executing[idx] = executing
	}
}

// SetActionResult sets the result status for the given action index.
func (d *DetailView) SetActionResult(idx int, status string) {
	if idx >= 0 && idx < len(d.lastResult) {
		d.lastResult[idx] = actionResult{status: status, timestamp: time.Now()}
	}
}

// LastExecutedIndex returns the index of the last executed action.
func (d *DetailView) LastExecutedIndex() int {
	return d.lastExecutedIndex
}

// View renders the detail view.
func (d *DetailView) View() string {
	var b strings.Builder

	b.WriteString(detailTitleStyle.Render(d.iconProvider.ForModule(d.manifest.Name) + " " + d.manifest.Name))
	b.WriteString("\n\n")

	if len(d.manifest.Actions) == 0 {
		b.WriteString("No actions available for this module.\n")
		return b.String()
	}

	for i, act := range d.manifest.Actions {
		cursor := "  "
		if d.cursor == i {
			cursor = detailCursorStyle.Render("> ")
		}

		line := d.renderActionLine(i, act)
		if d.cursor == i {
			b.WriteString(cursor + detailSelectedStyle.Render(line))
		} else {
			b.WriteString(cursor + detailItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	if d.warning != "" {
		b.WriteString("\n")
		warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warn))
		b.WriteString(warningStyle.Render("⚠ " + d.warning))
		b.WriteString("\n")
	}

	// Render confirm dialog inline if visible
	if d.showingConfirm && d.confirmDialog != nil {
		b.WriteString("\n")
		b.WriteString(d.confirmDialog.View())
	}

	return b.String()
}

func (d *DetailView) renderActionLine(idx int, act module.ActionConfig) string {
	ico := d.iconProvider
	width := ico.Width()
	pad := strings.Repeat(" ", width-1)

	var prefix strings.Builder

	// Loading state: prepend spinner icon
	if d.executing[idx] {
		prefix.WriteString(ico.ForWidget("spinner"))
		prefix.WriteString(pad)
		prefix.WriteString(" ")
	}

	// Success/error feedback: prepend status icon for 2 seconds
	if d.lastResult[idx].status != "" && time.Since(d.lastResult[idx].timestamp) < 2*time.Second {
		if d.lastResult[idx].status == "success" {
			prefix.WriteString(theme.SuccessStyle.Render(ico.ForWidget("success")))
		} else {
			prefix.WriteString(theme.ErrorStyle.Render(ico.ForWidget("error")))
		}
		prefix.WriteString(pad)
		prefix.WriteString(" ")
	}

	var body string
	switch act.Type {
	case "toggle":
		onOff := ico.ForWidget("toggle_off")
		stateText := "Off"
		if d.states[idx].toggleOn {
			onOff = ico.ForWidget("toggle_on")
			stateText = "On"
		}
		body = fmt.Sprintf("%s  %s %s", act.Label, onOff, stateText)

	case "select", "list":
		selected := ""
		if d.states[idx].selectIndex < len(act.Options) {
			selected = act.Options[d.states[idx].selectIndex]
		}
		body = fmt.Sprintf("%s  ◄ %s ►", act.Label, selected)

	case "text":
		placeholder := act.Label
		if d.textInputs[idx] != nil {
			placeholder = d.textInputs[idx].View()
		}
		body = fmt.Sprintf("%s  %s", act.Label, placeholder)

	case "confirm":
		body = fmt.Sprintf("%s  [Press Enter to confirm]", act.Label)

	case "execute":
		body = fmt.Sprintf("%s  [Press Enter to execute]", act.Label)

	default:
		body = act.Label
	}

	// Prepend widget type icon
	widgetIcon := ico.ForWidget(act.Type)
	return fmt.Sprintf("%s%s%s %s", prefix.String(), widgetIcon, pad, body)
}

// MergeDynamicOptions merges dynamic options from a module response into the view state.
func (d *DetailView) MergeDynamicOptions(options map[string][]string, values map[string]interface{}) {
	for i, act := range d.manifest.Actions {
		if opts, ok := options[act.Name]; ok && len(opts) > 0 {
			// Replace options for select/list actions
			if act.Type == "select" || act.Type == "list" {
				d.manifest.Actions[i].Options = opts
				// Reset selection to 0 or preserve if current value matches
				if val, vok := values[act.Name]; vok {
					for j, o := range opts {
						if o == fmt.Sprintf("%v", val) {
							d.states[i].selectIndex = j
							break
						}
					}
				} else {
					d.states[i].selectIndex = 0
				}
			}
		}
		if val, ok := values[act.Name]; ok {
			if act.Type == "text" {
				d.states[i].textValue = fmt.Sprintf("%v", val)
				if d.textInputs[i] != nil {
					d.textInputs[i].SetValue(d.states[i].textValue)
				}
			} else if act.Type == "toggle" {
				d.states[i].toggleOn = val == true || val == "true"
			}
		}
	}
}

// Manifest returns the current manifest (with any merged dynamic options).
func (d *DetailView) Manifest() module.Manifest {
	return d.manifest
}

// SetWarning sets a warning message to display in the view.
func (d *DetailView) SetWarning(msg string) {
	d.warning = msg
}
