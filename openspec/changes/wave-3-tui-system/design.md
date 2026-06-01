# Design: Wave 3 — TUI Interface + System Modules

## Technical Approach

Rewrite the flat TUI model into a component-based sub-model architecture using Bubble Tea's delegation pattern. Each view state (`categories`, `modules`, `moduleDetail`, `confirmDialog`) is a self-contained `tea.Model`. The parent `Model` delegates `Update`/`View` to the active sub-model. Module detail views are built dynamically from `manifest.actions[]` — the TUI reads action types and wraps `bubbles` components (textinput, list, viewport, help, key) into LambdaOS composite models. Four new system modules (keyboard, appearance, audio, defaults) follow the established JSON-over-stdout pattern with a shared `CLIExecutor` interface for testability. Schema migrates v1.0.0 → v1.1.0 in one atomic bump adding 7 sections.

## Architecture Decisions

| # | Decision | Options Considered | Chosen | Rationale |
|---|----------|-------------------|--------|-----------|
| 1 | TUI architecture | (A) Flat model + switch cases, (B) Sub-model delegation, (C) Composable view funcs | **B: Sub-model delegation** | Each view maintains independent state (cursor, widgets); flat model unmaintainable at 300+ lines |
| 2 | Widgets from manifest | (A) Hardcoded per module, (B) Dynamic type-switch, (C) Widget registry | **B: Dynamic type-switch** | 5 action types (toggle/select/text/confirm/execute) don't warrant a registry; `bubbles` components wrapped per type |
| 3 | Theme sync mechanism | (A) Push (appearance triggers), (B) Poll (modules check), (C) Inter-module calls | **A: Push via settings_delta** | Existing delta mechanism works. Modules read `use_global_theme` on each action; static mapping table maps appearance.theme → module theme |
| 4 | CLI testing | (A) Real exec only, (B) Interface mocking, (C) Shell script fixtures | **B: CLIExecutor interface** | `Run(cmd, args...) → (stdout, stderr, error)`. Prod wraps `os/exec`, test returns fixtures. All 4 modules inject it |
| 5 | Schema migration | (A) 7 incremental bumps, (B) One atomic bump, (C) No migration | **B: One atomic v1.0.0 → v1.1.0** | 7 sequential migrations would be 7 states to test; atomic bump uses existing `deepMerge` with defaults. Existing store already supports this |
| 6 | TUI-widget state | (A) Widgets own state, (B) TUI holds all state, (C) Module holds all state | **A: Widgets own state** | Each widget sub-model holds its own focused/blur/value state. TUI reads manifest, creates widgets, delegates updates. Widgets emit action requests up |

## Data Flow

### Toggle Setting (e.g., mute toggle in audio)

```
User presses Space → moduleDetail.Update → toggle widget flips local value
  → hub.ExecuteAction(name, action, params) → module binary runs
  → module writes settings_delta → hub.SaveDelta() → module returns response
  → TUI receives execMsg → status bar updates → toggle shows new state
```

### Theme Sync (appearance → neovim)

```
User selects "nord" → appearance module sets theme → delta: {appearance.theme: "nord"}
  → hub merges delta → (later) user opens neovim detail
  → neovim module loads settings → sees use_global_theme=true
  → maps "nord" → "nordic" via lookup table → returns effective theme in response
```

### Text Input (e.g., keyboard variant)

```
User types "dvorak" → textinput widget.Update → chars appear in bubbles textinput
  → User presses Enter → widget emits action(name, value)
  → hub.ExecuteAction("keyboard", "set-variant", "dvorak") → module runs setxkbmap
  → module returns settings_delta → TUI updates status bar
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/tui/model.go` | Modify | Parent model delegates to sub-models; viewState gains moduleDetail, confirmDialog |
| `internal/tui/update.go` | Modify | Parent Update delegates to active sub-model; execMsg for action results |
| `internal/tui/view.go` | Modify | Parent View delegates to active sub-model; status bar rendered by parent |
| `internal/tui/components/toggle.go` | Create | Toggle widget wrapping bubbles help+key styles |
| `internal/tui/components/textinput.go` | Create | Text input wrapping bubbles textinput |
| `internal/tui/components/confirmer.go` | Create | Confirm dialog wrapping bubbles viewport |
| `internal/tui/components/helpoverlay.go` | Create | Help overlay wrapping bubbles help |
| `internal/tui/components/statusbar.go` | Create | Persistent status bar (module, version, modified indicator) |
| `internal/tui/views/categories.go` | Create | Category list sub-model (extracted from current flat model) |
| `internal/tui/views/modules.go` | Create | Module list sub-model |
| `internal/tui/views/detail.go` | Create | Module detail sub-model — reads manifest.actions, builds widgets |
| `pkg/module/manifest.go` | Modify | Add `Actions []Action` field, `Action` struct with type/label/options/required/value |
| `internal/modules/system/executor.go` | Create | `CLIExecutor` interface + `RealExecutor` (os/exec) + `MockExecutor` (fixtures) |
| `internal/modules/keyboard/main.go` | Create | Keyboard module: setxkbmap layout/variant/options via CLIExecutor |
| `internal/modules/appearance/main.go` | Create | Appearance module: gsettings/feh theme/wallpaper via CLIExecutor |
| `internal/modules/audio/main.go` | Create | Audio module: pactl/wpctl volume/mute/sink via CLIExecutor |
| `internal/modules/defaults/main.go` | Create | Defaults module: xdg-mime assignments via CLIExecutor |
| `internal/settings/schema.go` | Modify | Bump CurrentVersion to "1.1.0"; add 7 new section structs; add UseGlobalTheme to Neovim/Qtile |
| `internal/hub/execution.go` | Modify | Add `ExecuteAction(name, action, params)` that passes LAMBDA_ENV_ACTION |

## Interfaces / Contracts

### Manifest Action Type

```go
type Action struct {
    Name        string   `json:"name"`        // "set-layout", "toggle-mute"
    Type        string   `json:"type"`        // toggle|select|text|confirm|execute
    Label       string   `json:"label"`       // "Keyboard Layout"
    Options     []string `json:"options,omitempty"` // for select type
    Required    bool     `json:"required,omitempty"`
    Value       string   `json:"value,omitempty"`   // current value path in settings
    Destructive bool     `json:"destructive,omitempty"` // triggers confirm dialog
}
```

### CLIExecutor Interface

```go
type CLIExecutor interface {
    Run(command string, args ...string) (stdout, stderr string, err error)
}
```

### Theme Mapping Table (appearance module)

```go
var themeMap = map[string]struct{ Neovim, Qtile string }{
    "dark":       {"tokyonight-night", "dracula"},
    "light":      {"tokyonight-day",   "nord-light"},
    "nord":       {"nord",             "nord"},
    "catppuccin": {"catppuccin",       "catppuccin"},
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | TUI sub-models (categories, modules, detail) | Bubbletea `tea.NewProgram(m).Send(msg); assert view output` |
| Unit | TUI widgets (toggle, textinput, confirm, help) | Bubbletea model tests with keypress sequences |
| Unit | CLIExecutor mock | Verify modules call correct commands with correct args |
| Unit | Schema migration v1.0.0 → v1.1.0 | Write v1.0.0 JSON fixture, load, assert 7 new sections with defaults |
| Unit | Theme mapping | Test each global theme → correct neovim/qtile theme |
| Integration | Module settings_delta flow | Write settings, execute action via mock executor, verify delta merge |
| Integration | Manifest action parsing | Parse manifest with actions field, assert widget types correct |

## Migration / Rollout

Schema v1.0.0 → v1.1.0 migration uses existing `deepMerge` infrastructure (already atomic). Migration adds 7 sections with defaults and `use_global_theme: false` to Neovim/Qtile. No user data is overwritten — `deepMerge` preserves existing keys. Feature flag: all new TUI components behind `enableWave3` boolean (default `true`) for rollback.

## Open Questions

None — all 8 questions from proposal were decided.
