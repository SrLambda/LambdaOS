# Spec: core/01-hub-plugin-system

## Intent

Establish the `lambda-env` Go binary as the TUI hub with a plugin system that discovers, validates, and executes modules via a JSON-over-stdout protocol. The hub renders a navigable menu organized by categories and manages module lifecycle from discovery through logging.

## Scope

### In Scope
- Go module initialization at `src/lambda-env/` with bubbletea and bubbles dependencies
- Module discovery from system and user paths (`/usr/share/lambda-env/modules/`, `~/.local/share/lambda-env/modules/`)
- `manifest.json` parsing and validation against required fields, including the `actions` array for interactive widgets
- Menu rendering with categories: System, Apps, Ops, Setup (3-level navigation: categories → modules → module detail)
- Keyboard navigation: arrows, Enter, Esc, q (context-sensitive based on view state)
- Module execution with environment variables (`LAMBDA_ENV_ACTION` per-action, `LAMBDA_ENV_SETTINGS`, `LAMBDA_ENV_HUB_VERSION`, `LAMBDA_ENV_LOCALE`)
- JSON stdout parsing from module responses (ok/error/warning), with settings_delta merge for interactive modules
- Error handling with logging to `/var/log/lambda-env/modules.log`
- Dependency checking before module execution (`pacman -Q`)
- Root detection for modules requiring sudo
- Settings delta merging from module responses
- Component-based TUI architecture with sub-model delegation (categories, modules, detail, confirm views)
- Dynamic widget rendering from manifest actions (toggle, select, text, confirm, execute types)
- CLIExecutor interface for module testing (RealExecutor wraps os/exec, MockExecutor returns fixtures)

### Out of Scope
- Module implementations (screen, audio, network, etc.) — Wave 2+; system modules (keyboard, appearance, audio, defaults) are Wave 3
- PKGBUILD for `lambdaos-tui` — Wave 2
- HTML/GUI prototypes — Wave 1 is infrastructure only

## Requirements

### Requirement 1: Go Module Structure

The system SHALL initialize a Go module at `src/lambda-env/` with module path `lambdaos.dev/lambda-env` and include `bubbletea`, `lipgloss`, and `bubbles` as dependencies.

#### Scenario: Go module initializes successfully

- GIVEN the directory `src/lambda-env/` does not contain `go.mod`
- WHEN `go mod init lambdaos.dev/lambda-env` is executed
- THEN `go.mod` exists with the correct module path
- AND `go get` for bubbletea, lipgloss, and bubbles succeeds

#### Scenario: Go project builds without errors

- GIVEN `go.mod` and source files exist under `src/lambda-env/`
- WHEN `go build ./...` is executed from `src/lambda-env/`
- THEN the build completes with exit code 0
- AND no compilation errors are reported

### Requirement 2: Module Discovery

The hub SHALL scan `/usr/share/lambda-env/modules/` (system) and `~/.local/share/lambda-env/modules/` (user) for directories containing a valid `manifest.json`.

#### Scenario: System modules are discovered

- GIVEN three module directories exist under `/usr/share/lambda-env/modules/` each with a valid `manifest.json`
- WHEN the hub starts
- THEN all three modules are loaded into the module registry

#### Scenario: User modules override system modules

- GIVEN a module `screen` exists in both system and user paths
- WHEN the hub performs discovery
- THEN the user `screen` module is used and the system version is ignored

#### Scenario: Invalid modules are skipped with warning

- GIVEN a module directory contains a malformed `manifest.json` (invalid JSON)
- WHEN the hub scans that directory
- THEN the module is skipped
- AND a warning is logged to stderr

#### Scenario: Empty module directories are ignored

- GIVEN a directory exists under a module path but contains no `manifest.json`
- WHEN the hub scans
- THEN the directory is ignored and no error is raised

### Requirement 3: Manifest Validation

The hub SHALL validate every `manifest.json` against required fields: `name`, `version`, `description`, `description_es`, `category`, `requires_root`, `dependencies`, `min_hub_version`.

#### Scenario: Valid manifest passes validation

- GIVEN a `manifest.json` contains all required fields with correct types
- WHEN the hub validates the manifest
- THEN validation passes and the module is registered

#### Scenario: Missing required field fails validation

- GIVEN a `manifest.json` is missing the `category` field
- WHEN the hub validates the manifest
- THEN validation fails
- AND the module is not registered

#### Scenario: Invalid category value fails validation

- GIVEN a `manifest.json` has `category` set to `"invalid"`
- WHEN the hub validates the manifest
- THEN validation fails (valid values: `system`, `apps`, `ops`, `setup`)

#### Scenario: Name format validation

- GIVEN a `manifest.json` has `name` set to `"My Module"` (contains spaces)
- WHEN the hub validates the manifest
- THEN validation fails (name must be lowercase with hyphens only)

### Requirement 4: Menu Rendering by Category

The hub SHALL render a navigable menu grouping modules by category: System, Apps, Ops, Setup. Modules within each category are sorted alphabetically by name. The menu SHALL support 3-level navigation with category selection as the first level.

#### Scenario: Menu displays all categories with modules

- GIVEN four modules exist, one per category
- WHEN the hub renders the main menu
- THEN four category sections are displayed
- AND each section contains its module sorted alphabetically

#### Scenario: Empty categories are hidden

- GIVEN no modules exist in the `apps` category
- WHEN the hub renders the main menu
- THEN the `apps` category section is not displayed

#### Scenario: Module count per category is shown

- GIVEN three modules exist in the `system` category
- WHEN the hub renders the menu
- THEN the system category header shows "(3)" or equivalent count indicator

#### Scenario: Category selection navigates to module list

- GIVEN the category list is displayed
- WHEN the user presses Enter on a category
- THEN the hub navigates to the module list for that category
- AND the category name is shown as the current context

### Requirement 5: Keyboard Navigation

The hub SHALL support navigation via arrow keys (up/down), Enter (select/confirm), Esc (back/quit current view), and q (quit application). Key bindings SHALL be context-sensitive based on the current view state.

#### Scenario: Arrow keys move selection

- GIVEN the menu is displayed with 5 items
- WHEN the user presses down arrow 3 times
- THEN the selection highlight moves to the 4th item

#### Scenario: Enter selects or confirms based on context

- GIVEN a module is highlighted in the module list
- WHEN the user presses Enter
- THEN the hub navigates to the module detail view

#### Scenario: Esc returns to previous view

- GIVEN the user is viewing a module execution result
- WHEN the user presses Esc
- THEN the hub returns to the main menu

#### Scenario: q quits the application

- GIVEN the hub is running at any view
- WHEN the user presses q
- THEN the application exits with code 0

#### Scenario: Context-specific key bindings shown in help

- GIVEN the user is on the module detail view with toggle widgets
- WHEN the user presses `?` for help
- THEN the help overlay shows Space/Enter for toggling, not just navigation

### Requirement 6: Module Execution Protocol

The hub SHALL execute module actions by passing the action name via `LAMBDA_ENV_ACTION` environment variable. For interactive modules, the hub SHALL execute actions individually as the user triggers them, rather than executing a single `run` action.

The hub SHALL execute the selected module's `module` executable with environment variables: `LAMBDA_ENV_ACTION=<action-name>`, `LAMBDA_ENV_SETTINGS` (path to settings.json), `LAMBDA_ENV_HUB_VERSION`, `LAMBDA_ENV_LOCALE`.

#### Scenario: Module executes with specific action

- GIVEN a module `keyboard` has actions `set-layout` and `set-variant`
- WHEN the user triggers the `set-layout` action
- THEN the hub executes the module with `LAMBDA_ENV_ACTION=set-layout`
- AND the module receives the action-specific context

#### Scenario: Module executes with correct environment

- GIVEN a module `screen` is selected
- WHEN the hub executes the module
- THEN the process receives `LAMBDA_ENV_ACTION=run`
- AND `LAMBDA_ENV_SETTINGS` points to `~/.config/lambdaos/settings.json`
- AND `LAMBDA_ENV_HUB_VERSION` matches the hub's version

#### Scenario: Module timeout is enforced

- GIVEN a module has `"timeout": 10` in its manifest
- WHEN the module runs longer than 10 seconds
- THEN the hub terminates the process
- AND an error is displayed to the user

#### Scenario: Module default timeout is 30 seconds

- GIVEN a module manifest does not specify a `timeout` field
- WHEN the module runs
- THEN the hub enforces a 30-second default timeout

### Requirement 7: JSON Response Parsing

The hub SHALL parse JSON from module stdout and handle three response types: `ok` (exit 0), `error` (exit 1), `warning` (exit 2). For interactive modules, the response MAY include `settings_delta` that the hub merges immediately.

#### Scenario: Success response is rendered

- GIVEN a module outputs `{"status":"ok","action":"run","data":{"outputs":[]}}` to stdout and exits 0
- WHEN the hub parses the response
- THEN the data is extracted and rendered in the TUI

#### Scenario: Error response triggers error display

- GIVEN a module outputs `{"status":"error","code":"FAIL","message":"something broke"}` to stdout and exits 1
- WHEN the hub parses the response
- THEN an error view is displayed with the message
- AND the error is logged to `/var/log/lambda-env/modules.log`

#### Scenario: Warning response shows confirmation prompt

- GIVEN a module outputs `{"status":"warning","code":"DEP_MISSING","message":"xrandr not installed","suggestion":"pacman -S xorg-xrandr"}` and exits 2
- WHEN the hub parses the response
- THEN a warning is displayed with the suggestion
- AND the user is prompted to continue or return to menu

#### Scenario: Non-JSON stdout is handled gracefully

- GIVEN a module outputs plain text (not valid JSON) to stdout
- WHEN the hub attempts to parse
- THEN the hub displays a parse error
- AND logs the raw stdout content

#### Scenario: Interactive action response updates widget state

- GIVEN a toggle action returns `{"status":"ok","action":"toggle-setting","settings_delta":{"keyboard":{"layout":"us"}}}`
- WHEN the hub processes the response
- THEN the settings_delta is merged into settings.json
- AND the toggle widget state updates to reflect the new value

### Requirement 8: Error Logging

The hub SHALL log all module execution errors to `/var/log/lambda-env/modules.log` with timestamp, module name, action, exit code, stdout, stderr, and environment variables.

#### Scenario: Error is logged with full context

- GIVEN a module `screen` fails with exit code 1
- WHEN the hub handles the error
- THEN `/var/log/lambda-env/modules.log` contains a line with timestamp, `module=screen`, `exit_code=1`
- AND the stdout JSON is logged
- AND the stderr text is logged

#### Scenario: Log directory is created if missing

- GIVEN `/var/log/lambda-env/` does not exist
- WHEN the hub attempts to log an error
- THEN the directory is created with appropriate permissions
- AND the log file is written

### Requirement 9: Dependency Checking

The hub SHALL check module dependencies before execution by verifying each package in `dependencies` is installed via `pacman -Q <package>`.

#### Scenario: All dependencies are satisfied

- GIVEN a module declares `dependencies: ["xorg-xrandr"]`
- AND `xorg-xrandr` is installed
- WHEN the hub checks dependencies
- THEN execution proceeds normally

#### Scenario: Missing dependency blocks execution

- GIVEN a module declares `dependencies: ["xorg-xrandr"]`
- AND `xorg-xrandr` is NOT installed
- WHEN the hub checks dependencies
- THEN execution is blocked
- AND the user is shown which packages are missing

#### Scenario: Multiple missing dependencies are listed

- GIVEN a module declares `dependencies: ["xorg-xrandr", "xorg-xprop"]`
- AND neither package is installed
- WHEN the hub checks dependencies
- THEN both missing packages are listed in the error message

### Requirement 10: Root Detection

The hub SHALL check `requires_root` in the module manifest and verify the current user has sudo privileges before executing modules that require root.

#### Scenario: Root module executes with sudo

- GIVEN a module has `"requires_root": true`
- AND the user has valid sudo access
- WHEN the module is selected
- THEN the hub executes it via `sudo` or equivalent privilege escalation

#### Scenario: Root module blocked without sudo

- GIVEN a module has `"requires_root": true`
- AND the user does NOT have sudo access
- WHEN the module is selected
- THEN execution is blocked
- AND an error message explains root is required

#### Scenario: Non-root module runs without sudo

- GIVEN a module has `"requires_root": false`
- WHEN the module is selected
- THEN the hub executes it directly without privilege escalation

### Requirement 11: Settings Delta Merging

The hub SHALL extract `settings_delta` from successful module responses and merge it into `~/.config/lambdaos/settings.json` via atomic write.

#### Scenario: Settings delta is merged on success

- GIVEN a module response contains `"settings_delta": {"display": {"active_profile": "home"}}`
- WHEN the hub processes the response
- THEN the delta is merged into settings.json
- AND only the `display.active_profile` field is updated (other fields preserved)

#### Scenario: Response without delta does not modify settings

- GIVEN a module response has no `settings_delta` field
- WHEN the hub processes the response
- THEN settings.json is not modified

### Requirement 12: Three-Level Navigation Support

The hub SHALL support navigation across 3+ levels: categories → modules → module detail view. The viewState enum SHALL expand to include `categories`, `modules`, `moduleDetail`, and `confirmDialog` states.

#### Scenario: Hub transitions to module detail view

- GIVEN a module is selected from the module list
- WHEN the user presses Enter
- THEN the hub transitions to `moduleDetail` viewState
- AND the module detail is rendered from the module's manifest actions

#### Scenario: Hub displays confirm dialog state

- GIVEN a module action is marked as requiring confirmation
- WHEN the user triggers that action
- THEN the hub transitions to `confirmDialog` viewState
- AND the confirm dialog is rendered with action details

#### Scenario: Hub maintains view stack for back navigation

- GIVEN the user has navigated: categories → modules → moduleDetail
- WHEN the user presses Esc at moduleDetail
- THEN the hub returns to modules view
- AND the module list selection is preserved

### Requirement 13: Manifest Actions Field

The manifest.json SHALL support an `actions` field that declares supported actions with their types, labels, and metadata for dynamic widget rendering.

#### Scenario: Manifest with actions is parsed

- GIVEN a manifest.json contains an `actions` array with typed entries
- WHEN the hub parses the manifest
- THEN each action is validated for required fields (name, type, label)
- AND the actions are stored in the module registry

#### Scenario: Action types are validated

- GIVEN a manifest action has an invalid type (not one of: toggle, select, text, confirm, execute)
- WHEN the hub validates the manifest
- THEN validation fails with an error about invalid action type
- AND the module is not registered

#### Scenario: Manifest without actions field is backward compatible

- GIVEN a manifest.json does not include the `actions` field
- WHEN the hub parses the manifest
- THEN the module is still registered (backward compatible)
- AND the module detail view shows a "no interactive actions" message

### Requirement 14: Dynamic Widget Rendering from Actions

The hub SHALL render appropriate TUI widgets based on the action types declared in the module manifest.

#### Scenario: Toggle widget rendered for boolean action

- GIVEN a module action has `type: "toggle"`
- WHEN the module detail view is rendered
- THEN a toggle widget is displayed for that action
- AND the current value is shown (on/off)

#### Scenario: Select list rendered for enum action

- GIVEN a module action has `type: "select"` with an `options` array
- WHEN the module detail view is rendered
- THEN a selectable list is displayed with the options
- AND the current value is highlighted

#### Scenario: Text input rendered for string action

- GIVEN a module action has `type: "text"`
- WHEN the module detail view is rendered
- THEN a text input widget is displayed
- AND the current value is pre-filled

#### Scenario: Confirm button rendered for confirm action

- GIVEN a module action has `type: "confirm"`
- WHEN the module detail view is rendered
- THEN a button/trigger is displayed that opens a confirm dialog when activated

## Technical Details

### Package Structure
```
src/lambda-env/
├── go.mod
├── go.sum
├── cmd/
│   └── lambda-env/
│       └── main.go              # Entry point, CLI flags
├── internal/
│   ├── hub/
│   │   ├── hub.go               # Hub struct, main loop
│   │   ├── discovery.go         # Module scanning, manifest parsing
│   │   ├── execution.go         # Module execution, ExecuteAction, JSON parsing, logging
│   │   └── version.go           # Hub version constant
│   ├── tui/
│   │   ├── model.go             # Parent model with sub-model delegation
│   │   ├── update.go            # Parent update dispatches to active sub-model
│   │   ├── view.go              # Parent view renders active sub-model + status bar
│   │   ├── components/
│   │   │   ├── toggle.go        # Boolean toggle widget
│   │   │   ├── textinput.go     # Text input widget (wraps bubbles/textinput)
│   │   │   ├── confirm.go       # Confirm dialog widget
│   │   │   ├── help.go          # Full-screen help overlay (custom)
│   │   │   └── statusbar.go     # Persistent status bar widget
│   │   └── views/
│   │       ├── categories.go    # Category list sub-model
│   │       ├── modules.go       # Module list sub-model
│   │       └── detail.go        # Module detail sub-model (dynamic widgets from manifest.actions)
│   └── module/
│       ├── types.go             # JSON protocol types (Manifest, Response)
│       └── logger.go            # Module log writer
└── pkg/
    ├── module/
    │   ├── manifest.go          # Manifest + ActionConfig structs, validation
    │   └── executor.go          # CLIExecutor interface, RealExecutor, MockExecutor
    └── version/
        └── version.go           # Hub version constant (moved from internal/)
```

### Manifest JSON Schema
```json
{
  "name": "string (lowercase, hyphens)",
  "version": "string (semver)",
  "description": "string",
  "description_es": "string",
  "category": "enum: system|apps|ops|setup",
  "requires_root": "boolean",
  "dependencies": "string[]",
  "min_hub_version": "string (semver)",
  "actions": [
    {
      "name": "string (hyphenated action name)",
      "label": "string (display label for TUI)",
      "type": "enum: toggle|select|text|confirm|execute",
      "field": "string (settings field path)",
      "options": ["string[] (optional, for select type)"],
      "requires_root": "boolean (optional)"
    }
  ]
}
```

### Module Response JSON Schema
```json
{
  "status": "enum: ok|error|warning",
  "action": "string",
  "data": "object (optional)",
  "code": "string (optional, for error/warning)",
  "message": "string (optional)",
  "message_es": "string (optional)",
  "suggestion": "string (optional)",
  "settings_delta": "object (optional)"
}
```

### Log Format
```
2026-05-31T00:00:00Z [LEVEL] module=<name> action=<action> exit_code=<n>
  stdout: <raw stdout>
  stderr: <raw stderr>
  env: LAMBDA_ENV_ACTION=<action>, LAMBDA_ENV_LOCALE=<locale>
```

## Dependencies

- Go 1.21+ (already in `packages.x86_64`)
- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/charmbracelet/lipgloss` — styling
- `github.com/charmbracelet/bubbles` — list/input components
- `pacman` — dependency checking via `pacman -Q`
- `sudo` — privilege escalation for root modules

## Verification Steps

```bash
# 1. Go module initializes and builds
cd src/lambda-env && go mod init lambdaos.dev/lambda-env
cd src/lambda-env && go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/bubbles
cd src/lambda-env && go build ./...

# 2. Go linting passes
cd src/lambda-env && go vet ./...

# 3. Go unit tests pass
cd src/lambda-env && go test ./... -v

# 4. Hub binary runs and exits cleanly
cd src/lambda-env && go build -o /tmp/lambda-env ./cmd/lambda-env
/tmp/lambda-env --help  # Should show help text

# 5. Module discovery with test fixtures
mkdir -p /tmp/test-modules/system/screen
echo '{"name":"screen","version":"0.1.0","description":"test","description_es":"test","category":"system","requires_root":false,"dependencies":[],"min_hub_version":"0.1.0"}' > /tmp/test-modules/system/screen/manifest.json
# Hub should discover and list the screen module

# 6. Manifest validation rejects invalid input
echo '{"name":"bad"}' > /tmp/test-bad/manifest.json
# Hub should skip the invalid module with warning

# 7. Dependency check works
# Module with dependencies: ["nonexistent-pkg-123"] should be blocked

# 8. Log file is created on error
ls -la /var/log/lambda-env/modules.log  # Should exist after error scenario
```
