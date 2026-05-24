# Proposal: Corregir Paquetes Faltantes

## Intent

The README advertises features (VPN, office suite, browser, gaming, Docker, etc.) that are NOT in packages.x86_64. The ISO builds without them, producing a system that doesn't match its own documentation. We must close this gap by adding official packages, enabling multilib for 32-bit gaming, and creating a clear AUR install path for packages archiso cannot deliver.

## Scope

### In Scope
- Uncomment `[multilib]` repo in `pacman.conf`
- Add all missing official packages to `packages.x86_64` (tailscale, libreoffice-fresh, thunderbird, vlc, okular, qalculate-gtk, keepassxc, chromium, yazi, virtualbox, docker, docker-compose, lazydocker, steam, wine, wine-mono, winetricks)
- Create `scripts/aur-packages.sh` with documented AUR installs (spotify, obsidian, megasync, bluetui, impala)
- Update README.md to clearly label what's in-ISO vs what's post-install AUR

### Out of Scope
- Automating AUR installation into the ISO (archiso constraint)
- GUI AUR helper installation (user installs yay/paru post-boot)
- Desktop entries, theming, or configuration for new packages
- Snapper setup (separate change)

## Capabilities

### New Capabilities
- `aur-install-script`: Script to install AUR-only packages post-boot with clear instructions

### Modified Capabilities
- `iso-packages`: Package list expanded with 16+ official packages including multilib
- `iso-pacman-config`: Multilib repository enabled for 32-bit compatibility

## Approach

1. **pacman.conf**: Uncomment lines 93-94 (`[multilib]` and its Include). Simple 2-line edit.
2. **packages.x86_64**: Add missing official packages in alphabetical order (maintain sort). Group by purpose with comments.
3. **aur-packages.sh**: Create self-documenting shell script with `# AUR packages` header, install commands using `yay -S`, and echo instructions.
4. **README.md**: Add section marking AUR packages as "Post-installation" with a reference to the script.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `pacman.conf` | Modified | Uncomment multilib repo (2 lines) |
| `packages.x86_64` | Modified | Add ~16 packages (~16 lines) |
| `scripts/aur-packages.sh` | New | AUR install script (~30 lines) |
| `README.md` | Modified | Add AUR/post-install section (~15 lines) |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| ISO size grows significantly with LibreOffice + Chromium | Medium | Monitor ISO size; consider lighter alternatives (firefox, abiword) |
| Multilib conflicts with existing x86_64 packages | Low | archiso handles multilib natively; test build required |
| AUR packages break on install | Medium | Script uses `yay --needed` and includes per-package error handling |
| Package removed from official repos | Low | Verify in archlinux.org/packages before adding; CI can catch |

## Rollback Plan

1. `pacman.conf`: Re-comment the 2 multilib lines.
2. `packages.x86_64`: Remove added packages (git revert the file).
3. `scripts/aur-packages.sh`: Delete the file.
4. `README.md`: Git revert changes.

All changes are in 4 files with no dependency cascades — clean revert via `git revert`.

## Dependencies

- archiso build environment (existing)
- No new external dependencies (all packages are in official repos or documented as AUR)

## Success Criteria

- [ ] `pacman.conf` has `[multilib]` uncommented with Include line
- [ ] `packages.x86_64` contains: tailscale, libreoffice-fresh, thunderbird, vlc, okular, qalculate-gtk, keepassxc, chromium, yazi, virtualbox, docker, docker-compose, lazydocker, steam, wine
- [ ] `scripts/aur-packages.sh` exists and documents: spotify, obsidian, megasync, bluetui, impala
- [ ] README distinguishes ISO-included packages from AUR post-install packages
- [ ] ISO builds successfully with `mkarchiso`