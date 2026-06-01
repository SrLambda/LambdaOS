# Spec: infra-01-repo-pacman-setup

## Intent

Establish a local pacman repository at `/srv/repo/lambdaos/` on the LambdaOS ISO, configure `pacman.conf` with a `[lambdaos]` section, provide a `repo-add` script for database regeneration, and set up GPG signing for package verification. This infrastructure enables packaging and distribution of `lambdaos-tui` and its modules.

## Scope

### In Scope
- Repository directory structure at `/srv/repo/lambdaos/x86_64/`
- `pacman.conf` configuration adding `[lambdaos]` section with `file://` server
- `scripts/repo-update.sh` for regenerating the repo database with `repo-add`
- GPG signing key generation and configuration for the LambdaOS repo
- PKGBUILD template for `lambdaos-tui` package
- Verification that `pacman -Sl lambdaos` lists available packages

### Out of Scope
- Actual PKGBUILD for `lambdaos-tui` — Wave 2 (`infra-02-repo-package-tui`)
- Remote repository hosting — Wave 1 is local-only (file://)
- CI pipeline for package building — Wave 2

## Requirements

### Requirement 1: Repository Directory Structure

The system SHALL create the directory structure `/srv/repo/lambdaos/x86_64/` on the ISO with appropriate permissions for pacman to read.

#### Scenario: Directory structure is created

- GIVEN the ISO is being built
- WHEN the airootfs overlay is applied
- THEN `/srv/repo/lambdaos/x86_64/` exists on the live system
- AND the directory is readable by all users (755)

#### Scenario: Repository is empty initially

- GIVEN the ISO is freshly built
- WHEN no packages have been added to the repo
- THEN `/srv/repo/lambdaos/x86_64/` contains no `.pkg.tar.zst` files
- AND no database files exist (`lambdaos.db.tar`, `lambdaos.files.tar`)

### Requirement 2: Pacman Configuration

The system SHALL add a `[lambdaos]` section to `pacman.conf` with `Server = file:///srv/repo/lambdaos/$arch` and `SigLevel = Required`.

#### Scenario: Lambdaos repo section is present and enabled

- GIVEN pacman.conf is modified
- WHEN the ISO is built
- THEN `[lambdaos]` section exists and is uncommented
- AND `Server = file:///srv/repo/lambdaos/$arch` is configured

#### Scenario: Signature level is set to Required

- GIVEN the `[lambdaos]` section is configured
- WHEN pacman.conf is parsed
- THEN `SigLevel = Required` is set for the lambdaos repo
- AND packages without valid signatures are rejected

#### Scenario: Lambdaos repo is positioned correctly

- GIVEN pacman.conf contains all repositories
- WHEN the file is read top-to-bottom
- THEN `[lambdaos]` appears after `[multilib]` and before any `[custom]` section

### Requirement 3: Repo Update Script

The system SHALL provide `scripts/repo-update.sh` that regenerates the pacman database using `repo-add --sign`.

#### Scenario: Script regenerates database

- GIVEN one or more `.pkg.tar.zst` files exist in `/srv/repo/lambdaos/x86_64/`
- WHEN `scripts/repo-update.sh` is executed
- THEN `lambdaos.db.tar` is created/updated in `/srv/repo/lambdaos/`
- AND `lambdaos.db.tar.sig` is created (signed database)

#### Scenario: Script handles empty repository

- GIVEN no `.pkg.tar.zst` files exist in `/srv/repo/lambdaos/x86_64/`
- WHEN `scripts/repo-update.sh` is executed
- THEN the script exits with code 0
- AND no error is raised (empty repo is valid)

#### Scenario: Script requires root privileges

- WHEN `scripts/repo-update.sh` is executed as a non-root user
- THEN the script exits with a non-zero code
- AND an error message indicates root is required

#### Scenario: Script uses repo-add with signing

- GIVEN the script is executed
- WHEN `repo-add` is invoked
- THEN the `--sign` flag is passed
- AND the database is signed with the LambdaOS GPG key

### Requirement 4: GPG Signing Key

The system SHALL generate and configure a GPG signing key for the LambdaOS repository.

#### Scenario: GPG key is generated

- GIVEN the ISO build process includes key generation
- WHEN the build completes
- THEN a GPG key exists with UID containing "LambdaOS"
- AND the key has signing capability

#### Scenario: Key is trusted in pacman keyring

- GIVEN the GPG key is generated
- WHEN the ISO is built
- THEN the key is added to the pacman keyring (`pacman-key --add`)
- AND the key is locally signed (`pacman-key --lsign-key`)

#### Scenario: Packages are signed during repo-add

- GIVEN a package is added to the repo via `repo-update.sh`
- WHEN the database is regenerated
- THEN the package signature is included in the database
- AND `pacman -S lambdaos-tui` verifies the signature successfully

### Requirement 5: PKGBUILD Template

The system SHALL provide a PKGBUILD template for `lambdaos-tui` at a documented location that can be used as a starting point for Wave 2.

#### Scenario: PKGBUILD template exists

- GIVEN the repo infrastructure is set up
- WHEN the template is needed
- THEN a PKGBUILD template file exists at `templates/PKGBUILD.lambdaos-tui` or equivalent
- AND the template includes standard Arch packaging fields (pkgname, pkgver, pkgrel, arch, depends)

#### Scenario: Template has correct package metadata

- GIVEN the PKGBUILD template is read
- WHEN the metadata is inspected
- THEN `pkgname=lambdaos-tui`
- AND `arch=('x86_64')`
- AND `depends=('glibc')` (minimum, expanded in Wave 2)

#### Scenario: Template includes Go build steps

- GIVEN the PKGBUILD template is read
- WHEN the build() function is inspected
- THEN it contains `go build` commands for `src/lambda-env/`
- AND the binary is installed to `/usr/bin/lambda-env`

### Requirement 6: Repository Verification

The system SHALL enable verification that the lambdaos repo is functional via `pacman -Sl lambdaos`.

#### Scenario: Pacman lists repo packages

- GIVEN the lambdaos repo is configured and has packages
- WHEN `pacman -Sl lambdaos` is executed
- THEN packages in the repo are listed with their versions
- AND each package shows as `[installed]` or `[available]`

#### Scenario: Pacman syncs repo database

- GIVEN the lambdaos repo is configured in pacman.conf
- WHEN `pacman -Sy` is executed
- THEN the lambdaos database is downloaded/synced
- AND no errors are reported for the lambdaos repo

#### Scenario: Package from lambdaos repo is installable

- GIVEN `lambdaos-tui` is in the lambdaos repo
- WHEN `pacman -S lambdaos-tui` is executed
- THEN the package is found in the lambdaos repo
- AND it installs successfully with signature verification

## Technical Details

### Repository Directory Structure
```
/srv/repo/lambdaos/
├── x86_64/
│   ├── lambdaos-tui-0.1.0-1-x86_64.pkg.tar.zst    # Package file
│   └── lambdaos-tui-0.1.0-1-x86_64.pkg.tar.zst.sig # Package signature
├── lambdaos.db.tar                                  # Package database
├── lambdaos.db.tar.sig                              # Database signature
└── lambdaos.files.tar                               # File list database
```

### Pacman.conf Addition
```ini
[lambdaos]
Server = file:///srv/repo/lambdaos/$arch
SigLevel = Required
```

### Repo Update Script (scripts/repo-update.sh)
```bash
#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="/srv/repo/lambdaos"
ARCH="x86_64"

if [ "$(id -u)" -ne 0 ]; then
    echo "Error: this script requires root privileges" >&2
    exit 1
fi

cd "$REPO_DIR"

# Check if there are packages to add
shopt -s nullglob
packages=("$ARCH"/*.pkg.tar.zst)
shopt -u nullglob

if [ ${#packages[@]} -eq 0 ]; then
    echo "No packages found in $ARCH/. Repository is empty."
    exit 0
fi

echo "Adding ${#packages[@]} package(s) to lambdaos repo..."
repo-add --sign "$REPO_DIR/lambdaos.db.tar" "${packages[@]}"
echo "Repository database updated."
```

### GPG Key Generation
```bash
# Generate signing key (during ISO build)
gpg --batch --gen-key <<EOF
%no-protection
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: LambdaOS Package Signing
Name-Email: packages@lambdaos.local
Expire-Date: 0
%commit
EOF

# Add to pacman keyring
pacman-key --add <(gpg --export packages@lambdaos.local)
pacman-key --lsign-key packages@lambdaos.local
```

### PKGBUILD Template (templates/PKGBUILD.lambdaos-tui)
```bash
# Maintainer: LambdaOS Team
pkgname=lambdaos-tui
pkgver=0.1.0
pkgrel=1
pkgdesc="LambdaOS TUI hub and module system"
arch=('x86_64')
url="https://github.com/lambdaos/lambdaos"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz::https://github.com/lambdaos/lambdaos/archive/refs/tags/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
    cd "lambdaos-${pkgver}/src/lambda-env"
    go build -o lambda-env ./cmd/lambda-env
}

package() {
    cd "lambdaos-${pkgver}/src/lambda-env"
    install -Dm755 lambda-env "${pkgdir}/usr/bin/lambda-env"
}
```

### Modified Files
- `pacman.conf` — add `[lambdaos]` section after `[multilib]`
- `scripts/repo-update.sh` — new file, repo database regeneration
- `Makefile` — add `repo-update` target

## Dependencies

- `pacman` — package manager, repo-add tool
- `gnupg` — GPG key generation and signing
- `repo-add` — pacman database tool (part of pacman package)
- Root privileges — required for repo-update.sh and pacman operations

## Verification Steps

```bash
# 1. Repository directory structure exists
ls -la /srv/repo/lambdaos/x86_64/
# Expect: directory exists, readable by all

# 2. Pacman.conf has lambdaos section
grep -A2 '\[lambdaos\]' /etc/pacman.conf
# Expect: [lambdaos] section with Server and SigLevel

# 3. Repo update script exists and is executable
test -x scripts/repo-update.sh && echo "Script is executable"
# Expect: "Script is executable"

# 4. Repo update script runs shellcheck
shellcheck scripts/repo-update.sh
# Expect: no warnings or errors

# 5. Repo update script runs shfmt
shfmt -d scripts/repo-update.sh
# Expect: no formatting differences

# 6. GPG key exists in pacman keyring
pacman-key --list-keys | grep -i lambdaos
# Expect: LambdaOS Package Signing key listed

# 7. Pacman syncs lambdaos repo
pacman -Sy lambdaos 2>&1 | grep -i lambdaos
# Expect: lambdaos database synced

# 8. Pacman lists repo (empty is valid)
pacman -Sl lambdaos
# Expect: no error, may show "no packages" if repo is empty

# 9. PKGBUILD template exists
test -f templates/PKGBUILD.lambdaos-tui && echo "Template exists"
# Expect: "Template exists"

# 10. PKGBUILD passes namcap (if available)
namcap templates/PKGBUILD.lambdaos-tui 2>/dev/null || echo "namcap not available"
# Expect: no critical errors
```
