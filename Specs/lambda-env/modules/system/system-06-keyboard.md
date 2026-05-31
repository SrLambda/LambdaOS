# lambda-env: Keyboard Manager (system-06)

## Intent

Módulo TUI para gestionar configuración del teclado: layout, variante, repeat rate, delay, y atajos globales del sistema.

## Scope

- Seleccionar layout de teclado (us, es, latam, etc.)
- Configurar variante (dvorak, colemak, etc.)
- Repeat rate y repeat delay
- Toggle numlock/capslock behavior
- Ver atajos globales configurados

## Requirements

1. Lista de layouts comunes con búsqueda
2. Preview del layout seleccionado
3. Slider para repeat rate y repeat delay
4. Aplicar cambios persiste en `settings.json`
5. Backend: `localectl`, `xkb`, `setxkbmap`

## Scenarios

### Escenario 1: Cambiar layout a español
- Usuario busca "spanish" o "es"
- Selecciona "es" → variante "latam"
- Aplica → el teclado cambia inmediatamente

### Escenario 2: Ajustar repeat rate
- Usuario siente que las teclas repiten muy lento
- Sube repeat rate de 25 a 50
- Reduce repeat delay de 600 a 300

## Technical Notes

- `localectl set-x11-keymap layout variant`
- `setxkbmap` para cambios en caliente
- Persistir en `/etc/X11/xorg.conf.d/00-keyboard.conf` o en settings.json
- Para Qtile, la config va en `config.py`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
