# lambda-env: CI Workflow (ci-01)

## Intent

Configurar el workflow de CI (Continuous Integration) que se ejecuta en cada push y PR.

## Scope

- Linting: shellcheck, shfmt, black, luacheck
- Tests unitarios: pytest
- Validación de specs: estructura y formato
- Build de ISO (solo en main, no en PRs para ahorrar recursos)

## Requirements

1. CI bloquea merge si lint o tests fallan
2. Resultados visibles en GitHub Checks
3. Build de ISO solo en main (no en PRs)
4. Artifacts de build disponibles para descarga

## Technical Notes

- Workflow: `.github/workflows/ci.yml`
- Jobs:
  - `lint` → shellcheck, shfmt, black, isort, luacheck
  - `test-unit` → pytest tests/unit/
  - `validate-specs` → validar Specs/ estructura
  - `build-iso` → mkarchiso (solo en main)
- Dependencias entre jobs: lint + test → build-iso
- Timeout: 60 minutos para build-iso
- Runner: ubuntu-latest con qemu/kvm support

## Dependencies

- Ninguno

## Verification

- PR con error de lint → CI falla, no se puede mergear
- Push a main → CI ejecuta todo incluyendo build
- ISO artifact descargable
