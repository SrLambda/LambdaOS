# README Accuracy Specification

## Purpose

Ensures the README.md accurately distinguishes between packages included in the ISO versus those requiring post-boot installation, so users have correct expectations.

## Requirements

### Requirement: Package Installation Source Is Clearly Labeled

The system MUST clearly label each advertised package or feature as either "included in ISO" or "requires post-boot AUR install" in README.md.

#### Scenario: ISO-included packages are marked

- GIVEN README.md lists available software
- WHEN a user reads the package list
- THEN packages in packages.x86_64 are visually distinguished (e.g., section heading, badge, or label) as included in the ISO

#### Scenario: AUR packages are marked as post-install

- GIVEN README.md lists available software
- WHEN a user reads the package list
- THEN AUR-only packages (spotify, obsidian, megasync, bluetui, impala) are labeled as requiring post-boot installation
- AND a reference to `scripts/aur-packages.sh` is provided

#### Scenario: Multilib packages are labeled

- GIVEN README.md lists gaming/compatibility features
- WHEN a user reads Steam or Wine entries
- THEN they are labeled as requiring multilib (included in ISO once multilib is enabled)

### Requirement: README Matches Actual ISO Content

The system MUST not advertise any package as "included" unless it is present in packages.x86_64 or enabled via pacman.conf configuration.

#### Scenario: No phantom packages in README

- GIVEN README.md claims a package is included
- WHEN packages.x86_64 and pacman.conf are checked
- THEN the package is either in packages.x86_64 or enabled via repository configuration

#### Scenario: All packages.x86_64 packages are documented

- GIVEN packages.x86_64 contains a user-facing package
- WHEN README.md is reviewed
- THEN the package or its feature category is mentioned in the README

### Requirement: Post-Install Instructions Are Actionable

The system SHOULD provide clear, copy-pasteable instructions for post-boot AUR package installation.

#### Scenario: User can follow post-install steps without ambiguity

- GIVEN a user just booted the LambdaOS ISO for the first time
- WHEN they read the README post-install section
- THEN the steps to install yay/paru and run aur-packages.sh are explicit and complete
