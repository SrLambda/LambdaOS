# lambda-env: Settings Schema Unificado

## Intent

Definir un schema JSON unificado (`~/.config/lambdaos/settings.json`) que centralice toda la configuración del sistema. Este archivo es leído tanto por la TUI como por Qtile, Neovim y cualquier app que necesite configuración.

## Scope

- Reemplazar `tui_settings.json` de Neovim
- Reemplazar `os_theme.json` de Qtile
- Unificar settings de audio, red, pantalla, apariencia, etc.
- Schema versionado para migraciones futuras

## Requirements

1. El archivo debe ser un JSON válido con estructura predecible
2. Debe tener un campo `version` para migraciones
3. Cada módulo lee solo la sección que le corresponde
4. La TUI es la única que ESCRIBE en el archivo (las apps solo leen)
5. Debe tener valores por defecto si el archivo no existe o faltan campos
6. Debe ser compatible con el `tui_bridge.lua` existente de Neovim

## Scenarios

### Escenario 1: Primer boot (no existe settings.json)
- La TUI crea el archivo con valores por defecto
- Qtile lee el tema default (catppuccin)
- Neovim lee los defaults del bridge

### Escenario 2: Usuario cambia tema desde TUI
- TUI escribe `{"theme": "gruvbox"}` en la sección `appearance`
- Qtile detecta el cambio (file watcher o signal) y recarga
- Neovim no se ve afectado (lee otra sección)

### Escenario 3: Migración de schema
- Nueva versión de lambda-env agrega campo `audio.default_output`
- El hub detecta que `version` es menor que la actual
- Ejecuta migración automática agregando defaults para campos nuevos

## Schema Propuesto

```json
{
  "version": 1,
  "theme": "catppuccin",
  "appearance": {
    "theme": "catppuccin",
    "wallpaper": "/usr/share/lambdaos/wallpapers/default.png",
    "icon_theme": "Papirus",
    "cursor_theme": "Bibata-Modern-Ice",
    "font": "Monoid Nerd Font",
    "font_size": 12
  },
  "display": {
    "profiles": [
      {
        "name": "default",
        "outputs": [
          { "name": "HDMI-1", "mode": "1920x1080@60", "position": "0x0", "primary": true },
          { "name": "eDP-1", "mode": "1920x1080@60", "position": "1920x0" }
        ]
      }
    ],
    "active_profile": "default"
  },
  "audio": {
    "default_output": "alsa_output.pci-0000_00_1f.3.analog-stereo",
    "default_input": "alsa_input.pci-0000_00_1f.3.analog-stereo",
    "volume": 75
  },
  "network": {
    "wifi_enabled": true,
    "known_networks": []
  },
  "bluetooth": {
    "enabled": true,
    "paired_devices": []
  },
  "keyboard": {
    "layout": "us",
    "variant": "",
    "repeat_rate": 25,
    "repeat_delay": 600
  },
  "neovim": {
    "enable_lsp": true,
    "enable_copilot": false,
    "enable_neotree": true
  },
  "qtile": {
    "terminal": "kitty",
    "browser": "chromium",
    "file_manager": "yazi"
  },
  "services": {
    "enabled": ["docker", "tailscaled"],
    "disabled": []
  }
}
```

## Technical Notes

- La TUI debe validar el JSON antes de escribir
- Usar escritura atómica (escribir a temp file + rename) para evitar corrupción
- Qtile puede usar `watchdog` o `inotify` para detectar cambios
- Neovim bridge debe actualizarse para leer de este schema en vez de `tui_settings.json`

## Dependencies

- `01-hub-plugin-system` (el hub usa el schema para cargar contexto)
