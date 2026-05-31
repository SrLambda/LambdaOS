# lambda-env: Security Manager (system-14)

## Intent

Módulo TUI para gestionar seguridad básica del sistema: firewall, SSH keys, GPG, permisos.

## Scope

- Firewall: habilitar/deshabilitar, agregar reglas (ufw)
- SSH: generar keys, ver keys autorizadas, config de sshd
- GPG: listar keys, generar nueva key, firmar archivos
- Sudo: ver configuración de sudoers
- Permisos de archivos sensibles

## Requirements

1. Firewall status: enabled/disabled, reglas activas
2. Agregar regla: puerto, protocolo, dirección
3. SSH keygen: rsa, ed25519
4. GPG: listar, generar, exportar
5. Backend: `ufw`, `ssh-keygen`, `gpg`, `sudo`

## Scenarios

### Escenario 1: Habilitar firewall
- Usuario abre Security → Firewall
- Ve: "Firewall: inactive"
- Enable → "Allow SSH (22/tcp)" → "Allow HTTP (80/tcp)"
- Firewall activo

### Escenario 2: Generar SSH key
- Security → SSH Keys → "Generate new key"
- Tipo: ed25519
- Email: usuario@lambdaos
- Key generada en `~/.ssh/id_ed25519`

## Technical Notes

- `ufw` como firewall (simple, TUI-friendly)
- `ssh-keygen -t ed25519 -C "email"`
- `gpg --full-generate-key`
- `sudo visudo` para editar sudoers
- Agregar ufw a packages.x86_64

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- ufw (agregar a packages.x86_64)
- Requiere root para firewall
