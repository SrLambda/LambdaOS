# Accumulated Learnings — Wave 3: TUI Interface + System Modules

## 1. Hyphen Convention for Action Names

**Problem**: Mixed hyphen/underscore conventions across modules caused action dispatch failures. Neovim module used `toggle-lsp` in manifest but some code paths expected `toggle_lsp`.

**Solution**: Standardized on hyphens everywhere (`toggle-lsp`, `set-layout`, `check-conflicts`). Applied retroactively to Wave 2 modules (neovim, qtile, dotfiles) in PR #4b.

**Why hyphens win**: Manifest JSON keys follow JSON convention (hyphens), Go identifiers use PascalCase/camelCase, the cross-boundary action names are the communication contract between the TUI hub and module executables — hyphens are shell-safe and consistent with CLI flag conventions.

**Enforcement**: All manifest action names are validated by `ActionConfig.Validate()` to reject names with underscores.

## 2. Option C — Dynamic Options Merge Pattern

**Problem**: Module actions with `type: "select"` need option lists. Hardcoding in manifest.json is static and doesn't reflect system state. Querying at every render is expensive.

**Solution**: Hybrid approach (Option C):
- Manifest has static `options` as fallback
- Module can emit dynamic options via `data.available_options` in the `run` response
- TUI detail view merges both: static manifest options are base, dynamic options from module execution are layered on top
- Implementation: `detail.go` calls `executeRun()` on module entry, merges `available_options` into each action's options list

**Tradeoff**: Requires an extra module execution when entering the detail view. Acceptable because module executions are fast (<100ms for CLI queries).

## 3. Component-Based Sub-Model TUI Architecture

**Problem**: Original 300-line flat model (model.go + update.go + view.go) with 2-level navigation and switch-case logic. Adding interactive widgets would have made it unmanageable.

**Solution**: Bubble Tea sub-model delegation pattern:
- Parent `Model` holds a `currentView SubModel` field
- Each view state is a self-contained `tea.Model` (categories, modules, detail, confirm)
- `Update()` delegates to `currentView.Update()`, `View()` delegates to `currentView.View()`
- Widgets are their own sub-models (toggle, textinput, confirm, help, statusbar)
- Status bar is rendered by parent across ALL views (not per-sub-model)

**Key insight**: Widget sub-models need pointer receivers because `bubbles` components like `textinput.Model` mutate internal state on `SetValue()`/`Focus()`/`Blur()`.

**File structure**:
```
internal/tui/
├── model.go           # Parent model with SubModel interface
├── update.go           # Delegation dispatch
├── view.go             # Parent view + status bar
├── components/         # Reusable widgets
│   ├── toggle.go
│   ├── textinput.go
│   ├── confirm.go
│   ├── help.go         # Custom full-screen overlay
│   └── statusbar.go
└── views/              # View-level sub-models
    ├── categories.go
    ├── modules.go
    └── detail.go
```

## 4. Theme Sync via settings_delta Push

**Problem**: When user changes the global theme (appearance module), Neovim and Qtile modules need to reflect that change if `use_global_theme` is enabled.

**Solution**: Push-based sync via existing `settings_delta` mechanism:
1. Appearance module sets `appearance.theme` → delta emitted to hub
2. Hub merges delta into settings.json
3. When user opens Neovim/Qtile detail view, module loads fresh settings
4. If `use_global_theme=true` (default: `true`), module maps `appearance.theme` through a lookup table to its own theme name
5. Theme mapping table: 4 presets (dark→tokyonight, light→tokyonight-light, nord→nord, catppuccin→catppuccin-mocha)

**Note**: This is push-on-access, not push-on-change. Modules read the current global theme when they execute, not when the theme changes. This avoids cross-module coupling.

**Debt**: Theme mapping table is hardcoded (4 entries). Pre-release wave should make it configurable via `/usr/share/lambda-env/themes.json`.

## 5. CLIExecutor Interface for Testability

**Problem**: System modules (keyboard, appearance, audio, defaults) call CLI tools (`setxkbmap`, `feh`, `pactl`, `xdg-mime`). Testing with real CLI calls would depend on system state and be non-deterministic.

**Solution**: `CLIExecutor` interface in `pkg/module/executor.go`:
```go
type CLIExecutor interface {
    Run(command string, args ...string) (stdout, stderr string, exitCode int, err error)
}
```

- `RealExecutor`: wraps `os/exec.Command`
- `MockExecutor`: returns pre-configured fixtures (stdout, stderr, exitCode)
- Modules receive the executor via dependency injection (constructor parameter)

**Impact**: All 4 system modules are fully testable without system dependencies. 292 tests pass with zero real CLI calls.

## 6. Manifest Action Type Safety

**Problem**: Action types (toggle, select, text, confirm, execute) need runtime validation to prevent malformed manifests from crashing the TUI.

**Solution**: `ActionConfig.Validate()` checks:
- `name` must be non-empty, lowercase, hyphens only
- `type` must be one of 5 valid types
- `label` must be non-empty
- `type: "select"` requires non-empty `options` array
- Unknown fields are ignored (forward compatibility)

**Bug caught**: Two manifests (audio, defaults) had `type: "select"` with empty `options` arrays, causing validation failure. Caught during PR #5 testing.

## 7. Help Overlay — Custom vs Library

**Decision**: Custom help overlay in `internal/tui/components/help.go` instead of wrapping `bubbles/help`.

**Rationale**: `bubbles/help` is a bottom-bar style component showing 1-2 lines of shortcuts. Wave 3 needed a full-screen overlay with context-sensitive bindings organized by category (Navigation, Actions, Global). Custom implementation at ~150 lines was simpler than adapting bubbles/help.

## 8. Testing Patterns for Bubble Tea TUI

**Patterns established**:
- **Unit tests**: Create widget model, `Update()` with key messages, assert view output or state changes
- **E2E tests**: `tea.NewProgram(model)` with `tea.WithInput()` pipe, send key sequences, assert final state
- **Integration tests**: Wire real hub with mock executor, verify module action → delta → settings flow
- **Safe pattern**: Always call `model.Init()` before `model.Update()` (bubbletea contract) — several early bugs were from omitting Init

## 9. Schema Migration — Atomic Deep Merge

**Key design**: Schema v1.0.0 → v1.1.0 migration uses `deepMerge`:
- Reads existing settings.json
- Creates default v1.1.0 object
- Deep-merges: user values overwrite defaults (not the other way)
- Writes atomically (temp file + rename)

**Why not additive only**: Simple field addition doesn't work when struct shapes change. `deepMerge` handles nested object merging correctly.

**Why atomic**: Temp file in same directory + `os.Rename()` ensures the write either completes fully or leaves the original intact.

## 10. Branch Chain Strategy

**Pattern**: 7 PRs in a feature-branch chain:
```
develop → wave-3/feat-1 → wave-3/feat-2 → wave-3/feat-3a → wave-3/feat-3b → wave-3/feat-4a → wave-3/feat-4b → wave-3/feat-5 → wave-3-tracker
```

- Each PR targets the previous PR's branch (except #1 targets `wave-3-tracker`)
- Final merge: `wave-3-tracker` → `develop`
- Each PR is independently reviewable (~200-500 lines)
- Critical: GitHub shows the diff from the base branch — if a previous slice leaks into a child's diff, retarget/rebase

**Workflow**: Implement → commit → push → create PR (ask user first) → wait for approval → merge → rebase next branch → continue.
