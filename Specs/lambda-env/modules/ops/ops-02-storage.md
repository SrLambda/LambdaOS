# lambda-env: Storage Manager (ops-02)

## Intent

Módulo TUI para gestionar almacenamiento: ver discos, particiones, uso, montar/desmontar, formatear.

## Scope

- Listar discos y particiones con uso
- Montar/desmontar particiones
- Ver uso de directorios principales
- Formatear partición (con advertencia)
- Configurar mounts en fstab

## Requirements

1. Listar discos: nombre, tamaño, tipo (SSD/HDD/NVMe), particiones
2. Para cada partición: filesystem, tamaño, usado, disponible, mount point
3. Montar: seleccionar partición → mount point → montar
4. Desmontar: seleccionar mount → desmontar
5. Formatear: seleccionar partición → filesystem → confirmar (destructivo)
6. Backend: `lsblk`, `df`, `mount`, `umount`, `blkid`, `fstab`

## Scenarios

### Escenario 1: Ver uso de disco
- Usuario abre Storage → ve:
  - nvme0n1: 512GB (SSD)
    - nvme0n1p1: EFI, 512MB, 30% usado
    - nvme0n1p2: BTRFS, 511GB, 45% usado, mounted at /

### Escenario 2: Montar USB
- Conecta USB → abre Storage
- Ve "sda1: FAT32, 32GB, unmounted"
- Selecciona → Mount → elige `/media/usb`
- Montado

## Technical Notes

- `lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINT,LABEL` para listar
- `df -h` para uso
- `mount /dev/X /path` para montar
- `udisksctl mount -b /dev/X` para montar como user
- `blkid` para info de filesystem
- Requiere root para formatear y editar fstab

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- udisks2 (para mount como user)
- Requiere root para formatear/fstab
