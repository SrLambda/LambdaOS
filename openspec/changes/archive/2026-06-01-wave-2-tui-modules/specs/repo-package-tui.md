# Spec: repo-package-tui

## Intent

Define the PKGBUILD for the `lambdaos-tui` pacman package that installs the hub binary, all Wave 2 modules, default settings, and Go templates to correct system paths.

## Requirements

### Requirement 1: PKGBUILD Structure

The system SHALL provide a valid PKGBUILD at `packages/lambdaos-tui/PKGBUILD` that builds the `lambdaos-tui` package from source using Go.

#### Scenario: PKGBUILD contains required fields

- GIVEN the PKGBUILD at `packages/lambdaos-tui/PKGBUILD`
- WHEN the PKGBUILD is parsed
- THEN it contains: `pkgname=lambdaos-tui`, `pkgver`, `pkgrel`, `pkgdesc`, `arch=('x86_64')`
- AND `depends` includes: `go`, `stow`, `qtile`, `neovim`
- AND `makedepends` includes: `go`

#### Scenario: PKGBUILD source points to correct location

- GIVEN the PKGBUILD `source` array
- WHEN the source is evaluated
- THEN it references the local source tree (no remote URL for repo package)
- AND uses `source=("$pkgname-$pkgver.tar.gz"::"$srcdir/../../")` or equivalent local reference

### Requirement 2: Build Process

The system SHALL compile all Go binaries (hub + modules) during the `build()` phase using `go build`.

#### Scenario: Build compiles hub binary

- GIVEN source is extracted to `$srcdir/`
- WHEN `build()` runs
- THEN `go build -o lambda-env ./cmd/lambda-env` executes from `src/lambda-env/`
- AND the `lambda-env` binary is produced

#### Scenario: Build compiles all modules

- GIVEN source contains modules under `src/lambda-env/internal/modules/`
- WHEN `build()` runs
- THEN each module is built as a separate binary:
  - `lambda-env-module-neovim`
  - `lambda-env-module-qtile`
  - `lambda-env-module-dotfiles`
- AND all binaries are produced without errors

#### Scenario: Build runs Go tests

- GIVEN the source contains test files
- WHEN `check()` runs
- THEN `go test ./...` executes from `src/lambda-env/`
- AND all tests pass (non-zero exit fails the build)

### Requirement 3: Install Paths

The system SHALL install files to the correct paths during the `package()` phase.

| Source | Destination |
|--------|-------------|
| `lambda-env` binary | `/usr/bin/lambda-env` |
| Module binaries | `/usr/share/lambda-env/modules/<name>/module` |
| Module manifests | `/usr/share/lambda-env/modules/<name>/manifest.json` |
| Go templates | `/usr/share/lambda-env/templates/` |
| Default settings | `/etc/lambdaos/settings.json` |

#### Scenario: Hub binary installed to /usr/bin

- GIVEN the `lambda-env` binary is built
- WHEN `package()` runs
- THEN `install -Dm755 lambda-env "$pkgdir/usr/bin/lambda-env"` executes
- AND the binary is at `/usr/bin/lambda-env` in the package

#### Scenario: Module binaries installed to modules directory

- GIVEN module binaries are built
- WHEN `package()` runs
- THEN each module is installed to `/usr/share/lambda-env/modules/<name>/module`
- AND each module's `manifest.json` is installed alongside it

#### Scenario: Default settings installed to /etc/lambdaos

- GIVEN a default `settings.json` exists in the source
- WHEN `package()` runs
- THEN `install -Dm644 settings.json "$pkgdir/etc/lambdaos/settings.json"` executes
- AND the default config is at `/etc/lambdaos/settings.json` in the package

#### Scenario: Templates installed to shared directory

- GIVEN Go templates exist in `src/lambda-env/pkg/templates/`
- WHEN `package()` runs
- THEN templates are installed to `/usr/share/lambda-env/templates/`
- AND the directory structure is preserved (neovim/, qtile/ subdirs)

### Requirement 4: Post-Install Hooks

The system SHALL provide a `.install` file at `packages/lambdaos-tui/lambdaos-tui.install` with post-install and pre-remove hooks.

#### Scenario: Post-install creates user config directory

- GIVEN the package is installed via `pacman -S lambdaos-tui`
- WHEN the `post_install()` hook runs
- THEN it creates `~/.config/lambdaos/` if it does not exist (for each user)
- AND copies default settings if `settings.json` does not exist

#### Scenario: Pre-remove preserves user settings

- GIVEN the package is being removed via `pacman -R lambdaos-tui`
- WHEN the `pre_remove()` hook runs
- THEN it does NOT delete `~/.config/lambdaos/settings.json`
- AND logs a message that user settings are preserved

#### Scenario: Post-upgrade is no-op

- GIVEN the package is being upgraded via `pacman -Syu`
- WHEN the `post_upgrade()` hook runs
- THEN it performs the same actions as `post_install()` (idempotent)

### Requirement 5: makepkg Verification

The system SHALL build successfully with `makepkg -s` in a clean environment.

#### Scenario: makepkg builds without errors

- GIVEN the PKGBUILD and all source files are in `packages/lambdaos-tui/`
- WHEN `makepkg -s` is executed from that directory
- THEN the build completes with exit code 0
- AND a `.pkg.tar.zst` file is produced

#### Scenario: makepkg fails on missing dependencies

- GIVEN a required build dependency is not installed
- WHEN `makepkg -s` is executed
- THEN makepkg attempts to install the dependency via `pacman -S`
- AND fails if the dependency cannot be resolved

#### Scenario: Package contents are correct

- GIVEN a `.pkg.tar.zst` file is produced
- WHEN `tar -tf lambdaos-tui-*.pkg.tar.zst` is executed
- THEN the archive contains:
  - `/usr/bin/lambda-env`
  - `/usr/share/lambda-env/modules/*/module`
  - `/usr/share/lambda-env/modules/*/manifest.json`
  - `/usr/share/lambda-env/templates/`
  - `/etc/lambdaos/settings.json`

## Technical Details

- PKGBUILD path: `packages/lambdaos-tui/PKGBUILD`
- Install file path: `packages/lambdaos-tui/lambdaos-tui.install`
- Default settings: `packages/lambdaos-tui/settings.json` (copy of Wave 1 defaults)
- Package name: `lambdaos-tui`
- Version: `0.2.0` (Wave 2 bump from Wave 1's `0.1.0`)
- Compression: `.pkg.tar.zst` (zstandard)

## Dependencies

- `infra-01-repo-pacman-setup` (Wave 1) — pacman repository configuration
- `core/01-hub-plugin-system` — hub binary and module structure
- All Wave 2 modules (neovim, qtile, dotfiles) — module binaries

## Verification Steps

```bash
# 1. PKGBUILD syntax is valid
cd packages/lambdaos-tui && namcap PKGBUILD

# 2. makepkg builds successfully
cd packages/lambdaos-tui && makepkg -s --noconfirm

# 3. Package contents are correct
tar -tf packages/lambdaos-tui/lambdaos-tui-*.pkg.tar.zst | sort

# 4. Install package
sudo pacman -U packages/lambdaos-tui/lambdaos-tui-*.pkg.tar.zst

# 5. Binary runs
lambda-env --help

# 6. Modules are discoverable
ls /usr/share/lambda-env/modules/

# 7. Default settings exist
cat /etc/lambdaos/settings.json | python3 -m json.tool

# 8. Remove package (verify pre-remove hook)
sudo pacman -R lambdaos-tui
# Verify: ~/.config/lambdaos/settings.json still exists
```
