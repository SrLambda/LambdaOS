# lambda-env: Notifications Manager (system-16)

## Intent

Módulo TUI para configurar el servidor de notificaciones Dunst: timeouts, colores, posición, reglas por app.

## Scope

- Configurar timeout de notificaciones (urgency low/normal/critical)
- Configurar posición en pantalla
- Configurar colores por urgencia
- Reglas por aplicación (silenciar, timeout custom)
- Ver historial de notificaciones recientes
- Test de notificación

## Requirements

1. Leer/escribir `~/.config/dunst/dunstrc`
2. Configurar timeout por urgencia
3. Configurar posición: top-right, top-center, bottom-right, etc.
4. Configurar colores: frame_color, background, foreground por urgencia
5. Reglas: por app_name, summary, o body
6. Test: enviar notificación de prueba con `notify-send`

## Scenarios

### Escenario 1: Cambiar posición de notificaciones
- Usuario abre Notifications → Position
- Cambia de "top-right" a "bottom-right"
- Recarga Dunst → notificaciones aparecen abajo

### Escenario 2: Silenciar notificaciones de una app
- Abre Rules → "Add rule"
- App: "Spotify"
- Action: "Silent" (no mostrar)
- Aplica

## Technical Notes

- Config de Dunst: `~/.config/dunst/dunstrc`
- `dunstctl reload` para recargar config
- `notify-send` para test
- `dunstctl history` para ver recientes
- Dunst ya está en packages.x86_64

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- dunst (ya en packages.x86_64)
