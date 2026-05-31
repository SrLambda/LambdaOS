# lambda-env: Appearance Manager (system-09)

## Intent

Módulo TUI para gestionar la apariencia visual del sistema: tema de colores, wallpaper, icon theme, cursor theme, fuente del sistema.

## Scope

- Selector de tema de colores (Catppuccin, Gruvbox, TokyoNight, Nord, OneDark)
- Selector de wallpaper
- Selector de icon theme
- Selector de cursor theme
- Selector de fuente del sistema y tamaño
- Preview del tema antes de aplicar

## Requirements

1. Mostrar los 5 temas disponibles con preview de colores en TUI
2. Al seleccionar tema, actualizar `settings.json` y disparar reload de Qtile
3. Wallpaper: seleccionar desde `/usr/share/lambdaos/wallpapers/` o path custom
4. Icon theme y cursor theme: listar disponibles via `fc-list` y `~/.icons`
5. Fuente: listar fuentes monospace disponibles
6. Todos los cambios persisten en `settings.json` bajo `appearance`

## Scenarios

### Escenario 1: Cambiar tema de Catppuccin a Gruvbox
- Usuario abre Appearance → Themes
- Ve preview de cada tema
- Selecciona Gruvbox → aplica
- Qtile recarga con nuevos colores
- Terminal cambia colores (si está configurado)

### Escenario 2: Cambiar wallpaper
- Abre Appearance → Wallpaper
- Ve thumbnails (ASCII o lista) de wallpapers disponibles
- Selecciona uno → aplica
- Qtile setea el wallpaper de fondo

## Technical Notes

- Temas definidos en `theme.py` de Qtile (ya existen 5)
- Wallpaper: `feh --bg-fill` o `xsetroot` en Qtile autostart
- Icon themes: listar dirs en `/usr/share/icons/` y `~/.local/share/icons/`
- Cursor themes: listar dirs en `/usr/share/icons/` con `index.theme` que tenga `Inherits=`
- Fuentes: `fc-list :spacing=mono` para monospace

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- feh (agregar a packages.x86_64 para wallpaper)
