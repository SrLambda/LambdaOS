# lambda-env: System Profiles (setup-03)

## Intent

Módulo TUI para gestionar perfiles de sistema: cambiar entre configuraciones predefinidas (dev, gaming, rescue, minimal) que ajustan servicios, apps, y comportamiento.

## Scope

- Perfiles predefinidos:
  - **Dev**: Docker, lenguajes, LSP, tmux, servicios de dev habilitados
  - **Gaming**: Steam, Wine, GPU performance, servicios de gaming
  - **Rescue**: mínimo, solo herramientas de rescate, sin GUI
  - **Minimal**: solo lo esencial, sin servicios extra
- Aplicar perfil: habilita/deshabilita servicios, configura Qtile, ajusta servicios
- Crear perfil custom
- Exportar/importar perfiles

## Requirements

1. Cada perfil define: servicios a habilitar, servicios a deshabilitar, apps al inicio, modo de Qtile (full/minimal)
2. Aplicar perfil es idempotente
3. Mostrar diff antes de aplicar: "Esto habilitará X, deshabilitará Y"
4. Confirmar antes de aplicar
5. Perfil activo visible en el hub principal
6. Backend: `systemctl`, settings.json, Qtile config

## Scenarios

### Escenario 1: Cambiar a perfil Gaming
- Usuario abre Profiles → ve "Current: Dev"
- Selecciona "Gaming"
- Ve diff:
  - + Enable: docker, steam, gamemode
  - - Disable: tailscaled, dnsmasq
  - Qtile: agregar keybindings de gaming
- Confirma → aplica → reboot de Qtile

### Escenario 2: Crear perfil custom
- Abre Profiles → "Create custom"
- Nombre: "Streaming"
- Habilita: OBS, nginx (para streaming server)
- Deshabilita: docker
- Qtile: layout optimizado para OBS
- Guarda

## Technical Notes

- Perfiles definidos como JSON en `/usr/share/lambdaos/profiles/`
- Perfil custom en `~/.config/lambdaos/profiles/`
- Cada perfil: `{ "name", "services_enable", "services_disable", "autostart", "qtile_mode" }`
- Aplicar perfil: iterar servicios + regenerar Qtile config + reload
- Perfil activo en `settings.json` bajo `active_profile`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- `system-12-services`
- `apps-02-qtile`
- Requiere root para servicios
