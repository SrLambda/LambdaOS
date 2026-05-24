# Design: CI/CD Pipeline for LambdaOS

## Technical Approach

Add GitHub Actions workflows (`.github/workflows/`) for CI (push/PR) and release (tag push), a `Makefile` for local parity, and git-tag-based semver replacing date-based versioning in `profiledef.sh`. All additive — existing build scripts (`build_and_test.sh`, `run_qemu.sh`) remain untouched. Follows the project's pip-based dependency model (no pyproject.toml).

## Architecture Decisions

| Decision | Options | Tradeoff | Choice |
|----------|---------|----------|--------|
| Workflow split | One monolith vs. `ci.yml` + `release.yml` | Monolith: simpler but mixed triggers and longer config. Split: trigger-specific, faster CI, clearer rollback | **Split**: `ci.yml` (push/PR to main) and `release.yml` (tag push `v*`) |
| Python dep install | pip from `requirements-dev.txt` vs. pacman packages | Pacman is Arch-only, not portable. pip works on all runners and mirrors local dev | **pip**: `python -m pip install -r requirements-dev.txt` |
| Shell lint scope | All `.sh` vs. `scripts/` + root only | All catches everything; narrow is predictable | **All**: `shellcheck **/*.sh` and `shfmt -d **/*.sh` — CI is the gate |
| Pacman caching | actions/cache vs. no cache | No cache: build 30+ min every time. Cache: first run slow, subsequent ~10-15 min saved | **actions/cache** on `/var/cache/pacman/pkg` |
| Version source | env var → git tag → describe → short hash | Proposal asked for `GITHUB_REF_NAME` first. Local dev needs fallbacks | **Priority chain**: `LAMBDAOS_VERSION` env → `git describe --tags --exact-match` → `git describe --tags` → `git rev-parse --short HEAD` |
| Makefile build target | Require sudo vs. check and warn | Requiring sudo silently breaks local use; warn is friendlier | **Warn** if not root, then `sudo mkarchiso -v -w work/ -o out/ .` |

## Data Flow

**CI Pipeline (push/PR to main):**
```
  push/PR event
       │
       ├──→ checkout@v4
       │       │
       │       ├──→ setup-python (3.12+)
       │       │       ├── pip install -r requirements-dev.txt
       │       │       └── lint (parallel with tests)
       │       │             ├── black --check .
       │       │             ├── isort --check .
       │       │             ├── shellcheck **/*.sh
       │       │             └── shfmt -d **/*.sh
       │       │
       │       └──→ tests (parallel with lint)
       │               └── pytest tests/unit/ -v
       │
       └──→ status ✅/❌
```

**Release Pipeline (tag push v*):**
```
  tag push v1.0.0
       │
       └──→ checkout@v4
             │
             ├──→ install archiso + base-devel (pacman)
             ├──→ cache pacman pkg dir
             ├──→ export SOURCE_DATE_EPOCH=$(git log -1 --format=%ct $TAG)
             ├──→ sudo mkarchiso -v -w work/ -o out/ .
             ├──→ sha256sum out/*.iso > out/checksums.sha256
             └──→ softprops/action-gh-release@v2
                   └── files: out/*.iso, out/checksums.sha256
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `.github/workflows/ci.yml` | Create | CI workflow: setup → lint (black/isort/shellcheck/shfmt) + test (pytest unit) in parallel |
| `.github/workflows/release.yml` | Create | Release workflow: install deps → build ISO → checksums → GitHub Release, 45 min timeout |
| `Makefile` | Create | Targets: `lint`, `test`, `build`, `release`, `clean`. Mirrors CI commands |
| `profiledef.sh` | Modify | Replace `date`-based `iso_version` and `iso_label` with git-tag-driven version logic |
| `docs/BRANCHING.md` | Create | Git Flow light strategy: main/develop/feature/release/hotfix branches |

## Interfaces / Contracts

**Makefile targets** (consistent with existing `build_and_test.sh` conventions):
```makefile
.PHONY: lint test build release clean

lint:
	black --check . && isort --check . && shellcheck **/*.sh && shfmt -d **/*.sh

test:
	python -m pytest tests/unit/ -v

build:
	@[ "$$(id -u)" -eq 0 ] || { echo "Need sudo for mkarchiso"; exit 1; }
	mkarchiso -v -w work/ -o out/ .

clean:
	sudo rm -rf work/ out/ .venv/
```

**Version resolution contract** (in `profiledef.sh`):
```bash
if [[ -n "${LAMBDAOS_VERSION:-}" ]]; then
    iso_version="${LAMBDAOS_VERSION}"
elif tag=$(git describe --tags --exact-match 2>/dev/null); then
    iso_version="${tag#v}"
else
    iso_version="$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
fi
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit (existing) | 68 tests in `tests/unit/` — TUI configurator, pacman cfg, Qtile AST, aur-packages | `pytest tests/unit/ -v` runs in CI as gate; unchanged by this change |
| Workflow validation | CI triggers correctly on push/PR; release triggers on tag; Makefile targets work | Manual: push to branch, open PR, push tag. No automated workflow testing framework added (out of scope) |
| Version resolution | `profiledef.sh` version logic for env var, tag, describe, hash fallbacks | Shell script unit test (add to `test` target or verify locally) |
| E2E (existing) | QEMU boot — `tests/qemu/test_live_boot.py` | Local only; explicitly excluded from CI (per exploration) |

## Migration / Rollout

No migration required. Rollback per proposal: delete `.github/workflows/`, revert `profiledef.sh` line 8 to `date` command, remove `Makefile` and `docs/BRANCHING.md`. No airootfs overlay implications — version changes are metadata only.

## Open Questions

- [ ] Tag format for pre-releases: propose `v1.0.0-beta.1` blocked by release `v*` trigger. Should we use `-` suffix tags as pre-release convention with manual release upload until semver pre-release support is added?
- [ ] GitHub Release already-exists behavior: fail (strict) vs skip (idempotent). Spec says fail; confirm this with team.
