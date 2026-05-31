# lambda-env: Docs Content (infra-05)

## Intent

Crear el contenido de documentación de LambdaOS: guía de instalación, guía de TUI, FAQ, troubleshooting.

## Scope

- Guía de instalación paso a paso
- Guía de uso de la TUI lambda-env
- FAQ: preguntas frecuentes
- Troubleshooting: problemas comunes y soluciones
- Referencia de módulos de la TUI
- Guía de dotfiles y stow

## Requirements

1. Documentación completa y actualizada
2. Formato HTML para servir con darkhttpd
3. Markdown source en el repo
4. Index navegable con links entre secciones

## Technical Notes

- Source: Markdown en `docs/` del repo
- Build: pandoc o mkdocs para generar HTML
- Estructura:
  ```
  docs/
  ├── index.html
  ├── install/
  │   ├── index.html
  │   └── live-usb.html
  ├── tui/
  │   ├── index.html
  │   ├── modules.html
  │   └── settings.html
  ├── faq.html
  ├── troubleshooting.html
  └── dotfiles.html
  ```
- Estilo: CSS con tema Catppuccin
- Incluir screenshots o diagrams ASCII

## Dependencies

- `infra-04-docs-local` (infraestructura)

## Verification

- Docs accesibles en `http://localhost:8080`
- Todas las secciones navegable
- Links internos funcionan
- Contenido actualizado
