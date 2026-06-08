# lambda-env: Services Default (polish-02)

## Intent

Configurar qué servicios systemd están habilitados por defecto en LambdaOS.

## Scope

- Definir lista de servicios habilitados por defecto
- Definir lista de servicios deshabilitados por defecto
- Configurar en el perfil de archiso

## Requirements

1. Servicios esenciales habilitados: NetworkManager, docker, tailscaled, sshd
2. Servicios no esenciales deshabilitados: se activan on-demand via TUI
3. Servicios habilitados al boot de la ISO live
4. Servicios habilitados tras la instalación

## Technical Notes

- Servicios a habilitar:
  - `NetworkManager.service` (o `iwd.service` + `systemd-networkd`)
  - `docker.service`
  - `tailscaled.service`
  - `sshd.service`
  - `pipewire.service` (user)
  - `wireplumber.service` (user)
  - `dunst.service` (user)
  - `ly.service` (display manager)
- Servicios a deshabilitar:
  - `dnsmasq.service` (solo si se usa)
  - `modemmanager.service` (solo si hay modem)
- Archiso: configurar en `airootfs/etc/systemd/system/` con symlinks
- Post-install: Calamares shellprocess habilita servicios

## Dependencies

- Ninguno

## Verification

- `systemctl list-unit-files --state=enabled` → muestra servicios correctos
- ISO bootea con red funcionando
- Docker disponible tras boot
