# lambda-env: Terminal Config (apps-05)

## Intent

Módulo TUI para configurar Kitty: tema de colores, fuente, tamaño, keymaps, comportamiento.

## Scope

- Tema de colores (sincronizado con tema del sistema)
- Fuente y tamaño
- Opacidad del fondo
- Keymaps de Kitty
- Comportamiento: tabs, splits, scrollback
- Configurar shell por defecto

## Requirements

1. Leer/escribir `~/.config/kitty/kitty.conf`
2. Tema sincronizado con `appearance.theme` de settings.json
3. Fuente: selector de fuentes monospace instaladas
4. Tamaño: 8-24pt con preview
5. Opacidad: 0.5-1.0
6. Backend: kitty conf generation

## Scenarios

### Escenario 1: Cambiar tema de terminal
- Usuario abre Terminal Config → Theme
- Selecciona "tokyonight" (sincronizado con sistema)
- Aplica → regenera kitty.conf → kitty recarga con `kitty @ set-colors`

### Escenario 2: Cambiar tamaño de fuente
- Abre Terminal → Font size
- Cambia de 12 a 14
- Aplica → cambio inmediato en instancias existentes

## Technical Notes

- Kitty config: `~/.config/kitty/kitty.conf`
- `kitty @ set-colors` para cambiar colores en caliente
- `kitty @ set-font-size` para cambiar tamaño
- Temas de Kitty: generar desde los 5 temas del sistema
- Kitty ya está en packages.x86_64

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- kitty (ya en packages.x86_64)
