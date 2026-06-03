# lambda-env: MOTD Personalizado (branding-01)

## Intent

Reemplazar el MOTD genérico de Arch Linux por un MOTD personalizado de LambdaOS con logo ASCII, versión, y links útiles.

## Scope

- Crear MOTD con logo ASCII de LambdaOS
- Mostrar versión de la ISO (desde git tag/describe)
- Links a documentación y repo
- Colores ANSI consistentes con el tema

## Requirements

1. Reemplazar `airootfs/etc/motd` con contenido custom
2. Logo ASCII reconocible de LambdaOS
3. Versión dinámica o placeholder que se resuelve en build time
4. Links: docs, GitHub, wiki

## Technical Notes

- Archivo: `airootfs/etc/motd`
- El MOTD actual es el default de Arch Linux (sin personalizar)
- Versión: usar `LAMBDAOS_VERSION` env var o git describe en profiledef.sh
- Colores ANSI: usar los del tema Catppuccin (default)
- Tamaño: máximo 20 líneas para no saturar la terminal

## Dependencies

- Ninguno

## Verification

- Boot de la ISO → MOTD muestra logo LambdaOS
- Versión visible
- Links correctos
