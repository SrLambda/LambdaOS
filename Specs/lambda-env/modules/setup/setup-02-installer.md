# lambda-env: Installer (setup-02)

## Intent

Módulo TUI para instalar LambdaOS en disco: particionado, filesystem, boot loader, usuario, y copia del sistema.

## Scope

- Seleccionar disco de destino
- Particionado: automático (BTRFS con subvolumes) o manual
- Configurar bootloader (GRUB + systemd-boot)
- Configurar hostname
- Crear usuario principal
- Copiar sistema al disco
- Configurar fstab
- Instalar bootloader
- Reboot

## Requirements

1. Detectar discos disponibles con tamaño y tipo
2. Particionado automático:
   - EFI partition (512MB, FAT32)
   - BTRFS root (restante, con subvolumes: @, @home, @log, @snapshots)
3. Particionado manual: herramienta de partición TUI
4. Bootloader: GRUB para BIOS+UEFI, systemd-boot para UEFI
5. Usuario: nombre, password, hostname
6. Progreso visible durante la instalación
7. Validación: confirmar disco destino (es destructivo)

## Scenarios

### Escenario 1: Instalación automática
- Usuario abre Installer
- Selecciona disco: "nvme0n1 (512GB)"
- "Automatic partitioning (BTRFS)" → confirma
- Hostname: "lambdaos"
- Usuario: "lambda", password, confirma
- "Begin installation" → progreso
- "Installation complete" → Reboot

### Escenario 2: Instalación con particionado manual
- Selecciona disco
- "Manual partitioning"
- Crea: EFI 512MB, / 100GB BTRFS, /home restante BTRFS
- Continúa con resto del flujo

## Technical Notes

- Opción A: Calamares (GUI, ya planificado) — este módulo sería un launcher
- Opción B: Installer TUI propio — más trabajo pero coherente con la distro
- Para v1.0: recomendación es Calamares con branding LambdaOS
- Para v2.0: installer TUI propio
- Herramientas necesarias: `arch-install-scripts`, `grub`, `efibootmgr`, `btrfs-progs`
- Ya están en packages.x86_64

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Requiere root
- Branding (setup-02 depende de branding para Calamares)
