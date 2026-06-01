# Spec: qtile-module

## Intent

TUI module for Qtile configuration: set default terminal and browser, regenerate config.py from Go template, and reload Qtile via `qtile cmd-obj`.

## Requirements

### Requirement 1: Qtile Settings Schema

The system SHALL read and write the `qtile` section of `settings.json` with the following fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `default_terminal` | string | `"kitty"` | Terminal emulator for keybindings |
| `default_browser` | string | `"firefox"` | Default web browser |
| `default_file_manager` | string | `"thunar"` | Default file manager |
| `bar_position` | string | `"top"` | Bar position: top/bottom/none |
| `bar_size` | int | `24` | Bar height in pixels |
| `layouts` | []string | `["columns","monadtall"]` | Active layout names |
| `groups` | []GroupConfig | `[{name:"1"},...{name:"9"}]` | Workspace definitions |

#### Scenario: Module reads qtile settings

- GIVEN `settings.json` exists with `qtile.default_terminal: "foot"`
- WHEN the qtile module loads settings
- THEN it returns `DefaultTerminal: "foot"` in the typed struct

#### Scenario: Module writes qtile settings delta

- GIVEN `qtile.default_browser` is `"firefox"`
- WHEN the user changes browser to `"brave"` via TUI
- THEN the module emits `settings_delta: {"qtile":{"default_browser":"brave"}}`
- AND other qtile fields are preserved

#### Scenario: Missing qtile section uses defaults

- GIVEN `settings.json` exists but has no `qtile` section
- WHEN the module loads settings
- THEN it returns default values for all qtile fields

### Requirement 2: Set Default Applications

The system SHALL support actions to set `default_terminal`, `default_browser`, and `default_file_manager`. Each action writes to settings.json and triggers config regeneration.

#### Scenario: Change default terminal

- GIVEN `qtile.default_terminal` is `"kitty"`
- WHEN the user selects "Set Terminal" → "foot" in TUI
- THEN `default_terminal` is updated to `"foot"` in settings.json
- AND `config.py` is regenerated with `terminal = "foot"`
- AND the module returns `{"status":"ok","action":"set_terminal","data":{"terminal":"foot"}}`

#### Scenario: Change default browser

- GIVEN `qtile.default_browser` is `"firefox"`
- WHEN the user selects "Set Browser" → "chromium" in TUI
- THEN `default_browser` is updated to `"chromium"` in settings.json
- AND `config.py` is regenerated with the new browser keybinding

#### Scenario: Invalid terminal is rejected

- GIVEN the user enters a terminal name not in the allowed list
- WHEN the module validates the input
- THEN validation fails with error `"unknown terminal: <name>"`
- AND settings.json is NOT modified

### Requirement 3: config.py Regeneration

The system SHALL regenerate `~/.config/qtile/config.py` from a Go template using current `settings.json` qtile values. The generated config SHALL be valid Python before reload.

#### Scenario: Regenerate with custom terminal

- GIVEN `default_terminal: "foot"` and `default_browser: "brave"`
- WHEN `RegenerateConfigPy()` is called
- THEN `config.py` contains `terminal = "foot"` and browser keybinding uses `"brave"`

#### Scenario: Python syntax validation before reload

- GIVEN `config.py` is generated from template
- WHEN the module validates the output
- THEN it runs `python3 -m py_compile config.py`
- AND proceeds to reload only if compilation succeeds

#### Scenario: Backup before regeneration

- GIVEN `~/.config/qtile/config.py` exists
- WHEN `RegenerateConfigPy()` is called
- THEN the existing file is copied to `config.py.bak`
- AND the new content is written to `config.py`

### Requirement 4: Qtile Reload

The system SHALL reload Qtile after successful config regeneration using `qtile cmd-obj -o cmd -f reload_config`. If reload fails, the system SHALL restore the backup config.

#### Scenario: Successful reload

- GIVEN `config.py` is valid Python
- WHEN the module calls reload
- THEN `qtile cmd-obj -o cmd -f reload_config` executes with exit code 0
- AND the module returns `{"status":"ok","action":"reload"}`

#### Scenario: Reload failure restores backup

- GIVEN `config.py.bak` exists from previous backup
- AND `qtile cmd-obj` returns non-zero exit code
- WHEN the module detects reload failure
- THEN it restores `config.py.bak` to `config.py`
- AND returns `{"status":"error","code":"RELOAD_FAILED","message":"Qtile reload failed, config restored from backup"}`

#### Scenario: Qtile not running

- GIVEN Qtile is not the active window manager
- WHEN the module attempts reload
- THEN the module returns a warning: `{"status":"warning","code":"QTILE_NOT_RUNNING","message":"Qtile is not running, config saved but not reloaded"}`
- AND the config.py is still written (reload skipped)

## Technical Details

- Go package: `src/lambda-env/internal/modules/qtile/`
- Template: `src/lambda-env/pkg/templates/qtile/config.py.tmpl`
- Config path: `~/.config/qtile/config.py`
- Backup path: `~/.config/qtile/config.py.bak`
- Reload command: `qtile cmd-obj -o cmd -f reload_config`
- Module manifest category: `apps`

## Dependencies

- `core/01-hub-plugin-system` — module discovery and execution
- `core/02-settings-schema` — settings.json read/write
- Qtile installed on target system
- Python 3 for config validation

## Verification Steps

```bash
# 1. Module compiles
cd src/lambda-env && go build ./internal/modules/qtile/...

# 2. Unit tests pass
cd src/lambda-env && go test ./internal/modules/qtile/... -v -cover

# 3. Change terminal and verify config.py
# Set qtile.default_terminal: "foot" in settings.json
# Run module with LAMBDA_ENV_ACTION=set_terminal
# Verify config.py contains terminal = "foot"

# 4. Generated config is valid Python
# Run: python3 -m py_compile ~/.config/qtile/config.py
# Expect: exit code 0

# 5. Qtile reload command executes
# Run: qtile cmd-obj -o cmd -f reload_config
# Expect: exit code 0 (if Qtile is running)
```
