# Tasks: Fix Missing Packages

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~140-170 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-always |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

## Phase 1: Multilib Config

- [x] 1.1 Write grep-based test: `[multilib]` section uncommented in `pacman.conf`
- [x] 1.2 Uncomment lines 93-94 in `pacman.conf` (`[multilib]` header + Include line)

## Phase 2: Official Package List

- [x] 2.1 Write grep test: all 17 packages present in `packages.x86_64`
- [x] 2.2 Write sort-order test: file matches `sort` output
- [x] 2.3 Insert 17 packages at correct alphabetical positions in `packages.x86_64`

## Phase 3: AUR Install Script

- [x] 3.1 Write shellcheck integration test for `scripts/aur-packages.sh`
- [x] 3.2 Write unit test: AUR helper detection (mock `command -v`, assert exit code 0 vs 1)
- [x] 3.3 Write unit test: per-package error handling + `--needed` idempotency (exit 0 on re-run)
- [x] 3.4 Create `scripts/aur-packages.sh` — yay/paru fallback, per-package error handling, summary

## Phase 4: README Documentation

- [x] 4.1 Label official packages as "included in ISO" in README
- [x] 4.2 Add "requires post-boot AUR install" section with script reference and copy-paste instructions

## Phase 5: Integration Verification

- [x] 5.1 Run full test suite (unit + grep validations)
- [x] 5.2 `shellcheck` all scripts, verify changed-line count under budget
