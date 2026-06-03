# lambda-env: Wallpaper Default (branding-02)

## Intent

Agregar un wallpaper default de LambdaOS a la ISO y configurarlo en Qtile.

## Scope

- Crear/incluir wallpaper de LambdaOS
- Agregar a la ISO en `/usr/share/lambdaos/wallpapers/`
- Configurar Qtile para usarlo como fondo de pantalla
- Agregar `feh` o `nitrogen` para gestión de wallpaper

## Requirements

1. Wallpaper incluido en la ISO
2. Qtile aplica el wallpaper en autostart
3. Wallpaper configurable desde TUI (system-09-appearance)

## Technical Notes

- Directorio: `airootfs/usr/share/lambdaos/wallpapers/default.png`
- Herramienta: `feh --bg-fill` (agregar a packages.x86_64)
- Qtile autostart: `subprocess.Popen(["feh", "--bg-fill", wallpaper_path])`
- Wallpaper path en `settings.json` bajo `appearance.wallpaper`
- Diseño del wallpaper: logo LambdaOS con fondo del tema Catppuccin

## Dependencies

- feh (agregar a packages.x86_64)
- `core/02-settings-schema` (para path en settings.json)

## Verification

- Wallpaper visible al iniciar Qtile
- `feh` instalado
- Path correcto en settings.json
