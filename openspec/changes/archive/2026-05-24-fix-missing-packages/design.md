# Design: Fix Missing Packages

## Technical Approach

Four independent, low-risk file changes that close the gap between README promises and ISO content. No cross-dependencies between changes — each spec maps to a single file. pacman.conf and packages.x86_64 are build-time; aur-packages.sh and README are documentation/runtime.

> **Delivery strategy**: `single-pr` (total delta ~55 lines across all files, well under 400-line budget).

## Architecture Decisions

### Decision: AUR script location

**Choice**: `scripts/aur-packages.sh` in project root (not airootfs/)
**Alternatives**: `airootfs/root/` (copied into ISO but increases image size and complexity), `airootfs/usr/local/bin/` (convention for live-ISO tools but this is a post-install script, not a live tool)
**Rationale**: The script is for users AFTER booting the installed system, not for the live ISO environment. Keeping it in the repo root as documentation-and-tool hybrid makes it discoverable without bloating the airootfs overlay. Users clone the repo or read it from GitHub.

### Decision: AUR helper detection strategy

**Choice**: Check `yay` first, fall back to `paru`. Print install instructions if neither found.
**Alternatives**: Hardcode a single helper, require the user to export `AUR_HELPER`
**Rationale**: yay and paru cover 95%+ of AUR helper users. A simple `command -v` check is more user-friendly than requiring env var setup. The script exits with install instructions — does not auto-install helpers (security boundary).

### Decision: Package list organization

**Choice**: Insert new packages in alphabetical order with category comments
**Alternatives**: Append to end, group by function irrespective of sort
**Rationale**: The existing 172-line file is strictly alphabetical. Maintaining sort prevents merge conflicts in future changes and matches the file's existing convention exactly.

## Data Flow

```
build time:
  pacman.conf ──[multilib uncommented]──→ mkarchiso ──→ ISO
  packages.x86_64 ──[+15 packages]──────→ mkarchiso ──→ ISO

post-boot (user triggered):
  README.md ──→ User reads AUR section
                   │
                   ├── 1. Install yay/paru
                   ├── 2. Run scripts/aur-packages.sh
                   └── 3. AUR packages installed
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `pacman.conf` | Modify | Uncomment `[multilib]` section (lines 93-94) |
| `packages.x86_64` | Modify | Add 15 official packages at sorted positions |
| `scripts/aur-packages.sh` | Create | Post-boot AUR install script with helper detection |
| `README.md` | Modify | Separate ISO-included vs AUR-post-install packages |

## Interfaces / Contracts

### aur-packages.sh interface

```bash
# Exit codes: 0=all installed, 1=missing AUR helper, 2=partial failures
./scripts/aur-packages.sh
```

No arguments. Reads no environment variables. Detects AUR helper via `command -v`. Installs each package with `--needed --noconfirm`, logs per-package results, prints summary.

### packages.x86_64 additions (sorted positions)

```text
chromium        # after chntpw, before clang
docker          # after dnsmasq, before dosfstools
docker-compose  # after docker, before dosfstools
keepassxc       # after jfsutils, before kitty
lazydocker      # after lazydocker line removed (none exists); after lazygit
libreoffice-fresh # after less, before lftp
okular          # after nvme-cli, before open-iscsi
qalculate-gtk   # after pv, before python
steam           # after squashfs-tools, before stow
tailscale       # after systemd-resolvconf, before tcpdump
thunderbird     # after terminus-font, before testdisk
virtualbox      # after usbmuxd, before usbutils
vlc             # after virtualbox-guest-utils-nox, before vpnc
wine            # after wget, before wireless-regdb
wine-mono       # after wine, before wireless-regdb
winetricks      # after wine-mono, before wireless-regdb
yazi            # after xorg-xsetroot, before zathura
```

Note: `docker-compose` is a separate package from `docker` in Arch repos. `wine-mono` is a dependency of wine but listed explicitly per spec.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Shell | aur-packages.sh syntax | `shellcheck scripts/aur-packages.sh` (already installed) |
| Shell | aur-packages.sh AUR helper detection | Unit test: mock `command -v`, assert exit codes per scenario |
| Config | pacman.conf multilib uncommented | grep-based validation: assert `[multilib]` not prefixed with `#` |
| Config | packages.x86_64 contains all additions | grep-based validation: assert each package name present |
| Config | packages.x86_64 maintains sort | `diff <(cat packages.x86_64) <(sort packages.x86_64)` returns empty |
| E2E | ISO builds with new packages | `build_and_test.sh` (existing pipeline) |

## Migration / Rollout

No migration required. This is a greenfield addition to an ISO definition. Users rebuild the ISO to get new packages. The AUR script is opt-in.

Rollback: revert 4 files — no state, no data migration, no cascading changes.

## Open Questions

- [ ] Does `wine-mono` need to be listed separately, or does `wine` pull it as a dependency? (Proposal lists all three wine packages explicitly — verify with `pacman -Si wine` before finalizing)
- [ ] ISO size impact: LibreOffice (~500MB compressed) + Chromium (~100MB) — monitor build output; may need to switch to `firefox` + `abiword` if ISO exceeds 8GB target
