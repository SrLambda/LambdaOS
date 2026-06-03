# lambda-env: Repo Pacman Setup (infra-01)

## Intent

Crear la estructura de un repositorio pacman local para LambdaOS y configurar `pacman.conf` para usarlo.

## Scope

- Crear estructura de repo: directorio, signing key, base de datos
- Configurar `pacman.conf` para incluir el repo de LambdaOS
- Documentar cómo agregar paquetes al repo
- Script para actualizar el repo

## Requirements

1. Repo pacman funcional con estructura correcta
2. `pacman.conf` configurado para usar el repo
3. Script para agregar paquetes y regenerar la base de datos
4. Signing key para paquetes del repo

## Technical Notes

- Estructura del repo:
  ```
  /srv/repo/lambdaos/
  ├── lambdaos.db.tar
  ├── lambdaos.files.tar
  └── x86_64/
      ├── lambdaos-tui-0.1.0-1-x86_64.pkg.tar.zst
      └── lambdaos-tui-0.1.0-1-x86_64.pkg.tar.zst.sig
  ```
- `pacman.conf`:
  ```
  [lambdaos]
  Server = https://repo.lambdaos.dev/$arch
  SigLevel = Required
  ```
- Herramientas: `repo-add`, `gpg` para signing
- Script: `scripts/repo-update.sh` → `repo-add lambdaos.db.tar *.pkg.tar.zst`
- Para la ISO live: repo local en la ISO o mirror interno

## Dependencies

- Ninguno (pero se usa en Fase 5.2 y 5.3)

## Verification

- `pacman -Sl lambdaos` → lista paquetes del repo
- `pacman -S lambdaos-tui` → instala desde el repo
- Signing verifica correctamente
