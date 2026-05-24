# AUR Install Script Specification

## Purpose

Defines a post-boot script for installing AUR-only packages that cannot be included in the archiso build, with clear documentation and error resilience.

## Requirements

### Requirement: AUR Install Script Exists

The system MUST provide a shell script at `scripts/aur-packages.sh` that installs all AUR-only packages documented in the README.

#### Scenario: Script file exists and is executable

- GIVEN the LambdaOS repository is cloned
- WHEN the user checks `scripts/aur-packages.sh`
- THEN the file exists and has executable permissions

#### Scenario: Script documents all AUR packages

- GIVEN aur-packages.sh exists
- WHEN the user reads the script
- THEN it lists: `spotify`, `obsidian`, `megasync`, `bluetui`, `impala` with comments describing each package's purpose

### Requirement: Script Requires AUR Helper

The system SHALL document that an AUR helper (yay or paru) MUST be installed before running the script, and SHALL provide installation instructions if missing.

#### Scenario: AUR helper check and install guidance

- GIVEN a fresh LambdaOS installation without an AUR helper
- WHEN the user runs aur-packages.sh
- THEN the script detects the missing helper and prints instructions to install yay or paru
- AND the script exits with a non-zero status code

#### Scenario: AUR helper is available

- GIVEN yay or paru is installed on the system
- WHEN the user runs aur-packages.sh
- THEN the script proceeds with package installation using the detected helper

### Requirement: Script Continues on Individual Package Failure

The system MUST not abort the entire script when a single AUR package fails to build or install.

#### Scenario: One AUR package fails, others continue

- GIVEN an AUR package fails to build (e.g., upstream source unavailable)
- WHEN aur-packages.sh runs
- THEN the script logs the failure for that package
- AND the script continues installing remaining packages
- AND the script exits with a summary of successful and failed installations

#### Scenario: All AUR packages install successfully

- GIVEN all AUR packages are buildable
- WHEN aur-packages.sh completes
- THEN all listed packages are installed
- AND the script exits with status code 0

### Requirement: Script Uses --needed Flag

The system SHOULD use the `--needed` flag (or equivalent) when invoking the AUR helper to avoid reinstalling already-present packages.

#### Scenario: Re-running script skips installed packages

- GIVEN some AUR packages are already installed
- WHEN the user runs aur-packages.sh again
- THEN the AUR helper skips packages that are already up-to-date
- AND only missing or outdated packages are processed
