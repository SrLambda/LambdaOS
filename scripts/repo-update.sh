#!/bin/bash
#
# LambdaOS Repository Update Script
#
# Regenerates the pacman database for the local lambdaos repository.
# Must be run as root. Handles empty repositories gracefully.
#
# Usage: ./scripts/repo-update.sh
# Exit codes: 0=success/empty, 1=not root, 2=repo-add failed

set -euo pipefail

REPO_DIR="/srv/repo/lambdaos"
ARCH="x86_64"
DB_NAME="lambdaos.db.tar"

# Require root privileges
if [[ "$(id -u)" -ne 0 ]]; then
    printf 'Error: this script requires root privileges\n' >&2
    exit 1
fi

# Ensure repository directory exists
if [[ ! -d "${REPO_DIR}/${ARCH}" ]]; then
    printf 'Error: repository directory %s/%s does not exist\n' "${REPO_DIR}" "${ARCH}" >&2
    exit 1
fi

cd "${REPO_DIR}"

# Collect packages using nullglob to handle empty directories safely
shopt -s nullglob
packages=("${ARCH}"/*.pkg.tar.zst)
shopt -u nullglob

# Handle empty repository gracefully
if ((${#packages[@]} == 0)); then
    printf 'No packages found in %s/. Repository is empty.\n' "${ARCH}"
    exit 0
fi

printf 'Adding %d package(s) to lambdaos repo...\n' "${#packages[@]}"

# Regenerate database with GPG signing
if ! repo-add --sign "${DB_NAME}" "${packages[@]}"; then
    printf 'Error: repo-add failed\n' >&2
    exit 2
fi

printf 'Repository database updated successfully.\n'
exit 0
