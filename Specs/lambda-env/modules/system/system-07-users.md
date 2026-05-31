# lambda-env: Users Manager (system-07)

## Intent

Módulo TUI para gestión básica de cuentas de usuario: crear, modificar, eliminar usuarios, cambiar password, gestionar grupos.

## Scope

- Listar usuarios del sistema
- Crear nuevo usuario con home directory
- Cambiar password de usuario
- Agregar/quitar usuario de grupos (wheel, docker, etc.)
- Configurar shell por defecto
- Auto-login toggle

## Requirements

1. Requiere privilegios de root (sudo)
2. Listar usuarios con uid, shell, grupos principales
3. Crear usuario: nombre, password, grupos, shell
4. Cambiar password de usuario existente
5. Toggle auto-login en display manager (Ly)

## Scenarios

### Escenario 1: Crear usuario nuevo
- Admin abre módulo → "Create user"
- Ingresa: username, password, confirma password
- Selecciona grupos: wheel, docker
- Usuario creado con home directory

### Escenario 2: Agregar usuario a grupo docker
- Selecciona usuario existente
- "Manage groups" → agrega "docker"
- Aplica

## Technical Notes

- `useradd`, `usermod`, `passwd`, `userdel` como backend
- `getent passwd` para listar usuarios
- Auto-login: configurar Ly (`/etc/ly/config.ini`)
- Siempre requerir sudo

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Requiere root
