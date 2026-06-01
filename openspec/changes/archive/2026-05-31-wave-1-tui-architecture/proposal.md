# Proposal: Wave 1 — TUI Architecture

**Change name**: `wave-1-tui-architecture`
**Status**: `complete`
**Date**: 2026-05-31
**Author**: LambdaOS Team

---

## Executive Summary

Wave 1 establishes the architectural foundation for the LambdaOS TUI: a Go-based hub binary using bubbletea, a unified JSON settings schema with atomic writes, and a local pacman repository for packaging the TUI and its modules. This replaces the existing Python/Textual prototype (`src/os_tui_configurator/`) with the production Go implementation defined in the module interface contract.

**Wave 0 is complete** (CI builds ISO, boots in QEMU). Wave 1 delivers 3 specs in an estimated 3-5 days.

---

## Intent

Establish the TUI hub, settings schema, and pacman repo infrastructure so that:

1. `lambda-env` binary runs from terminal and displays a navigable menu
2. `~/.config/lambdaos/settings.json` exists with the unified schema, replacing `tui_settings.json` and `os_theme.json`
3. Local pacman repo structure is ready for packaging `lambdaos-tui`

---

## Scope

### In Scope

| Spec | What it covers |
|------|---------------|
| `core/01-hub-plugin-system` | Go module init, bubbletea TUI, hub binary, module discovery via manifest.json |
| `core/02-settings-schema` | Settings reader/writer with atomic writes, schema validation, default values, migration support |
| `infra-01-repo-pacman-setup` | Pacman repo structure, signing key, `repo-add` script, `pacman.conf` update |

### Out of Scope

- Module implementations (screen, audio, network, etc.) — Wave 2+
- PKGBUILD for `lambdaos-tui` — Wave 2 (`infra-02-repo-package-tui`)
- HTML prototypes — Wave 1 is infrastructure, no TUI views to prototype yet
- CI Go linting/test jobs — added as part of the Go module setup

---

## Current State Analysis

### What Exists Today

- **Python/Textual prototype**: `src/os_tui_configurator/` — a working but incomplete TUI using Textual (Python). Reads `tui_settings.json` from nvim config dir. This is a Wave 0 artifact that will be **replaced**, not migrated.
- **No Go module**: No `go.mod` or `go.sum` exists anywhere in the repo.
- **ISO build**: Standard mkarchiso with `profiledef.sh` at root, `airootfs/` overlay directory. Builds in Docker Arch container via CI.
- **CI**: GitHub Actions with Python lint (black, isort), shellcheck, luacheck, pytest, and ISO build. No Go tooling.
- **Settings**: Two separate files — `airootfs/etc/skel/dotfiles/nvim/.config/nvim/tui_settings.json` and `os_theme.json` at home root.
- **Pacman**: `pacman.conf` has `[custom]` repo commented out. No lambdaos repo.
- **Packages**: `go` is already in `packages.x86_64` (line 38).

### What Needs to Be Created

| Item | Path | Purpose |
|------|------|---------|
| Go module | `src/lambda-env/go.mod` | Root of Go workspace |
| Hub binary | `src/lambda-env/cmd/lambda-env/main.go` | Entry point |
| Hub package | `src/lambda-env/internal/hub/` | Plugin system, discovery |
| Settings package | `src/lambda-env/internal/settings/` | Schema reader/writer |
| TUI package | `src/lambda-env/internal/tui/` | Bubbletea models/views |
| Module contract | `src/lambda-env/internal/module/` | JSON protocol types |
| Pacman repo | `/srv/repo/lambdaos/` (on ISO) | Local package repository |
| Repo script | `scripts/repo-update.sh` | Regenerate repo database |

---

## Approach

### Step 1: Initialize Go Module with Bubbletea

```
src/lambda-env/
├── go.mod          ← module lambdaos.dev/lambda-env
├── go.sum
└── ...
```

- Initialize Go module at `src/lambda-env/`
- Add dependencies: `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss` (for styling), `github.com/charmbracelet/bubbles` (for list/input components)
- Add dev dependencies: `golangci-lint` (via Makefile target)
- Update Makefile: add `lint-go`, `test-go`, `build-go` targets
- Update CI: add Go lint + test jobs before ISO build

**Why bubbletea**: Elm architecture, excellent terminal rendering via `tcell`, battle-tested in production (charm.sh ecosystem), pure Go (no C dependencies), works in TTY without X11.

### Step 2: Create Settings Schema Reader/Writer

Package: `src/lambda-env/internal/settings/`

- Define Go structs matching the proposed schema from `core/02-settings-schema.md`
- `Load()` — reads `~/.config/lambdaos/settings.json`, returns defaults if missing
- `Save()` — atomic write (temp file + `os.Rename`)
- `Validate()` — JSON schema validation (required fields, valid types)
- `Migrate()` — detects version mismatch, adds missing fields with defaults
- `GetDelta()` — merges a settings delta from module response

**Atomic write pattern**:
```go
func (s *Store) Save(data Settings) error {
    tmp, _ := os.CreateTemp(dir, "settings-*.tmp")
    json.NewEncoder(tmp).Encode(data)
    tmp.Close()
    os.Rename(tmp.Name(), path)  // atomic on same filesystem
}
```

### Step 3: Create Hub Binary with Module Discovery

Package: `src/lambda-env/cmd/lambda-env/` (entry point)
Package: `src/lambda-env/internal/hub/` (plugin system)

- Hub binary scans `/usr/share/lambda-env/modules/` and `~/.local/share/lambda-env/modules/`
- Parses `manifest.json` for each module directory
- Validates required fields (name, version, description, category, requires_root, dependencies, min_hub_version)
- Merges: user modules override system modules with same name
- Sorts by category, then name
- Renders main menu via bubbletea with categories: System, Apps, Ops, Setup
- On module selection: executes `module run` with env vars (`LAMBDA_ENV_ACTION`, `LAMBDA_ENV_SETTINGS`, `LAMBDA_ENV_HUB_VERSION`, `LAMBDA_ENV_LOCALE`)
- Parses JSON from stdout, text from stderr
- Handles exit codes (0=success, 1=error, 2=warning)
- Merges `settings_delta` from module response atomically
- Logs errors to `/var/log/lambda-env/modules.log`

### Step 4: Set Up Pacman Repo Structure

- Create directory structure: `/srv/repo/lambdaos/x86_64/`
- Generate GPG signing key for LambdaOS repo
- Create `scripts/repo-update.sh`:
  ```bash
  #!/usr/bin/env bash
  cd /srv/repo/lambdaos
  repo-add --sign lambdaos.db.tar x86_64/*.pkg.tar.zst
  ```
- Update `pacman.conf` to include:
  ```
  [lambdaos]
  Server = file:///srv/repo/lambdaos/$arch
  SigLevel = Required
  ```
- For ISO live: repo is local (file://), not remote

### Step 5: Update CI for Go

- Add Go setup step to CI workflow
- Add `golangci-lint` run
- Add `go test` run
- Ensure Go build happens before ISO build (binary must be in airootfs)

---

## Framework Decision: Go + Bubbletea

**Already decided** in `module-interface-contract.md`. Confirmed after analysis:

| Criterion | Bubbletea (Go) | Textual (Python) | whiptail (bash) |
|-----------|---------------|-----------------|-----------------|
| TTY support | ✅ Native | ❌ Requires Python runtime | ✅ Native |
| Binary size | ~10-15MB | N/A (interpreted) | ~50KB |
| Dependencies | None (static) | Python + pip packages | None |
| Testability | ✅ Unit testable | ✅ Unit testable | ❌ Hard to test |
| Learning curve | Moderate | Moderate | Low |
| Ecosystem | Rich (charm.sh) | Growing | Minimal |
| Performance | Fast | Good | Fast |

**Decision**: Go + bubbletea. The module interface contract already mandates Go for the hub and default modules. Bubbletea is the natural choice for Go TUIs.

---

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Go binary size (~10-15MB) | Low — ISO is ~800MB+, 15MB is negligible | Acceptable; static binary means no runtime deps |
| Bubbletea learning curve | Medium — team may not know Elm architecture | Start with simple list model, iterate |
| Packaging for Arch | Medium — need PKGBUILD for Go project | Use `go build` in PKGBUILD, standard pattern |
| Atomic write on overlay filesystem | Low — tmp + rename works on same fs | Ensure temp file is in same directory as target |
| CI build time increase | Low — Go build is fast (~10s) | Parallel jobs, cached Go modules |

---

## Estimated Lines

| Component | Estimated Lines | Notes |
|-----------|----------------|-------|
| `go.mod` + `go.sum` | ~50 | Dependencies |
| `cmd/lambda-env/main.go` | ~80 | Entry point, cobra/cli flags |
| `internal/hub/discovery.go` | ~150 | Module scanning, manifest parsing |
| `internal/hub/execution.go` | ~120 | Module execution, JSON parsing |
| `internal/hub/hub.go` | ~100 | Hub struct, main loop |
| `internal/settings/schema.go` | ~200 | Go structs for full schema |
| `internal/settings/store.go` | ~150 | Load, Save, Validate, Migrate |
| `internal/tui/menu.go` | ~200 | Bubbletea main menu model |
| `internal/tui/module_view.go` | ~150 | Module execution view |
| `internal/module/types.go` | ~100 | JSON protocol types |
| `internal/module/logger.go` | ~80 | Module log writer |
| `scripts/repo-update.sh` | ~30 | Repo database regeneration |
| `pacman.conf` (updated) | ~5 | Added [lambdaos] section |
| CI workflow updates | ~40 | Go lint + test jobs |
| Makefile updates | ~20 | Go targets |
| **Total** | **~1,475 lines** | |

---

## File Changes Summary

### New Files
- `src/lambda-env/go.mod`
- `src/lambda-env/go.sum`
- `src/lambda-env/cmd/lambda-env/main.go`
- `src/lambda-env/internal/hub/hub.go`
- `src/lambda-env/internal/hub/discovery.go`
- `src/lambda-env/internal/hub/execution.go`
- `src/lambda-env/internal/settings/schema.go`
- `src/lambda-env/internal/settings/store.go`
- `src/lambda-env/internal/settings/store_test.go`
- `src/lambda-env/internal/tui/menu.go`
- `src/lambda-env/internal/tui/module_view.go`
- `src/lambda-env/internal/module/types.go`
- `src/lambda-env/internal/module/logger.go`
- `scripts/repo-update.sh`

### Modified Files
- `pacman.conf` — add `[lambdaos]` repo section
- `Makefile` — add Go targets (lint-go, test-go, build-go)
- `.github/workflows/ci.yml` — add Go lint + test jobs
- `packages.x86_64` — no changes needed (go already included)

### Not Touched (Yet)
- `src/os_tui_configurator/` — existing Python TUI, will be removed in Wave 2 when PKGBUILD replaces it
- `airootfs/etc/skel/dotfiles/nvim/.config/nvim/tui_settings.json` — migration happens in Wave 2
- `airootfs/etc/skel/dotfiles/nvim/.config/nvim/lua/core/tui_bridge.lua` — updated in Wave 2 to read from new schema path

---

## Next Recommended: `specs`

The orchestrator should now launch spec writing for the 3 specs:
1. `core/01-hub-plugin-system` — detailed requirements and scenarios
2. `core/02-settings-schema` — schema validation, atomic writes, migration
3. `infra-01-repo-pacman-setup` — repo structure, signing, scripts
