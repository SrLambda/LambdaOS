# lambda-env: Neovim Config (apps-01)

## Intent

Módulo TUI para configurar Neovim desde la interfaz central: toggles de LSP, Copilot, Neo-tree, gestión de plugins, keymaps.

## Scope

- Toggles: LSP on/off, Copilot on/off, Neo-tree on/off
- Lista de plugins instalados con toggle enable/disable
- Configurar tema de Neovim (sincronizado con tema del sistema)
- Configurar LSP servers activos
- Keymaps: ver y modificar atajos principales

## Requirements

1. Leer/escribir en `settings.json` bajo `neovim`
2. Actualizar `tui_settings.json` para compatibilidad con bridge existente
3. Al cambiar toggle, regenerar config de Neovim y recargar
4. Lista de plugins con estado (installed, disabled, update available)
5. Tema de Neovim sincronizado con `appearance.theme` del settings

## Scenarios

### Escenario 1: Desactivar Copilot
- Usuario abre Neovim Config → Toggles
- Toggle "Copilot" off
- Escribe en settings.json → regenera lazy.lua
- Próximo nvim abre sin Copilot

### Escenario 2: Cambiar tema
- Abre Neovim Config → Theme
- Selecciona "gruvbox" (sincronizado con tema del sistema)
- Aplica → nvim usa gruvbox

## Technical Notes

- Settings en `settings.json` → `neovim` section
- Bridge existente: `lua/core/tui_bridge.lua` lee `tui_settings.json`
- Actualizar bridge para leer de settings.json unificado
- Plugins gestionados por lazy.nvim
- Temas: catppuccin, gruvbox, tokyonight, nord, onedark (mismos que Qtile)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- `core/01-hub-plugin-system` (actualizar tui_bridge.lua)
