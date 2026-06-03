# lambda-env: File Manager Config (apps-06)

## Intent

Módulo TUI para configurar Yazi: tema, keymaps, plugins, comportamiento de preview.

## Scope

- Tema de Yazi (sincronizado con tema del sistema)
- Keymaps personalizados
- Plugins de Yazi: instalar, habilitar, deshabilitar
- Configuración de preview: imágenes, PDFs, videos
- Configurar opener por tipo de archivo

## Requirements

1. Leer/escribir `~/.config/yazi/` (theme.toml, keymap.toml, yazi.toml)
2. Tema sincronizado con `appearance.theme` de settings.json
3. Preview: toggle para imágenes (sixel/kitty graphics)
4. Openers: configurar qué app abre cada tipo de archivo
5. Backend: Yazi config files

## Scenarios

### Escenario 1: Cambiar tema de Yazi
- Usuario abre File Manager Config → Theme
- Selecciona "nord"
- Aplica → regenera theme.toml → Yazi usa nuevo tema

### Escenario 2: Configurar opener para PDFs
- Abre File Manager → Openers
- PDF: cambia de "zathura" a "okular"
- Aplica

## Technical Notes

- Yazi config: `~/.config/yazi/theme.toml`, `keymap.toml`, `yazi.toml`
- Temas de Yazi: generar desde los 5 temas del sistema
- Preview de imágenes: requiere kitty con graphics protocol o sixel
- Yazi ya está en packages.x86_64

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- yazi (ya en packages.x86_64)
