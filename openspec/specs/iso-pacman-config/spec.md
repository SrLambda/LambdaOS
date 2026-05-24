# ISO Pacman Configuration Specification

## Purpose

Defines the pacman.conf configuration required for the LambdaOS ISO, including repository selection and multilib support for 32-bit compatibility.

## Requirements

### Requirement: Multilib Repository Enabled

The system MUST enable the `[multilib]` repository in pacman.conf to support 32-bit packages required by Steam, Wine, and other multilib software.

#### Scenario: Multilib section is uncommented and functional

- GIVEN pacman.conf exists in the archiso profile
- WHEN the ISO is built with mkarchiso
- THEN the `[multilib]` section header is uncommented (not prefixed with `#`)
- AND the `Include = /etc/pacman.d/mirrorlist` line under `[multilib]` is uncommented

#### Scenario: Multilib repository is accessible during ISO build

- GIVEN `[multilib]` is enabled in pacman.conf
- WHEN pacman synchronizes package databases during ISO build
- THEN multilib packages are resolvable (e.g., `pacman -Ss lib32-*` returns results)

#### Scenario: Core and extra repositories remain enabled

- GIVEN pacman.conf is modified to enable multilib
- WHEN the ISO is built
- THEN `[core]` and `[extra]` sections remain uncommented and functional
- AND no existing repository configuration is removed or disabled

### Requirement: Repository Order and Priority

The system SHALL maintain repository order with `[core]` before `[extra]` before `[multilib]` to follow Arch Linux best practices.

#### Scenario: Repositories are in correct order

- GIVEN pacman.conf contains all three repositories
- WHEN pacman resolves package dependencies
- THEN `[core]` packages take priority over `[extra]`, and `[extra]` over `[multilib]`
