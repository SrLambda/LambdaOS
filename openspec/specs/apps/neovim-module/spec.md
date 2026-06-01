# Spec: neovim-module

## Intent

TUI module for Neovim configuration toggles (LSP, Copilot, Neo-tree), lazy.lua regeneration from Go templates, and tui_bridge.lua update to read unified settings.json.

## Requirements

### Requirement 1: Neovim Settings Schema

The system SHALL read and write the `neovim` section of `settings.json` with the following fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enable_lsp` | bool | `true` | Enable LSP servers |
| `enable_copilot` | bool | `true` | Enable GitHub Copilot |
| `enable_neotree` | bool | `true` | Enable Neo-tree file explorer |
| `theme` | string | `"tokyonight"` | Color theme name |
| `font` | string | `"JetBrainsMono"` | Font family |
| `lines` | int | `40` | Default window lines |
| `columns` | int | `120` | Default window columns |
| `lsp_servers` | []string | `["gopls","pyright"]` | Active LSP server list |

#### Scenario: Module reads neovim settings

- GIVEN `settings.json` exists with `neovim.enable_lsp: true`
- WHEN the neovim module loads settings
- THEN it returns a typed `NeovimSettings` struct with `EnableLSP: true`

#### Scenario: Module writes neovim settings delta

- GIVEN `settings.json` has `neovim.enable_copilot: true`
- WHEN the user toggles Copilot off via TUI
- THEN the module emits `settings_delta: {"neovim":{"enable_copilot":false}}`
- AND the hub merges the delta preserving other neovim fields

#### Scenario: Missing neovim section uses defaults

- GIVEN `settings.json` exists but has no `neovim` section
- WHEN the module loads settings
- THEN it returns default values for all neovim fields

### Requirement 2: Toggle Actions

The system SHALL support three toggle actions: `toggle_lsp`, `toggle_copilot`, `toggle_neotree`. Each toggle flips the boolean value, writes to settings.json, and triggers config regeneration.

#### Scenario: Toggle LSP off

- GIVEN `neovim.enable_lsp` is `true`
- WHEN the user selects "Toggle LSP" in TUI
- THEN `enable_lsp` is set to `false` in settings.json
- AND `lazy.lua` is regenerated without LSP plugin entries
- AND the module returns `{"status":"ok","action":"toggle_lsp","data":{"enabled":false}}`

#### Scenario: Toggle Copilot on

- GIVEN `neovim.enable_copilot` is `false`
- WHEN the user selects "Toggle Copilot" in TUI
- THEN `enable_copilot` is set to `true` in settings.json
- AND `lazy.lua` is regenerated with Copilot plugin entry
- AND the module returns `{"status":"ok","action":"toggle_copilot","data":{"enabled":true}}`

#### Scenario: Toggle Neo-tree

- GIVEN `neovim.enable_neotree` is `true`
- WHEN the user selects "Toggle Neo-tree" in TUI
- THEN `enable_neotree` is set to `false` in settings.json
- AND `lazy.lua` is regenerated without Neo-tree plugin entry

### Requirement 3: lazy.lua Regeneration

The system SHALL regenerate `~/.config/nvim/lua/plugins/lazy.lua` from a Go template using current `settings.json` neovim values. The template SHALL conditionally include plugin entries based on toggle states.

#### Scenario: Regenerate with all toggles on

- GIVEN all toggles are `true`
- WHEN `RegenerateLazyLua()` is called
- THEN `lazy.lua` contains LSP, Copilot, and Neo-tree plugin entries

#### Scenario: Regenerate with LSP off

- GIVEN `enable_lsp` is `false`
- WHEN `RegenerateLazyLua()` is called
- THEN `lazy.lua` does NOT contain `nvim-lspconfig` or LSP server entries
- AND `lazy.lua` still contains Copilot and Neo-tree entries

#### Scenario: Template validation before write

- GIVEN a Go template renders `lazy.lua` content
- WHEN the content is generated
- THEN the module validates it is non-empty before writing
- AND backs up the previous `lazy.lua` to `lazy.lua.bak`

### Requirement 4: tui_bridge.lua Unified Settings

The system SHALL update `airootfs/etc/skel/dotfiles/nvim/lua/core/tui_bridge.lua` to read from `settings.json` instead of `tui_settings.json`. The bridge SHALL parse the `neovim` section and apply settings on Neovim startup.

#### Scenario: Bridge reads settings.json

- GIVEN `settings.json` exists with `neovim.theme: "gruvbox"`
- WHEN Neovim starts and tui_bridge.lua executes
- THEN the bridge reads `settings.json` and applies `gruvbox` theme

#### Scenario: Bridge falls back to defaults

- GIVEN `settings.json` does not exist
- WHEN tui_bridge.lua executes
- THEN the bridge uses hardcoded default values
- AND no error is raised

#### Scenario: Bridge applies LSP toggle

- GIVEN `neovim.enable_lsp` is `false` in settings.json
- WHEN tui_bridge.lua executes
- THEN the bridge does NOT call `require("lspconfig")` setup

## Technical Details

- Go package: `src/lambda-env/internal/modules/neovim/`
- Template: `src/lambda-env/pkg/templates/neovim/lazy.lua.tmpl`
- Settings path: `~/.config/nvim/lua/plugins/lazy.lua`
- Bridge path: `airootfs/etc/skel/dotfiles/nvim/lua/core/tui_bridge.lua`
- Backup suffix: `.bak` (previous config preserved before regeneration)
- Module manifest category: `apps`

## Dependencies

- `core/01-hub-plugin-system` — module discovery and execution
- `core/02-settings-schema` — settings.json read/write
- Neovim + lazy.nvim installed on target system

## Verification Steps

```bash
# 1. Module compiles
cd src/lambda-env && go build ./internal/modules/neovim/...

# 2. Unit tests pass
cd src/lambda-env && go test ./internal/modules/neovim/... -v -cover

# 3. Toggle LSP off and verify lazy.lua
# Create settings.json with neovim.enable_lsp: true
# Run module with LAMBDA_ENV_ACTION=toggle_lsp
# Verify lazy.lua does NOT contain lspconfig entries

# 4. Template renders valid Lua
# Run: lua -l ~/.config/nvim/lua/plugins/lazy.lua
# Expect: no syntax errors

# 5. tui_bridge.lua reads settings.json
# Create settings.json with neovim.theme: "catppuccin"
# Run: nvim --headless -c "lua print(require('core.tui_bridge').theme)" -c q
# Expect: prints "catppuccin"
```
