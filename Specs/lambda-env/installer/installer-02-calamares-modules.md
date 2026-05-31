# lambda-env: Calamares Modules (installer-02)

## Intent

Configurar todos los módulos de Calamares para una instalación completa de LambdaOS.

## Scope

- Configurar módulo de particionado: BTRFS con subvolumes automático
- Configurar módulo de users: crear usuario con sudo
- Configurar módulo de bootloader: GRUB + systemd-boot
- Configurar módulo de summary: resumen personalizado
- Configurar módulo de shellprocess: post-install scripts

## Requirements

1. Particionado automático con BTRFS subvolumes (@, @home, @log, @snapshots)
2. Usuario creado con grupo wheel y sudo habilitado
3. Bootloader instalado correctamente (BIOS + UEFI)
4. Post-install: configurar fstab, generar initramfs, habilitar servicios

## Technical Notes

- Partición: `btrfs` con subvolumes
  - `@` → /
  - `@home` → /home
  - `@log` → /var/log
  - `@snapshots` → /.snapshots
- Users: `auto-login=false`, `sudoers-group=wheel`, `require-root-password=true`
- Bootloader: `grub` para BIOS, `systemd-boot` para UEFI
- Post-install shellprocess:
  - `pacman -S --noconfirm linux linux-firmware`
  - `mkinitcpio -P`
  - `grub-install` o `bootctl install`
  - Habilitar servicios: NetworkManager, docker, tailscaled, etc.

## Dependencies

- `installer-01-calamares-scaffold`

## Verification

- Instalación completa en VM → sistema booteable
- BTRFS subvolumes creados correctamente
- Usuario puede hacer sudo
- Bootloader funcional
