# Tasks: Wave 2 — TUI Modules

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~950 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 (schema) → PR 2 (neovim) → PR 3 (qtile) → PR 4 (dotfiles) → PR 5 (package+CI) |
| Delivery strategy | ask-on-risk |
| Chain strategy | feature-branch-chain |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Settings schema extensions (neovim/qtile fields + defaults + validation) | PR 1 | Base = feature/wave-2-tui-modules |
| 2 | Neovim module: toggles, lazy.lua template, tui_bridge update + tests | PR 2 | Base = PR 1 branch |
| 3 | Qtile module: defaults, config.py template, reload safety + tests | PR 3 | Base = PR 2 branch |
| 4 | Dotfiles module: stow/unstow, conflict detection, backup + tests | PR 4 | Base = PR 3 branch |
| 5 | PKGBUILD, .install, default settings, CI updates | PR 5 | Base = PR 4 branch; merges to main |

## Phase 1: Settings Schema Extensions

- [x] 1.1 Add `EnableLSP`, `EnableCopilot`, `EnableNeotree` (bool), `LspServers` ([]string) to `NeovimSettings` in `src/lambda-env/internal/settings/schema.go`; add defaults `true, true, true, ["gopls","pyright"]` in `Defaults()`
- [x] 1.2 Add `Terminal`, `Browser` (string), `DefaultFileManager` (string), `Groups` ([]GroupConfig) to `QtileSettings` in `src/lambda-env/internal/settings/schema.go`; add defaults `"kitty", "firefox", "thunar", [{name:"1"}..{name:"9"}]` in `Defaults()`
- [x] 1.3 Add validation in `schema.go Validate()`: reject unknown terminal names (allowlist: kitty, foot, alacritty, st, xterm), reject unknown browser names (allowlist: firefox, chromium, brave, chrome); add `TestValidateNeovimDefaults` and `TestValidateQtileDefaults` in `store_test.go`

## Phase 2: Neovim Module

- [x] 2.1 Create `src/lambda-env/internal/modules/neovim/main.go`: read `LAMBDA_ENV_ACTION`, `LAMBDA_ENV_SETTINGS` env vars; dispatch to `toggle_lsp`, `toggle_copilot`, `toggle_neotree`, `run` actions; emit `module.Response` JSON to stdout
- [x] 2.2 Create `src/lambda-env/internal/modules/neovim/config.go`: `GenerateLazyLua(settings.NeovimSettings) (string, error)` using embedded Go template with `{{if .EnableLSP}}` conditional blocks; backup existing `lazy.lua` to `.bak` before write; validate non-empty output
- [x] 2.3 Create `src/lambda-env/internal/modules/neovim/templates/lazy_lua.go`: embedded `const lazyLuaTemplate` string rendering lazy.nvim bootstrap + conditional imports for LSP, Copilot, Neo-tree (per design.md sketch)
- [x] 2.4 Create `src/lambda-env/internal/modules/neovim/bridge.go`: `UpdateTuiBridgeLua(nvimConfigPath string) error` rewrites `lua/core/tui_bridge.lua` to read `~/.config/lambdaos/settings.json` via `vim.json.decode` instead of `tui_settings.json`; backup before rewrite
- [x] 2.5 Create `src/lambda-env/internal/modules/neovim/manifest.json`: name `neovim`, version `0.1.0`, category `apps`, deps `["neovim"]`, min_hub_version `1.0.0`
- [x] 2.6 Write `src/lambda-env/internal/modules/neovim/config_test.go`: golden file tests for `GenerateLazyLua()` with all toggles on, LSP off, Copilot off; assert template output contains/omits expected plugin entries

## Phase 3: Qtile Module

- [x] 3.1 Create `src/lambda-env/internal/modules/qtile/main.go`: read env vars; dispatch to `set_terminal`, `set_browser`, `set_file_manager`, `reload`, `run` actions; validate input against allowlists; emit `module.Response` JSON
- [x] 3.2 Create `src/lambda-env/internal/modules/qtile/config.go`: `GenerateConfigPy(settings.QtileSettings) (string, error)` using embedded Go template with `{{.Terminal}}`, `{{.Browser}}`, `{{range .Layouts}}` variables; backup `config.py.bak` before write
- [x] 3.3 Create `src/lambda-env/internal/modules/qtile/templates/config_py.go`: embedded `const configPyTemplate` string rendering config.py with parameterized terminal, browser, layouts (per design.md sketch)
- [x] 3.4 Create `src/lambda-env/internal/modules/qtile/reload.go`: `ValidateConfigPy(path string) error` runs `python3 -m py_compile`; `ReloadQtile() error` runs `qtile cmd-obj -o cmd -f reload_config`; `SafeApply()` orchestrates: backup → generate → validate → reload → restore on failure
- [x] 3.5 Create `src/lambda-env/internal/modules/qtile/manifest.json`: name `qtile`, version `0.1.0`, category `apps`, deps `["qtile"]`, min_hub_version `1.0.0`
- [x] 3.6 Write `src/lambda-env/internal/modules/qtile/config_test.go`: golden file tests for `GenerateConfigPy()` with custom terminal/browser; `TestValidateConfigPy` with valid/invalid Python temp files

## Phase 4: Dotfiles Module

- [x] 4.1 Create `src/lambda-env/internal/modules/dotfiles/main.go`: read env vars; dispatch to `list`, `stow`, `unstow`, `check_conflicts`, `backup`, `run` actions; emit `module.Response` JSON
- [x] 4.2 Create `src/lambda-env/internal/modules/dotfiles/stow.go`: `ListModules(dir string) ([]StowModule, error)` lists directories; `StowModule(dir, target, name string) error` runs `stow -t <target> <module>`; `UnstowModule(...)` runs `stow -D -t <target> <module>`
- [x] 4.3 Create `src/lambda-env/internal/modules/dotfiles/conflicts.go`: `DetectConflicts(dir, target, module string) ([]Conflict, error)` walks module tree, computes SHA-256 of source vs target files, returns mismatches with paths and checksums
- [x] 4.4 Create `src/lambda-env/internal/modules/dotfiles/backup.go`: `BackupConfig(sourceDir, backupDir, module string) (string, error)` copies files to `~/.config/lambdaos/backups/<timestamp>/<module>/` preserving structure; skips identical files
- [x] 4.5 Create `src/lambda-env/internal/modules/dotfiles/manifest.json`: name `dotfiles`, version `0.1.0`, category `ops`, deps `["stow"]`, min_hub_version `1.0.0`
- [x] 4.6 Write `src/lambda-env/internal/modules/dotfiles/stow_test.go`: temp dir tests for `ListModules`, `DetectConflicts` (matching/differing checksums), `BackupConfig` (copies only changed files)

## Phase 5: PKGBUILD and CI

- [x] 5.1 Create `packages/lambdaos-tui/PKGBUILD`: pkgname `lambdaos-tui`, pkgver `0.2.0`, pkgrel `1`, arch `x86_64`, depends `go stow qtile neovim`; `build()` compiles hub + 3 module binaries; `package()` installs to `/usr/bin/lambda-env`, `/usr/share/lambda-env/modules/*/`, `/etc/lambdaos/settings.json`
- [x] 5.2 Create `packages/lambdaos-tui/lambdaos-tui.install`: `post_install()` creates `~/.config/lambdaos/` per user, copies default settings if missing; `pre_remove()` preserves user settings.json with warning message
- [x] 5.3 Create `packages/lambdaos-tui/settings.json`: default settings file with Wave 2 schema (all neovim/qtile fields populated with defaults)
- [x] 5.4 Update `.github/workflows/ci.yml`: add `test-go-modules` job that runs `go test ./internal/modules/... -v -cover` in `src/lambda-env`; ensure `build-go` job builds new module packages — VERIFIED: existing jobs already cover all Go packages recursively, no changes needed

## Phase 6: Technical Debt Fixes

- [x] 6.1 **Fix keys.py terminal update in config.go**: `updateKeysPyTerminal` verified + 2 test cases added (TestUpdateKeysPyTerminal, TestUpdateKeysPyTerminalNoExistingTerminal)
- [x] 6.2 **Add EnableNeotree conditional to lazy.lua template**: Conditional block added + 2 tests (TestGenerateLazyLuaNeotreeOff, TestGenerateLazyLuaNeotreeOn)
- [x] 6.3 **Update tui_bridge.lua skeleton to read settings.json**: Updated to read `~/.config/lambdaos/settings.json`, extract `neovim` sub-object, map enable_* flags
- [x] 6.4 **Add Groups field to config.py template**: Groups section added with `{{range .Groups}}Group("{{.Name}}")` + test (TestGenerateConfigPyWithGroups)
- [x] 6.5 **Replace sha256sums=('SKIP') in PKGBUILD**: Changed to git source pattern `source=("${pkgname}::git+file://${PWD}/../..#tag=v${pkgver}")` with SKIP (valid for git sources)
- [x] 6.6 **Generate .SRCINFO for AUR compatibility**: `.SRCINFO` generated via `makepkg --printsrcinfo`
- [x] 6.7 **Add integration tests for module → settings → regeneration flow**: 3 integration tests added (TestIntegrationToggleLsp, TestIntegrationToggleCopilot, TestIntegrationSetTerminal)
