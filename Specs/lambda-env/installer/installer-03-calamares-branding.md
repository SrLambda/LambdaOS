# lambda-env: Calamares Branding (installer-03)

## Intent

Personalizar el branding de Calamares con logo, colores y textos de LambdaOS.

## Scope

- Crear package de branding de Calamares para LambdaOS
- Logo en la ventana de Calamares
- Colores consistentes con el tema de LambdaOS
- Textos personalizados (bienvenida, slides)
- Slides informativos durante la instalación

## Requirements

1. Calamares muestra branding de LambdaOS
2. Colores del tema Catppuccin (o configurable)
3. Slides informativos sobre LambdaOS durante la instalación
4. Logo visible en la ventana principal

## Technical Notes

- Branding package: `/usr/share/calamares/branding/lambdaos/`
  - `branding.desc` (config del branding)
  - `logo.png` (logo de LambdaOS)
  - `show.png` (imagen de presentación)
  - `slide*.html` (slides informativos)
- Colores: Catppuccin palette
  - Background: `#1e1e2e`
  - Text: `#cdd6f4`
  - Accent: `#cba6f7`
- Slides:
  1. Bienvenida a LambdaOS
  2. TUI-first design
  3. Qtile window manager
  4. Neovim integrado
  5. Herramientas de desarrollo incluidas
  6. Instalación en progreso

## Dependencies

- `installer-01-calamares-scaffold`
- `branding-02-wallpaper` (para assets)

## Verification

- Calamares abre con branding LambdaOS
- Logo visible
- Slides muestran durante instalación
- Colores consistentes
