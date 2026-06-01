# Design: Wave 2 — TUI Modules

## Technical Approach

Four new Go module binaries extend the Wave 1 hub infrastructure. Each module is a standalone Go binary placed in its own `internal/modules/<name>/` package, built as a separate binary, and discovered by the hub via `manifest.json` + JSON-over-stdout protocol. Modules read settings through the shared `internal/settings` package, perform system actions (regenerate configs, run stow, reload Qtile), and emit `settings_delta` back to the hub for atomic merge into `settings.json`.

Config regeneration uses Go `text/template` for Lua (lazy.lua) and `text/template` for Python (config.py). Templates live in `internal/modules/<name>/templates/` as embedded strings — no external template files at runtime.

## Architecture Decisions

### Decision: Module binary structure

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Single binary with subcommands | Simpler build, but violates hub plugin contract | Rejected |
| Separate binary per module | Matches hub discovery, isolates failures, clean JSON protocol | **Chosen** |

Each module binary has a `main.go` that: reads env vars → loads settings via shared package → executes action → emits JSON to stdout.

### Decision: Template embedding

| Option | Tradeoff | Decision |
|--------|----------|----------|
| External `.tmpl` files on disk | Editable post-install, but fragile paths | Rejected |
| `go:embed` string constants | Single binary, no runtime deps, versioned with code | **Chosen** |

### Decision: Qtile reload safety

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Reload blindly, hope it works | Simple, but can lock user out | Rejected |
| `python -m py_compile` + backup + reload | Validates before reload, rollback on failure | **Chosen** |

### Decision: Neovim tui_bridge.lua migration

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Keep reading `tui_settings.json` alongside `settings.json` | Two sources of truth, migration complexity | Rejected |
| Rewrite `tui_bridge.lua` to read `settings.json` directly | Single source, but Lua JSON parsing needed | **Chosen** |

Lua's `vim.json.decode` already exists in tui_bridge.lua — just change the file path from `~/.config/nvim/tui_settings.json` to `~/.config/lambdaos/settings.json` and read the `neovim` sub-object.

### Decision: Stow conflict detection

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Let stow fail on conflicts | Simple, but error messages are opaque | Rejected |
| Pre-scan target files, compare checksums, report conflicts before stow | Clear UX, safe by default | **Chosen** |

## Data Flow

```
TUI (Bubbletea) ──→ User selects module action
                       │
                       ▼
                  Hub.ExecuteModule()
                       │
                       ├── env: LAMBDA_ENV_SETTINGS → settings.json path
                       ├── env: LAMBDA_ENV_ACTION → action name
                       │
                       ▼
              Module binary (stdout JSON)
                       │
                       ├── Reads settings.json via settings.Load()
                       ├── Performs system action (regen config, stow, reload)
                       └── Emits Response{status, action, data, settings_delta}
                       │
                       ▼
                  Hub merges settings_delta
                       │
                       ▼
              settings.SaveDelta() → atomic write settings.json
```

## Module Structure

### Neovim Module

**Package**: `src/lambda-env/internal/modules/neovim/`

| File | Description |
|------|-------------|
| `main.go` | Entry point: reads env, dispatches action, emits JSON |
| `config.go` | Reads settings, generates lazy.lua from template |
| `bridge.go` | Updates tui_bridge.lua to read from settings.json |
| `templates/lazy_lua.go` | Go template string for lazy.lua |

**Key functions**:
- `GenerateLazyLua(settings.NeovimSettings) (string, error)` — renders template
- `UpdateTuiBridgeLua(nvimConfigPath string) error` — rewrites tui_bridge.lua to read `~/.config/lambdaos/settings.json` instead of `~/.config/nvim/tui_settings.json`
- `Apply(settingsPath string) error` — orchestrates: load settings → regen lazy.lua → update bridge → emit delta

**lazy.lua template approach**: The template renders the existing lazy.lua structure with `tui_flags` replaced by settings-driven values. Three toggle sections (LSP, Copilot, Neo-tree) use `{{if .EnableLSP}}...{{end}}` blocks. The template output replaces `lua/core/lazy.lua` in the dotfiles directory.

**tui_bridge.lua update strategy**: The module rewrites `lua/core/tui_bridge.lua` to:
1. Read `~/.config/lambdaos/settings.json` (path from `LAMBDA_ENV_SETTINGS` env or default)
2. Parse JSON with `vim.json.decode`
3. Extract `neovim` section fields: `enable_lsp`, `enable_copilot`, `enable_neotree`
4. Return flags map (same interface as current `get_flags()`)

**Settings extension**: Add `EnableLSP`, `EnableCopilot`, `EnableNeotree` bool fields to `NeovimSettings` in `schema.go`.

### Qtile Module

**Package**: `src/lambda-env/internal/modules/qtile/`

| File | Description |
|------|-------------|
| `main.go` | Entry point: reads env, dispatches action, emits JSON |
| `config.go` | Reads settings, generates config.py from template |
| `reload.go` | Validates and reloads Qtile safely |
| `templates/config_py.go` | Go template string for config.py |

**Key functions**:
- `GenerateConfigPy(settings.QtileSettings) (string, error)` — renders config.py template
- `ValidateConfigPy(path string) error` — runs `python -m py_compile <path>` to validate syntax
- `ReloadQtile() error` — runs `qtile cmd-obj -o cmd -f reload`
- `Apply(settingsPath string) error` — orchestrates: load settings → backup config.py → regen → validate → reload → emit delta

**config.py template approach**: The template renders the existing `config.py` structure. Template variables: `.Terminal`, `.Browser`, `.BarPosition`, `.BarSize`, `.Layouts`. The template imports remain static (`groups`, `keys`, `screens`, `theme`), only the configurable values are templated.

**Qtile reload safety**:
1. Backup current `config.py` → `config.py.bak`
2. Write generated `config.py`
3. Run `python -m py_compile config.py` — if fails, restore backup, return error
4. Run `qtile cmd-obj -o cmd -f reload` — if fails, restore backup, return error

**Settings extension**: Add `Terminal`, `Browser` string fields to `QtileSettings` in `schema.go`.

### Dotfiles Module

**Package**: `src/lambda-env/internal/modules/dotfiles/`

| File | Description |
|------|-------------|
| `main.go` | Entry point: reads env, dispatches action, emits JSON |
| `stow.go` | Wraps `stow` CLI with conflict detection |
| `conflicts.go` | Pre-scan and checksum-based conflict detection |
| `backup.go` | Backup current configs before stow operations |

**Key functions**:
- `ListModules(dotfilesDir string) ([]StowModule, error)` — lists directories in dotfiles repo
- `DetectConflicts(dotfilesDir, targetDir, module string) ([]Conflict, error)` — scans what stow would link, checks if target files exist with different checksums
- `StowModule(dotfilesDir, targetDir, module string) error` — runs `stow -t <target> <module>` after conflict check
- `UnstowModule(dotfilesDir, targetDir, module string) error` — runs `stow -D -t <target> <module>`
- `BackupConfig(sourcePath, backupDir string) (string, error)` — copies file to backup dir with timestamp

**Conflict detection algorithm**:
1. Walk module directory tree, compute relative paths
2. For each relative path, check if `~/.config/<path>` exists
3. If exists, compute SHA256 of both source and target
4. If checksums differ → conflict (user modified the file)
5. Return list of conflicts with paths and checksums

**Backup strategy**: Before any stow operation that has conflicts, copy conflicting files to `~/.config/lambdaos/backups/<timestamp>/<module>/`. Backup path is returned in the module response data.

### PKGBUILD

**Location**: `packages/lambdaos-tui/`

| File | Description |
|------|-------------|
| `PKGBUILD` | Main build script |
| `lambdaos-tui.install` | Post-install hooks |

**Build process**:
1. Build hub binary: `go build -o lambda-env ./cmd/lambda-env`
2. Build each module binary: `go build -o <module-name> ./internal/modules/<module>/`
3. Install hub to `/usr/bin/lambda-env`
4. Install modules to `/usr/share/lambda-env/modules/<name>/` (each with `manifest.json` + `module` binary)
5. Install default settings to `/etc/lambdaos/settings.json`
6. Create `/var/log/lambda-env/` directory

**Install paths**:
```
/usr/bin/lambda-env                          (hub binary)
/usr/share/lambda-env/modules/neovim/manifest.json
/usr/share/lambda-env/modules/neovim/module  (neovim module binary)
/usr/share/lambda-env/modules/qtile/manifest.json
/usr/share/lambda-env/modules/qtile/module   (qtile module binary)
/usr/share/lambda-env/modules/dotfiles/manifest.json
/usr/share/lambda-env/modules/dotfiles/module (dotfiles module binary)
/etc/lambdaos/settings.json                  (default settings)
/var/log/lambda-env/                         (log directory)
```

**Post-install hooks** (`lambdaos-tui.install`):
- `post_install`: Copy `/etc/lambdaos/settings.json` to `~/.config/lambdaos/settings.json` if not exists (per-user)
- `post_upgrade`: Migrate settings if schema version changed

**Clean chroot compatibility**: PKGBUILD uses only `go` and `glibc` as build/runtime deps. No network access needed during build (all deps vendored via `go mod vendor` or `go mod download` in `prepare()`).

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `src/lambda-env/internal/modules/neovim/main.go` | Create | Neovim module entry point |
| `src/lambda-env/internal/modules/neovim/config.go` | Create | lazy.lua generation |
| `src/lambda-env/internal/modules/neovim/bridge.go` | Create | tui_bridge.lua update |
| `src/lambda-env/internal/modules/neovim/templates/lazy_lua.go` | Create | Embedded lazy.lua template |
| `src/lambda-env/internal/modules/qtile/main.go` | Create | Qtile module entry point |
| `src/lambda-env/internal/modules/qtile/config.go` | Create | config.py generation |
| `src/lambda-env/internal/modules/qtile/reload.go` | Create | Qtile reload with validation |
| `src/lambda-env/internal/modules/qtile/templates/config_py.go` | Create | Embedded config.py template |
| `src/lambda-env/internal/modules/dotfiles/main.go` | Create | Dotfiles module entry point |
| `src/lambda-env/internal/modules/dotfiles/stow.go` | Create | Stow wrapper |
| `src/lambda-env/internal/modules/dotfiles/conflicts.go` | Create | Conflict detection |
| `src/lambda-env/internal/modules/dotfiles/backup.go` | Create | Config backup |
| `src/lambda-env/internal/settings/schema.go` | Modify | Add EnableLSP, EnableCopilot, EnableNeotree to NeovimSettings; add Terminal, Browser to QtileSettings |
| `src/lambda-env/cmd/lambda-env/main.go` | Modify | No changes needed — hub discovers modules automatically |
| `packages/lambdaos-tui/PKGBUILD` | Modify | Add module build steps, install paths |
| `packages/lambdaos-tui/lambdaos-tui.install` | Create | Post-install hooks |
| `airootfs/etc/skel/dotfiles/nvim/.config/nvim/lua/core/tui_bridge.lua` | Modify | Read from settings.json instead of tui_settings.json |
| `.github/workflows/ci.yml` | Modify | Add test jobs for new modules |

## Interfaces / Contracts

### NeovimSettings (extended)

```go
type NeovimSettings struct {
    Theme        string `json:"theme"`
    Font         string `json:"font"`
    Lines        int    `json:"lines"`
    Columns      int    `json:"columns"`
    EnableLSP    bool   `json:"enable_lsp"`
    EnableCopilot bool  `json:"enable_copilot"`
    EnableNeotree bool  `json:"enable_neotree"`
}
```

### QtileSettings (extended)

```go
type QtileSettings struct {
    BarPosition string   `json:"bar_position"`
    BarSize     int      `json:"bar_size"`
    Layouts     []string `json:"layouts"`
    Terminal    string   `json:"terminal"`
    Browser     string   `json:"browser"`
}
```

### Module Response (existing, reused)

```go
type Response struct {
    Status        string                 `json:"status"`
    Action        string                 `json:"action"`
    Data          map[string]interface{} `json:"data,omitempty"`
    Message       string                 `json:"message,omitempty"`
    SettingsDelta map[string]interface{} `json:"settings_delta,omitempty"`
}
```

### lazy.lua Go template (sketch)

```go
const lazyLuaTemplate = `local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
{{/* lazy.nvim bootstrap — always present */}}
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({"git", "clone", "--filter=blob:none", "https://github.com/folke/lazy.nvim.git", "--branch=stable", lazypath})
end
vim.opt.rtp:prepend(lazypath)

require("lazy").setup({
  spec = {
    { import = "plugins" },
{{if .EnableLSP}}    { import = "plugins.lsp" },{{end}}
{{if .EnableCopilot}}    { import = "plugins.ai" },{{end}}
{{if .EnableNeotree}}    { "nvim-tree/nvim-tree.lua", config = function() require("nvim-tree").setup() end },{{end}}
  },
  defaults = { lazy = false, version = false },
  install = { colorscheme = { vim.g.nvim_theme or "catppuccin" } },
  checker = { enabled = true, notify = false },
  change_detection = { notify = false },
})`
```

### config.py Go template (sketch)

```go
const configPyTemplate = `import os
import subprocess
from pathlib import Path

from groups import groups
from keys import keys
from libqtile import hook, qtile
from libqtile.config import Key, Screen
from libqtile.layout import Columns, Max, MonadTall
from theme import load_theme

colors = load_theme()

terminal = "{{.Terminal}}"
browser = "{{.Browser}}"

layouts = [
{{range .Layouts}}    {{.}},
{{end}}]

# ... rest of config from existing config.py, with terminal/browser substituted ...

@hook.subscribe.startup_once
def autostart():
    home = Path.home()
    dotfiles_dir = home / "dotfiles"
    if dotfiles_dir.is_dir():
        subprocess.Popen(["stow", "*/"], cwd=dotfiles_dir)
    subprocess.Popen(["xsetroot", "-solid", colors["bg"]])
    subprocess.Popen(["picom", "--experimental-backends"])
    subprocess.Popen(["flameshot"])`
```

## Error Handling

| Level | Error | Handling |
|-------|-------|----------|
| Module binary | Settings file not found | Return `status: "error"`, message with path |
| Module binary | Template render fails | Return `status: "error"`, include template error |
| Module binary | External command fails (stow, qtile reload) | Return `status: "error"`, capture stderr in message |
| Hub | Module executable not found | Skip module, log warning |
| Hub | Module JSON parse fails | Return error, log raw output |
| Hub | Settings delta merge fails | Return error, settings unchanged (atomic write) |
| Qtile module | py_compile fails | Restore backup, return error |
| Qtile module | reload fails | Restore backup, return error |
| Dotfiles module | Conflicts detected | Return `status: "ok"`, list conflicts in `data`, do NOT stow |
| Dotfiles module | Stow fails | Return error, no partial state (stow is atomic per module) |

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `GenerateLazyLua()` output correctness | Golden file tests — compare rendered template against expected `.lua` output |
| Unit | `GenerateConfigPy()` output correctness | Golden file tests — compare rendered template against expected `.py` output |
| Unit | `ValidateConfigPy()` with valid/invalid Python | Write temp `.py` files, run validation, assert pass/fail |
| Unit | `DetectConflicts()` with matching/differing checksums | Create temp dir trees with known file contents, assert conflict detection |
| Unit | Settings schema defaults for new fields | Assert `Defaults()` returns correct bool/string values |
| Integration | Module binary emits valid JSON | Build module, run with test env, parse stdout JSON |
| Integration | Hub discovers and executes modules | Use `setupTestHome` pattern from existing tests |
| Integration | Settings delta merge after module execution | Existing `TestSettingsDeltaMerge` pattern extended |

**Mock strategy**: No external mocking needed — use temp directories for filesystem operations, shell scripts for external command simulation (stow, qtile reload).

## Migration Notes

**tui_settings.json → settings.json**: The Neovim module updates `tui_bridge.lua` to read from `~/.config/lambdaos/settings.json`. The old `~/.config/nvim/tui_settings.json` is no longer read. No migration of existing `tui_settings.json` values is needed — the TUI will write the correct values to `settings.json` on first use. The old file is left in place (harmless).

**os_theme.json**: Unchanged — still read by `theme.py` (Qtile) and `env.lua` (Neovim) for theme colors. Out of scope for Wave 2.

**Existing lazy.lua**: Replaced by template-generated version. The template preserves the existing structure — only the conditional import blocks change.

**Existing config.py**: Replaced by template-generated version. The template preserves the existing structure — only `terminal`, `browser`, and layout values are parameterized.

## Risks and Mitigations

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Qtile reload locks user out of WM | Medium | py_compile validation + backup + restore on failure |
| Stow overwrites user-modified configs | Medium | Pre-scan conflict detection with checksums, never stow without explicit confirmation |
| Template produces invalid Lua/Python | Medium | Golden file tests + `lua -l` validation (optional) + `python -m py_compile` (mandatory) |
| PKGBUILD fails in clean chroot | Low | Test with `extra-x86_64-build`; vendor Go deps or use `go mod download` in `prepare()` |
| Settings schema drift between modules | Low | Single shared `internal/settings` package; all modules import it |
| Module binary size bloats package | Low | `CGO_ENABLED=0` static builds; each module is a small Go binary (~5-10MB) |
| tui_bridge.lua rewrite breaks existing nvim | Low | Backup tui_bridge.lua before rewrite; template preserves same `get_flags()` interface |

## Open Questions

- [ ] Should the Neovim module also manage the `plugins/` directory imports, or only toggle the three flags (LSP, Copilot, Neo-tree)?
- [ ] Should the Qtile template also parameterize `keys.py` (terminal variable) or only `config.py`?
- [ ] Should the PKGBUILD include a `.SRCINFO` file for AUR compatibility, or is this repo-only?
