# system-audio-module Specification

## Purpose

Define the audio configuration module for `lambda-env`: PipeWire/PulseAudio volume control, sink selection, and mute toggle via pactl/wpctl commands.

## Requirements

### Requirement: Audio Backend Detection

The system SHALL detect whether PipeWire or PulseAudio is the active audio backend at runtime and use the appropriate commands.

#### Scenario: PipeWire backend is detected

- GIVEN PipeWire is running as the audio server
- WHEN the module queries the backend via `pactl info`
- THEN the server name contains "PipeWire"
- AND the module uses wpctl commands for operations

#### Scenario: PulseAudio backend is detected

- GIVEN PulseAudio is running as the audio server
- WHEN the module queries the backend via `pactl info`
- THEN the server name contains "PulseAudio"
- AND the module uses pactl commands for operations

#### Scenario: No audio backend is running

- GIVEN neither PipeWire nor PulseAudio is running
- WHEN the module attempts to detect the backend
- THEN the module returns an error status
- AND the TUI displays a message indicating no audio server is available

### Requirement: Volume Control

The system SHALL adjust the master volume level in 5% increments, with a range of 0-100%.

#### Scenario: Volume is increased

- GIVEN the current volume is at 50%
- WHEN the user triggers volume up
- THEN the volume is set to 55% via the appropriate backend command
- AND the settings_delta includes `audio.volume: 55`

#### Scenario: Volume is decreased

- GIVEN the current volume is at 50%
- WHEN the user triggers volume down
- THEN the volume is set to 45% via the appropriate backend command
- AND the settings_delta includes `audio.volume: 45`

#### Scenario: Volume caps at 100%

- GIVEN the current volume is at 98%
- WHEN the user triggers volume up
- THEN the volume is set to 100% (not 103%)
- AND the settings_delta includes `audio.volume: 100`

#### Scenario: Volume floors at 0%

- GIVEN the current volume is at 3%
- WHEN the user triggers volume down
- THEN the volume is set to 0% (not -2%)
- AND the settings_delta includes `audio.volume: 0`

#### Scenario: Volume is set to specific value

- GIVEN the user enters a specific volume value via text input
- WHEN the value is between 0 and 100
- THEN the volume is set to that exact value
- AND the settings_delta includes the new volume

#### Scenario: Invalid volume value is rejected

- GIVEN the user enters a volume value outside 0-100 range
- WHEN the module validates the input
- THEN the module returns an error status
- AND the volume is not changed

### Requirement: Mute Toggle

The system SHALL toggle the mute state of the default audio sink.

#### Scenario: Mute is toggled on

- GIVEN the audio is currently unmuted
- WHEN the user triggers mute toggle
- THEN the default sink is muted via the backend command
- AND the settings_delta includes `audio.muted: true`

#### Scenario: Mute is toggled off

- GIVEN the audio is currently muted
- WHEN the user triggers mute toggle
- THEN the default sink is unmuted via the backend command
- AND the settings_delta includes `audio.muted: false`

#### Scenario: Current mute state is detected

- GIVEN the module queries the current audio state
- WHEN it reads the sink properties
- THEN the mute state is parsed and displayed in the TUI
- AND the volume slider shows muted state visually

### Requirement: Sink Selection

The system SHALL list available audio sinks and allow the user to select the default sink.

#### Scenario: Available sinks are listed

- GIVEN multiple audio sinks are available
- WHEN the module queries sinks via the backend
- THEN each sink is listed with its name and description
- AND the current default sink is indicated

#### Scenario: Default sink is changed

- GIVEN the user selects a different sink from the list
- WHEN the module sets the default sink
- THEN the backend command to set default sink is executed
- AND the settings_delta includes `audio.default_sink: <new-sink-name>`

#### Scenario: No sinks available

- GIVEN no audio sinks are detected
- WHEN the module queries sinks
- THEN the module returns a warning status
- AND the TUI displays a message indicating no output devices

### Requirement: Volume State Persistence

The system SHALL persist volume, mute, and default sink settings in the settings schema.

#### Scenario: Settings are persisted after change

- GIVEN the user changes volume to 80%
- WHEN the module emits the settings_delta
- THEN the hub merges the delta into settings.json
- AND subsequent loads show volume as 80%

#### Scenario: Settings are restored on module reload

- GIVEN settings.json has `audio.volume: 60` and `audio.muted: true`
- WHEN the module loads
- THEN the TUI displays volume at 60% with muted indicator
- AND the actual audio state matches the settings

## Non-Functional Requirements

- **Dependencies**: The module SHALL declare `pipewire` or `pulseaudio` and `pactl` as dependencies
- **Performance**: Volume changes SHALL be reflected within 500ms
- **Backend agnostic**: The module SHALL work with both PipeWire and PulseAudio without user configuration
- **Testing**: All CLI tool interactions SHALL use interface-based mocking (`CLIExecutor`) for unit tests
