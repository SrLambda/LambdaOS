# lambda-env: Mkinitcpio Custom (polish-03)

## Intent

Configurar hooks y módulos custom de mkinitcpio para optimizar el initramfs de LambdaOS.

## Scope

- Configurar hooks de mkinitcpio
- Agregar módulos necesarios para hardware soporte
- Optimizar tamaño del initramfs
- Configurar para BTRFS + Plymouth

## Requirements

1. mkinitcpio.conf optimizado para LambdaOS
2. Hooks necesarios: base, udev, autodetect, modconf, kms, block, filesystems, keyboard, fsck
3. Plymouth hook incluido (si branding-05 está activo)
4. BTRFS modules incluidos
5. Initramfs tamaño razonable (< 100MB)

## Technical Notes

- Archivo: `airootfs/etc/mkinitcpio.conf`
- HOOKS:
  ```
  HOOKS=(base udev autodetect modconf kms block keyboard keymap consolefont btrfs plymouth filesystems fsck)
  ```
- MODULES:
  ```
  MODULES=(btrfs nvme)
  ```
- COMPRESSION: `zstd` (ya configurado)
- Plymouth hook requiere `plymouth` instalado
- BTRFS hook para soporte de subvolumes en boot
- Verificar orden de hooks: plymouth antes de kms

## Dependencies

- `branding-05-plymouth` (si se usa Plymouth)

## Verification

- `mkinitcpio -P` → genera initramfs sin errores
- Initramfs tamaño < 100MB
- Boot funciona con BTRFS root
- Plymouth muestra durante boot (si activo)
