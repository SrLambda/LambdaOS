#!/usr/bin/env bash
# All Profiles — installs every post-install profile.
# Runs all profile scripts in sequence: dev, gaming, office, vm, media, rescue.
#
# Usage: sudo bash scripts/post-install/all.sh
# Idempotent: safe to run multiple times (each profile script is idempotent).

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "============================================"
echo "  LambdaOS — Installing ALL profiles"
echo "============================================"

profiles=(dev gaming office vm media rescue)

for profile in "${profiles[@]}"; do
    echo ""
    echo "--- Profile: ${profile} ---"
    bash "${SCRIPT_DIR}/${profile}.sh"
done

echo ""
echo "============================================"
echo "  All profiles installed successfully."
echo "============================================"
