# lambda-env: CI/CD Configuration

## Intent

Configurar el pipeline de CI/CD completo para LambdaOS: build automático de ISO, tests, linting, y deployment de releases.

## Scope

- CI: linting, tests unitarios, validación de specs, build de ISO
- CD: publicación de releases, upload de ISO, generación de changelog
- Artifacts: ISO, checksums, logs de build
- Notificaciones: estado del pipeline

## Requirements

1. CI se ejecuta en cada push a main y en PRs
2. CD se ejecuta solo en tags `v*`
3. Tests unitarios pasan antes de build
4. Build de ISO en runner con soporte de virtualización
5. ISO publicada en GitHub Releases con checksums
6. Changelog automático generado desde commits

## Technical Notes

### Pipeline CI (en cada push/PR)

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - checkout
      - shellcheck scripts/*.sh build_and_test.sh
      - shfmt check
      - python black/isort check
      - lua check (luacheck)

  test-unit:
    runs-on: ubuntu-latest
    steps:
      - checkout
      - setup python
      - pip install -r requirements-dev.txt
      - pytest tests/unit/ -v

  validate-specs:
    runs-on: ubuntu-latest
    steps:
      - checkout
      - validate Specs/ structure
      - validate markdown format
      - check for broken links/references

  build-iso:
    runs-on: ubuntu-latest
    needs: [lint, test-unit]
    steps:
      - checkout
      - install archiso dependencies
      - sudo mkarchiso -v -w work/ -o out/ .
      - upload artifact: ISO
```

### Pipeline CD (en tags v*)

```yaml
name: CD
on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - checkout
      - build ISO
      - generate sha256sum
      - generate changelog (git-cliff)
      - create GitHub Release
      - upload ISO + checksum + changelog
```

### Estructura de archivos CI/CD

```
.github/
├── workflows/
│   ├── ci.yml              ← CI pipeline
│   ├── cd.yml              ← CD pipeline
│   └── nightly.yml         ← Build nightly (opcional)
├── ISSUE_TEMPLATE/
│   ├── bug.md
│   └── feature.md
└── PULL_REQUEST_TEMPLATE.md
```

### Nightly Builds (opcional)

- Trigger: schedule daily at 00:00 UTC
- Build ISO con versión `nightly-YYYYMMDD`
- Upload como artifact (no como release)
- Útil para testing continuo

### Release Process

1. Developer crea tag: `git tag v0.1.0 && git push origin v0.1.0`
2. CI/CD se ejecuta automáticamente
3. ISO se build y sube a GitHub Release
4. Changelog se genera automáticamente
5. Release publicado con notas

## Dependencies

- `polish-04-release-tag` (versionado)
- GitHub Actions (ya configurado en `.github/`)

## Verification

- Push a main → CI ejecuta lint + tests + build
- PR → CI ejecuta lint + tests
- Push tag `v0.1.0` → CD ejecuta release
- ISO disponible en GitHub Releases
- Checksums correctos
- Changelog generado
