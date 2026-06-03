# lambda-env: Version Milestones

## Overview

This document maps each wave to a semantic version, defines the release process, and establishes the versioning scheme for LambdaOS pre-1.0.

**Este documento mapea cada wave a una versión semántica, define el proceso de release y establece el esquema de versionado para LambdaOS pre-1.0.**

---

## Version Scheme: `v0.wave.patch`

**EN**: For all pre-1.0 releases, we use `v0.wave.patch` where:
- `0` = pre-1.0 (no stable API yet)
- `wave` = the wave number that was completed (1-9)
- `patch` = bugfixes within that wave (0, 1, 2...)

This scheme is simple, traceable, and directly maps to the implementation plan. After v1.0.0, we switch to standard semver (`vMAJOR.MINOR.PATCH`).

**ES**: Para todos los releases pre-1.0, usamos `v0.wave.patch` donde:
- `0` = pre-1.0 (API no estable aún)
- `wave` = el número de wave completada (1-9)
- `patch` = bugfixes dentro de esa wave (0, 1, 2...)

Este esquema es simple, trazable, y mapea directamente al plan de implementación. Después de v1.0.0, cambiamos a semver estándar (`vMAJOR.MINOR.PATCH`).

### Version Mapping

| Wave | Version | Release Name | Key Deliverable |
|---|---|---|---|
| 0 | `v0.0.1` | Pipeline | CI buildea ISO, Flameshot, nombre correcto |
| 1 | `v0.1.0` | Foundation | Framework decidido, hub abre, repo pacman |
| 2 | `v0.2.0` | First Modules | TUI controla Neovim, Qtile, dotfiles |
| 3 | `v0.3.0` | Identity | Temas, MOTD, wallpaper, iconos |
| 4 | `v0.4.0` | Hardware | Pantalla, audio, energía, OBS |
| 5 | `v0.5.0` | Connectivity | Red, BT, teclado, todos los paquetes |
| 6 | `v0.6.0` | System | Servicios, updates, docs |
| 7 | `v0.7.0` | Complete Apps | Todas las apps + ops |
| 8 | `v0.9.0` | Installer | Wizard + Calamares funcional |
| 9 | `v1.0.0` | **RELEASE** | Distro completa con CD automático |

### Why v0.9.0 for Wave 8 (not v0.8.0)?

**EN**: Wave 8 is the installer — it's the last major feature before v1.0.0. Using v0.9.0 signals "almost there" and creates a natural progression: v0.9.0 (installer) → v1.0.0 (release).

**ES**: Wave 8 es el installer — es la última feature mayor antes de v1.0.0. Usar v0.9.0 señala "casi listo" y crea una progresión natural: v0.9.0 (installer) → v1.0.0 (release).

---

## Release Process

### Automatic Tags

**EN**: Each wave generates an automatic release tag when all specs in the wave are complete and the ISO passes all tests. The CI/CD pipeline (wave 9) automates this. Until then, tags are created manually.

**ES**: Cada wave genera un tag de release automático cuando todas las specs de la wave están completas y la ISO pasa todos los tests. El pipeline CI/CD (wave 9) automatiza esto. Hasta entonces, los tags se crean manualmente.

### Manual Release Process (Waves 0-8)

```bash
# 1. Complete all specs in the wave
# 2. Verify ISO builds and passes tests
./build_and_test.sh

# 3. Run QEMU tests
pytest tests/qemu/ -v

# 4. Create tag
git tag v0.1.0 -m "Wave 1: Foundation - hub, settings schema, repo pacman"

# 5. Push tag
git push origin v0.1.0

# 6. Build release ISO
LAMBDAOS_VERSION=0.1.0 sudo mkarchiso -v -w work/ -o out/ .

# 7. Create GitHub Release (manual until wave 9)
gh release create v0.1.0 out/LambdaOS-0.1.0-x86_64.iso \
  --title "LambdaOS v0.1.0 - Foundation" \
  --notes "Wave 1 complete: hub opens, settings schema, repo pacman configured"
```

### Automatic Release Process (Wave 9+)

```bash
# Just push the tag — CI/CD does everything
git tag v1.0.0 -m "LambdaOS v1.0.0 - First stable release"
git push origin v1.0.0

# CI/CD pipeline:
# 1. Builds ISO with LAMBDAOS_VERSION=1.0.0
# 2. Runs smoke + feature + install tests
# 3. Generates changelog (git-cliff)
# 4. Creates GitHub Release
# 5. Uploads ISO + SHA256SUMS + CHANGELOG.md
```

---

## Post-1.0 Versioning

After v1.0.0, we switch to standard semver:

| Version | Meaning | Example |
|---|---|---|
| `v1.1.0` | New feature (new module, new wave) | Add profiles module |
| `v1.1.1` | Bugfix in existing feature | Fix screen module crash |
| `v2.0.0` | Breaking change (API change, major refactor) | Switch hub language |

---

## ISO Naming Convention

The ISO filename follows the version:

```
LambdaOS-<version>-x86_64.iso
```

Examples:
- `LambdaOS-0.0.1-x86_64.iso` (Wave 0)
- `LambdaOS-0.1.0-x86_64.iso` (Wave 1)
- `LambdaOS-1.0.0-x86_64.iso` (Release)

This is configured in `profiledef.sh`:

```bash
iso_name="LambdaOS"
iso_version="${LAMBDAOS_VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
```

---

## Changelog Convention

**EN**: Each release has a changelog generated from conventional commits between tags. The changelog is bilingual (English/Spanish).

**ES**: Cada release tiene un changelog generado desde conventional commits entre tags. El changelog es bilingüe (Inglés/Español).

### Changelog Format

```markdown
# LambdaOS v0.1.0 - Foundation

## Features / Funcionalidades

- feat(hub): implement plugin system with module discovery
- feat(settings): unified JSON settings schema
- feat(repo): local pacman repository setup

## Bug Fixes / Correcciones

- fix(build): unmount chroot mounts before cleanup

## Documentation / Documentación

- docs(specs): add 60 spec files for lambda-env suite
```
