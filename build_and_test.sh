#!/usr/bin/env bash
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
