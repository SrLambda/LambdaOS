# lambda-env: Agregar Flameshot (pkg-01)

## Intent

Agregar Flameshot a la ISO como herramienta de captura de pantalla y configurar un keybinding en Qtile para acceso rápido.

## Scope

- Agregar `flameshot` a `packages.x86_64`
- Configurar keybinding en Qtile: `Mod+Shift+S` → `flameshot gui`
- Configurar keybinding: `Mod+Ctrl+S` → `flameshot full`
- Agregar Flameshot al autostart de Qtile (para clipboard integration)

## Requirements

1. Flameshot disponible en la ISO live
2. Keybindings funcionales en Qtile
3. Directorio default de capturas: `~/Pictures/Screenshots/`
4. Copiar al portapapeles automáticamente

## Technical Notes

- Paquete: `flameshot` (repositorio extra de Arch)
- Qtile keybinding: `Key([mod, "Shift"], "s", lazy.spawn("flameshot gui"))`
- Qtile keybinding: `Key([mod, "Control"], "s", lazy.spawn("flameshot full"))`
- Crear directorio `~/Pictures/Screenshots/` en skel

## Dependencies

- Ninguno

## Verification

- `pacman -Q flameshot` → instalado
- `flameshot gui` → abre selector de área
- Keybindings responden correctamente
