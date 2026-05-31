#!/bin/bash
#
# AUR Package Installer for LambdaOS
#
# Installs AUR-only packages that cannot be included in the archiso build.
# Requires yay or paru to be installed beforehand.
#
# Usage: ./scripts/aur-packages.sh
# Exit codes: 0=all installed, 1=missing AUR helper, 2=partial failures

set -euo pipefail

# AUR packages to install (must exist in AUR)
AUR_PACKAGES=(
    spotify  # Music streaming client
    obsidian # Knowledge base / notes
    megasync # Mega.nz cloud sync
    bluetui  # Bluetooth TUI manager
    impala   # WiFi TUI manager
)

# Detect AUR helper
detect_aur_helper() {
    if command -v yay > /dev/null 2>&1; then
        echo "yay"
        return 0
    elif command -v paru > /dev/null 2>&1; then
        echo "paru"
        return 0
    fi
    return 1
}

helper=$(detect_aur_helper) || {
    echo "Error: No AUR helper found."
    echo ""
    echo "Install one of the following before running this script:"
    echo "  yay:    pacman -S yay"
    echo "  paru:   pacman -S paru"
    echo ""
    echo "Or build from source:"
    echo "  git clone https://aur.archlinux.org/yay.git && cd yay && makepkg -si"
    exit 1
}

echo "Using AUR helper: ${helper}"
echo ""

failed=()
success=()

for pkg in "${AUR_PACKAGES[@]}"; do
    echo "Installing ${pkg}..."
    if "${helper}" -S --needed --noconfirm "${pkg}"; then
        success+=("${pkg}")
        echo "  OK: ${pkg}"
    else
        failed+=("${pkg}")
        echo "  FAILED: ${pkg} (continuing...)"
    fi
    echo ""
done

echo "========================================"
echo "AUR Package Installation Summary"
echo "========================================"
echo "Successful (${#success[@]}): ${success[*]}"
echo "Failed     (${#failed[@]}): ${failed[*]}"
echo "========================================"

if ((${#failed[@]} > 0)); then
    exit 2
fi

exit 0
