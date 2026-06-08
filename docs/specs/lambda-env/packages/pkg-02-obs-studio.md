# lambda-env: Agregar OBS Studio (pkg-02)

## Intent

Agregar OBS Studio a la ISO como herramienta de grabación y streaming de pantalla.

## Scope

- Agregar `obs-studio` a `packages.x86_64`
- Crear directorio default de grabaciones: `~/Videos/Recordings/`
- Configurar perfil básico de grabación (1080p30)

## Requirements

1. OBS Studio disponible en la ISO live
2. Funcional con PipeWire (ya incluido)
3. Perfil básico pre-configurado para grabación de pantalla completa

## Technical Notes

- Paquete: `obs-studio` (repositorio extra de Arch)
- PipeWire ya está incluido → OBS lo detecta automáticamente
- Config default en `~/.config/obs-studio/` via skel
- Dependencias de OBS: pipewire, xdg-desktop-portal, xdg-desktop-portal-gtk

## Dependencies

- Ninguno (PipeWire ya está en packages.x86_64)

## Verification

- `pacman -Q obs-studio` → instalado
- `obs` → abre sin errores
- PipeWire source disponible en OBS
