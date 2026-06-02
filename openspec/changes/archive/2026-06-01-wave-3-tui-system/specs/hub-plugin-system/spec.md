# Delta for hub-plugin-system

## ADDED Requirements

### Requirement: Three-Level Navigation Support

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

### Requirement: Manifest Actions Field

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

### Requirement: Dynamic Widget Rendering from Actions

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

## MODIFIED Requirements

### Requirement: Module Execution Protocol

The hub SHALL execute module actions by passing the action name via `LAMBDA_ENV_ACTION` environment variable. For interactive modules, the hub SHALL execute actions individually as the user triggers them, rather than executing a single `run` action.

The hub SHALL execute the selected module's `module` executable with environment variables: `LAMBDA_ENV_ACTION=<action-name>`, `LAMBDA_ENV_SETTINGS` (path to settings.json), `LAMBDA_ENV_HUB_VERSION`, `LAMBDA_ENV_LOCALE`.
(Previously: Module execution always used `LAMBDA_ENV_ACTION=run` for a single execution)

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

### Requirement: Menu Rendering by Category

The hub SHALL render a navigable menu grouping modules by category: System, Apps, Ops, Setup. Modules within each category are sorted alphabetically by name. The menu SHALL support 3-level navigation with category selection as the first level.
(Previously: Menu was a flat 2-level navigation: categories → module execution)

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

### Requirement: Keyboard Navigation

The hub SHALL support navigation via arrow keys (up/down), Enter (select/confirm), Esc (back/quit current view), and q (quit application). Key bindings SHALL be context-sensitive based on the current view state.
(Previously: Key bindings were fixed and not context-sensitive)

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

### Requirement: JSON Response Parsing

The hub SHALL parse JSON from module stdout and handle three response types: `ok` (exit 0), `error` (exit 1), `warning` (exit 2). For interactive modules, the response MAY include `settings_delta` that the hub merges immediately.
(Previously: Response parsing was only for single-run module execution)

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
