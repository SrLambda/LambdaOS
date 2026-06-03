## Exploration: ci-01-ci-workflow

### Current State

The CI workflow at `.github/workflows/ci.yml` currently has **2 jobs**:
- **lint** — runs black, isort, shellcheck, shfmt (Python + Shell only)
- **test** — runs pytest on `tests/unit/`

A separate `.github/workflows/release.yml` exists (from polish-04) that builds the ISO and creates a GitHub release on tag pushes (`v*`). It uses `archlinux:latest` container with `--privileged` and runs `mkarchiso`.

**What exists vs what the spec requires:**

| Job | Spec Requires | Current CI | Gap |
|-----|--------------|------------|-----|
| lint: black | Yes | Yes | None |
| lint: isort | Yes | Yes | None |
| lint: shellcheck | Yes | Yes | None |
| lint: shfmt | Yes | Yes | None |
| lint: luacheck | Yes | **Missing** | Need to add |
| test-unit: pytest | Yes | Yes | None (job named `test`, spec says `test-unit`) |
| validate-specs | Yes | **Missing** | Need to add |
| build-iso (main only) | Yes | **Missing** | Exists in release.yml but tag-only, not main-push |

### Affected Areas

- `.github/workflows/ci.yml` — main target; needs 3 new jobs + luacheck step
- `requirements-dev.txt` — needs `luacheck` (or install via apt/npm)
- `Specs/` directory — 48 spec files across 10 subdirectories; validation script needed
- `Makefile` — currently mirrors lint/test but missing luacheck; could be updated for local parity
- `airootfs/etc/skel/dotfiles/nvim/.config/nvim/lua/` — 21 Lua files that luacheck would lint

### Approaches

#### 1. **Extend ci.yml in-place (Recommended)**
Add luacheck step to existing lint job, add validate-specs and build-iso as new jobs. Reuse the Arch Linux container pattern from release.yml for the ISO build.

- **Pros**: Single source of truth for CI, clear job dependencies, minimal duplication
- **Cons**: ci.yml grows from 53 to ~120 lines
- **Effort**: Low

#### 2. **Split into separate workflow files**
Keep ci.yml for lint+test, create `build-iso.yml` for main-only builds, `validate-specs.yml` as independent check.

- **Pros**: Each file is focused, easier to reason about individually
- **Cons**: More files to maintain, harder to see full CI picture at a glance
- **Effort**: Low

#### 3. **Reuse release.yml ISO build logic**
Move the ISO build steps into a reusable composite action or workflow_call, then both ci.yml and release.yml reference it.

- **Pros**: DRY, single place to update ISO build logic
- **Cons**: Over-engineering for a project this size; adds indirection
- **Effort**: Medium

### Recommendation

**Approach 1** — extend ci.yml in-place. The project is small enough that a single workflow file is the right complexity level. The release.yml can keep its tag-triggered release logic (changelog + GitHub release) while ci.yml handles the main-push ISO build. The ISO build artifact in CI would be uploaded as a workflow artifact (not a GitHub release), which is the correct distinction.

Key implementation details:
- **luacheck**: Install via `luarocks` or `apt` (`apt install lua-check`). The 21 Lua files are Neovim config files under `airootfs/.../nvim/lua/` — standard Lua 5.1 with Neovim globals. Need to configure luacheck for Neovim globals (`vim`).
- **validate-specs**: Create a simple script (e.g., `scripts/validate-specs.sh`) that checks:
  - All spec files exist under `Specs/`
  - Required frontmatter/headers present (e.g., `# `, `## ` sections)
  - Directory structure matches expected categories
- **build-iso**: Use `if: github.ref == 'refs/heads/main'` condition. Use `archlinux:latest` container with `--privileged` (same as release.yml). Upload ISO as workflow artifact via `actions/upload-artifact@v4`. Set timeout to 60 minutes per spec.

### Risks

- **ISO build time on GitHub runners**: Arch Linux ISO builds can take 15-30 minutes. The 60-minute timeout should be sufficient but monitor actual times.
- **Luacheck Neovim globals**: The Lua files use `vim.*` extensively. Luacheck will flag these as undefined globals unless configured with `globals = { vim = true }` in `.luacheckrc`.
- **Privileged container for ISO build**: The `--privileged` flag is required for mkarchiso (chroot/mount operations). GitHub Actions supports this but some enterprise environments restrict it.
- **Spec validation scope**: The spec says "validar Specs/ estructura" but doesn't define what "valid" means. Need clarification on validation rules (frontmatter? naming convention? required sections?).

### Ready for Proposal

**Yes.** The codebase has sufficient structure to implement this change. The main unknown is the exact validation rules for specs — the proposal phase should define what "valid" means for the Specs/ directory (minimal validation: file existence + markdown headers; or stricter: required sections per spec type).
