# Release Pipeline Specification

## Purpose

Automated ISO build and GitHub Release creation triggered by semantic version tags, producing downloadable release artifacts with integrity checksums.

## Requirements

### Requirement: Release Trigger on Version Tag

The release pipeline SHALL execute only when a git tag matching the `v*` pattern is pushed.

#### Scenario: Valid semver tag pushed

- GIVEN a tag matching `v{major}.{minor}.{patch}` (e.g., `v1.0.0`) is pushed
- WHEN the tag push event occurs
- THEN the release workflow SHALL start automatically

#### Scenario: Non-version tag pushed

- GIVEN a tag NOT matching the `v*` pattern is pushed
- WHEN the tag push event occurs
- THEN the release workflow SHALL NOT trigger

#### Scenario: Tag deleted

- GIVEN an existing version tag is deleted
- WHEN the tag deletion event occurs
- THEN the release workflow SHALL NOT trigger

### Requirement: Reproducible ISO Build

The release pipeline SHALL build the ISO with `SOURCE_DATE_EPOCH` set from the tagged commit timestamp for reproducible builds.

#### Scenario: SOURCE_DATE_EPOCH set from tag

- GIVEN a version tag points to a specific commit
- WHEN the release workflow starts the build step
- THEN `SOURCE_DATE_EPOCH` SHALL be set to that commit's timestamp
- AND `mkarchiso` SHALL use this value for deterministic output

#### Scenario: Build requires sudo

- GIVEN the release workflow runs on a GitHub Actions runner
- WHEN the ISO build step executes
- THEN `sudo mkarchiso` SHALL be used
- AND the build SHALL complete without interactive prompts

### Requirement: ISO Build Completes Successfully

The release pipeline SHALL produce a valid ISO file in the output directory.

#### Scenario: Successful ISO build

- GIVEN all build dependencies are installed
- WHEN `sudo mkarchiso` completes
- THEN exactly one `.iso` file SHALL exist in the output directory
- AND the ISO filename SHALL contain the version from the tag

#### Scenario: Build dependency missing

- GIVEN required packages (archiso, base-devel) are not installed
- WHEN the build step attempts to run `mkarchiso`
- THEN the workflow SHALL install dependencies before building
- AND the build SHALL proceed after installation

#### Scenario: Build failure

- GIVEN a configuration error causes `mkarchiso` to fail
- WHEN the build step executes
- THEN the workflow SHALL fail immediately
- AND no GitHub Release SHALL be created

### Requirement: SHA256 Checksum Generation

The release pipeline SHALL generate SHA256 checksums for the built ISO.

#### Scenario: Checksum file created

- GIVEN a valid ISO file exists in the output directory
- WHEN the checksum generation step runs
- THEN a `.sha256` file SHALL be created alongside the ISO
- AND the checksum file SHALL contain the correct SHA256 hash and filename

#### Scenario: Checksum verification

- GIVEN an ISO and its corresponding `.sha256` file
- WHEN `sha256sum --check` is run against the checksum file
- THEN verification SHALL pass

### Requirement: GitHub Release Creation

The release pipeline SHALL create a GitHub Release with the ISO and checksum file as assets.

#### Scenario: Release with assets

- GIVEN a successful ISO build and checksum generation
- WHEN the release step executes
- THEN a GitHub Release SHALL be created with the tag name as the release title
- AND the ISO file SHALL be uploaded as a release asset
- AND the `.sha256` file SHALL be uploaded as a release asset

#### Scenario: Release already exists

- GIVEN a GitHub Release already exists for the tag
- WHEN the release step executes
- THEN the workflow SHALL fail or skip (no duplicate releases)

### Requirement: Build Timeout Protection

The release pipeline SHALL enforce a maximum build duration to prevent hung builds.

#### Scenario: Build completes within timeout

- GIVEN the ISO build completes in under 45 minutes
- WHEN the build step finishes
- THEN the workflow SHALL proceed to release creation

#### Scenario: Build exceeds timeout

- GIVEN the ISO build exceeds 45 minutes
- WHEN the timeout is reached
- THEN the workflow SHALL be cancelled
- AND no partial release SHALL be created
