#!/bin/bash

# 1. Define la ruta a tu ISO
ISO_PATH=$(find out/ -maxdepth 1 -name 'LambdaOS-*-x86_64.iso' -print | tail -n 1)
echo "Usando ISO: $ISO_PATH"

# 2. Obtener el UUID del sistema de archivos de la ISO
ISO_UUID=$(blkid -s UUID -o value "$ISO_PATH")

# 3. Extraer el kernel y el initramfs temporalmente
mkdir -p /tmp/iso_extract
echo "Extrayendo kernel e initramfs..."
7z x -aoa -o/tmp/iso_extract "$ISO_PATH" arch/boot/x86_64/vmlinuz-linux arch/boot/x86_64/initramfs-linux.img -y > /dev/null

# 4. Lanzar QEMU directamente en esta misma terminal
echo "Iniciando QEMU..."
qemu-system-x86_64 \
    -M pc \
    -m 2G \
    -enable-kvm \
    -nographic \
    -device virtio-rng-pci \
    -kernel /tmp/iso_extract/arch/boot/x86_64/vmlinuz-linux \
    -initrd /tmp/iso_extract/arch/boot/x86_64/initramfs-linux.img \
    -append "archisobasedir=arch archisosearchuuid=$ISO_UUID console=ttyS0,115200" \
    -cdrom "$ISO_PATH"
