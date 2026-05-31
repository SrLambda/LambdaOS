# lambda-env: Development Conventions

## Overview

This document defines the development conventions for LambdaOS: branch strategy, commit messages, code structure, and documentation standards.

**Este documento define las convenciones de desarrollo para LambdaOS: estrategia de branches, mensajes de commit, estructura de código y estándares de documentación.**

---

## 1. Branch Strategy

### Branches

| Branch | Purpose | Protected |
|---|---|---|
| `main` | Stable releases only. Each commit is a tagged release. | ✅ Yes |
| `develop` | Integration branch. All wave work merges here first. | ✅ Yes |
| `wave-N-*` | Feature branches for each wave. Created from `develop`. | ❌ No |

### Current State

**EN**: As of the planning phase, only `main` exists. The `develop` branch needs to be created before Wave 0 work begins.

**ES**: Al momento de la planificación, solo existe `main`. La branch `develop` debe crearse antes de comenzar el trabajo de Wave 0.

### Setup

```bash
# Create develop branch from main
git checkout main
git checkout -b develop
git push -u origin develop
```

### Workflow

```
main ────────────────────────────────────────────────────── v0.0.1 ─── v0.1.0 ─── v1.0.0
  ↑                                                          ↑          ↑          ↑
develop ────────────────────────────────────────────────────┼──────────┼──────────┼
  ↑                                                          │          │          │
wave-0-pipeline ────────────────────────────────────────────┘          │          │
  ↑                                                                     │          │
wave-1-core ───────────────────────────────────────────────────────────┘          │
  ↑                                                                                │
wave-2-modules ───────────────────────────────────────────────────────────────────┘
```

### Merge Process

```bash
# 1. Create wave branch from develop
git checkout develop
git checkout -b wave-0-pipeline

# 2. Implement specs (one commit per spec)
git add <files>
git commit -m "feat(ci): add CI workflow for ISO build and lint"

# 3. Push and create PR to develop
git push -u origin wave-0-pipeline
gh pr create --base develop --title "Wave 0: Pipeline + ISO mínima funcional"

# 4. After PR review and merge to develop
git checkout develop
git pull

# 5. Tag release on main
git checkout main
git merge develop
git tag v0.0.1
git push origin main --tags

# 6. Merge main back to develop (sync)
git checkout develop
git merge main
git push origin develop
```

---

## 2. Commit Messages

### Format

We use Conventional Commits with scope:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | When to use | Example |
|---|---|---|
| `feat` | New feature or spec implementation | `feat(hub): add plugin system` |
| `fix` | Bug fix | `fix(build): unmount chroot before rm` |
| `docs` | Documentation changes | `docs(specs): add testing strategy` |
| `test` | Test additions or improvements | `test(qemu): add smoke test for wave 0` |
| `chore` | Maintenance, config, tooling | `chore(ci): add GitHub Actions workflow` |
| `refactor` | Code restructuring without behavior change | `refactor(airootfs): reorganize skel structure` |
| `style` | Formatting, whitespace, no code change | `style(qtile): format config.py with black` |
| `perf` | Performance improvements | `perf(hub): use parallel module discovery` |

### Scopes

| Scope | What it covers |
|---|---|
| `hub` | Core hub and plugin system |
| `settings` | Settings schema and management |
| `ci` | CI/CD workflows |
| `build` | ISO build process (mkarchiso, profiledef.sh) |
| `branding` | MOTD, wallpaper, boot theme, icons |
| `pkg` | Package additions (flameshot, obs, etc.) |
| `qtile` | Qtile window manager config |
| `nvim` | Neovim configuration |
| `docs` | Documentation (specs, README, guides) |
| `repo` | Pacman repository setup |
| `installer` | Calamares installer |
| `tui` | TUI modules (general) |

### Examples

```
feat(hub): add plugin system with module discovery
feat(settings): implement unified JSON settings schema
feat(pkg): add flameshot to packages.x86_64 + Qtile keybinding
fix(build): unmount chroot procfs before rm -rf
docs(specs): add module interface contract
test(qemu): add smoke test with 3-minute activity timeout
chore(ci): configure GitHub Actions for ISO build
refactor(airootfs): reorganize skel directory structure
style(qtile): format config.py with black
perf(hub): cache module manifest parsing
```

---

## 3. Bilingual Documentation

**EN**: All user-facing documentation (specs, README, guides, changelog) is written in both English and Spanish. Technical documentation (code comments, module interface contract) uses English as primary with Spanish summaries.

**ES**: Toda la documentación orientada al usuario (specs, README, guías, changelog) se escribe en inglés y español. La documentación técnica (comentarios de código, contrato de interfaz de módulos) usa inglés como primario con resúmenes en español.

### Spec Files

Each spec file has bilingual sections:

```markdown
# lambda-env: Module Name

## Intent / Intención

**EN**: One sentence describing what this module does.

**ES**: Una oración describiendo qué hace este módulo.

## Requirements / Requisitos

1. **EN**: First requirement.
   **ES**: Primer requisito.

2. **EN**: Second requirement.
   **ES**: Segundo requisito.

## Technical Notes / Notas Técnicas

**EN**: Technical details in English (for broader audience).

**ES**: Resumen en español si hay detalles importantes.
```

### Code Comments

```go
// Module discovery scans two directories for manifest.json files.
// User modules override system modules with the same name.
// Descubrimiento de módulos escanea dos directorios por manifest.json.
// Los módulos de usuario sobrescriben los del sistema con el mismo nombre.
func (h *Hub) DiscoverModules() ([]Module, error) {
    // ...
}
```

### README

The main README.md has bilingual sections:

```markdown
# LambdaOS

**EN**: The TUI-first Linux distribution.

**ES**: La distribución Linux donde la TUI es primero.

## Quick Start / Inicio Rápido

**EN**: Build the ISO with `sudo mkarchiso -v .`

**ES**: Construí la ISO con `sudo mkarchiso -v .`
```

---

## 4. Monorepo Structure

**EN**: All code lives in this repository until v1.0.0. After v1.0.0, the TUI may be split into its own repository if it becomes a standalone product.

**ES**: Todo el código vive en este repositorio hasta v1.0.0. Después de v1.0.0, la TUI puede separarse en su propio repositorio si se convierte en un producto independiente.

### Directory Structure

```
LambdaOS/
├── .github/
│   └── workflows/
│       ├── ci.yml                    ← CI pipeline (lint, build, test)
│       ├── cd.yml                    ← CD pipeline (releases)
│       └── nightly.yml               ← Nightly builds
│
├── Specs/
│   └── lambda-env/
│       ├── 00-suite-overview.md      ← Suite vision
│       ├── implementation-waves.md   ← Wave plan
│       ├── module-interface-contract.md
│       ├── testing-strategy.md
│       ├── version-milestones.md
│       ├── development-conventions.md
│       ├── core/                     ← Core specs
│       ├── packages/                 ← Package specs
│       ├── modules/                  ← Module specs
│       ├── branding/                 ← Branding specs
│       ├── installer/                ← Installer specs
│       ├── infrastructure/           ← Infrastructure specs
│       ├── polish/                   ← Polish specs
│       └── ci-cd/                    ← CI/CD specs
│
├── src/
│   └── lambda-env/                   ← TUI source code (hub + modules)
│       ├── cmd/
│       │   └── lambda-env/           ← Hub binary
│       ├── internal/
│       │   ├── hub/                  ← Plugin system
│       │   ├── settings/             ← Settings management
│       │   └── tui/                  ← TUI rendering
│       ├── modules/                  ← Module source code
│       │   ├── system/
│       │   ├── apps/
│       │   ├── ops/
│       │   └── setup/
│       ├── go.mod
│       └── go.sum
│
├── airootfs/                         ← Archiso root filesystem overlay
│   ├── etc/
│   ├── root/
│   └── usr/
│
├── grub/                             ← GRUB boot config
├── efiboot/                          ← UEFI boot config
├── syslinux/                         ← BIOS boot config
│
├── scripts/
│   ├── aur-packages.sh               ← AUR post-boot installer
│   ├── headless-install.sh           ← CI install test
│   └── repo-update.sh                ← Repo database update
│
├── tests/
│   ├── unit/                         ← Unit tests
│   └── qemu/                         ← QEMU integration tests
│
├── docs/                             ← Local documentation (served by darkhttpd)
│
├── profiledef.sh                     ← Archiso profile definition
├── packages.x86_64                   ← ISO package list
├── pacman.conf                       ← Pacman configuration
├── build_and_test.sh                 ← Build + test script
├── run_qemu.sh                       ← QEMU runner
├── Makefile                          ← Build shortcuts
├── README.md                         ← Project documentation (bilingual)
└── .gitignore
```

---

## 5. Code Quality

### Go Code

```bash
# Format
gofmt -w src/lambda-env/

# Lint
golangci-lint run src/lambda-env/

# Test
go test ./src/lambda-env/...
```

### Bash Scripts

```bash
# Format
shfmt -w scripts/*.sh build_and_test.sh

# Lint
shellcheck scripts/*.sh build_and_test.sh
```

### Python

```bash
# Format
black tests/

# Sort imports
isort tests/

# Lint
flake8 tests/
```

### Lua (Neovim config)

```bash
# Lint
luacheck airootfs/etc/skel/dotfiles/nvim/.config/nvim/
```

---

## 7. TUI View Prototypes (HTML)

**EN**: Every TUI interface design (view/screen) MUST have an HTML prototype before implementation. The prototype lives in `src/lambda-env/prototypes/` and is a single self-contained HTML file that renders the intended layout using basic HTML + CSS only — no JavaScript frameworks, no build step.

**ES**: Cada diseño de interfaz TUI (vista/pantalla) DEBE tener un prototipo en HTML antes de la implementación. El prototipo vive en `src/lambda-env/prototypes/` y es un archivo HTML autocontenido que renderiza el layout usando HTML + CSS básico — sin frameworks JavaScript, sin build step.

### Why HTML Prototypes

**EN**: HTML prototypes let you:
- Iterate on layout and visual hierarchy in seconds (browser refresh vs recompiling TUI)
- Share designs with stakeholders who can open a file in any browser
- Validate spacing, colors, and readability before locking implementation
- Serve as the visual spec that the TUI implementation must match

**ES**: Los prototipos HTML permiten:
- Iterar sobre layout y jerarquía visual en segundos (refresh del browser vs recompilar TUI)
- Compartir diseños con stakeholders que pueden abrir un archivo en cualquier browser
- Validar espaciado, colores y legibilidad antes de fijar la implementación
- Servir como especificación visual que la implementación TUI debe respetar

### File Naming

```
src/lambda-env/prototypes/
├── hub-menu.html          ← Main hub menu
├── settings-schema.html   ← Settings browser
├── screen-config.html     ← Screen/display config module
├── audio-config.html      ← Audio module
├── network-config.html    ← Network module
└── ...
```

### Prototype Requirements

1. **Single file** — all CSS inline or in `<style>` tag, no external dependencies
2. **Terminal aesthetic** — use monospace fonts, dark background, terminal-like colors
3. **No interactivity required** — static layout is enough; hover effects are optional
4. **Labeled sections** — each UI region has a comment explaining its purpose
5. **Updated when views change** — if the TUI view changes significantly, update the prototype

### Example Skeleton

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Hub Menu — LambdaOS TUI Prototype</title>
    <style>
        body {
            background: #1a1a2e;
            color: #e0e0e0;
            font-family: 'Courier New', monospace;
            padding: 2rem;
            max-width: 80ch;
            margin: 0 auto;
        }
        .header { color: #00d4ff; border-bottom: 1px solid #333; padding-bottom: 0.5rem; }
        .menu-item { padding: 0.25rem 0; }
        .menu-item.selected { background: #16213e; color: #00d4ff; }
        .footer { color: #666; margin-top: 1rem; }
    </style>
</head>
<body>
    <!-- Main menu listing -->
    <div class="header">LambdaOS — Configuration Hub</div>
    <div class="menu-item selected">▸ System</div>
    <div class="menu-item">  Applications</div>
    <div class="menu-item">  Operations</div>
    <div class="menu-item">  Setup</div>
    <div class="footer">[↑↓] Navigate  [Enter] Select  [q] Quit</div>
</body>
</html>
```

### PR Checklist Addition

- [ ] HTML prototype exists in `src/lambda-env/prototypes/` for each new TUI view
- [ ] Prototype reflects the final layout before TUI implementation begins

---

## 6. PR Guidelines

### PR Size

**EN**: Each PR should be reviewable. If the diff exceeds 400 lines, split into smaller PRs.

**ES**: Cada PR debe ser revisable. Si el diff supera 400 líneas, dividir en PRs más pequeños.

### PR Checklist

- [ ] All specs for this wave are implemented
- [ ] ISO builds successfully
- [ ] Smoke test passes (QEMU boot + 3-min activity check)
- [ ] Feature tests pass (if wave 2+)
- [ ] HTML prototype exists in `src/lambda-env/prototypes/` for each new TUI view
- [ ] Commit messages follow conventional commits format
- [ ] Documentation is bilingual (EN/ES)
- [ ] No linting errors

### PR Title Format

```
Wave N: <short description>

Examples:
Wave 0: Pipeline + ISO mínima funcional
Wave 1: Core TUI - hub, settings schema, repo pacman
Wave 2: First modules - Neovim, Qtile, dotfiles
```
