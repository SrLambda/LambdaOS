# lambda-env: System Monitor (ops-01)

## Intent

Módulo TUI que integra un monitor de sistema: CPU, RAM, disco, red, procesos.

## Scope

- Vista de CPU: uso por core, load average, temperatura
- Vista de RAM: uso, swap, por proceso
- Vista de disco: uso por partición, I/O
- Vista de red: throughput por interfaz
- Lista de procesos: ordenar por CPU, RAM, nombre
- Matar procesos

## Requirements

1. Integrar `htop` o `btop` como subprocesso O re-implementar vista básica
2. Mostrar resumen en el menú principal del hub
3. Permitir matar procesos (seleccionar de lista → kill)
4. Backend: `htop`, `btop`, o lectura directa de `/proc`

## Scenarios

### Escenario 1: Ver uso de sistema
- Usuario abre Monitor → ve dashboard:
  - CPU: 23% (4 cores)
  - RAM: 4.2G / 16G
  - Disk: 45% /home
  - Net: ↓ 2.3 MB/s ↑ 0.5 MB/s

### Escenario 2: Matar proceso colgado
- Abre Monitor → Processes
- Busca "chromium" → ve proceso usando 95% CPU
- Selecciona → Kill → SIGTERM

## Technical Notes

- Opción A: lanzar `btop` como subprocesso fullscreen (más fácil, mejor UX)
- Opción B: leer `/proc` directamente y renderizar en TUI (más control, más trabajo)
- Recomendación: Opción A para v1.0, Opción B para futuro
- btop ya es más completo que htop visualmente
- Agregar btop a packages.x86_64 si no está

## Dependencies

- `core/01-hub-plugin-system`
- btop o htop (agregar a packages.x86_64)
