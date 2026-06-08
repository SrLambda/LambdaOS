# lambda-env TUI — Architecture Guide

## Overview

```
cmd/lambda-env/main.go          ← Entry point
    │
    ├─ internal/hub/             ← Discovers modules, executes binaries
    │   ├── hub.go               ← New(), BuildMenu()
    │   ├── discovery.go         ← Scan() — reads manifest.json from filesystem
    │   └── execution.go         ← ExecuteModule(), ExecuteAction()
    │
    ├─ internal/tui/             ← Bubbletea TUI (the framework)
    │   ├── model.go             ← Model struct, viewState, SubModel interface
    │   ├── update.go            ← Message router, delegation to sub-models
    │   ├── view.go              ← Render delegated to sub-models + statusBar
    │   │
    │   ├── components/          ← Reusable widgets (5)
    │   │   ├── toggle.go        ← On/Off with ✓/○
    │   │   ├── textinput.go     ← Wrap bubbles/textinput with validation
    │   │   ├── confirm.go       ← Full-screen Yes/No dialog
    │   │   ├── help.go          ← Key bindings overlay (? / esc)
    │   │   └── statusbar.go     ← Persistent bottom bar
    │   │
    │   └── views/               ← Sub-models per view (3)
    │       ├── categories.go    ← Category list (system, apps, ops, setup)
    │       ├── modules.go       ← Module list per category
    │       └── detail.go        ← Detail view with dynamic widgets
    │
    ├─ internal/modules/         ← 7 independent modules (Go binaries)
    │   ├── neovim/
    │   ├── qtile/
    │   ├── dotfiles/
    │   ├── keyboard/
    │   ├── appearance/
    │   ├── audio/
    │   └── defaults/
    │
    └── pkg/module/              ← Shared types
        ├── manifest.go          ← Manifest, ActionConfig, Response
        └── executor.go          ← CLIExecutor interface + mock
```

## Execution Flow

### 1. Startup (`cmd/lambda-env/main.go`)

```go
h, err := hub.New(settingsPath)   // Load settings.json + discover modules
m := tui.NewModel(h)              // Create Model with categories from Hub
p := tea.NewProgram(m)            // Bubbletea takes over
p.Run()
```

The Hub in `New()`:
1. Loads `~/.config/lambdaos/settings.json` (with automatic v1.0.0→v1.1.0 migration)
2. Scans `/usr/share/lambda-env/modules/` and `~/.local/share/lambda-env/modules/`
3. Reads `manifest.json` from each module and parses it

### 2. Navigation (`internal/tui/`)

The `Model` has a `viewState` that controls which screen is displayed:

```go
type viewState string
const (
    viewCategories     // Category list
    viewModules        // Module list in a category
    viewModuleDetail   // Detail view with widgets
    viewConfirmDialog  // Confirmation dialog
)
```

`Update()` is a **router**: receives `tea.KeyMsg` and delegates to the active sub-model:

```
Update(msg):
  ├─ viewCategories   → categoriesView.Update(msg)
  ├─ viewModules      → modulesView.Update(msg)
  ├─ viewModuleDetail → detailView.Update(msg)
  └─ viewConfirmDialog → confirmDialog.Update(msg)
```

**View transitions**:
```
Categories ──Enter──→ Modules ──Enter──→ Detail
     ↑                   ↑                   │
     └──── Esc ──────────┘                   │ Action
                                             ↓
                                       ConfirmDialog
```

### 3. Sub-models (`internal/tui/views/`)

Each view is a struct implementing `SubModel` (embeds `tea.Model`):

```go
type SubModel interface {
    tea.Model  // Init(), Update(msg tea.Msg) (tea.Model, tea.Cmd), View() string
}
```

**CategoriesView**: Lists the 4 categories (`system`, `apps`, `ops`, `setup`) with module count. On Enter emits `CategorySelectedMsg`.

**ModulesView**: Lists modules in the selected category with `name — description`. On Enter emits `ModuleSelectedMsg`.

**DetailView**: The most complex view. On receiving `ModuleSelectedMsg`:
1. Reads `manifest.actions[]` from the module
2. For each action, creates a widget based on its `Type`:

| `action.Type` | Widget Created | Interaction |
|---------------|---------------|-------------|
| `"toggle"` | `components.Toggle` | Space/Enter flip |
| `"select"` | Option list with cursor | ↑↓ navigate, Enter select |
| `"text"` | `components.TextInput` | Type, Enter submit, Esc cancel |
| `"confirm"` | `components.Confirm` | ←→ Yes/No, Enter confirm |
| `"execute"` | Simple button | Enter executes |

3. On load, runs the module in background (`run` action) to get dynamic options (Option C)
4. On pressing Enter on an action, emits `ActionExecuteMsg`

### 4. Module Communication

The Hub does **not call Go functions** — it executes external binaries and communicates via **JSON on stdout**:

```
TUI (detail view)
  │  User presses Enter on "Enable LSP"
  │
  ├─→ Emits ActionExecuteMsg{Name: "toggle-lsp", Action: "toggle-lsp"}
  │
  ├─→ hub.ExecuteAction(manifest, "toggle-lsp", params)
  │     │
  │     └─→ os/exec: ./neovim/module
  │           env: LAMBDA_ENV_ACTION=toggle-lsp
  │           env: LAMBDA_ENV_SETTINGS=/home/user/.config/lambdaos/settings.json
  │           │
  │           └─→ STDOUT: {"status":"ok","action":"toggle-lsp",
  │                         "settings_delta":{"neovim":{"enable_lsp":true}},
  │                         "message":"LSP enabled"}
  │
  ├─→ Hub parses JSON, merges settings_delta into settings.json
  │
  └─→ execMsg → TUI updates status bar + widget state
```

**Why external binaries and not imports**: Each module is an independent `package main`. The Hub has no code dependency on modules. A new module only needs a `manifest.json` and a `module` binary — no Hub recompilation needed.

### 5. Settings Delta

Modules **never write settings.json directly**. They emit a partial `settings_delta` and the Hub does the atomic merge:

```json
// Module emits:
{"settings_delta": {"audio": {"volume": 80, "muted": false}}}

// Hub merges with existing settings.json:
{"audio": {"volume": 80, "muted": false, "default_sink": "alsa_output.pci-0000"}}
//         ↑ new             ↑ new               ↑ preserved
```

`SaveDelta()` uses deep merge + atomic rename (temp file → rename).

### 6. Schema Migration

When loading `settings.json`, if the version is lower than `1.1.0`:

```
v1.0.0 settings.json
    │
    ├─→ detects version < "1.1.0"
    ├─→ adds 7 new sections with defaults
    ├─→ adds use_global_theme: true to neovim and qtile
    ├─→ bumps version to "1.1.0"
    └─→ saves atomically
```

## Key Types

### ActionConfig (`pkg/module/manifest.go`)

```go
type ActionConfig struct {
    Name         string   `json:"name"`          // "toggle-lsp", "set-browser"
    Label        string   `json:"label"`         // "Enable LSP"
    Type         string   `json:"type"`          // toggle|select|text|confirm|execute
    Field        string   `json:"field"`         // "neovim.enable_lsp"
    Options      []string `json:"options,omitempty"` // for select type
    RequiresRoot bool     `json:"requires_root,omitempty"`
}
```

### Response (`pkg/module/manifest.go`)

```go
type Response struct {
    Status        string                 `json:"status"`         // "ok" | "error" | "warning"
    Action        string                 `json:"action"`
    Data          map[string]interface{} `json:"data,omitempty"` // available_options, current_value
    Message       string                 `json:"message,omitempty"`
    SettingsDelta map[string]interface{} `json:"settings_delta,omitempty"`
}
```

### CLIExecutor (`pkg/module/executor.go`)

```go
type CLIExecutor interface {
    Run(command string, args ...string) (stdout, stderr string, exitCode int, err error)
}
```

## Conventions

| Convention | Example |
|-----------|---------|
| Action names | **hyphens** — `toggle-lsp`, `set-terminal`, `check-conflicts` |
| Struct names | `ActionConfig`, `CLIExecutor`, `SubModel` |
| Messages | `ToggleChangedMsg`, `ConfirmResultMsg`, `ActionExecuteMsg` |
| TextInput | **Pointer receiver** — `*TextInput` (bubbles mutates internal state) |
| CLI testing | `MockExecutor` — no real tools required |

## How to Add a New Module

```bash
# 1. Create directory
mkdir -p internal/modules/my-module/

# 2. Create manifest.json
cat > internal/modules/my-module/manifest.json << 'EOF'
{
  "name": "my-module",
  "version": "0.1.0",
  "description": "Configure something",
  "description_es": "Configurar algo",
  "category": "system",
  "dependencies": [],
  "min_hub_version": "1.0.0",
  "actions": [
    {"name": "run", "label": "Refresh", "type": "execute"},
    {"name": "my-action", "label": "My Action", "type": "toggle", "field": "my_module.my_field"}
  ]
}
EOF

# 3. Create main.go (package main, func main(), reads LAMBDA_ENV_ACTION)
# 4. Build: go build -o module .
# 5. Copy to ~/.local/share/lambda-env/modules/my-module/
# 6. Done — the Hub discovers it automatically
```

## Test File Conventions

| Type | Location | Example |
|------|----------|---------|
| Unit test | Same package, `_test.go` | `components/toggle_test.go` |
| Integration test | `test/` or same package | `test/settings_delta_flow_test.go` |
| E2E test | `internal/tui/` with `tea.NewProgram` | `e2e_navigation_test.go` |
| Build test | `test/build_test.go` | Verifies `go build` works |

## Theme Sync Architecture

The `use_global_theme` flag (default: `true`) controls whether a module follows the global theme:

```json
{
  "appearance": {"theme": "dark"},
  "neovim": {"use_global_theme": true},   // derives theme from appearance
  "qtile": {"use_global_theme": false}    // uses its own color_scheme
}
```

**Mapping table**:
| appearance.theme | neovim theme | qtile color_scheme |
|-----------------|-------------|---------------------|
| dark | tokyonight | dracula |
| light | tokyonight-light | nord-light |
| nord | nord | nord |
| catppuccin | catppuccin-mocha | catppuccin |

Sync is **push-based**: when the appearance module changes the theme, it emits a `settings_delta` that includes the neovim and qtile theme fields (if their `use_global_theme` is `true`).

## Option C — Dynamic Options

Select actions have two sources of options:

1. **Static** (manifest.json): `"options": ["firefox", "chromium", "brave", "chrome"]`
2. **Dynamic** (module `run` action): `data.available_options: {"browsers": ["firefox", "chromium", "brave", "chrome", "edge"]}`

When entering a detail view, the TUI:
1. Shows static options immediately
2. Runs module `run` in background
3. Merges dynamic options with static (module wins)
4. Shows merged options in the widget

If the module fails to respond, static options are used as fallback.
