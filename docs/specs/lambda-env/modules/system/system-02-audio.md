# lambda-env: Audio Manager (system-02)

## Intent

Módulo TUI para gestionar audio: volumen, dispositivos de salida/entrada, perfiles de sonido, usando PipeWire como backend.

## Scope

- Control de volumen master (salida y entrada)
- Selección de dispositivo de salida por defecto
- Selección de dispositivo de entrada por defecto
- Per-files de audio (ej: "auriculares", "parlantes", "HDMI")
- Test de audio (play sample)

## Requirements

1. Mostrar volumen actual con slider TUI (0-100%)
2. Listar dispositivos de salida disponibles con nombre amigable
3. Listar dispositivos de entrada disponibles
4. Cambiar dispositivo default escribe en `settings.json`
5. Mute/unmute con toggle
6. Backend: `wpctl` (wireplumber) o `pactl` (pipewire-pulse)

## Scenarios

### Escenario 1: Conectar auriculares Bluetooth
- El dispositivo aparece en la lista
- Usuario lo selecciona como output default
- El volumen se ajusta al último usado para ese dispositivo

### Escenario 2: Micrófono no funciona
- Usuario va a Input devices
- Ve que el default es el micrófono de la webcam
- Cambia al micrófono interno
- Test de input: ve el nivel de audio en tiempo real

## Technical Notes

- `wpctl status` para listar dispositivos
- `wpctl set-volume` para cambiar volumen
- `wpctl set-default` para cambiar dispositivo default
- `pactl` como fallback si wireplumber no está disponible

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- pipewire, wireplumber (ya en packages.x86_64)
