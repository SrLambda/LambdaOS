# Spec: core/02-settings-schema

## Intent

Define a unified settings schema at `~/.config/lambdaos/settings.json` with a Go-based reader/writer that supports atomic writes, schema validation, default values, delta merging, and version migration. This replaces the separate `tui_settings.json` and `os_theme.json` files from Wave 0.

## Scope

### In Scope
- JSON schema definition with `version` field at root level
- Schema sections: appearance, display, audio, network, bluetooth, keyboard, neovim, qtile, services
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

The system SHALL define the following top-level sections in the settings schema: `appearance`, `display`, `audio`, `network`, `bluetooth`, `keyboard`, `neovim`, `qtile`, `services`.

#### Scenario: All sections are present in default settings

- GIVEN no settings.json exists
- WHEN the hub loads settings (returns defaults)
- THEN the settings object contains all nine sections as empty or default objects

#### Scenario: Section structure matches defined schema

- GIVEN settings.json is loaded
- WHEN each section is validated
- THEN `appearance` contains: theme, font_size, opacity, wallpaper
- AND `display` contains: active_profile, profiles (array of output configs)
- AND `audio` contains: default_sink, volume, muted
- AND `network` contains: wifi_enabled, known_networks (array)
- AND `bluetooth` contains: enabled, paired_devices (array)
- AND `keyboard` contains: layout, variant, options
- AND `neovim` contains: theme, font, lines, columns
- AND `qtile` contains: bar_position, bar_size, layouts (array)
- AND `services` contains: enabled (array of service names)

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

## Technical Details

### Go Struct Definitions (simplified)
```go
type Settings struct {
    Version    string            `json:"version"`
    Appearance AppearanceSettings `json:"appearance"`
    Display    DisplaySettings    `json:"display"`
    Audio      AudioSettings      `json:"audio"`
    Network    NetworkSettings    `json:"network"`
    Bluetooth  BluetoothSettings  `json:"bluetooth"`
    Keyboard   KeyboardSettings   `json:"keyboard"`
    Neovim     NeovimSettings     `json:"neovim"`
    Qtile      QtileSettings      `json:"qtile"`
    Services   ServicesSettings   `json:"services"`
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
```

### Default Settings (v1.0.0)
```json
{
  "version": "1.0.0",
  "appearance": { "theme": "dark", "font_size": 14, "opacity": 100, "wallpaper": "" },
  "display": { "active_profile": "default", "profiles": [] },
  "audio": { "default_sink": "", "volume": 75, "muted": false },
  "network": { "wifi_enabled": true, "known_networks": [] },
  "bluetooth": { "enabled": true, "paired_devices": [] },
  "keyboard": { "layout": "us", "variant": "", "options": "" },
  "neovim": { "theme": "tokyonight", "font": "JetBrainsMono", "lines": 40, "columns": 120 },
  "qtile": { "bar_position": "top", "bar_size": 24, "layouts": [] },
  "services": { "enabled": [] }
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
