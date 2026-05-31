#!/usr/bin/env bash
# Profile: Media
# Installs media playback and streaming tools.
# Packages: vlc
#
# Usage: sudo bash scripts/post-install/media.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing Media profile packages..."

pacman -S --noconfirm \
    vlc

echo "==> Media profile installed successfully."
