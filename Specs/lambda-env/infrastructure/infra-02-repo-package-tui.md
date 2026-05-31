# lambda-env: Empaquetar TUI (infra-02)

## Intent

Empaquetar la TUI lambda-env como paquete pacman (`lambdaos-tui`) para distribución via repo propio.

## Scope

- Crear PKGBUILD para lambda-env
- Configurar build y install del paquete
- Publicar en repo pacman de LambdaOS
- Dependencias del paquete

## Requirements

1. PKGBUILD funcional que construye el paquete
2. Paquete instala binario `lambda-env` en `/usr/bin/`
3. Instala módulos en `/usr/share/lambda-env/modules/`
4. Instala config default en `/etc/lambdaos/`

## Technical Notes

- PKGBUILD:
  ```
  pkgname=lambdaos-tui
  pkgver=0.1.0
  pkgrel=1
  pkgdesc="TUI configuration suite for LambdaOS"
  arch=('x86_64')
  depends=('python' 'python-textual' 'jq')
  source=("$pkgname-$pkgver.tar.gz::https://github.com/...")
  ```
- Install file: `lambdaos-tui.install` → post-install hooks
- Files:
  - `/usr/bin/lambda-env`
  - `/usr/share/lambda-env/modules/*`
  - `/etc/lambdaos/settings.json` (config default)
  - `/usr/share/lambdaos/wallpapers/` (si aplica)

## Dependencies

- `infra-01-repo-pacman-setup`

## Verification

- `makepkg -s` → construye sin errores
- `pacman -U lambdaos-tui-*.pkg.tar.zst` → instala
- `lambda-env` → ejecuta correctamente
- `pacman -Q lambdaos-tui` → muestra instalado
