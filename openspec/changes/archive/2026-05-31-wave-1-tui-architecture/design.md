# Design: Wave 1 TUI Architecture

## Overview

Establish the `lambda-env` Go binary as a bubbletea-based TUI hub. Three specs converge into a single Go module at `src/lambda-env/`: hub plugin system (core/01), settings schema with atomic writes (core/02), and local pacman repo (infra/01). The existing Python/Textual prototype (`src/os_tui_configurator/`) is replaced, not migrated.

## Architecture Diagram

```
┌──────────────────── lambda-env TUI ────────────────────┐
│  ┌──────────┐   ┌──────────────┐   ┌────────────────┐  │
│  │ Bubbletea│──▶│  Hub         │──▶│ Settings Store │  │
│  │ UI       │◀──│  Controller  │◀──│  (singleton)   │  │
│  └──────────┘   │              │   └───────┬────────┘  │
│                 │  ┌──────────┐│           │           │
│                 │  │Discovery ││   ~/.config/lambdaos/ │
│                 │  │Scanner   ││   settings.json      │
│                 │  └────┬─────┘│                       │
│                 │       │      │                       │
│                 │  ┌────▼─────┐│                       │
│                 │  │ Executor ││                       │
│                 │  └────┬─────┘│                       │
│                 └───────┼──────┘                       │
├─────────────────────────┼──────────────────────────────┤
│         Module Execution (JSON over stdout)             │
│  /usr/share/lambda-env/modules/ + ~/.local/share/...   │
│  ┌──────────────────────────────────────────────────┐  │
│  │ /var/log/lambda-env/modules.log                   │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

## Go Module Structure

```
src/lambda-env/
├── go.mod                          # module lambdaos.dev/lambda-env
├── go.sum
├── cmd/lambda-env/main.go          # Entry point, CLI flags
├── internal/
│   ├── hub/
│   │   ├── hub.go                  # Hub struct, main loop, dependency check
│   │   ├── discovery.go            # Module scanning, manifest parsing
│   │   └── execution.go            # Module exec, JSON parsing, logging
│   ├── settings/
│   │   ├── schema.go               # Go structs, defaults, validation
│   │   ├── store.go                # Load, Save(Settings), SaveDelta(map), Migrate
│   │   └── store_test.go           # Unit tests: atomic write, delta merge, migration
│   └── tui/
│       ├── model.go                # Bubbletea model (Elm: Model+Init)
│       ├── update.go               # Message handling (Update)
│       └── view.go                 # Menu/category rendering (View)
├── pkg/module/
│   ├── manifest.go                 # Manifest JSON types, validation
│   └── manifest_test.go            # Manifest validation unit tests
├── pkg/version/
│   └── version.go                  # Hub version constant (semver)
└── test/                           # Integration tests (mock modules)
    ├── hub_integration_test.go
    └── fixtures/
        └── modules/                # Test module directories
```

## Architecture Decisions

| Decision | Choice | Rejected | Rationale |
|----------|--------|----------|-----------|
| TUI framework | bubbletea (Go) | Textual (Python), whiptail (bash) | Static binary, no runtime deps, TTY-native, testable. Already mandated in module-interface-contract.md |
| UI pattern | Elm (Model/Update/View) | MVC, event-bus | bubbletea native pattern; single-direction data flow simplifies state management |
| Settings access | Singleton, lazy-loaded | Dependency injection, global var | Single file, single writer. DI adds ceremony with no benefit. Lazy-load avoids reading on startup if no module needs settings |
| Module discovery | Scan every menu open | Cache with TTL, inotify watch | Module count is small (~40). No-cache avoids staleness, invalidation bugs. Per contract.md — scan is <10ms |
| Error handling | Hub-centralized; modules emit JSON only | Modules write to TUI directly (control codes) | Separation of concerns — modules are testable without a TUI; hub owns all rendering decisions |
| Atomic writes | os.CreateTemp(dir) + os.Rename | Write directly, fsync, flock | POSIX guarantees rename atomicity on same fs. Temp file in same dir ensures same-fs operation |
| Logging | `log` package with file writer | logrus, zap, zerolog | stdlib only — no external deps. Structured format (timestamp, module, action, exit_code) is sufficient for audit trail |
| Test location | `*_test.go` alongside source + `test/` for integration | Only `*_test.go` or only top-level `test/` | Go convention for unit tests; separate `test/fixtures/` needed for mock module directories |
| Repo signing | GPG batch keygen during ISO build | Pre-generated key in repo | Per-ISO keys avoid key lifecycle management; `SigLevel = Required` enforces verification |

## Data Flow

### core/01 — Module lifecycle
```
Startup → Discovery.Scan(system + user paths)
  → manifest.Parse() → validate required fields → filter invalid
  → Hub.BuildMenu() → group by category, sort by name
  → TUI.Update(menu) → View renders

User selects module:
  → Hub.CheckDeps(pacman -Q) → block if missing
  → Hub.CheckRoot(sudo) → escalate if required
  → Executor.Run(module bin, env vars) → exec.Command
  → Parse stdout (JSON) + stderr (text)
  → Log to /var/log/lambda-env/modules.log
  → If settings_delta: Settings.SaveDelta() → atomic write
  → TUI.renderResponse(status, data, message)
```

### core/02 — Settings lifecycle
```
Hub.Startup → Store.Load()
  → if file missing → return Defaults()
  → if file exists → json.Decode(Settings{})
  → Migrate(loaded.Version, currentVersion) → add missing fields
  → Validate() → check types, enums

Module response has settings_delta:
  → Store.SaveDelta(delta map[string]interface{})
  → Deep merge delta into current Settings
  → Validate() merged result
  → Atomic write: CreateTemp → Encode → Close → Rename
```

### infra/01 — Repo lifecycle
```
ISO build → mkdir -p /srv/repo/lambdaos/x86_64/
  → gpg --batch --gen-key (LambdaOS signing key)
  → pacman-key --add + --lsign-key
  → pacman.conf: [lambdaos] Server = file:///srv/repo/lambdaos/$arch

Post-build: scripts/repo-update.sh
  → repo-add --sign lambdaos.db.tar x86_64/*.pkg.tar.zst
  → generates .db.tar + .db.tar.sig + .files.tar
```

## Error Handling Strategy

Three-tier approach:

| Tier | Where | What |
|------|-------|------|
| Module | JSON on stdout | `status: "error"` / `"warning"` with code, message, suggestion |
| Hub | execution.go | Parses exit codes (0=ok, 1=error, 2=warning), extracts message, renders TUI error view |
| System | modules.log | Full context: timestamp, module name, action, exit code, raw stdout, raw stderr, env vars |

Modules NEVER write errors to the TUI — they always emit structured JSON. The hub decides how to render. Non-JSON stdout is treated as a parse error and logged with raw content.

## Testing Strategy

| Layer | Scope | Approach | Location |
|-------|-------|----------|----------|
| Unit | Manifest validation | Table-driven: valid/invalid manifests, all edge cases | `pkg/module/manifest_test.go` |
| Unit | Settings store | Atomic write failure simulation, delta merge, migration paths | `internal/settings/store_test.go` |
| Unit | TUI model | bubbletea `Model.Update(msg)` without terminal — test state transitions | `internal/tui/model_test.go` |
| Integration | Hub execution | Mock module scripts emit known JSON; verify hub parsing + logging | `test/hub_integration_test.go` |
| CI | Go build + lint | `go build ./...`, `go vet ./...`, `golangci-lint` | GitHub Actions |

Mock modules in `test/fixtures/modules/` provide controlled stdout/stderr for integration tests.

## Branching Strategy

GitFlow adapted for SDD waves:
```
main ─────  develop ── wave-1 ── wave-1/core-01-hub
                                      ├── wave-1/core-02-settings
                                      └── wave-1/infra-01-pacman
```

Each spec maps to a feature branch off `wave-1`. PRs merge feat → wave-1 → develop → main.

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `go` | 1.21+ | Already in packages.x86_64:38 |
| `charmbracelet/bubbletea` | latest | TUI framework (Elm architecture) |
| `charmbracelet/lipgloss` | latest | Terminal styling |
| `charmbracelet/bubbles` | latest | List/input widgets |
| `pacman` | system | Dependency checking (`pacman -Q`) |
| `sudo` | system | Root escalation |
| `gnupg` | system | GPG key generation for repo signing |
| `repo-add` | system | Pacman database tool (included with pacman) |

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Go binary ~15MB | Low — ISO is 800MB+ | Acceptable; static binary with zero runtime deps |
| bubbletea learning curve | Medium — team unfamiliar with Elm pattern | Start with simple list model; bubbletea tutorials available |
| Atomic write on overlay fs | Low — tmp+Rename on same fs | Ensure temp file in same directory; test on overlayfs in CI |
| GPG key per ISO build | Low — keys not persistent across rebuilds | Package signing is local-only in Wave 1; remote repo in Wave 2 |
| CI build time increase | Low — Go build ~10s | Cache Go modules in CI; run Go jobs in parallel with Python |