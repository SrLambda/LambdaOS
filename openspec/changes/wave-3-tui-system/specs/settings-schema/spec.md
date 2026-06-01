# Delta for settings-schema

## ADDED Requirements

### Requirement: Schema Version 1.1.0

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

### Requirement: Power Section

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

### Requirement: Defaults Section

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

### Requirement: Autostart Section

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

### Requirement: Updates Section

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

### Requirement: Security Section

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

### Requirement: Fonts Section

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

### Requirement: Notifications Section

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

### Requirement: use_global_theme in Neovim and Qtile Settings

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

#### Scenario: use_global_theme defaults to false

- GIVEN no settings.json exists
- WHEN the hub loads default settings
- THEN `neovim.use_global_theme` defaults to false
- AND `qtile.use_global_theme` defaults to false

## MODIFIED Requirements

### Requirement: Schema Sections

The system SHALL define the following top-level sections in the settings schema: `appearance`, `display`, `audio`, `network`, `bluetooth`, `keyboard`, `neovim`, `qtile`, `services`, `power`, `defaults`, `autostart`, `updates`, `security`, `fonts`, `notifications`.
(Previously: Schema had 9 sections — appearance, display, audio, network, bluetooth, keyboard, neovim, qtile, services)

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

### Requirement: Schema Migration Support

The system SHALL detect version mismatches between the loaded settings file and the current schema version, and migrate by adding missing fields with defaults. Migration from v1.0.0 to v1.1.0 SHALL add all 7 new sections and the `use_global_theme` field to Neovim and Qtile settings.
(Previously: Migration only handled single-section additions between minor versions)

#### Scenario: Migration adds missing fields

- GIVEN settings.json has `"version": "1.0.0"` with only `appearance` and `display` sections
- AND the current schema version is `"1.1.0"` which adds 7 new sections
- WHEN the hub loads settings
- THEN all 7 new sections are added with default values
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
- THEN `neovim.use_global_theme` is added with default value `false`
- AND `qtile.use_global_theme` is added with default value `false`
- AND existing neovim/qtile fields are preserved
