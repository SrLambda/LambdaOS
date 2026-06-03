# lambda-env: Nightly Builds (ci-03)

## Intent

Configurar builds nocturnos automáticos de la ISO para testing continuo sin necesidad de tags de release.

## Scope

- Trigger: schedule daily at 00:00 UTC
- Build de ISO con versión `nightly-YYYYMMDD`
- Upload como artifact (no como release)
- Limpieza de artifacts antiguos (mantener últimos 7)

## Requirements

1. Build se ejecuta automáticamente cada día
2. Versión de la ISO: `nightly-YYYYMMDD`
3. Artifacts disponibles por 7 días
4. Notificación de fallo (si el build falla)

## Technical Notes

- Workflow: `.github/workflows/nightly.yml`
- Trigger: `on.schedule: [{ cron: '0 0 * * *' }]`
- Version: `LAMBDAOS_VERSION=nightly-$(date +%Y%m%d)`
- Artifacts: GitHub Actions artifacts (auto-expire en 7 días)
- Notification: GitHub Actions email o webhook a Discord/Slack
- Limpiar artifacts viejos: action `c-hive/gha-remove-artifacts`

### Workflow

```yaml
name: Nightly Build
on:
  schedule:
    - cron: '0 0 * * *'
  workflow_dispatch:  # Manual trigger también

jobs:
  nightly:
    runs-on: ubuntu-latest
    steps:
      - checkout
      - install dependencies
      - build ISO with LAMBDAOS_VERSION=nightly-$(date +%Y%m%d)
      - upload artifact
      - cleanup old artifacts
```

## Dependencies

- `ci-01-ci-workflow`

## Verification

- Build ejecuta a medianoche UTC
- ISO artifact disponible con nombre `LambdaOS-nightly-YYYYMMDD-x86_64.iso`
- Artifacts viejos (>7 días) eliminados
- Manual trigger funciona
