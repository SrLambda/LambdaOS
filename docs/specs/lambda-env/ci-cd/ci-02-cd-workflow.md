# lambda-env: CD Workflow (ci-02)

## Intent

Configurar el workflow de CD (Continuous Deployment) que publica releases automáticamente.

## Scope

- Trigger: push de tag `v*`
- Build de ISO de release
- Generación de checksums (sha256)
- Generación de changelog automático
- Creación de GitHub Release
- Upload de ISO + checksums + changelog

## Requirements

1. Solo tags `v*` triggeran el release
2. ISO de release es build limpia (sin dirty flag)
3. Checksums publicados junto a la ISO
4. Changelog generado desde commits entre tags
5. Release publicado como draft primero (revisión manual antes de publicar)

## Technical Notes

- Workflow: `.github/workflows/cd.yml`
- Trigger: `on.push.tags: ['v*']`
- Steps:
  1. Checkout con fetch-depth: 0 (para git log completo)
  2. Install archiso dependencies
  3. Build ISO con `LAMBDAOS_VERSION` del tag
  4. Generate sha256sum: `sha256sum LambdaOS-*.iso > SHA256SUMS`
  5. Generate changelog: `git-cliff --tag <tag>` o `git log`
  6. Create GitHub Release (draft=true)
  7. Upload ISO + SHA256SUMS + CHANGELOG.md
- Herramientas:
  - `git-cliff` para changelog automático
  - `gh release create` para crear release
  - `gh release upload` para subir artifacts

### git-cliff config

```toml
# cliff.toml
[changelog]
body = """
{% for group, commits in commits | group_by(attribute="group") %}
### {{ group | upper_first }}
{% for commit in commits %}
- {{ commit.message | upper_first }} ({{ commit.id | truncate(length=7, end="") }})
{% endfor %}
{% endfor %}
"""
[git]
conventional_commits = true
filter_unconventional = true
commit_parsers = [
    { message = "^feat", group = "Features" },
    { message = "^fix", group = "Bug Fixes" },
    { message = "^chore", group = "Chores" },
    { message = "^docs", group = "Documentation" },
    { message = "^test", group = "Tests" },
]
```

## Dependencies

- `ci-01-ci-workflow`
- `polish-04-release-tag`

## Verification

- Push tag `v0.1.0` → CD ejecuta
- Release draft creado en GitHub
- ISO + SHA256SUMS + CHANGELOG.md subidos
- Release listo para publicación manual
