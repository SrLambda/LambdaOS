# lambda-env: Services Manager (system-12)

## Intent

Módulo TUI para gestionar servicios del sistema: habilitar, deshabilitar, iniciar, detener, reiniciar servicios systemd.

## Scope

- Listar servicios activos, habilitados, fallidos
- Start/stop/restart/reload servicios
- Enable/disable servicios
- Ver logs de un servicio (integración con ops-03-logs)
- Filtrar por estado (active, inactive, failed)

## Requirements

1. Lista de servicios con estado visual (● running, ○ stopped, ✗ failed)
2. Acciones: start, stop, restart, enable, disable
3. Requiere root para servicios system; user para servicios --user
4. Mostrar descripción del servicio
5. Backend: `systemctl`

## Scenarios

### Escenario 1: Habilitar Docker al inicio
- Usuario busca "docker"
- Ve: "docker.service - Docker Application Container Engine"
- Estado: inactive
- Enable + Start

### Escenario 2: Investigar servicio fallido
- Ve "tailscaled.service" en rojo (failed)
- Selecciona → "View logs" → ve journalctl del servicio
- Selecciona → "Restart" → se recupera

## Technical Notes

- `systemctl list-units --type=service --all`
- `systemctl status <service>`
- `systemctl start/stop/restart/enable/disable <service>`
- `journalctl -u <service>` para logs
- Separar servicios system de servicios user

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Requiere root para servicios system
