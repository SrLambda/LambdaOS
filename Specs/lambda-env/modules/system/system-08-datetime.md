# lambda-env: Date & Time Manager (system-08)

## Intent

Módulo TUI para gestionar fecha, hora, zona horaria y sincronización NTP.

## Scope

- Ver fecha y hora actual
- Cambiar zona horaria
- Toggle NTP (sincronización automática)
- Configurar hora manual (si NTP está desactivado)
- Formato de hora (12h / 24h)

## Requirements

1. Mostrar hora actual en tiempo real
2. Selector de zona horaria con búsqueda por ciudad/país
3. Toggle NTP on/off
4. Si NTP off, permitir setear hora manualmente
5. Formato 12h/24h persiste en settings.json
6. Backend: `timedatectl`

## Scenarios

### Escenario 1: Cambiar zona horaria
- Usuario viaja → abre módulo
- Busca "Buenos Aires" → selecciona "America/Argentina/Buenos_Aires"
- Aplica → hora se actualiza

### Escenario 2: Desactivar NTP y setear hora manual
- Toggle NTP off
- Setea hora: 14:30:00
- Setea fecha: 2026-05-30

## Technical Notes

- `timedatectl` como backend único
- `timedatectl list-timezones` para lista
- `timedatectl set-timezone ZONE`
- `timedatectl set-ntp true/false`
- `timedatectl set-time "YYYY-MM-DD HH:MM:SS"`
- Formato 12h/24h es un setting de UI (Qtile widgets, terminal)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Requiere root para cambiar timezone/hora
