# lambda-env: Autostart Manager (system-11)

## Intent

Módulo TUI para gestionar aplicaciones que se ejecutan al inicio de sesión.

## Scope

- Listar apps al inicio (systemd user units + autostart dirs)
- Habilitar/deshabilitar apps
- Agregar nueva app al inicio
- Configurar delay de inicio

## Requirements

1. Listar servicios user habilitados: `systemctl --user list-unit-files`
2. Listar entries en `~/.config/autostart/`
3. Toggle on/off para cada entry
4. Agregar nuevo: seleccionar app de lista de `.desktop` disponibles
5. Backend: `systemctl --user`, XDG autostart spec

## Scenarios

### Escenario 1: Deshabilitar app al inicio
- Usuario ve "dunst.service" en la lista
- Toggle off → `systemctl --user disable dunst.service`

### Escenario 2: Agregar app al inicio
- "Add autostart" → lista de apps instaladas
- Selecciona "picom" → agrega como systemd user unit o .desktop

## Technical Notes

- Preferir systemd user units sobre XDG autostart (más control)
- `systemctl --user enable/disable <unit>`
- `systemctl --user start/stop <unit>`
- Crear units en `~/.config/systemd/user/` para apps custom

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
