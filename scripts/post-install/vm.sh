#!/usr/bin/env bash
# Profile: VM
# Installs virtualization and guest agent support for multiple hypervisors.
# Packages: virtualbox, virtualbox-guest-utils-nox, qemu-guest-agent,
#           open-vm-tools, hyperv
#
# Usage: sudo bash scripts/post-install/vm.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing VM profile packages..."

pacman -S --noconfirm \
    virtualbox \
    virtualbox-guest-utils-nox \
    qemu-guest-agent \
    open-vm-tools \
    hyperv

echo "==> VM profile installed successfully."
