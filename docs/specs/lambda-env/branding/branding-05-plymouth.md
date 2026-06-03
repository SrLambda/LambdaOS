# lambda-env: Plymouth Boot Splash (branding-05)

## Intent

Agregar Plymouth como boot splash animado para LambdaOS, mostrando logo durante el boot.

## Scope

- Instalar Plymouth
- Crear tema de Plymouth con logo de LambdaOS
- Configurar initramfs para usar Plymouth
- Animación de boot con branding

## Requirements

1. Plymouth instalado y configurado
2. Tema custom con logo LambdaOS
3. Funciona con mkinitcpio
4. Boot muestra splash en vez de texto de kernel

## Technical Notes

- Paquete: `plymouth` (extra de Arch)
- Tema: crear en `/usr/share/plymouth/themes/lambdaos/`
  - `lambdaos.plymouth` (config del tema)
  - Logo animado (script o animación de frames)
- mkinitcpio: agregar `plymouth` hook en `airootfs/etc/mkinitcpio.conf`
- Kernel parameter: `quiet splash` en GRUB/systemd-boot
- Hook order: `plymouth` debe ir antes de `udev` y después de `base`

## Dependencies

- plymouth (agregar a packages.x86_64)
- mkinitcpio hook configuration

## Verification

- Boot → splash animado visible
- No se ve texto de kernel (quiet splash)
- Plymouth hook en mkinitcpio.conf
