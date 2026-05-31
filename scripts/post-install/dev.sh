#!/usr/bin/env bash
# Profile: Dev
# Installs development toolchains, language servers, and build tools.
# Packages: jdk-openjdk, docker, docker-compose, texlive-basic,
#           texlive-latex, texlive-latexextra, shellcheck, shfmt,
#           clang, go, nodejs, npm, python-black, python-isort, stylua
#
# Usage: sudo bash scripts/post-install/dev.sh
# Idempotent: safe to run multiple times (pacman skips installed packages).

set -euo pipefail

echo "==> Installing Dev profile packages..."

pacman -S --noconfirm \
    jdk-openjdk \
    docker \
    docker-compose \
    texlive-basic \
    texlive-latex \
    texlive-latexextra \
    shellcheck \
    shfmt \
    clang \
    go \
    nodejs \
    npm \
    python-black \
    python-isort \
    stylua

echo "==> Dev profile installed successfully."
