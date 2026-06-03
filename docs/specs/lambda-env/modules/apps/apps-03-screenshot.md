# lambda-env: Screenshot Config (apps-03)

## Intent

Módulo TUI para configurar Flameshot: atajos de teclado, directorio de capturas, comportamiento de copiado al portapapeles.

## Scope

- Configurar directorio donde se guardan las capturas
- Toggle: copiar al portapapeles automáticamente
- Toggle: mostrar notificación tras captura
- Configurar atajos: fullscreen, area selection, window capture
- Configurar delay antes de captura

## Requirements

1. Configurar Flameshot via su config file o CLI
2. Directorio de capturas: selector de path
3. Atajos: integrar con Qtile keybindings
4. Backend: `flameshot config`, `flameshot gui`, `flameshot full`
5. Persistir settings en `settings.json`

## Scenarios

### Escenario 1: Configurar directorio de capturas
- Usuario abre Screenshot Config → Save directory
- Cambia de `~/Pictures` a `~/Screenshots`
- Aplica → próxima captura se guarda ahí

### Escenario 2: Configurar atajo de captura de área
- Abre Shortcuts → "Area selection"
- Asigna: Mod+Shift+S
- Integra con Qtile keybindings

## Technical Notes

- Flameshot: `flameshot config` para settings
- `flameshot gui` para selección de área
- `flameshot full -p <path>` para fullscreen
- `flameshot full -c` para copiar al clipboard
- Atajos se integran en Qtile `keys.py`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- `apps-02-qtile` (para integrar keybindings)
- flameshot (agregar a packages.x86_64)
