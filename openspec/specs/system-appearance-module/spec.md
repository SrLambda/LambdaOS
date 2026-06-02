# system-appearance-module Specification

## Purpose

Define the appearance configuration module for `lambda-env`: theme selection, wallpaper setting, icon/cursor theme, and font selection, with support for cross-module theme synchronization.

## Requirements

### Requirement: Theme Selection

The system SHALL provide a list of available GTK themes and allow the user to select one via the TUI.

#### Scenario: Available themes are listed

- GIVEN theme directories exist in `/usr/share/themes/` and `~/.themes/`
- WHEN the appearance module loads
- THEN it scans both directories for valid themes
- AND presents them as a selectable list in the TUI

#### Scenario: Theme is applied

- GIVEN the user selects a theme (e.g., "Adwaita-dark")
- WHEN the module applies the theme
- THEN `gsettings set org.gnome.desktop.interface gtk-theme` is executed
- AND the settings_delta includes `appearance.theme: "Adwaita-dark"`

#### Scenario: No themes found

- GIVEN no theme directories exist or are empty
- WHEN the module scans for themes
- THEN it returns a warning status
- AND the TUI displays a message suggesting theme installation

### Requirement: Wallpaper Setting

The system SHALL set the desktop wallpaper using feh and persist the wallpaper path in settings.

#### Scenario: Wallpaper is set from file path

- GIVEN a valid image file path is provided
- WHEN the module applies the wallpaper
- THEN `feh --bg-scale <path>` is executed
- AND the settings_delta includes `appearance.wallpaper: <path>`

#### Scenario: Wallpaper file does not exist

- GIVEN the specified wallpaper file does not exist
- WHEN the module attempts to apply it
- THEN the module returns an error status
- AND no settings_delta is emitted

#### Scenario: Wallpaper is cleared

- GIVEN a wallpaper is currently set
- WHEN the user selects "None" or clears the wallpaper
- THEN the wallpaper path in settings_delta is set to empty string
- AND feh is not executed (no-op for clearing)

### Requirement: Icon and Cursor Theme Selection

The system SHALL allow selection of icon themes and cursor themes independently.

#### Scenario: Icon theme is applied

- GIVEN the user selects an icon theme (e.g., "Papirus")
- WHEN the module applies the icon theme
- THEN `gsettings set org.gnome.desktop.interface icon-theme` is executed
- AND the settings_delta includes the icon theme name

#### Scenario: Cursor theme is applied

- GIVEN the user selects a cursor theme (e.g., "Bibata-Modern-Ice")
- WHEN the module applies the cursor theme
- THEN `gsettings set org.gnome.desktop.interface cursor-theme` is executed
- AND the settings_delta includes the cursor theme name

#### Scenario: Icon theme directory not found

- GIVEN the selected icon theme does not exist in theme paths
- WHEN the module attempts to apply it
- THEN the module returns an error status
- AND settings are not modified

### Requirement: Font Selection

The system SHALL allow selection of font family and font size for the desktop environment.

#### Scenario: Font family is changed

- GIVEN the user selects a font family (e.g., "Inter")
- WHEN the module applies the font
- THEN `gsettings set org.gnome.desktop.interface font-name` is executed with the new family
- AND the settings_delta includes the font family

#### Scenario: Font size is changed

- GIVEN the user changes font size via text input widget
- WHEN the module applies the size
- THEN the font-name gsetting is updated with the new size
- AND the settings_delta includes the new font size

#### Scenario: Invalid font family is rejected

- GIVEN the user enters a font family that is not installed
- WHEN the module attempts to apply it
- THEN the module returns a warning status
- AND the user is prompted to confirm or cancel

### Requirement: Theme Sync Flag (use_global_theme)

The system SHALL expose a `use_global_theme` boolean that, when enabled, maps `appearance.theme` to module-specific themes for Neovim and Qtile.

#### Scenario: Theme sync is enabled

- GIVEN `use_global_theme` is set to true in Neovim settings
- WHEN the appearance module changes the global theme
- THEN the Neovim module maps the appearance theme to its equivalent Neovim theme
- AND the settings_delta for Neovim includes the mapped theme

#### Scenario: Theme sync is disabled

- GIVEN `use_global_theme` is set to false in Neovim settings
- WHEN the appearance module changes the global theme
- THEN the Neovim module theme is NOT changed
- AND Neovim retains its independently configured theme

#### Scenario: Theme mapping has no equivalent

- GIVEN `use_global_theme` is true and the appearance theme has no Neovim equivalent
- WHEN the theme is changed
- THEN the module falls back to a default Neovim theme
- AND a warning is logged indicating no direct mapping exists

## Non-Functional Requirements

- **Dependencies**: The module SHALL declare `feh`, `gsettings` (via glib2) as dependencies
- **Performance**: Theme application SHALL complete within 3 seconds per change
- **Safety**: The module SHALL NOT modify theme files directly; only gsettings and feh commands
- **Theme mapping**: The appearance → Neovim/Qtile theme mapping SHALL be defined as a lookup table in the module code
