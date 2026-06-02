# Tasks: Wave 3 — TUI Interface + System Modules

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 2,400–2,800 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 → PR 2 → PR 3a → PR 3b → PR 4a → PR 4b → PR 5 |
| Delivery strategy | ask-on-risk |
| Chain strategy | feature-branch-chain |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Foundation: schema v1.1.0, manifest actions, CLIExecutor, deps | PR #1 | base = feature/tracker branch |
| 2 | TUI Components: toggle, textinput, confirm, help, statusbar | PR #2 | base = PR #1 |
| 3 | TUI Sub-models: model/update/view refactor, categories/modules views | PR #3a | base = PR #2; ~450 lines |
| 4 | Module Detail + ExecuteAction: detail.go with dynamic widgets | PR #3b | base = PR #3a; fits budget |
| 5 | System Modules Part 1: keyboard + appearance + theme sync | PR #4a | base = PR #3b; ~500 lines |
| 6 | System Modules Part 2: audio + defaults | PR #4b | base = PR #4a; ~460 lines |
| 7 | Integration + E2E Tests | PR #5 | base = PR #4b; fits budget |

## Phase 1: Foundation (Infrastructure)

- [x] 1.1 Add `go get github.com/charmbracelet/bubbles` to go.mod (+1)
- [x] 1.2 Add `ActionConfig` struct + `Actions []ActionConfig` to `pkg/module/manifest.go` with type validation
- [x] 1.3 Bump `CurrentVersion` to "1.1.0" in `internal/settings/schema.go`; add 7 section structs (Power, Defaults, Autostart, Updates, Security, Fonts, Notifications)
- [x] 1.4 Add `UseGlobalTheme bool` to NeovimSettings + QtileSettings; default `true` in `Defaults()`
- [x] 1.5 Create `pkg/module/executor.go` — CLIExecutor interface + RealExecutor + MockExecutor; returns `(stdout, stderr, exitCode int, err error)`
- [x] 1.6 Tests: schema v1.0.0→1.1.0 migration preserves existing values; manifest action parsing validates types

## Phase 2: TUI Components (Reusable Widgets)

- [x] 2.1 Create `internal/tui/components/toggle.go` — boolean toggle; renders ✓/○; emits `ToggleChangedMsg`
- [x] 2.2 Create `internal/tui/components/textinput.go` — wraps bubbles/textinput; pointer receiver; allowlist/regex/numeric validation
- [x] 2.3 Create `internal/tui/components/confirm.go` — Yes/No dialog; emits `ConfirmResultMsg`
- [x] 2.4 Create `internal/tui/components/help.go` — custom full-screen overlay (NOT bubbles/help); context-sensitive bindings
- [x] 2.5 Create `internal/tui/components/statusbar.go` — persistent bar (context, module, settings state, modified indicator)
- [x] 2.6 Tests: each component with keypress sequences; verify state transitions, disabled state, boundary behavior

## Phase 3a: TUI Sub-model Architecture

- [x] 3a.1 Expand viewState: add `viewModuleDetail`, `viewConfirmDialog` to `internal/tui/model.go`; add sub-model fields + SubModel interface
- [x] 3a.2 Extract `internal/tui/views/categories.go` — category list sub-model from existing flat logic
- [x] 3a.3 Extract `internal/tui/views/modules.go` — module list sub-model with preserved selection state
- [x] 3a.4 Rewrite `internal/tui/update.go` — parent delegates to active sub-model; route execMsg to status bar component
- [x] 3a.5 Rewrite `internal/tui/view.go` — parent delegates rendering; status bar component rendered by parent across all views
- [x] 3a.6 Tests: sub-model navigation (categories→modules→back), selection preservation, empty states
- [x] 3a.7 **Debt fix**: Update design.md to match implementation conventions (`ActionConfig`, `pkg/module/executor.go`, pointer receivers, custom help overlay)

## Phase 3b: Module Detail View + ExecuteAction

- [x] 3b.1 Create `internal/tui/views/detail.go` — reads manifest.actions, renders widgets per type (toggle/select/text/confirm/execute)
- [x] 3b.2 Add `ExecuteAction(name, action, params)` to `internal/hub/execution.go` — passes LAMBDA_ENV_ACTION env var; backward-compatible with existing `run` action
- [x] 3b.3 Wire detail widget actions → hub.ExecuteAction → execMsg → widget state update + status bar
- [x] 3b.4 **Option C**: Detail view executes module "run" on enter, merges `data.available_options` into widget options (static manifest options → dynamic system options)
- [x] 3b.5 Tests: detail view renders all 5 action types; ExecuteAction sets correct env; widget state syncs with response; dynamic options merge

## Phase 4a: Keyboard + Appearance Modules

- [x] 4a.1 Create `internal/modules/keyboard/main.go` — setxkbmap layout/variant/options via CLIExecutor; emit `data.available_layouts` from `setxkbmap -layout` query
- [x] 4a.2 Create `internal/modules/appearance/main.go` — gsettings theme/wallpaper, feh wallpaper, icon/cursor/font via CLIExecutor; emit `data.available_themes` from filesystem scan
- [x] 4a.3 Add theme sync: `use_global_theme` mapping in appearance module; themeMap lookup table (dark→tokyonight, light→tokyonight-light, nord→nord, catppuccin→catppuccin-mocha)
- [x] 4a.4 Add `use_global_theme` support to `internal/modules/neovim/config.go` — map appearance.theme → neovim theme when `use_global_theme: true`
- [x] 4a.5 Add `use_global_theme` support to `internal/modules/qtile/config.go` — map appearance.theme → qtile color scheme when `use_global_theme: true`
- [x] 4a.6 Tests: keyboard layout apply, invalid layout error, theme mapping for all 4 presets, use_global_theme toggle, dynamic options discovery

## Phase 4b: Audio + Defaults Modules + Debt Fix

- [x] 4b.1 Create `internal/modules/audio/main.go` — pactl/wpctl volume 0–100, 5% steps, mute toggle, sink selection; emit `data.available_sinks` from `pactl list short sinks`; detect pipewire vs pulseaudio
- [x] 4b.2 Create `internal/modules/defaults/main.go` — xdg-mime browser/terminal/editor/file-manager; emit `data.available_apps` from `/usr/share/applications/` scan; batch apply with confirm
- [x] 4b.3 Tests: audio backend detection (pipewire/pulse), volume capping, mute state, dynamic sink discovery; defaults desktop file validation, batch partial failure, dynamic app discovery
- [x] 4b.4 **Debt fix**: Fix hyphen mismatch — update neovim/qtile/dotfiles `main.go` to accept hyphenated action names (e.g., `toggle-lsp` not `toggle_lsp`), matching their `manifest.json` actions. New standard: hyphens everywhere.

## Phase 5: Integration + E2E Tests

- [x] 5.1 Integration tests: settings_delta flow — write settings, execute module action via MockExecutor, verify delta merge
- [x] 5.2 Integration tests: manifest action parsing — parse manifest with all 5 action types, assert widget rendering
- [x] 5.3 E2E tests: TUI navigation — programmatic tea.NewProgram key sequences through all 3 levels + confirm dialog
