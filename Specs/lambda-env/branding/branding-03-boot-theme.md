# lambda-env: Boot Theme (branding-03)

## Intent

Personalizar el tema de boot: GRUB para BIOS y systemd-boot para UEFI, con logo y colores de LambdaOS.

## Scope

- Tema de GRUB con logo LambdaOS y colores del tema
- Fondo de systemd-boot con logo
- Mensajes de boot con branding

## Requirements

1. GRUB muestra logo y colores de LambdaOS
2. systemd-boot muestra fondo custom
3. Funciona tanto en BIOS como UEFI

## Technical Notes

- GRUB: crear tema en `grub/` directory del perfil archiso
  - `theme.txt` con colores, fuentes, posiciones
  - Logo PNG para GRUB
  - Fuente terminus o personalizada
- systemd-boot: splash screen con `splash` kernel parameter
  - Fondo en `efiboot/loader/`
- Archiso ya tiene estructura `grub/` y `efiboot/` en el repo

## Dependencies

- Ninguno (GRUB y systemd-boot ya configurados en el perfil)

## Verification

- Boot BIOS → GRUB con tema LambdaOS
- Boot UEFI → systemd-boot con fondo LambdaOS
- Logo visible en ambos casos
