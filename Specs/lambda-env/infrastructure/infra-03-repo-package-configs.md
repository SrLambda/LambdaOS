# lambda-env: Empaquetar Configs (infra-03)

## Intent

Empaquetar las configuraciones base de LambdaOS como paquete pacman (`lambdaos-configs`) para distribución via repo propio.

## Scope

- Crear PKGBUILD para lambdaos-configs
- Incluir: Qtile config, Neovim config, Kitty config, Yazi config, dotfiles
- Configurar como paquete de configuración (no sobrescribe user configs)
- Publicar en repo pacman de LambdaOS

## Requirements

1. PKGBUILD funcional que construye el paquete
2. Instala configs en `/etc/skel/` para nuevos usuarios
3. No sobrescribe configs existentes de usuarios
4. Incluye todos los dotfiles del skel actual

## Technical Notes

- PKGBUILD:
  ```
  pkgname=lambdaos-configs
  pkgver=0.1.0
  pkgrel=1
  pkgdesc="Default configuration files for LambdaOS"
  arch=('any')
  ```
- Files:
  - `/etc/skel/.config/qtile/*`
  - `/etc/skel/.config/nvim/*`
  - `/etc/skel/.config/kitty/*`
  - `/etc/skel/.config/yazi/*`
  - `/etc/skel/dotfiles/*`
  - `/etc/skel/.zprofile`
  - `/etc/skel/.bash_profile`
- Usar `.pacnew` mechanism para no sobrescribir configs existentes
- Package type: `any` (no architecture-specific)

## Dependencies

- `infra-01-repo-pacman-setup`

## Verification

- `makepkg -s` → construye sin errores
- `pacman -U lambdaos-configs-*.pkg.tar.zst` → instala
- Nuevos usuarios reciben configs default
- Usuarios existentes no pierden sus configs
