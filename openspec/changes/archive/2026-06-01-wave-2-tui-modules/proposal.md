# Proposal: Wave 2 — TUI Modules

## Intent

Wave 1 delivered the hub infrastructure (Go + Bubbletea TUI, plugin system, settings schema, CI, dotfiles skeleton). Wave 2 closes the gap: the TUI cannot actually control the system. This change implements 4 modules that make the TUI functional — Neovim config, Qtile config, dotfiles management, and pacman packaging — so a user can install `lambdaos-tui` and configure their system from the TUI.

## Scope

### In Scope
- **apps-01-neovim**: TUI toggles for LSP/Copilot/Neo-tree; writes to `settings.json` `neovim` section; regenerates `lazy.lua`; updates `tui_bridge.lua` to read unified settings
- **apps-02-qtile**: TUI sets terminal/browser defaults; writes to `settings.json` `qtile` section; regenerates `config.py` from template; reloads Qtile via `qtile cmd-obj`
- **ops-05-dotfiles**: TUI stow/unstow modules; lists stowed/unstowed state; detects file conflicts; backup current configs to dotfiles repo
- **infra-02-repo-package-tui**: PKGBUILD for `lambdaos-tui`; installs binary to `/usr/bin/lambda-env`; installs modules to `/usr/share/lambda-env/modules/`; installs default config to `/etc/lambdaos/`

### Out of Scope
- Remaining system modules (screen, audio, network, bluetooth, etc.) — Wave 3 (14 system modules + TUI interface development)
- TUI interactive views (forms, toggles, sub-navigation) — Wave 3 Track A
- Installer/Calamares modules — Wave 8+
- Branding/polish modules (MOTD, wallpaper, icons) — Wave 9
- Migration of existing `tui_settings.json` to unified schema (deferred, settings.json is source of truth)

## Capabilities

### New Capabilities
- `neovim-module`: TUI module for Neovim config toggles (LSP, Copilot, Neo-tree), lazy.lua regeneration, tui_bridge.lua unified settings reader
- `qtile-module`: TUI module for Qtile defaults (terminal, browser), config.py template regeneration, Qtile reload
- `dotfiles-module`: TUI module for GNU Stow operations (stow/unstow), conflict detection, config backup to dotfiles repo
- `repo-package-tui`: PKGBUILD for lambdaos-tui pacman package with correct install paths and post-install hooks

### Modified Capabilities
- `hub-plugin-system`: Out of scope field updates — module implementations now in scope (was Wave 2+)
- `settings-schema`: Out of scope field updates — tui_settings.json migration deferred (was Wave 2)

## Approach

Each module is a Go binary following the hub plugin contract (manifest.json + JSON-over-stdout protocol). Modules read/write `settings.json` via the existing settings package, then execute system actions (regenerate configs, run stow, reload Qtile). The PKGBUILD packages the hub binary + all 4 modules for pacman distribution.

```
hub (lambda-env) → discovers modules → executes via JSON protocol → modules emit settings_delta → hub merges to settings.json
```

Config regeneration uses Go templates: Neovim `lazy.lua` and Qtile `config.py` are rendered from templates + settings.json values. Dotfiles module wraps `stow` CLI with conflict detection via file checksums.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `src/lambda-env/internal/modules/neovim/` | New | Neovim module: toggles, lazy.lua template, tui_bridge update |
| `src/lambda-env/internal/modules/qtile/` | New | Qtile module: defaults, config.py template, reload |
| `src/lambda-env/internal/modules/dotfiles/` | New | Dotfiles module: stow/unstow, conflict detection, backup |
| `src/lambda-env/internal/settings/` | Modified | Add neovim/qtile section structs if not complete |
| `src/lambda-env/pkg/templates/` | New | Go templates for lazy.lua and config.py |
| `packages/lambdaos-tui/` | New | PKGBUILD, .install file, default settings.json |
| `airootfs/etc/skel/dotfiles/nvim/` | Modified | tui_bridge.lua reads from settings.json |
| `airootfs/etc/skel/dotfiles/qtile/` | Modified | config.py generated from template |
| `.github/workflows/ci.yml` | Modified | Add Go module test jobs for new modules |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Qtile reload fails if config.py has syntax errors | Medium | AST-validate generated config before reload; rollback to previous config on failure |
| Stow conflicts corrupt user configs | Medium | Detect conflicts before stow; offer backup option; never overwrite without confirmation |
| PKGBUILD dependencies mismatch on Arch | Low | Pin to stable packages; test in clean chroot with `extra-x86_64-build` |
| Template rendering produces invalid Lua/Python | Medium | Unit tests for template output; validate with `lua -l` and `python -m py_compile` |
| Settings schema drift between modules | Low | Single settings package shared by all modules; strict typing in Go structs |

## Rollback Plan

1. **Package level**: `pacman -R lambdaos-tui` removes binary and modules; user settings.json preserved in `~/.config/lambdaos/`
2. **Config level**: Each module backs up the target config file before regeneration (e.g., `config.py.bak`); manual restore via `mv config.py.bak config.py`
3. **Dotfiles level**: `stow -D <module>` unstows any module; backup export available before any stow operation
4. **Settings level**: settings.json atomic writes preserve previous version on failure; manual edit always possible (plain JSON)

## Dependencies

- Wave 1 complete: hub plugin system, settings schema, CI pipeline, repo pacman setup
- Go 1.21+ (in packages.x86_64)
- GNU Stow (in packages.x86_64)
- Qtile (in packages.x86_64)
- Neovim + lazy.nvim (in packages.x86_64)

## Success Criteria

- [ ] `lambda-env` TUI shows Neovim module → toggle LSP off → nvim opens without LSP
- [ ] `lambda-env` TUI shows Dotfiles module → stow/unstow works without errors
- [ ] `pacman -S lambdaos-tui` installs successfully on clean Arch system
- [ ] `lambda-env` TUI shows Qtile module → change terminal default → Mod+Enter opens new terminal
- [ ] All Go module tests pass (`go test ./... -v`)
- [ ] CI pipeline passes with new module jobs
- [ ] PKGBUILD builds with `makepkg -s` in clean chroot
