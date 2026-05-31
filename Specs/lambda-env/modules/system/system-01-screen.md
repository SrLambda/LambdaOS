# lambda-env: Screen Manager (system-01)

## Intent

Módulo TUI para gestionar monitores: detectar displays conectados, cambiar resolución, refresh rate, posición, modo espejo/extendido, y guardar perfiles de configuración.

## Scope

- Wrapper TUI de `xrandr`
- Detección automática de monitores conectados
- Selección de resolución y refresh rate por monitor
- Modo espejo / extendido
- Perfiles guardables (ej: "casa", "trabajo", "solo-laptop")
- Aplicar configuración con preview de 10 segundos (revertir si no se confirma)

## Requirements

1. Listar todos los outputs detectados con su estado (connected/disconnected)
2. Para cada monitor conectado, mostrar resoluciones disponibles ordenadas por tamaño
3. Permitir arrastrar monitores visualmente para definir posición (o seleccionar posición numérica)
4. Guardar perfiles en `settings.json` bajo `display.profiles`
5. Al aplicar, mostrar countdown de 10s para revertir si queda sin imagen
6. Funcionar con X11 (xrandr); Wayland queda para futuro

## Scenarios

### Escenario 1: Usuario conecta segundo monitor
- Abre lambda-env → System → Screen
- Ve: "eDP-1 (connected, 1920x1080), HDMI-1 (connected, no mode set)"
- Selecciona HDMI-1 → elige 1920x1080@60
- Elige posición: "a la derecha de eDP-1"
- Aplica → countdown 10s → confirma

### Escenario 2: Cambiar a modo espejo
- Selecciona "Mirror mode"
- Elige resolución común
- Aplica

### Escenario 3: Queda sin imagen
- Aplica configuración incorrecta
- No confirma en 10s
- Se revierte automáticamente a la configuración anterior

## Technical Notes

- Usar `xrandr --query` para detectar
- Usar `xrandr --output ... --mode ... --pos ...` para aplicar
- El fallback usa `xrandr --size` anterior guardado en variable
- Para el preview visual en TUI, mostrar un diagrama ASCII de la disposición

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- xorg-xrandr (ya en packages.x86_64)
