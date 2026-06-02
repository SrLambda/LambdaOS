# tui-interactive-views Specification

## Purpose

Define the interactive TUI component layer for `lambda-env`: 3-level navigation, toggle widget, text input, confirm dialog, help overlay, and persistent status bar. These components replace the flat 2-level navigation with a component-based architecture.

## Requirements

### Requirement: Three-Level Navigation

The system SHALL support navigation across at least 3 levels: categories → modules → module detail view. Each level SHALL maintain independent scroll position and selection state.

#### Scenario: Navigate from categories to modules

- GIVEN the hub displays the category list
- WHEN the user presses Enter on a category with modules
- THEN the view transitions to the module list for that category
- AND the category name is shown in a breadcrumb or header

#### Scenario: Navigate from modules to module detail

- GIVEN the module list view is displayed
- WHEN the user presses Enter on a module
- THEN the module detail view is rendered based on the module's manifest actions
- AND the module name is shown in the header

#### Scenario: Navigate back through levels

- GIVEN the user is on the module detail view
- WHEN the user presses Esc
- THEN the view returns to the module list
- AND the previous selection is preserved

#### Scenario: Back from module list returns to categories

- GIVEN the user is on the module list view
- WHEN the user presses Esc
- THEN the view returns to the category list
- AND the previous category selection is preserved

### Requirement: Toggle Widget

The system SHALL render a toggle widget for boolean settings that displays current state (on/off) and flips state on user interaction.

#### Scenario: Toggle flips boolean value

- GIVEN a toggle widget displays "Off" for a boolean setting
- WHEN the user presses Space or Enter on the toggle
- THEN the toggle displays "On"
- AND the setting value is flipped internally

#### Scenario: Toggle preserves state across navigation

- GIVEN a toggle has been set to "On" in module detail view
- WHEN the user navigates away and returns to the same module
- THEN the toggle displays "On" (state preserved)

#### Scenario: Toggle shows disabled state

- GIVEN a toggle is associated with a setting that cannot be changed (e.g., missing dependency)
- WHEN the toggle is rendered
- THEN it is displayed in a disabled/grayed-out visual state
- AND user input does not change the value

### Requirement: Text Input Widget

The system SHALL render a text input widget for string settings that allows typing, editing, and submitting values.

#### Scenario: User types and submits text

- GIVEN a text input widget is focused
- WHEN the user types characters and presses Enter
- THEN the input is submitted and stored
- AND the widget exits edit mode

#### Scenario: Input respects max length constraint

- GIVEN a text input has a maximum length of 32 characters
- WHEN the user attempts to type the 33rd character
- THEN the input rejects the additional character
- AND the displayed text remains at 32 characters

#### Scenario: Cancel input with Esc

- GIVEN a text input widget is focused and partially edited
- WHEN the user presses Esc
- THEN the input is cancelled
- AND the original value is restored

### Requirement: Confirm Dialog

The system SHALL display a confirmation dialog before executing destructive or irreversible module actions. The dialog SHALL show the action description and require explicit confirmation.

#### Scenario: Confirm dialog appears before destructive action

- GIVEN the user triggers an action marked as destructive in the manifest
- WHEN the action is about to execute
- THEN a confirm dialog is displayed with the action description
- AND the user must choose "Confirm" or "Cancel"

#### Scenario: User confirms destructive action

- GIVEN a confirm dialog is displayed
- WHEN the user selects "Confirm" and presses Enter
- THEN the action is executed
- AND the dialog is dismissed

#### Scenario: User cancels destructive action

- GIVEN a confirm dialog is displayed
- WHEN the user selects "Cancel" or presses Esc
- THEN the action is NOT executed
- AND the view returns to the previous state

### Requirement: Help Overlay

The system SHALL display a context-sensitive help overlay showing available key bindings for the current view. The overlay SHALL be toggled with a dedicated key (e.g., `?`).

#### Scenario: Help overlay shows current view key bindings

- GIVEN the user is on the module detail view
- WHEN the user presses `?`
- THEN a help overlay is displayed
- AND it shows key bindings relevant to the module detail view

#### Scenario: Help overlay dismisses on key press

- GIVEN the help overlay is visible
- WHEN the user presses any key
- THEN the overlay is dismissed
- AND the underlying view is restored

#### Scenario: Help overlay does not block underlying view

- GIVEN the help overlay is visible
- WHEN the overlay is displayed
- THEN the underlying view content is still visible (dimmed or partially obscured)
- AND the overlay is rendered on top

### Requirement: Persistent Status Bar

The system SHALL render a persistent status bar at the bottom of every view showing the current module name, settings state indicator, and hub version.

#### Scenario: Status bar shows current context

- GIVEN the user is viewing a module detail
- WHEN the status bar is rendered
- THEN it shows the module name and current view level
- AND it shows the hub version string

#### Scenario: Status bar updates on settings change

- GIVEN a toggle or input changes a setting value
- WHEN the change is applied
- THEN the status bar updates to reflect the new state (e.g., shows "modified" indicator)

#### Scenario: Status bar is visible on all views

- GIVEN the user navigates through all view levels
- WHEN any view is displayed
- THEN the status bar is always visible at the bottom
- AND it occupies a fixed number of lines (1-2)

## Non-Functional Requirements

- **Performance**: View transitions SHALL complete within 100ms on a typical terminal
- **Compatibility**: All components SHALL work on terminals with at least 80x24 character dimensions
- **Accessibility**: Key bindings SHALL be discoverable via the help overlay without memorization
- **State isolation**: Each view level SHALL maintain independent state; navigating between levels SHALL NOT corrupt state
