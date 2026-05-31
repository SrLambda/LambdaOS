# Versioning Specification

## Purpose

Replace date-based ISO versioning with semantic versioning derived from git tags, enabling reproducible builds and clear release identification.

## Requirements

### Requirement: Semantic Version Format

All release versions SHALL follow the Semantic Versioning 2.0.0 scheme with a `v` prefix.

#### Scenario: Valid version tag

- GIVEN a developer creates a release tag
- WHEN the tag follows the format `v{major}.{minor}.{patch}` (e.g., `v1.0.0`, `v2.3.1`)
- THEN the tag SHALL be accepted by the release pipeline

#### Scenario: Invalid version tag

- GIVEN a tag does not follow the `v{major}.{minor}.{patch}` format
- WHEN the tag is pushed (e.g., `1.0.0`, `v1.0`, `v1.0.0-beta`)
- THEN the release pipeline SHALL NOT trigger for this tag

### Requirement: Version Derived from Git Tag

The ISO version SHALL be derived from the git tag name, not from the build timestamp.

#### Scenario: Version from tag name

- GIVEN a tag `v1.2.3` is checked out
- WHEN the build process determines the version
- THEN the version SHALL be `1.2.3` (without the `v` prefix)
- AND the ISO filename SHALL be `LambdaOS-1.2.3-x86_64.iso`

#### Scenario: Version from GITHUB_REF_NAME

- GIVEN the build runs in GitHub Actions with tag `v1.0.0`
- WHEN `profiledef.sh` reads the version
- THEN it SHALL use `$GITHUB_REF_NAME` (stripped of `v` prefix) as the version

#### Scenario: Local build fallback

- GIVEN the build runs locally without `GITHUB_REF_NAME`
- WHEN `profiledef.sh` determines the version
- THEN it SHALL fall back to `git describe --tags --abbrev=0`
- AND if no tags exist, it SHALL fall back to `git rev-parse --short HEAD`

### Requirement: SOURCE_DATE_EPOCH for Reproducibility

The build system SHALL set `SOURCE_DATE_EPOCH` to ensure reproducible ISO builds.

#### Scenario: Epoch from tag commit

- GIVEN a version tag points to a specific commit
- WHEN the build starts
- THEN `SOURCE_DATE_EPOCH` SHALL be set to the commit timestamp of the tagged commit
- AND the same tag on different machines SHALL produce identical ISO checksums

#### Scenario: Epoch from environment override

- GIVEN `SOURCE_DATE_EPOCH` is already set in the environment
- WHEN the build starts
- THEN the existing value SHALL be used (no override)

### Requirement: Profiledef Version Logic

The `profiledef.sh` file SHALL determine the ISO version using the priority: environment variable > git tag > git short hash.

#### Scenario: Environment variable takes priority

- GIVEN `ISO_VERSION` is set to `1.5.0` in the environment
- WHEN `profiledef.sh` executes
- THEN `iso_version` SHALL be `1.5.0` regardless of git state

#### Scenario: Git tag when no env var

- GIVEN no `ISO_VERSION` environment variable is set
- AND a git tag `v1.0.0` exists on the current commit
- WHEN `profiledef.sh` executes
- THEN `iso_version` SHALL be `1.0.0`

#### Scenario: Short hash fallback

- GIVEN no `ISO_VERSION` environment variable is set
- AND no git tags exist on the current commit
- WHEN `profiledef.sh` executes
- THEN `iso_version` SHALL be the short git commit hash

### Requirement: Date-Based Versioning Removed

The previous date-based versioning scheme SHALL be fully replaced.

#### Scenario: No date-based output

- GIVEN a build is triggered
- WHEN the ISO is produced
- THEN the filename SHALL NOT contain a date pattern like `2026.05.24`
- AND `profiledef.sh` SHALL NOT use `date` for version generation
