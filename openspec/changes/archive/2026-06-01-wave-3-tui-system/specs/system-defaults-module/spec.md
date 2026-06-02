# system-defaults-module Specification

## Purpose

Define the default application assignment module for `lambda-env`: xdg-mime-based default browser, terminal, editor, and file manager associations.

## Requirements

### Requirement: Default App Discovery

The system SHALL query current default application assignments for common MIME types and display them in the TUI.

#### Scenario: Current defaults are displayed

- GIVEN xdg-mime is available
- WHEN the module loads
- THEN it queries defaults for key MIME types (text/html, x-scheme-handler/http, inode/directory, text/plain)
- AND the current defaults are displayed in the TUI

#### Scenario: No default set for MIME type

- GIVEN no default application is set for a MIME type
- WHEN the module queries that MIME type
- THEN it displays "Not set" for that entry
- AND the user can assign a default

#### Scenario: xdg-mime not available

- GIVEN xdg-utils is not installed
- WHEN the module attempts to load
- THEN the module returns an error status
- AND the TUI indicates xdg-utils is required

### Requirement: Default Browser Assignment

The system SHALL allow the user to set the default web browser via xdg-mime for HTTP/HTTPS MIME types.

#### Scenario: Default browser is set

- GIVEN the user selects "firefox" from the available browser list
- WHEN the module applies the change
- THEN `xdg-mime default firefox.desktop x-scheme-handler/http` is executed
- AND `xdg-mime default firefox.desktop x-scheme-handler/https` is executed
- AND the settings_delta includes `defaults.browser: "firefox"`

#### Scenario: Selected browser desktop file not found

- GIVEN the user selects a browser whose .desktop file does not exist
- WHEN the module attempts to set it as default
- THEN the module returns an error status
- AND no xdg-mime command is executed

#### Scenario: Multiple browsers are available

- GIVEN multiple browser .desktop files exist on the system
- WHEN the module scans for available browsers
- THEN all browsers are listed as options
- AND the current default is indicated

### Requirement: Default Terminal Assignment

The system SHALL allow the user to set the default terminal emulator.

#### Scenario: Default terminal is set

- GIVEN the user selects "kitty" from available terminals
- WHEN the module applies the change
- THEN `xdg-mime default kitty.desktop x-terminal-emulator` is executed (or equivalent)
- AND the settings_delta includes `defaults.terminal: "kitty"`

#### Scenario: Terminal desktop file validation

- GIVEN the selected terminal's .desktop file exists but lacks the Terminal=true field
- WHEN the module validates the desktop file
- THEN the module returns a warning
- AND the user can confirm or select a different terminal

### Requirement: Default Editor Assignment

The system SHALL allow the user to set the default text editor via environment variable and xdg-mime.

#### Scenario: Default editor is set

- GIVEN the user selects "nvim" from available editors
- WHEN the module applies the change
- THEN the appropriate xdg-mime assignment is made for text/plain
- AND the settings_delta includes `defaults.editor: "nvim"`

#### Scenario: Editor not installed

- GIVEN the user selects an editor that is not installed
- WHEN the module validates the selection
- THEN the module returns an error status
- AND the user is prompted to install the editor first

### Requirement: Default File Manager Assignment

The system SHALL allow the user to set the default file manager for directory and inode MIME types.

#### Scenario: Default file manager is set

- GIVEN the user selects "thunar" from available file managers
- WHEN the module applies the change
- THEN `xdg-mime default thunar.desktop inode/directory` is executed
- AND the settings_delta includes `defaults.file_manager: "thunar"`

#### Scenario: File manager desktop file not found

- GIVEN the selected file manager's .desktop file does not exist
- WHEN the module attempts to set it
- THEN the module returns an error status
- AND no changes are made

### Requirement: Batch Default Assignment

The system SHALL support applying multiple default app changes in a single action with confirmation.

#### Scenario: Multiple defaults are applied together

- GIVEN the user has changed browser, terminal, and editor selections
- WHEN the user triggers "Apply All"
- THEN a confirm dialog is shown listing all changes
- AND upon confirmation, all xdg-mime commands are executed
- AND a single settings_delta with all defaults is emitted

#### Scenario: Partial failure in batch assignment

- GIVEN a batch of 3 default assignments is being applied
- WHEN the second assignment fails
- THEN the module returns a warning status
- AND the settings_delta includes only the successfully applied changes
- AND the failed assignment is logged with the error

## Non-Functional Requirements

- **Dependencies**: The module SHALL declare `xdg-utils` as a dependency
- **Safety**: The module SHALL NOT modify system-wide defaults; only user-level (~/.config/mimeapps.list)
- **Validation**: All .desktop file references SHALL be validated before assignment
- **Testing**: All xdg-mime interactions SHALL use interface-based mocking for unit tests
