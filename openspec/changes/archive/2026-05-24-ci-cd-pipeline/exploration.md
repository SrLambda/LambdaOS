# Exploration: CI/CD Automatizado con Entregas de Versiones del OS

## Current State

### ISO Build Process
- **Build orchestration**: `build_and_test.sh` (32 lines) runs 5 phases:
  1. Python venv setup + `pip install -r requirements-dev.txt`
  2. Clean previous builds (`sudo rm -rf work/ out/`)
  3. Verify pacman.conf (falls back to archiso releng default)
  4. `sudo mkarchiso -v -w work/ -o out/ .`
  5. Run QEMU E2E tests (`pytest tests/qemu/test_live_boot.py -v`)
- **No Makefile exists** — all commands are in shell scripts
- **No CI/CD** — zero `.github/workflows/` directory

### Versioning
- `profiledef.sh` line 8: `iso_version="$(date --date="@${SOURCE_DATE_EPOCH:-$(date +%s)}" +%Y.%m.%d)"`
- Date-based versioning (e.g., `2026.05.24`) — no git tags, no semantic versioning
- ISO output: `out/lambda-os-*-x86_64.iso`

### Test Infrastructure
- **Unit tests** (`tests/unit/`): 7 test files, 68/68 passing
  - pytest + pytest-asyncio + Textual `run_test()` + AST parsing
  - Fast, no system dependencies
- **QEMU E2E tests** (`tests/qemu/`): 3 test scenarios
  - Requires: qemu-system-x86_64, built ISO, sudo/KVM
  - Boots ISO headless via serial console with pexpect
  - Tests: boot to shell, stow symlinks, neovim init.lua
  - NOT feasible in GitHub Actions (nested KVM not available on standard runners, very slow without it)
- **Legacy backup tests** (`backup_test/`): not relevant for CI

### Existing Specs (openspec/specs/)
- `iso-packages/` — package list validation
- `iso-pacman-config/` — pacman.conf validation
- `aur-install-script/` — AUR package script
- `readme-accuracy/` — README consistency

### Key Technical Constraints
1. **mkarchiso requires sudo** — available in GitHub Actions (runner has sudo)
2. **QEMU E2E needs KVM** — GitHub Actions Linux runners have KVM available via `sudo modprobe kvm`, but full ISO boot is extremely slow (15-30 min for build + 5-10 min for QEMU boot)
3. **ISO size ~2-3 GB** — within GitHub Release limits (2GB per file, 5GB total per release)
4. **No pyproject.toml** — Python deps in `requirements-dev.txt` and `src/requirements.txt`
5. **Build time** — mkarchiso takes 10-30+ min depending on network and cache

## Gap Analysis

| Need | Current State | Gap |
|------|--------------|-----|
| Automated testing on push/PR | None | Need CI workflow for lint + unit tests |
| Automated ISO builds | Manual `sudo mkarchiso` | Need release workflow triggered by tags |
| Versioned releases | Date-based in profiledef.sh | Need git-tag-based versioning |
| Release artifacts | ISO in local `out/` | Need GitHub Release with ISO + checksums |
| Linting in CI | Linters available but not automated | Need lint step in CI |
| E2E tests in CI | QEMU (needs KVM) | Not feasible for standard CI runners |

## Feasibility: GitHub Actions

| Feature | Feasible? | Notes |
|---------|-----------|-------|
| Unit tests on push/PR | YES | Fast, no system deps |
| Lint on push/PR | YES | black, isort, shellcheck, shfmt |
| ISO build on tag | YES | sudo available, ~15-30 min build time |
| QEMU E2E in CI | RISKY | KVM now available on ubuntu-latest but build+boot time is prohibitive for every push |
| Release upload | YES | GitHub Releases API, assets up to 2GB |

## Recommended Architecture

```
.github/workflows/
  ci.yml          — lint + unit tests on push/PR
  release.yml     — build ISO + publish release on tag push
```

**CI workflow** (every push/PR):
1. Checkout code
2. Set up Python + deps
3. Run linters (black --check, isort --check, shellcheck, shfmt)
4. Run unit tests (`pytest tests/unit/ -v`)

**Release workflow** (tag push `v*`):
1. Checkout code at tag
2. Set SOURCE_DATE_EPOCH from git tag timestamp (reproducible builds)
3. Install archiso packages
4. Build ISO with `sudo mkarchiso`
5. Generate SHA256 checksums
6. Create GitHub Release with ISO + checksums

**Versioning approach**: 
- Git tags follow semver: `v1.0.0`, `v1.1.0`, etc.
- `profiledef.sh` updated to derive version from git tag or env var
- `SOURCE_DATE_EPOCH` set from tag commit for reproducibility

**E2E testing** remains manual/local for now — could later add a dedicated self-hosted runner.