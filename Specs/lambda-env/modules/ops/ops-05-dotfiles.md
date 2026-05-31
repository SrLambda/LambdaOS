# lambda-env: Dotfiles Manager (ops-05)

## Intent

Módulo TUI para gestionar dotfiles con GNU Stow: ver módulos disponibles, stow/unstow, detectar conflictos, sync con backup.

## Scope

- Listar módulos de dotfiles disponibles (nvim, qtile, kitty, etc.)
- Stow/unstow módulos individuales
- Detectar conflictos (archivos que existen en home pero no en dotfiles)
- Backup: exportar dotfiles actuales al repo
- Perfiles de dotfiles: "dev", "minimal", "gaming"

## Requirements

1. Listar módulos en `~/dotfiles/` con estado (stowed/unstowed)
2. Stow: selecciona módulo → `stow <module>`
3. Unstow: selecciona módulo → `stow -D <module>`
4. Conflictos: detectar archivos en home que colisionan con dotfiles
5. Backup: copiar configs actuales a `~/dotfiles/` (override)
6. Backend: `stow`, diff de archivos

## Scenarios

### Escenario 1: Stow módulo de Kitty
- Usuario abre Dotfiles → ve lista:
  - nvim: ✓ stowed
  - qtile: ✓ stowed
  - kitty: ○ unstowed
- Selecciona kitty → Stow
- Kitty config aplicada

### Escenario 2: Detectar conflicto
- Usuario modificó `~/.config/kitty/kitty.conf` manualmente
- Abre Dotfiles → "Check conflicts"
- Ve: "kitty.conf: modified in home, differs from dotfiles repo"
- Opciones: "Use dotfiles version", "Keep home version", "Merge"

## Technical Notes

- GNU Stow ya está en packages.x86_64
- `stow --adopt` para tomar archivos existentes
- `stow -D` para unstow
- Detectar conflictos: comparar timestamps o checksums
- Perfiles: subdirectorios en `~/dotfiles/` con sets de módulos

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- stow (ya en packages.x86_64)
