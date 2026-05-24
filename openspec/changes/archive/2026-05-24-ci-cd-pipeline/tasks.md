# Tasks: CI/CD Pipeline for LambdaOS

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~260 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-always |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

## Phase 1: Foundation — Makefile + Versioning

- [x] 1.1 Create `Makefile` with `lint`, `test`, `build`, `clean` targets mirroring CI commands
- [x] 1.2 Write shell test for version resolution: `tests/unit/test_versioning.sh` — env var → tag → describe → hash fallback chain (TDD RED)
- [x] 1.3 Modify `profiledef.sh` — replace `date`-based `iso_version`/`iso_label` with `LAMBDAOS_VERSION` → git tag → short hash priority (TDD GREEN)

## Phase 2: CI Pipeline

- [x] 2.1 Create `.github/workflows/ci.yml` — trigger on push/PR to main, Python setup, pip install
- [x] 2.2 Add parallel lint job: `black --check`, `isort --check`, `shellcheck **/*.sh`, `shfmt -d **/*.sh`
- [x] 2.3 Add parallel test job: `pytest tests/unit/ -v` — aggregate status reflects all jobs

## Phase 3: Release Pipeline

- [x] 3.1 Create `.github/workflows/release.yml` — trigger on tag `v*`, install `archiso` + `base-devel`
- [x] 3.2 Add pacman cache with `actions/cache` for `/var/cache/pacman/pkg` keyed by hash
- [x] 3.3 Add ISO build: `SOURCE_DATE_EPOCH` from tag commit, `sudo mkarchiso`, 45 min timeout
- [x] 3.4 Add `sha256sum` generation + `softprops/action-gh-release@v2` asset upload (ISO + checksums)

## Phase 4: Documentation & Verification

- [x] 4.1 Create `docs/BRANCHING.md` — Git Flow light: main/develop/feature/release/hotfix branches
- [x] 4.2 Run `pytest tests/unit/ -v` — confirm all existing tests pass (no regressions)
- [x] 4.3 Verify Makefile targets locally: `make lint`, `make test`, `make clean`
