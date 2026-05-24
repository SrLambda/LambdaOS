# Proposal: CI/CD Automatizado con Entregas de Versiones del OS

## Intent

LambdaOS has zero CI/CD — every build, test, and release is manual. We need automated testing on every push/PR to catch regressions early, and automated ISO builds + GitHub Releases on version tags so users can download versioned ISOs without a developer manually running `mkarchiso`.

## Scope

### In Scope
- GitHub Actions CI workflow (lint + unit tests on push/PR)
- GitHub Actions release workflow (ISO build + publish on tag push)
- Semantic versioning via git tags (replacing date-based versioning)
- SHA256 checksums for release artifacts
- Makefile for consistent build/test/lint targets

### Out of Scope
- QEMU E2E tests in CI (requires self-hosted runner or prohibitively slow)
- Automated AUR package installation testing
- Docker-based build isolation
- Nightly/scheduled builds
- Code signing or GPG verification of ISOs

## Capabilities

### New Capabilities
- `ci-pipeline`: Automated lint + unit test execution on every push and PR
- `release-pipeline`: Automated ISO build and GitHub Release creation on tag push
- `versioning`: Semantic versioning scheme derived from git tags with reproducible builds

### Modified Capabilities
- `iso-pacman-config`: No spec-level change (CI validates config but doesn't change behavior)

## Approach

Create two GitHub Actions workflows plus a Makefile:

1. **`.github/workflows/ci.yml`** — triggers on push/PR to main
   - Python setup + deps installation
   - Lint: `black --check`, `isort --check`, `shellcheck`, `shfmt`
   - Test: `pytest tests/unit/ -v`

2. **`.github/workflows/release.yml`** — triggers on tag push `v*`
   - Build ISO with `mkarchiso` (sudo available in GH Actions)
   - Set `SOURCE_DATE_EPOCH` from git tag for reproducibility
   - Generate SHA256 checksums
   - Create GitHub Release with ISO + checksums as assets

3. **`Makefile`** — consistent targets for local development
   - `make lint`, `make test`, `make build`, `make release`
   - Mirrors what CI does locally

4. **Versioning** — derive from git tag
   - Tag format: `v1.0.0`, `v1.1.0`, etc.
   - `profiledef.sh` reads version from `GITHUB_REF_NAME` or falls back to git describe

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `.github/workflows/ci.yml` | New | CI pipeline workflow |
| `.github/workflows/release.yml` | New | Release pipeline workflow |
| `Makefile` | New | Build/test/lint targets |
| `profiledef.sh` | Modified | Version derived from git tag/env |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| ISO build time exceeds 30 min in CI | Med | Cache pacman package cache between runs |
| Release ISO exceeds 2GB file limit | Low | Monitor size; most Arch ISOs are 1.5-2.5GB |
| QEMU E2E not feasible in CI | High | Keep E2E as local-only; add self-hosted runner later |
| mkarchiso fails due to missing deps | Low | Install archiso + base-devel in workflow |

## Rollback Plan

1. Delete `.github/workflows/` directory — CI is fully additive, no existing code changes
2. Revert `profiledef.sh` to date-based versioning (single line)
3. Remove `Makefile` — no runtime dependency on it

## Dependencies

- GitHub Actions (ubuntu-latest runner with sudo)
- `archiso` package (installable via pacman in workflow)
- Git tags must follow `v*` pattern for release workflow to trigger

## Success Criteria

- [ ] CI workflow runs on every push/PR and passes lint + unit tests
- [ ] Release workflow builds ISO on tag push and creates GitHub Release with assets
- [ ] `make lint` and `make test` work locally
- [ ] ISO version in filename matches git tag (e.g., `v1.0.0` → `lambda-os-1.0.0-x86_64.iso`)
- [ ] SHA256 checksums are generated and uploaded alongside ISO