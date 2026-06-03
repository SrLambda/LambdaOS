# lambda-env: Calamares Scaffold (installer-01)

## Intent

Agregar Calamares como installer gráfico para LambdaOS con configuración básica funcional.

## Scope

- Agregar `calamares` a `packages.x86_64`
- Crear estructura de configuración en `airootfs/etc/calamares/`
- Configurar módulos básicos: welcome, locale, keyboard
- Desktop entry para lanzar desde el live environment

## Requirements

1. Calamares instalado en la ISO
2. Configuración básica funcional
3. Lanzable desde el menú de aplicaciones o terminal

## Technical Notes

- Paquete: `calamares` (extra de Arch)
- Config principal: `/etc/calamares/settings.conf`
- Módulos mínimos:
  - `welcome` → bienvenida
  - `locale` → idioma
  - `keyboard` → teclado
  - `partition` → particionado
  - `users` → crear usuario
  - `summary` → resumen antes de instalar
- Desktop entry: `/usr/share/applications/calamares.desktop` (viene con el paquete)
- Requiere `sudo calamares` para ejecutar

## Dependencies

- Ninguno (pero branding-03 mejora la experiencia)

## Verification

- `pacman -Q calamares` → instalado
- `sudo calamares` → abre installer
- Módulos básicos cargan sin error
