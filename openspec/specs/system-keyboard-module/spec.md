# system-keyboard-module Specification

## Purpose

Define the keyboard configuration module for `lambda-env`: setxkbmap-based keyboard layout and variant selection, exposed through interactive TUI widgets.

## Requirements

### Requirement: Keyboard Layout Discovery

The system SHALL query available keyboard layouts and variants using setxkbmap and present them as selectable options in the TUI.

#### Scenario: Available layouts are listed

- GIVEN setxkbmap is installed and accessible
- WHEN the keyboard module loads
- THEN it queries available layouts via setxkbmap
- AND the layouts are presented as a selectable list in the TUI

#### Scenario: setxkbmap not found

- GIVEN setxkbmap is not installed on the system
- WHEN the keyboard module attempts to load
- THEN the module returns an error status
- AND the TUI displays a message indicating setxkbmap is required

#### Scenario: Variants listed for selected layout

- GIVEN a layout (e.g., "us") is selected
- WHEN the module queries variants for that layout
- THEN available variants are displayed (e.g., "intl", "dvorak", "colemak")
- AND an empty variant option is included for "no variant"

### Requirement: Layout Application

The system SHALL apply the selected keyboard layout using setxkbmap and verify the change was successful.

#### Scenario: Layout is applied successfully

- GIVEN the user selects "us" layout and presses Enter
- WHEN the module executes `setxkbmap us`
- THEN the command exits with code 0
- AND the module returns a settings_delta with `keyboard.layout: "us"`

#### Scenario: Layout with variant is applied

- GIVEN the user selects "us" layout with "dvorak" variant
- WHEN the module executes `setxkbmap -layout us -variant dvorak`
- THEN the command exits with code 0
- AND the settings_delta includes both layout and variant

#### Scenario: Invalid layout is rejected

- GIVEN the user provides a layout name that does not exist
- WHEN the module attempts to apply it
- THEN setxkbmap returns a non-zero exit code
- AND the module returns an error status with a descriptive message

### Requirement: Current Layout Detection

The system SHALL detect the currently active keyboard layout and variant to display the current state in the TUI.

#### Scenario: Current layout is detected

- GIVEN the system has an active keyboard configuration
- WHEN the module queries the current layout via setxkbmap -query
- THEN the current layout is parsed and displayed
- AND the current variant is parsed and displayed (or "none" if empty)

#### Scenario: Query output parsing handles variations

- GIVEN setxkbmap -query output format varies slightly across versions
- WHEN the module parses the output
- THEN it extracts layout and variant using regex with fallback patterns
- AND unexpected output is logged for diagnosis

### Requirement: Keyboard Options Support

The system SHALL support setting keyboard options (e.g., caps:swapescape, ctrl:nocaps) via setxkbmap -option.

#### Scenario: Option is applied

- GIVEN the user selects the "swap caps and escape" option
- WHEN the module executes `setxkbmap -option caps:swapescape`
- THEN the option is applied
- AND the settings_delta includes the option string

#### Scenario: Multiple options are combined

- GIVEN the user selects two keyboard options
- WHEN the module applies them
- THEN both options are passed as comma-separated values to setxkbmap -option
- AND the settings_delta includes the combined option string

#### Scenario: Options are cleared before applying new set

- GIVEN existing keyboard options are active
- WHEN new options are applied
- THEN the module first clears options with `setxkbmap -option ""`
- AND then applies the new option set

## Non-Functional Requirements

- **Dependency**: The module SHALL declare `xorg-setxkbmap` as a dependency in its manifest
- **Performance**: Layout switching SHALL complete within 2 seconds
- **Safety**: The module SHALL NOT modify XKB configuration files directly; all changes through setxkbmap only
- **Reversibility**: The module SHALL backup the current layout before applying changes for rollback
