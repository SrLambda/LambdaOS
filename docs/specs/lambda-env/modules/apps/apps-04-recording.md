# lambda-env: Recording Config (apps-04)

## Intent

Módulo TUI para configurar OBS Studio: perfiles de grabación, escenas, atajos de grabación rápida, directorio de videos.

## Scope

- Configurar perfil de grabación: calidad, formato, FPS
- Configurar escenas básicas: "Pantalla completa", "Ventana", "Webcam"
- Configurar directorio de grabaciones
- Configurar atajos: iniciar/detener grabación, pausar
- Toggle: grabar audio del micrófono
- Toggle: grabar audio del sistema

## Requirements

1. Configurar OBS via CLI o archivo de config
2. Perfiles: Low (720p30), Medium (1080p30), High (1080p60), Ultra (4K30)
3. Directorio de grabaciones: selector de path
4. Atajos: integrar con Qtile keybindings
5. Backend: OBS CLI, config files en `~/.config/obs-studio/`

## Scenarios

### Escenario 1: Configurar grabación rápida
- Usuario abre Recording Config → Quick Record
- Selecciona perfil: "Medium (1080p30)"
- Configura atajo: Mod+Shift+R para start/stop
- Aplica

### Escenario 2: Cambiar directorio de grabaciones
- Abre Recording → Save directory
- Cambia a `~/Videos/Recordings`
- Aplica

## Technical Notes

- OBS config: `~/.config/obs-studio/basic/scenes/` y `profiles/`
- OBS CLI: `obs --startrecording`, `obs --stoprecording`
- Para quick recording, considerar `wf-recorder` o `gpu-screen-recorder` como alternativa más ligera
- Atajos se integran en Qtile `keys.py`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- `apps-02-qtile` (para integrar keybindings)
- obs-studio (agregar a packages.x86_64)
