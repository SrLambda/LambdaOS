# lambda-env: Agregar DaVinci Resolve (pkg-03)

## Intent

Agregar DaVinci Resolve al script de paquetes AUR post-boot para edición de video profesional.

## Scope

- Agregar `davinci-resolve` al script `scripts/aur-packages.sh`
- Documentar requisitos: libs propietarias, espacio en disco (~2GB)
- Agregar nota en README.md sobre instalación post-boot

## Requirements

1. DaVinci Resolve listado en `AUR_PACKAGES` del script
2. Documentación clara de que requiere AUR y ~2GB de espacio
3. El script maneja fallo gracefully (continúa con otros paquetes)

## Technical Notes

- Paquete AUR: `davinci-resolve` (o `davinci-resolve-studio` para versión paga)
- Requiere: `ocl-icd`, `opencl-nvidia` (para GPU NVIDIA), `libxcvt`
- Tamaño: ~2GB descargado + instalado
- No puede ir en la ISO directa (AUR-only, licencia propietaria)
- El script `aur-packages.sh` ya maneja fallos individuales con `continue`

## Dependencies

- `scripts/aur-packages.sh` existente

## Verification

- `davinci-resolve` aparece en `AUR_PACKAGES` del script
- README.md menciona DaVinci Resolve como post-boot
- Script ejecuta sin romper otros paquetes si davinci-resolve falla
