# lambda-env: Installer Launcher (installer-04)

## Intent

Crear accesos para lanzar el installer desde la TUI y desde el desktop environment.

## Scope

- Desktop entry para lanzar Calamares desde el menú de Qtile
- Módulo TUI `setup-02-installer` que lanza Calamares
- Atajo de teclado opcional

## Requirements

1. Calamares lanzable desde el menú de aplicaciones de Qtile
2. Calamares lanzable desde la TUI lambda-env
3. Requiere sudo para ejecutar

## Technical Notes

- Desktop entry: ya viene con el paquete calamares
- Qtile: agregar al menú de rofi o keybinding
  - `Key([mod], "i", lazy.spawn("sudo calamares"))`
- TUI: módulo en `modules/setup/setup-02-installer`
  - Detecta si calamares está instalado
  - Si no está: "Installer no disponible en esta versión"
  - Si está: `sudo calamares`
- Verificar que se ejecuta desde live environment (no desde sistema instalado)

## Dependencies

- `installer-01-calamares-scaffold`
- `core/01-hub-plugin-system`

## Verification

- Menú de Qtile → "Install LambdaOS" → abre Calamares
- TUI → Setup → Installer → abre Calamares
- Atajo de teclado funciona
