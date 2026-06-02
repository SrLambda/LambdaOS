# Spec: core/02-settings-schema

## Intent

Define a unified settings schema at `~/.config/lambdaos/settings.json` with a Go-based reader/writer that supports atomic writes, schema validation, default values, delta merging, and version migration. This replaces the separate `tui_settings.json` and `os_theme.json` files from Wave 0.

## Scope

### In Scope
- JSON schema definition with `version` field at root level
- Schema sections: appearance, display, audio, network, bluetooth, keyboard, neovim, qtile, services, power, defaults, autostart, updates, security, fonts, notifications
- Atomic write pattern (temp file + rename in same directory)
- Default values when file doesn't exist or fields are missing
- Settings reader returning typed Go structs
- Settings writer accepting deltas and merging into existing data
- Schema migration support (version bumping with field additions)
- Validation before write (required fields, valid types, enum values)

### Out of Scope
- Settings UI/TUI views — handled by hub plugin system (core/01)
- Module-specific settings logic — modules emit deltas, hub applies them
- Migration of existing `tui_settings.json` — Wave 2

## Requirements

### Requirement 1: Settings File Location and Format

The system SHALL store all settings in a single JSON file at `~/.config/lambdaos/settings.json` with a `version` field at the root level.

#### Scenario: Settings file is created at correct path

- GIVEN the directory `~/.config/lambdaos/` does not exist
- WHEN the hub saves settings for the first time
- THEN the directory is created
- AND `settings.json` is written at the correct path

#### Scenario: Settings file contains version field

- GIVEN settings.json exists
- WHEN the file is read
- THEN the root object contains a `"version"` field with a semver string (e.g., `"1.0.0"`)

#### Scenario: Settings file is valid JSON

- GIVEN settings.json exists
- WHEN the file is parsed
- THEN it is valid JSON with no syntax errors

### Requirement 2: Schema Sections

The system SHALL define the following top-level sections in the settings schema: `appearance`, `display`, `audio`, `network`, `bluetooth`, `keyboard`, `neovim`, `qtile`, `services`, `power`, `defaults`, `autostart`, `updates`, `security`, `fonts`, `notifications`.

#### Scenario: All sections are present in default settings

- GIVEN no settings.json exists
- WHEN the hub loads settings (returns defaults)
- THEN the settings object contains all sixteen sections as empty or default objects

#### Scenario: Section structure matches defined schema

- GIVEN settings.json is loaded
- WHEN each section is validated
- THEN `appearance` contains: theme, font_size, opacity, wallpaper
- AND `display` contains: active_profile, profiles (array of output configs)
- AND `audio` contains: default_sink, volume, muted
- AND `network` contains: wifi_enabled, known_networks (array)
- AND `bluetooth` contains: enabled, paired_devices (array)
- AND `keyboard` contains: layout, variant, options
- AND `neovim` contains: theme, font, lines, columns, enable_lsp, enable_copilot, enable_neotree, lsp_servers, use_global_theme
- AND `qtile` contains: bar_position, bar_size, layouts (array), default_terminal, default_browser, default_file_manager, groups, use_global_theme
- AND `services` contains: enabled (array of service names)
- AND `power` contains: screen_timeout, sleep_timeout, lid_close_action
- AND `defaults` contains: browser, terminal, editor, file_manager
- AND `autostart` contains: enabled (array of app names)
- AND `updates` contains: auto_update, check_interval, exclude_packages
- AND `security` contains: firewall_enabled, sudo_timeout, screen_lock_timeout
- AND `fonts` contains: monospace, sans_serif, serif, font_size
- AND `notifications` contains: enabled, do_not_disturb, timeout_seconds

### Requirement 3: Atomic Writes

The system SHALL use atomic writes for all settings modifications: write to a temporary file in the same directory, then rename to `settings.json`.

#### Scenario: Atomic write completes successfully

- GIVEN settings.json exists with valid content
- WHEN the hub saves new settings
- THEN a temp file is created in `~/.config/lambdaos/`
- AND the temp file is renamed to `settings.json`
- AND the final file contains the new content

#### Scenario: Atomic write preserves content on failure

- GIVEN settings.json exists with valid content
- WHEN a write fails mid-operation (e.g., disk full during temp write)
- THEN the original settings.json is unchanged
- AND no corrupted temp file remains

#### Scenario: Temp file uses same directory

- GIVEN the hub performs an atomic write
- WHEN the temp file is created
- THEN it is in `~/.config/lambdaos/` (same directory as target)
- AND the rename is a same-filesystem operation

### Requirement 4: Default Values

The system SHALL return a complete settings object with default values when the file does not exist or when individual fields are missing.

#### Scenario: Missing file returns full defaults

- GIVEN `~/.config/lambdaos/settings.json` does not exist
- WHEN the hub calls Load()
- THEN a complete settings object with all default values is returned
- AND no error is raised

#### Scenario: Missing fields are filled with defaults

- GIVEN settings.json exists with only `{"version": "1.0.0"}`
- WHEN the hub calls Load()
- THEN all missing sections are populated with their default values
- AND the version field is preserved as `"1.0.0"`

#### Scenario: Partial section is merged with defaults

- GIVEN settings.json has `"display": {"active_profile": "home"}` but no `profiles` field
- WHEN the hub calls Load()
- THEN `display.active_profile` is `"home"`
- AND `display.profiles` is set to its default value (empty array)

### Requirement 5: Settings Reader Returns Typed Go Structs

The system SHALL provide a `Load()` function that reads settings.json and returns typed Go structs, not raw `map[string]interface{}`.

#### Scenario: Load returns typed struct

- GIVEN settings.json exists with valid content
- WHEN `settings.Load()` is called
- THEN it returns a `Settings` struct (not `map[string]interface{}`)
- AND each section is a typed struct (e.g., `DisplaySettings`, `AudioSettings`)

#### Scenario: Load returns error on invalid JSON

- GIVEN settings.json contains invalid JSON
- WHEN `settings.Load()` is called
- THEN an error is returned
- AND the caller can inspect the parse error

#### Scenario: Load returns error on invalid version

- GIVEN settings.json has `"version": "not-a-version"`
- WHEN `settings.Load()` is called
- THEN an error is returned indicating invalid version format

### Requirement 6: Settings Writer Accepts Deltas

The system SHALL provide a `SaveDelta(delta map[string]interface{})` function that merges a partial settings update into the existing settings and writes atomically.

#### Scenario: Delta updates only specified fields

- GIVEN settings.json has `{"version":"1.0.0","display":{"active_profile":"default","profiles":[]}}`
- WHEN `SaveDelta` is called with `{"display":{"active_profile":"home"}}`
- THEN settings.json contains `{"display":{"active_profile":"home","profiles":[]}}`
- AND `profiles` is preserved (not overwritten)

#### Scenario: Delta adds new fields

- GIVEN settings.json has `{"version":"1.0.0"}`
- WHEN `SaveDelta` is called with `{"audio":{"volume":75}}`
- THEN settings.json contains the new `audio` section
- AND `version` is preserved

#### Scenario: Delta with empty object is a no-op

- GIVEN settings.json exists
- WHEN `SaveDelta` is called with `{}`
- THEN settings.json is not modified

### Requirement 7: Schema Validation Before Write

The system SHALL validate settings before writing: required fields must be present, types must match, and enum values must be valid.

#### Scenario: Valid settings pass validation

- GIVEN a settings object with all required fields and valid types
- WHEN `Validate()` is called
- THEN validation passes with no errors

#### Scenario: Invalid enum value fails validation

- GIVEN a settings object has `display.active_profile` set to `"invalid-profile"`
- AND valid profiles are only those defined in `display.profiles`
- WHEN `Validate()` is called
- THEN validation fails with an error about invalid profile reference

#### Scenario: Wrong type fails validation

- GIVEN a settings object has `audio.volume` set to `"loud"` (string instead of int)
- WHEN `Validate()` is called
- THEN validation fails with a type mismatch error

#### Scenario: Missing required field fails validation

- GIVEN a settings object is missing the `version` field
- WHEN `Validate()` is called
- THEN validation fails indicating version is required

### Requirement 8: Schema Migration Support

The system SHALL detect version mismatches between the loaded settings file and the current schema version, and migrate by adding missing fields with defaults.

#### Scenario: Migration adds missing fields

- GIVEN settings.json has `"version": "1.0.0"` with only `appearance` and `display` sections
- AND the current schema version is `"1.1.0"` which adds `bluetooth`
- WHEN the hub loads settings
- THEN the `bluetooth` section is added with default values
- AND existing `appearance` and `display` values are preserved

#### Scenario: Migration does not overwrite user values

- GIVEN settings.json has `"version": "1.0.0"` with `audio.volume: 80`
- AND migration to `"1.1.0"` would set `audio.volume` default to `50`
- WHEN migration runs
- THEN `audio.volume` remains `80` (user value preserved)
- AND only truly missing fields receive defaults

#### Scenario: Downgrade is rejected

- GIVEN settings.json has `"version": "2.0.0"`
- AND the current schema version is `"1.0.0"`
- WHEN the hub loads settings
- THEN an error is returned (cannot downgrade)

#### Scenario: Same version requires no migration

- GIVEN settings.json has `"version": "1.0.0"`
- AND the current schema version is `"1.0.0"`
- WHEN the hub loads settings
- THEN no migration is performed
- AND settings are loaded as-is

#### Scenario: v1.0.0 to v1.1.0 adds use_global_theme

- GIVEN settings.json has `"version": "1.0.0"` with neovim and qtile sections
- WHEN migration to v1.1.0 runs
- THEN `neovim.use_global_theme` is added with default value `true`
- AND `qtile.use_global_theme` is added with default value `true`
- AND existing neovim/qtile fields are preserved

### Requirement 9: Schema Version 1.1.0

The system SHALL bump the schema version from `1.0.0` to `1.1.0` and support migration from v1.0.0 to v1.1.0 with all new sections added as optional fields with defaults.

#### Scenario: Migration from v1.0.0 to v1.1.0

- GIVEN settings.json has `"version": "1.0.0"` with existing sections populated
- WHEN the hub loads settings with schema v1.1.0
- THEN all 7 new sections are added with default values
- AND existing section values are preserved unchanged

#### Scenario: Migration is atomic

- GIVEN settings.json has v1.0.0 content
- WHEN migration to v1.1.0 runs
- THEN the migration writes to a temp file first
- AND renames to settings.json only after successful validation
- AND the original file is preserved if migration fails

#### Scenario: Already at v1.1.0 requires no migration

- GIVEN settings.json has `"version": "1.1.0"`
- WHEN the hub loads settings
- THEN no migration is performed
- AND settings are loaded as-is

### Requirement 10: Power Section

The system SHALL define a `power` section in the settings schema with fields for power management configuration.

#### Scenario: Power section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `power` section is present with default values
- AND includes fields for screen timeout, sleep timeout, and lid close action

#### Scenario: Power section delta is merged

- GIVEN settings.json exists with v1.1.0 schema
- WHEN a module emits a settings_delta with `power.screen_timeout: 300`
- THEN the delta is merged into the power section
- AND other power fields are preserved

### Requirement 11: Defaults Section

The system SHALL define a `defaults` section for default application assignments (browser, terminal, editor, file manager).

#### Scenario: Defaults section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `defaults` section is present with empty or default values
- AND includes fields for browser, terminal, editor, and file_manager

#### Scenario: Defaults section stores app assignments

- GIVEN the defaults module sets browser to "firefox"
- WHEN the settings_delta is merged
- THEN `defaults.browser` is set to "firefox"
- AND the value persists across hub restarts

### Requirement 12: Autostart Section

The system SHALL define an `autostart` section for managing applications that start on login.

#### Scenario: Autostart section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `autostart` section is present with an empty enabled list
- AND includes fields for `enabled` (array of service/app names)

#### Scenario: Autostart entries are added

- GIVEN the autostart section exists
- WHEN a module adds an entry to the enabled list
- THEN the entry is appended to the `enabled` array
- AND duplicate entries are not added

### Requirement 13: Updates Section

The system SHALL define an `updates` section for system update configuration.

#### Scenario: Updates section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `updates` section is present with default values
- AND includes fields for auto_update (bool), check_interval, and exclude_packages (array)

#### Scenario: Updates settings are persisted

- GIVEN the user configures update preferences
- WHEN the settings_delta is applied
- THEN the updates section reflects the new configuration
- AND the values are validated against allowed types

### Requirement 14: Security Section

The system SHALL define a `security` section for security-related configuration.

#### Scenario: Security section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `security` section is present with default values
- AND includes fields for firewall_enabled, sudo_timeout, and screen_lock_timeout

#### Scenario: Security settings are validated

- GIVEN a module attempts to set an invalid security value
- WHEN the settings_delta is validated
- THEN the invalid value is rejected
- AND an error is returned to the caller

### Requirement 15: Fonts Section

The system SHALL define a `fonts` section for font configuration beyond the appearance section.

#### Scenario: Fonts section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `fonts` section is present with default values
- AND includes fields for monospace, sans_serif, serif, and font_size

#### Scenario: Font settings are merged with appearance

- GIVEN the appearance module sets a font
- WHEN the settings_delta is applied
- THEN the fonts section is updated accordingly
- AND the appearance.font_size is synchronized if applicable

### Requirement 16: Notifications Section

The system SHALL define a `notifications` section for notification daemon configuration.

#### Scenario: Notifications section is present in defaults

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN the `notifications` section is present with default values
- AND includes fields for enabled (bool), do_not_disturb (bool), and timeout_seconds

#### Scenario: Notification settings are persisted

- GIVEN the user toggles do_not_disturb mode
- WHEN the settings_delta is applied
- THEN `notifications.do_not_disturb` is set to true
- AND the value persists across restarts

### Requirement 17: use_global_theme in Neovim and Qtile Settings

The system SHALL add a `use_global_theme` boolean field to both `NeovimSettings` and `QtileSettings` structs. When true, the module SHALL derive its theme from `appearance.theme`; when false, it SHALL use its own `theme` field.

#### Scenario: Neovim uses global theme

- GIVEN `neovim.use_global_theme` is true
- AND `appearance.theme` is "dark"
- WHEN the Neovim module loads settings
- THEN it maps "dark" to the equivalent Neovim theme (e.g., "tokyonight-night")
- AND the module's own `theme` field is ignored

#### Scenario: Neovim uses independent theme

- GIVEN `neovim.use_global_theme` is false
- AND `neovim.theme` is "catppuccin"
- WHEN the Neovim module loads settings
- THEN it uses "catppuccin" regardless of appearance.theme
- AND appearance.theme changes do not affect Neovim

#### Scenario: Qtile uses global theme

- GIVEN `qtile.use_global_theme` is true
- AND `appearance.theme` is "light"
- WHEN the Qtile module loads settings
- THEN it maps "light" to the equivalent Qtile theme
- AND the module's own `theme` field is ignored

#### Scenario: Qtile uses independent theme

- GIVEN `qtile.use_global_theme` is false
- WHEN the Qtile module loads settings
- THEN it uses its own `theme` field value
- AND appearance.theme changes do not affect Qtile

#### Scenario: use_global_theme defaults to true

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN `neovim.use_global_theme` defaults to true
- AND `qtile.use_global_theme` defaults to true

## Technical Details

### Go Struct Definitions (simplified)
```go
type Settings struct {
    Version       string              `json:"version"`
    Appearance    AppearanceSettings  `json:"appearance"`
    Display       DisplaySettings     `json:"display"`
    Audio         AudioSettings       `json:"audio"`
    Network       NetworkSettings     `json:"network"`
    Bluetooth     BluetoothSettings   `json:"bluetooth"`
    Keyboard      KeyboardSettings    `json:"keyboard"`
    Neovim        NeovimSettings      `json:"neovim"`
    Qtile         QtileSettings       `json:"qtile"`
    Services      ServicesSettings    `json:"services"`
    Power         PowerSettings       `json:"power"`
    Defaults      DefaultsSettings    `json:"defaults"`
    Autostart     AutostartSettings   `json:"autostart"`
    Updates       UpdatesSettings     `json:"updates"`
    Security      SecuritySettings    `json:"security"`
    Fonts         FontsSettings       `json:"fonts"`
    Notifications NotificationsSettings `json:"notifications"`
}

type DisplaySettings struct {
    ActiveProfile string          `json:"active_profile"`
    Profiles      []OutputProfile `json:"profiles"`
}

type AudioSettings struct {
    DefaultSink string `json:"default_sink"`
    Volume      int    `json:"volume"`      // 0-100
    Muted       bool   `json:"muted"`
}

type PowerSettings struct {
    ScreenTimeout  int    `json:"screen_timeout"`   // seconds
    SleepTimeout   int    `json:"sleep_timeout"`    // seconds
    LidCloseAction string `json:"lid_close_action"` // "suspend"|"hibernate"|"ignore"
}

type DefaultsSettings struct {
    Browser     string `json:"browser"`
    Terminal    string `json:"terminal"`
    Editor      string `json:"editor"`
    FileManager string `json:"file_manager"`
}

type AutostartSettings struct {
    Enabled []string `json:"enabled"`
}

type UpdatesSettings struct {
    AutoUpdate      bool     `json:"auto_update"`
    CheckInterval   int      `json:"check_interval"`    // hours
    ExcludePackages []string `json:"exclude_packages"`
}

type SecuritySettings struct {
    FirewallEnabled   bool `json:"firewall_enabled"`
    SudoTimeout       int  `json:"sudo_timeout"`        // minutes
    ScreenLockTimeout int  `json:"screen_lock_timeout"` // seconds
}

type FontsSettings struct {
    Monospace string `json:"monospace"`
    SansSerif string `json:"sans_serif"`
    Serif     string `json:"serif"`
    FontSize  int    `json:"font_size"`
}

type NotificationsSettings struct {
    Enabled        bool `json:"enabled"`
    DoNotDisturb   bool `json:"do_not_disturb"`
    TimeoutSeconds int  `json:"timeout_seconds"`
}
```

### Default Settings (v1.1.0)
```json
{
  "version": "1.1.0",
  "appearance": { "theme": "dark", "font_size": 14, "opacity": 100, "wallpaper": "" },
  "display": { "active_profile": "default", "profiles": [] },
  "audio": { "default_sink": "", "volume": 75, "muted": false },
  "network": { "wifi_enabled": true, "known_networks": [] },
  "bluetooth": { "enabled": true, "paired_devices": [] },
  "keyboard": { "layout": "us", "variant": "", "options": "" },
  "neovim": { "theme": "tokyonight", "font": "JetBrainsMono", "lines": 40, "columns": 120, "enable_lsp": true, "enable_copilot": true, "enable_neotree": true, "lsp_servers": ["gopls", "pyright"], "use_global_theme": true },
  "qtile": { "bar_position": "top", "bar_size": 24, "layouts": ["columns", "monadtall"], "default_terminal": "kitty", "default_browser": "firefox", "default_file_manager": "thunar", "groups": [{"name": "1"}, {"name": "2"}, {"name": "3"}, {"name": "4"}, {"name": "5"}, {"name": "6"}, {"name": "7"}, {"name": "8"}, {"name": "9"}], "use_global_theme": true },
  "services": { "enabled": [] },
  "power": { "screen_timeout": 300, "sleep_timeout": 1800, "lid_close_action": "suspend" },
  "defaults": { "browser": "", "terminal": "", "editor": "", "file_manager": "" },
  "autostart": { "enabled": [] },
  "updates": { "auto_update": true, "check_interval": 24, "exclude_packages": [] },
  "security": { "firewall_enabled": true, "sudo_timeout": 5, "screen_lock_timeout": 300 },
  "fonts": { "monospace": "JetBrainsMono Nerd Font", "sans_serif": "Inter", "serif": "Noto Serif", "font_size": 11 },
  "notifications": { "enabled": true, "do_not_disturb": false, "timeout_seconds": 5 }
}
```

### Atomic Write Implementation
```go
func (s *Store) Save(data Settings) error {
    dir := filepath.Dir(s.path)
    tmp, err := os.CreateTemp(dir, "settings-*.tmp")
    if err != nil { return err }
    defer os.Remove(tmp.Name()) // cleanup on failure

    if err := json.NewEncoder(tmp).Encode(data); err != nil {
        tmp.Close()
        return err
    }
    if err := tmp.Close(); err != nil { return err }
    return os.Rename(tmp.Name(), s.path)
}
```

### Package Structure
```
src/lambda-env/internal/settings/
├── schema.go      # Go struct definitions, defaults, validation
├── store.go       # Load, Save, SaveDelta, Migrate, Validate
└── store_test.go  # Unit tests for all operations
```

## Dependencies

- Go standard library: `encoding/json`, `os`, `path/filepath`, `io/ioutil`
- Settings file path: `~/.config/lambdaos/settings.json`
- Hub plugin system (core/01) — modules emit deltas that this package applies

## Verification Steps

```bash
# 1. Settings package compiles
cd src/lambda-env && go build ./internal/settings/...

# 2. Unit tests pass
cd src/lambda-env && go test ./internal/settings/... -v -cover

# 3. Load returns defaults when file missing
rm -f ~/.config/lambdaos/settings.json
# Run: go test -run TestLoadDefaults ./internal/settings/
# Expect: full default settings returned, no error

# 4. Atomic write preserves original on failure
# Run: go test -run TestAtomicWriteFailure ./internal/settings/
# Expect: original file unchanged after simulated failure

# 5. Delta merge preserves existing fields
# Run: go test -run TestSaveDeltaMerge ./internal/settings/
# Expect: only delta fields updated, others preserved

# 6. Validation catches invalid types
# Run: go test -run TestValidateInvalidType ./internal/settings/
# Expect: validation error for wrong type

# 7. Migration adds missing fields
# Run: go test -run TestMigrationAddsFields ./internal/settings/
# Expect: missing sections added with defaults, existing values preserved

# 8. Downgrade is rejected
# Run: go test -run TestDowngradeRejected ./internal/settings/
# Expect: error returned when file version > schema version

# 9. Settings file is valid JSON after write
cat ~/.config/lambdaos/settings.json | python3 -m json.tool > /dev/null && echo "Valid JSON"
```
