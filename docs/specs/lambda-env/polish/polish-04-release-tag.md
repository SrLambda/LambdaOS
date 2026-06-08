# lambda-env: Release Tags & CI/CD (polish-04)

## Intent

Configurar flujo de releases con tags git y CI/CD para builds automáticos de la ISO.

## Scope

- Convención de versionado semántico
- Tags git para releases
- CI/CD pipeline para build automático de ISO
- Upload de ISO a GitHub Releases
- Changelog automático

## Requirements

1. Versionado semántico: `vMAJOR.MINOR.PATCH`
2. Tags git crean release automático
3. CI/CD build de ISO en GitHub Actions
4. ISO subida como artifact de release
5. Changelog generado desde commits

## Technical Notes

- Versionado: `v0.1.0`, `v0.2.0`, `v1.0.0`
- `profiledef.sh` ya soporta `LAMBDAOS_VERSION` env var y git tags
- GitHub Actions workflow:
  - Trigger: push de tag `v*`
  - Build: `sudo mkarchiso -v -w work/ -o out/ .`
  - Upload: ISO a GitHub Release
  - Checksum: sha256sum de la ISO
- Changelog: `git log --oneline` entre tags o con `git-cliff`
- Artifacts: ISO + sha256sum + changelog

## Dependencies

- `.github/workflows/` existente (ya hay CI/CD)
- `profiledef.sh` con soporte de versionado

## Verification

- Push tag `v0.1.0` → CI/CD se ejecuta
- ISO build completa
- ISO subida a GitHub Release
- Changelog generado
