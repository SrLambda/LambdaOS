# lambda-env: Docs Local (infra-04)

## Intent

Servir documentación local en la ISO live usando `darkhttpd` (ya incluido en packages.x86_64).

## Scope

- Crear estructura de documentación HTML
- Configurar darkhttpd para servir en la ISO
- Auto-start del servidor al bootear
- Acceso via browser (Chromium)

## Requirements

1. Documentación accesible en `http://localhost:8080` en la ISO live
2. darkhttpd auto-start al bootear
3. Documentación incluye: guía de instalación, guía de TUI, FAQ, troubleshooting
4. Desktop entry para abrir docs en el browser

## Technical Notes

- Directorio: `/usr/share/lambdaos/docs/`
- Contenido: HTML estático (Markdown convertido con pandoc o hand-written)
- darkhttpd: `darkhttpd /usr/share/lambdaos/docs/ --port 8080`
- Systemd service: `lambdaos-docs.service` → auto-start
- Desktop entry: "LambdaOS Documentation" → abre `http://localhost:8080` en Chromium
- darkhttpd ya está en packages.x86_64

## Dependencies

- `infra-05-docs-content` (contenido)
- darkhttpd (ya en packages.x86_64)

## Verification

- Boot ISO → `http://localhost:8080` → muestra docs
- Desktop entry abre browser con docs
- darkhttpd service running
