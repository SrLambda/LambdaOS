# Branching Strategy — Git Flow Light

LambdaOS follows a **Git Flow Light** branching model. It is a simplified adaptation of [nvie/gitflow](https://nvie.com/posts/a-successful-git-branching-model/) that keeps the ceremony low while preserving release stability.

## Branches

| Branch | Purpose | Lifespan | Merge target |
|--------|---------|----------|--------------|
| `main` | Production-ready code. Every commit is deployable. | Permanent | — |
| `develop` | Integration branch for the next release. | Permanent | `main` |
| `feature/*` | New features or non-trivial improvements. | Ephemeral | `develop` |
| `release/*` | Release preparation (version bump, changelog, final QA). | Ephemeral | `main` + `develop` |
| `hotfix/*` | Critical fixes for `main` that cannot wait for the next release cycle. | Ephemeral | `main` + `develop` |

## Workflow

### 1. Daily development

```
main ← develop ← feature/awesome-thing
```

1. Create a feature branch from `develop`:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/short-description
   ```
2. Open PRs against `develop`.
3. CI must pass (lint + unit tests) before merge.
4. Delete the feature branch after merge.

### 2. Releasing

```
develop → release/v1.2.0 → main  → tag v1.2.0
                       ↘ develop
```

1. When `develop` is ready, create `release/vX.Y.Z` from `develop`.
2. Apply only fixes on the release branch (no new features).
3. Merge `release/vX.Y.Z` into `main` **and** `develop`.
4. Tag `main` with `vX.Y.Z`. The release workflow triggers automatically.

### 3. Hotfixes

```
main → hotfix/critical-fix → main  → tag vX.Y.(Z+1)
                        ↘ develop
```

1. Create `hotfix/*` directly from `main`.
2. Apply the minimal fix.
3. Merge into `main` and `develop`.
4. Tag `main` immediately.

## Tagging Convention

- Format: `vMAJOR.MINOR.PATCH` (e.g., `v1.0.0`, `v1.2.3`)
- Pre-releases: `v1.0.0-beta.1` (supported by the release trigger `v*`)
- The `v` prefix is stripped by the ISO build to produce a clean semver string.

## CI / Release Integration

- Every push or PR to `main` triggers the **CI workflow** (lint + unit tests).
- Every push of a tag matching `v*` triggers the **Release workflow** (ISO build + GitHub Release).

## Why not trunk-based?

LambdaOS is an OS distribution. Releases are ISO images that users download and flash. A broken release is expensive (users download GBs of bad data). The `main`/`develop` split gives us a staging area to validate the full archiso build before declaring a release ready.
