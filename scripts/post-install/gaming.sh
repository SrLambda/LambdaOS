#!/usr/bin/env bash
# Profile: Gaming
# Installs Steam and Wine compatibility layer.
# Packages: steam, wine, wine-mono, winetricks
#
# Usage: sudo bash scripts/post-install/gaming.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing Gaming profile packages..."

pacman -S --noconfirm \
    steam \
    wine \
    wine-mono \
    winetricks

echo "==> Gaming profile installed successfully."
