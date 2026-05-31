# lambda-env: Icon & Cursor Theme (branding-04)

## Intent

Agregar y configurar un icon theme y cursor theme por defecto en LambdaOS.

## Scope

- Seleccionar icon theme (Papirus recomendado)
- Seleccionar cursor theme (Bibata Modern Ice recomendado)
- Agregar a packages.x86_64
- Configurar como default en GTK y Qtile

## Requirements

1. Icon theme instalado y configurado como default
2. Cursor theme instalado y configurado como default
3. Persiste en settings de GTK y X11

## Technical Notes

- Icon theme: `papirus-icon-theme` (repositorio extra de Arch)
- Cursor theme: `bibata-cursor-theme` (AUR) o `phinger-cursors` (extra)
- GTK config: `~/.config/gtk-3.0/settings.ini` → `gtk-icon-theme-name=Papirus`
- Cursor config: `~/.icons/default/index.theme` → `Inherits=Bibata-Modern-Ice`
- X11 cursor: `~/.Xresources` → `Xcursor.theme: Bibata-Modern-Ice`
- Si cursor está en AUR, agregar al script post-boot

## Dependencies

- papirus-icon-theme (agregar a packages.x86_64)
- phinger-cursors o bibata-cursor-theme (agregar)

## Verification

- Iconos visibles en apps GTK
- Cursor visible y consistente
- Config persiste entre sesiones
