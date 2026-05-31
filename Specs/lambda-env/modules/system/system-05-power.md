# lambda-env: Power Manager (system-05)

## Intent

Módulo TUI para gestionar energía: brillo de pantalla, batería, suspensión, hibernación.

## Scope

- Control de brillo de pantalla
- Información de batería (porcentaje, tiempo restante, estado)
- Configurar timeouts de suspensión
- Suspender/hibernar/reiniciar/apagar
- Configurar comportamiento al cerrar tapa (laptop)

## Requirements

1. Slider de brillo (0-100%) con cambio en tiempo real
2. Mostrar info de batería: %, tiempo restante, charging/discharging
3. Configurar: suspensión tras X minutos inactivo, apagar pantalla tras Y minutos
4. Acciones rápidas: Suspend, Hibernate, Reboot, Power Off
5. Backend: `brightnessctl`, `upower`, `systemd-inhibit`, `systemctl`

## Scenarios

### Escenario 1: Ajustar brillo
- Usuario abre módulo → ve brillo actual: 75%
- Sube a 100% con flechas → cambio inmediato

### Escenario 2: Configurar suspensión
- Va a "Power settings"
- Cambia "Suspend after" de 30min a 15min
- Cambia "Screen off after" de 10min a 5min

## Technical Notes

- `brightnessctl` para brillo (ya en paquetes o agregar)
- `upower -d` para info de batería
- `systemd-inhibit` para prevenir suspensión
- `systemctl suspend/hibernate/reboot/poweroff`
- Logind.conf para timeouts de suspensión

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- brightnessctl (agregar a packages.x86_64)
- upower (agregar a packages.x86_64)
