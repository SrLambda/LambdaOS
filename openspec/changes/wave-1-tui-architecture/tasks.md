# Tasks: Wave 1 — TUI Architecture

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~1,515 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 infra → PR 2 settings → PR 3 hub-p1 → PR 4 hub-p2 |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Pacman repo, scripts, PKGBUILD template, pacman.conf | PR 1 (~80 lines) | Base: wave-1. Independent infra slice |
| 2 | Settings structs, store, migration, unit tests | PR 2 (~480 lines) | Base: wave-1. No TUI deps, just stdlib |
| 3 | Module types, discovery, hub controller, logger | PR 3 (~400 lines) | Base: settings branch. Imports settings |
| 4 | TUI model/update/view, execution, main.go | PR 4 (~350 lines) | Base: PR 3 branch. Completes TUI |
| 5 | CI Go jobs, Makefile targets | (merged into PR 3-4) | ~50 lines, inline in hub PRs |
| 6 | CI autofix workflow (lint-fix.yml) | (merged into PR 4) | ~40 lines, new workflow file |
| 6 | CI autofix workflow (lint-fix.yml) | (merged into PR 4) | ~40 lines, new workflow file |

## Phase 1: Foundation — Types + Constants

- [ ] 1.1 Init `src/lambda-env/go.mod` with module path `lambdaos.dev/lambda-env`
- [ ] 1.2 Add bubbletea, lipgloss, bubbles deps via `go get`
- [ ] 1.3 Create `pkg/module/manifest.go` — Manifest + Response structs, Validate(), Helper()
- [ ] 1.4 Create `pkg/module/manifest_test.go` — table-driven validation tests (valid/invalid, edge cases)
- [ ] 1.5 Create `pkg/version/version.go` — `Version = "0.1.0"` constant

## Phase 2: Infrastructure — Pacman Repo

- [x] 2.1 Create `airootfs/srv/repo/lambdaos/x86_64/` directory in airootfs overlay
- [x] 2.2 Create `scripts/repo-update.sh` — root check, `repo-add --sign`, empty-repo guard
- [x] 2.3 Create `templates/PKGBUILD.lambdaos-tui` — Go build + install to `/usr/bin/lambda-env`
- [x] 2.4 Modify `pacman.conf` — add `[lambdaos]` section after `[multilib]` with `SigLevel = Required`

## Phase 3: Core — Settings Schema

- [ ] 3.1 Create `internal/settings/schema.go` — 9 section structs + defaults + `Validate()` typed checks
- [ ] 3.2 Create `internal/settings/store.go` — `Load()` (with defaults+migration), `Save()` (atomic temp+rename), `SaveDelta()` (deep merge), `Migrate()` (version check)
- [ ] 3.3 Create `internal/settings/store_test.go` — atomic write, delta merge, field preservation, downgrade rejection, migration adds missing fields

## Phase 4: Core — Hub System

- [ ] 4.1 Create `internal/hub/discovery.go` — `Scan()` system+user module paths, manifest parse, validation, user-override-system merge
- [ ] 4.2 Create `internal/hub/hub.go` — Hub struct, `BuildMenu()` with category groups, `CheckDeps(pacman -Q)`, `CheckRoot(sudo)`
- [ ] 4.3 Create `internal/module/logger.go` — file writer for `/var/log/lambda-env/modules.log`, structured format with timestamp+module+exit+stdout+stderr+env

## Phase 5: TUI + Execution

- [ ] 5.1 Create `internal/tui/model.go` — Model struct (categories, items, cursor), `Init()`
- [ ] 5.2 Create `internal/tui/update.go` — `Update()`: arrow keys, Enter (select+exec), Esc (back), q (quit)
- [ ] 5.3 Create `internal/tui/view.go` — `View()`: category headers with count, sorted module list, error/warning overlay
- [ ] 5.4 Create `internal/hub/execution.go` — `exec.Command` with env vars (`LAMBDA_ENV_ACTION`, `LAMBDA_ENV_SETTINGS`, `LAMBDA_ENV_HUB_VERSION`, `LAMBDA_ENV_LOCALE`), timeout, JSON stdout parse, `settings_delta` merge
- [ ] 5.5 Create `cmd/lambda-env/main.go` — flag parsing (--help), init Hub+Store, run bubbletea program

## Phase 6: Integration + CI

- [ ] 6.1 Create `test/fixtures/modules/` — mock module dirs with known manifest.json + module scripts
- [ ] 6.2 Create `test/hub_integration_test.go` — mock module exec, JSON parse, log verification, delta merge
- [ ] 6.3 Modify `Makefile` — add `lint-go`, `test-go`, `build-go` targets
- [ ] 6.4 Modify `.github/workflows/ci.yml` — add Go setup, golangci-lint, go test, go build jobs
- [ ] 6.5 Create `.github/workflows/lint-fix.yml` — autofix CI for push to non-main branches (black, isort, shfmt -w), commits fixes with `[bot]` prefix; leaves shellcheck/luacheck as gate in existing ci.yml
