#!/usr/bin/env bash
# Profile: Office
# Installs office suite, email client, PDF viewer, and web browser.
# Packages: libreoffice-fresh, thunderbird, okular, chromium
#
# Usage: sudo bash scripts/post-install/office.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing Office profile packages..."

pacman -S --noconfirm \
    libreoffice-fresh \
    thunderbird \
    okular \
    chromium

echo "==> Office profile installed successfully."
