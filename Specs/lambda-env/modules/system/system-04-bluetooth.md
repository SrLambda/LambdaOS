# lambda-env: Bluetooth Manager (system-04)

## Intent

Módulo TUI para gestionar dispositivos Bluetooth: escanear, emparejar, conectar, desconectar, olvidar.

## Scope

- Toggle Bluetooth on/off
- Escanear dispositivos cercanos
- Emparejar nuevo dispositivo
- Conectar/desconectar dispositivos emparejados
- Olvidar dispositivo
- Ver dispositivos conectados y su tipo (audio, input, etc.)

## Requirements

1. Mostrar estado de Bluetooth (on/off)
2. Lista de dispositivos emparejados con nombre, tipo, estado (connected/disconnected)
3. Escanear nuevos dispositivos (timeout 30s)
4. Proceso de pairing con PIN si es necesario
5. Backend: `bluetoothctl` via D-Bus o comando directo

## Scenarios

### Escenario 1: Emparejar auriculares
- Usuario activa Bluetooth
- Escanea → ve "AirPods Pro" en la lista
- Selecciona → pairing → confirma
- Se conecta automáticamente

### Escenario 2: Reconectar dispositivo conocido
- Abre módulo → ve "Mouse MX Master (disconnected)"
- Selecciona → Connect → listo

## Technical Notes

- `bluetoothctl` como backend principal
- `bluetoothctl power on/off`
- `bluetoothctl scan on`
- `bluetoothctl pair <mac>`
- `bluetoothctl connect <mac>`
- `bluetoothctl trust <mac>` para reconexión automática

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- bluez (incluir en packages.x86_64 si no está)
