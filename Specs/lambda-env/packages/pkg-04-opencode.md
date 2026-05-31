# lambda-env: Agregar OpenCode (pkg-04)

## Intent

Investigar y agregar OpenCode (agente IA CLI) a la distribución, ya sea en la ISO directa o como instalación post-boot.

## Scope

- Investigar método de distribución: npm global, cargo, binario precompilado, o empaquetado pacman
- Agregar al método de distribución elegido
- Documentar configuración básica

## Requirements

1. OpenCode disponible en LambdaOS (ISO o post-boot)
2. Método de instalación documentado
3. Config básica funcional

## Technical Notes

- OpenCode no está en repositorios oficiales de Arch ni en AUR (verificar al momento de implementar)
- Opciones de distribución:
  - **npm global**: `npm install -g opencode` → requiere Node.js (ya en ISO)
  - **cargo**: `cargo install opencode` → requiere Rust (no en ISO, agregar `rustup`)
  - **Binario precompilado**: descargar de GitHub releases
  - **Empaquetado propio**: crear PKGBUILD y repo pacman local (Fase 5)
- Recomendación v1.0: npm global (Node.js ya está en la ISO)
- Recomendación v2.0: empaquetar en repo pacman local

## Dependencies

- nodejs (ya en packages.x86_64) si se usa npm
- rustup (agregar) si se usa cargo

## Verification

- `opencode --version` → responde con versión
- `opencode` → abre sin errores
- Config básica en `~/.config/opencode/`
