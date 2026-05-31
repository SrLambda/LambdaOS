# lambda-env: Qtile Config (apps-02)

## Intent

Módulo TUI para configurar Qtile: layouts, keybindings, bar, groups, terminal default, browser default.

## Scope

- Configurar layouts activos y su orden
- Configurar keybindings principales (terminal, browser, file manager)
- Configurar bar: posición, altura, widgets activos
- Configurar groups (workspaces): nombres, cantidad, icons
- Configurar terminal y browser por defecto para Qtile
- Reload de Qtile tras cambios

## Requirements

1. Leer/escribir `settings.json` bajo `qtile`
2. Generar `config.py` de Qtile desde template + settings
3. Keybindings: terminal, browser, file manager, screenshot, launcher
4. Bar: widgets (CPU, memory, network, clock, layout, window-name)
5. Groups: 1-9 workspaces con nombres custom
6. `qtile cmd-obj -o cmd -f reload_config` para recargar

## Scenarios

### Escenario 1: Cambiar terminal default
- Usuario abre Qtile Config → Apps
- Cambia terminal de "kitty" a "foot"
- Aplica → regenera config.py → reload Qtile
- Próximo Mod+Enter abre foot

### Escenario 2: Agregar widget al bar
- Abre Qtile Config → Bar → Widgets
- Agrega "Battery" widget
- Aplica → bar se regenera con el nuevo widget

## Technical Notes

- Config de Qtile: `~/.config/qtile/config.py` (ya existe, modular)
- Modules existentes: `keys.py`, `groups.py`, `screens.py`, `theme.py`
- Generar config desde template Python con interpolación de settings
- `qtile cmd-obj -o cmd -f reload_config` para reload sin restart
- Theme ya está sincronizado via `theme.py`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- qtile (ya en packages.x86_64)
