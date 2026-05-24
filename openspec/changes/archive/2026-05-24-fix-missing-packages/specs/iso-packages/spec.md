# ISO Packages Specification

## Purpose

Defines the complete package list for the LambdaOS ISO, ensuring all advertised features are included and installable from official repositories.

## Requirements

### Requirement: All Advertised Official Packages Included

The system MUST include every package advertised in the README that exists in official Arch Linux repositories (core, extra, multilib) in packages.x86_64.

#### Scenario: Networking packages are present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `tailscale` is included for VPN functionality

#### Scenario: Office suite packages are present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `libreoffice-fresh` is included for office productivity

#### Scenario: Email client is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `thunderbird` is included for email functionality

#### Scenario: Multimedia player is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `vlc` is included for media playback

#### Scenario: PDF viewer is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `okular` is included for PDF and document viewing

#### Scenario: Calculator is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `qalculate-gtk` is included for calculation functionality

#### Scenario: Password manager is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `keepassxc` is included for credential management

#### Scenario: Web browser is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `chromium` is included as a web browser

#### Scenario: Terminal file manager is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `yazi` is included for terminal-based file management

### Requirement: Virtualization and Container Tools Included

The system MUST include virtualization and container packages for development workflows.

#### Scenario: VirtualBox is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `virtualbox` is included for virtualization

#### Scenario: Docker toolchain is present

- GIVEN packages.x86_64 is the ISO package list
- WHEN the ISO is built
- THEN `docker`, `docker-compose`, and `lazydocker` are all included

### Requirement: Gaming and Compatibility Layer Packages Included

The system MUST include gaming and Windows compatibility packages from the multilib repository.

#### Scenario: Steam is present

- GIVEN `[multilib]` is enabled in pacman.conf AND packages.x86_64 is the package list
- WHEN the ISO is built
- THEN `steam` is included and installable

#### Scenario: Wine is present

- GIVEN `[multilib]` is enabled in pacman.conf AND packages.x86_64 is the package list
- WHEN the ISO is built
- THEN `wine`, `wine-mono`, and `winetricks` are included

### Requirement: Package List Maintains Alphabetical Order

The system SHOULD maintain alphabetical ordering in packages.x86_64 for readability and merge conflict reduction.

#### Scenario: New packages are inserted in sorted position

- GIVEN packages.x86_64 is alphabetically sorted
- WHEN new packages are added
- THEN each package is inserted at its correct alphabetical position

### Requirement: All Packages Are Installable

The system MUST ensure every package in packages.x86_64 exists in an enabled repository at build time.

#### Scenario: No package resolution failures

- GIVEN packages.x86_64 contains all listed packages
- WHEN mkarchiso resolves the package list
- THEN zero packages fail with "target not found" errors
