#!/usr/bin/env bash
# shellcheck disable=SC1091
# build_and_test.sh — Compila LambdaOS ISO y ejecuta pruebas E2E en QEMU
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== Phase 0: Setting up Python virtual environment ==="
if [ ! -d .venv ]; then
    python3 -m venv .venv
    echo "Virtual environment created at .venv/"
fi
source .venv/bin/activate
pip install -q -r requirements-dev.txt
echo "Dependencies installed"

echo "=== Phase 1: Cleaning previous builds ==="
# Unmount virtual filesystems from previous chroot before removing
# Handle shared/busy mounts by iterating in reverse order with retries
if [ -d work/x86_64/airootfs ]; then
    echo "Unmounting previous chroot filesystems..."
    for attempt in 1 2 3; do
        echo "  Unmount attempt ${attempt}/3..."
        local_mounts=$(grep "work/x86_64/airootfs" /proc/self/mountinfo 2>/dev/null | awk '{print $5}' | sort -r)
        if [ -z "$local_mounts" ]; then
            break
        fi
        while IFS= read -r mountpoint; do
            sudo umount -l "$mountpoint" 2>/dev/null || true
        done <<< "$local_mounts"
        sleep 1
    done
    # Final check: if mounts still exist, warn but continue
    remaining=$(grep -c "work/x86_64/airootfs" /proc/self/mountinfo 2>/dev/null || echo "0")
    if [ "$remaining" -gt 0 ]; then
        echo "WARNING: $remaining mounts still active under work/. rm may fail."
        echo "Run manually: sudo umount -l -R work/x86_64/airootfs/"
    fi
fi
sudo rm -rf work/ out/

echo "=== Phase 2: Verifying prerequisites ==="
if [ ! -f pacman.conf ]; then
    cp /usr/share/archiso/configs/releng/pacman.conf pacman.conf
    echo "pacman.conf copied from archiso releng default"
fi

echo "=== Phase 3: Building ISO with mkarchiso ==="
sudo mkarchiso -v -w work/ -o out/ .

echo "=== Phase 4: Running QEMU integration tests ==="
python -m pytest tests/qemu/test_live_boot.py -v

echo "=== Done ==="
