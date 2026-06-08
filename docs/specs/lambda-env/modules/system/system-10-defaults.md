# lambda-env: Default Apps Manager (system-10)

## Intent

Módulo TUI para configurar aplicaciones por defecto del sistema usando XDG MIME.

## Scope

- Browser por defecto
- Terminal por defecto
- File manager por defecto
- Text editor por defecto
- Image viewer por defecto
- PDF viewer por defecto
- Video player por defecto
- Email client por defecto

## Requirements

1. Listar categorías de apps con la app actual asignada
2. Para cada categoría, listar alternativas instaladas
3. Cambiar app por defecto persiste en `settings.json` y XDG
4. Backend: `xdg-mime`, `xdg-settings`, `update-alternatives`

## Scenarios

### Escenario 1: Cambiar browser de Chromium a Firefox
- Usuario abre Default Apps → Browser
- Ve: "Current: chromium"
- Selecciona "firefox" de la lista de alternativas
- Aplica

### Escenario 2: Ver todas las defaults actuales
- Abre módulo → ve resumen:
  - Browser: chromium
  - Terminal: kitty
  - Editor: nvim
  - Files: yazi
  - PDF: zathura

## Technical Notes

- `xdg-settings set default-web-browser chromium.desktop`
- `xdg-mime default nvim.desktop text/plain`
- `update-alternatives` para algunos casos
- Apps disponibles: escanear `.desktop` files en `/usr/share/applications/`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- xdg-utils (ya en packages.x86_64)
