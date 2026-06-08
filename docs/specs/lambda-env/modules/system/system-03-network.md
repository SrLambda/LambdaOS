# lambda-env: Network Manager (system-03)

## Intent

Módulo TUI para gestionar conectividad de red: WiFi, Ethernet, VPN.

## Scope

- Escanear y conectar a redes WiFi
- Ver estado de Ethernet
- Configurar VPN (Tailscale, OpenVPN)
- Ver IP, DNS, gateway actual

## Requirements

1. Listar redes WiFi disponibles con SSID, señal, seguridad
2. Conectar a red WiFi (ingresar contraseña si es WPA)
3. Mostrar estado de conexión activa (IP, gateway, DNS)
4. Toggle WiFi on/off
5. Mostrar estado de VPN conectada
6. Backend: `iwctl` (iwd) o `nmcli` (NetworkManager)

## Scenarios

### Escenario 1: Conectar a WiFi nuevo
- Usuario escanea redes
- Selecciona su red del listado
- Ingresa contraseña
- Confirma conexión

### Escenario 2: Ver info de red actual
- Abre módulo → ve "Conectado a: MiRed (WiFi)"
- IP: 192.168.1.42, Gateway: 192.168.1.1, DNS: 1.1.1.1

## Technical Notes

- `iwctl station wlan0 scan` + `iwctl station wlan0 get-networks` para WiFi
- `iwctl station wlan0 connect <ssid>` para conectar
- `ip addr`, `ip route` para info de red
- Tailscale ya está en la ISO (`tailscale status`)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- iwd o NetworkManager (definir cuál usar)
