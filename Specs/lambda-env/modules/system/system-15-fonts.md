# lambda-env: Fonts Manager (system-15)

## Intent

Módulo TUI para gestionar fuentes del sistema: listar, previsualizar, instalar, desinstalar fuentes.

## Scope

- Listar fuentes instaladas con preview
- Instalar fuente desde archivo (TTF, OTF)
- Desinstalar fuente
- Reconstruir cache de fuentes
- Configurar fuente default del sistema

## Requirements

1. Listar fuentes agrupadas: system, user, monospace, sans, serif
2. Preview de cada fuente con texto de ejemplo
3. Instalar: copiar a `~/.local/share/fonts/` o `/usr/share/fonts/`
4. `fc-cache -fv` después de instalar/desinstalar
5. Backend: `fc-list`, `fc-cache`

## Scenarios

### Escenario 1: Instalar nueva fuente
- Usuario descarga un TTF
- Abre Fonts → "Install font"
- Selecciona archivo → instala en `~/.local/share/fonts/`
- Rebuild cache → fuente disponible

### Escenario 2: Cambiar fuente de terminal
- Abre Fonts → "Set default monospace"
- Elige "JetBrainsMono Nerd Font"
- Persiste en `settings.json`

## Technical Notes

- `fc-list --format='%{family}\n'` para listar
- `fc-list :spacing=mono` para monospace
- `fc-cache -fv` para rebuild
- Instalar en `~/.local/share/fonts/` (user) o `/usr/share/fonts/` (system)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- fontconfig (ya incluido via dependencias)
