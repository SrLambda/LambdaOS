#!/usr/bin/env bash
# Profile: Rescue
# Installs disk recovery, partition management, and forensic tools.
# Packages: clonezilla, testdisk, partclone, fsarchiver, ddrescue, gpart
#
# Usage: sudo bash scripts/post-install/rescue.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing Rescue profile packages..."

pacman -S --noconfirm \
    clonezilla \
    testdisk \
    partclone \
    fsarchiver \
    ddrescue \
    gpart

echo "==> Rescue profile installed successfully."
