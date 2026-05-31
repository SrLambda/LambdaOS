# lambda-env: Logs Viewer (ops-03)

## Intent

Módulo TUI para visualizar logs del sistema via journalctl con filtros y búsqueda.

## Scope

- Ver logs del sistema (journalctl)
- Filtrar por servicio, prioridad, tiempo
- Buscar en logs (keyword search)
- Ver logs de un servicio específico
- Seguir logs en tiempo real (tail -f equivalente)
- Exportar logs a archivo

## Requirements

1. Lista de servicios con logs disponibles
2. Filtros: tiempo (última hora, hoy, custom), prioridad (err, warning, info, debug)
3. Búsqueda: texto libre en logs
4. Follow mode: ver logs en tiempo real
5. Exportar: guardar logs filtrados a archivo
6. Backend: `journalctl`

## Scenarios

### Escenario 1: Investigar error de Docker
- Usuario abre Logs → Services → docker.service
- Ve últimos 50 logs
- Filtra por "error" → ve 3 entradas
- Selecciona entrada → ve detalle completo

### Escenario 2: Seguir logs en tiempo real
- Abre Logs → Follow mode
- Selecciona servicio: tailscaled
- Ve logs aparecer en tiempo real
- Ctrl+C para salir del follow

## Technical Notes

- `journalctl -u <service> -n 100` para últimos logs
- `journalctl -p err` para errores
- `journalctl --since "1 hour ago"`
- `journalctl -f -u <service>` para follow
- `journalctl | grep <keyword>` para búsqueda
- Renderizar logs con colores por prioridad

## Dependencies

- `core/01-hub-plugin-system`
- systemd (ya incluido)
