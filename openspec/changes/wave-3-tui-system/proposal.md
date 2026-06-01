# Proposal: Wave 3 — TUI Interface + System Modules

## Intent

Wave 1 delivered the hub infrastructure (Go + Bubbletea TUI, plugin system, settings schema, CI). Wave 2 delivered 4 functional modules (Neovim, Qtile, Dotfiles, Repo Package). Wave 3 closes the interactive gap: the TUI currently only has flat 2-level navigation with immediate execution — no toggles, no forms, no sub-views. This change adds interactive TUI components AND 4 system modules that exercise them, making `lambda-env` a functional system configurator.

## Scope

### In Scope
- **Track A (TUI Foundation)**: Add `bubbles` dependency; 3-level navigation; toggle, text input, confirm dialog, help overlay, status bar components
- **Track B (System Modules)**: system-06-keyboard (setxkbmap), system-09-appearance (theme/wallpaper), system-02-audio (pipewire/pactl), system-10-defaults (xdg-mime)
- Schema sections for 7 missing domains (defaults, autostart, updates, security, fonts, notifications, power) — added upfront for migration readiness
- Manifest `actions` field for module action discovery

### Out of Scope
- Remaining 10 system modules (screen, network, bluetooth, autostart, services, updates, security, fonts, notifications, power) — Wave 4+
- Installer/Calamares — Wave 8+
- Branding/polish — Wave 9+
- Cross-module theme sync (appearance → neovim/qtile) — IN SCOPE: `use_global_theme` flag per module

## Capabilities

### New Capabilities
- `tui-interactive-views`: 3-level navigation, toggle widget, text input, confirm dialog, help overlay, persistent status bar — the TUI interactive component layer
- `system-keyboard-module`: setxkbmap keyboard layout/variant selection
- `system-appearance-module`: theme, wallpaper, icon/cursor theme, font selection
- `system-audio-module`: pipewire/pactl volume, sink selection, mute toggle
- `system-defaults-module`: xdg-mime default app assignments

### Modified Capabilities
- `hub-plugin-system`: Navigation expands from 2 to 3+ levels; manifest gains `actions` field; module execution model supports interactive views
- `settings-schema`: Version bumps 1.0.0 → 1.1.0; adds 7 new sections (power, defaults, autostart, updates, security, fonts, notifications); existing sections get field expansions

## Approach

**Track A**: Rewrite TUI model/view/update into component-based architecture with sub-models per view. `viewState` enum expands to `categories | modules | moduleDetail | confirmDialog`. Each module detail view is built from `manifest.actions` — the TUI reads actions and renders appropriate widgets (toggle for booleans, list for selections, textinput for strings). Uses `bubbles` components (textinput, list, viewport, help, key, spinner) wrapped in LambdaOS composite models.

**Track B**: Each system module follows established pattern (manifest.json + JSON-over-stdout). Modules read settings, execute CLI tools via `os/exec`, and emit `settings_delta`. All 4 modules are self-contained Go binaries under `src/lambda-env/internal/modules/`. Schema additions are batched in one migration to v1.1.0.

## Technical Debt / Open Questions

| # | Question | Decision |
|---|----------|----------|
| 1 | `bubbles` dependency — vendor or `go get`? | **A: `go get`** with version pin to latest stable `v0.20.x`. Vendoring not needed |
| 2 | Schema v1.1.0 — incremental or all-at-once? | **A: All-at-once**: add all 7 sections in one migration. Cleaner than 7 sequential bumps |
| 3 | TUI architecture — rewrite or incremental? | **A: Rewrite** into component-based sub-models. Current 300-line flat model cannot scale |
| 4 | Module action discovery — manifest or hardcoded? | **Manifest `actions` field** — each module declares supported actions with types. TUI renders widgets dynamically |
| 5 | Settings delta vs direct TUI writes? | **A: All through modules**. TUI never writes settings.json directly. Modules own their settings lifecycle |
| 6 | Testing strategy for CLI-tool modules? | **Interface-based mocking**: `CLIExecutor` interface, production impl calls `os/exec`, test impl returns fixtures |
| 7 | Appearance → neovim/qtile theme sync? | **IN SCOPE**: Add `use_global_theme` bool to NeovimSettings/QtileSettings. When true, map `appearance.theme` → module-specific theme. When false, use module's own `theme` field |
| 8 | Root permissions for system modules? | **Per-action pkexec**: modules use `pkexec` for specific actions requiring root, not global sudo |

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `src/lambda-env/internal/tui/` | Modified | Rewrite from flat model to component-based sub-models |
| `src/lambda-env/internal/tui/components/` | New | Toggle, textinput, confirm, help, statusbar components |
| `src/lambda-env/internal/modules/keyboard/` | New | Keyboard layout module |
| `src/lambda-env/internal/modules/appearance/` | New | Theme/wallpaper module |
| `src/lambda-env/internal/modules/audio/` | New | Audio volume module |
| `src/lambda-env/internal/modules/defaults/` | New | Default apps module |
| `src/lambda-env/internal/settings/schema.go` | Modified | Add 7 new sections, bump version to 1.1.0, add `use_global_theme` to Neovim/Qtile settings |
| `src/lambda-env/internal/modules/neovim/` | Modified | Support `use_global_theme` flag, map appearance.theme → neovim theme |
| `src/lambda-env/internal/modules/qtile/` | Modified | Support `use_global_theme` flag, map appearance.theme → qtile theme |
| `src/lambda-env/internal/settings/migrate.go` | Modified | Migration handler for v1.0.0 → v1.1.0 |
| `src/lambda-env/pkg/module/types.go` | Modified | Manifest struct gains `actions` field |
| `src/lambda-env/go.mod` | Modified | Add `charmbracelet/bubbles` dependency |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| TUI rewrite breaks existing Wave 2 flows | Medium | Comprehensive test coverage before rewrite; feature-flag to old model as fallback |
| `bubbles` API breaking changes | Low | Pin to specific minor version; TUI components wrap bubbles with our own interfaces |
| CLI tool output parsing varies across Arch versions | Medium | Parse with regex fallbacks; log unexpected output for diagnosis |
| Schema migration corrupts user settings | Low | Atomic writes (already implemented); dry-run migration validation before commit |
| Audio module: pipewire vs pulseaudio | Medium | Detect backend at runtime (`pactl info`); support both with feature detection |

## Rollback Plan

1. **Package level**: `pacman -R lambdaos-tui` removes all modules; user settings preserved
2. **TUI fallback**: Feature flag to revert to 2-level navigation if component model has issues
3. **Schema migration**: `settings.json` v1.1.0 has all new sections as optional; removing them reverts to v1.0.0 compatible state
4. **Module level**: Each system module is an independent binary — remove or disable individual modules via manifest without affecting others
5. **Config level**: Every module backs up target config before modification (`.bak` suffix)

## Dependencies

- Wave 1 + Wave 2 complete (hub, settings, neovim/qtile/dotfiles modules)
- `github.com/charmbracelet/bubbles` — TUI interactive components
- CLI tools: `setxkbmap`, `feh`, `pactl`/`wpctl`, `xdg-mime` (in packages.x86_64)
- Go 1.24+ (in go.mod)

## Success Criteria

- [ ] `lambda-env` TUI navigates 3 levels (categories → modules → module detail)
- [ ] Toggle widget flips boolean settings and shows state in real-time
- [ ] Text input widget allows typing values (e.g., keyboard variant)
- [ ] Confirm dialog appears before destructive actions
- [ ] Help overlay shows context-specific key bindings
- [ ] Status bar persists across all views showing current module/settings state
- [ ] Keyboard module changes keyboard layout via setxkbmap
- [ ] Appearance module sets theme and wallpaper
- [ ] Appearance → neovim/qtile theme sync works (`use_global_theme` flag)
- [ ] Audio module adjusts volume and toggles mute
- [ ] Defaults module sets default browser/terminal/editor
- [ ] All `go test ./... -v` pass
- [ ] Schema v1.0.0 → v1.1.0 migration preserves existing user settings