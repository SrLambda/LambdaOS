# Exploration: fix-missing-packages

## Current State of packages.x86_64

172 packages currently installed. Key categories present:
- **Base system**: base, base-devel, linux, linux-firmware, systemd
- **WM/Display**: qtile, xorg-server, ly, rofi, picom, dunst, kitty
- **Development**: neovim, git, lazygit, go, nodejs, npm, python, clang, jdk-openjdk, shellcheck, shfmt, stylua
- **Audio**: pipewire, wireplumber, pipewire-pulse, pipewire-alsa, pipewire-jack
- **Networking**: iw, iwd, wpa_supplicant, dhcpcd, openvpn, openssh, curl, wget
- **Tools**: fd, ripgrep, zathura, stow, tmux, nano, mc, nmap

## Missing Packages (README vs Reality)

| README Feature | Package | In pacman repos? | Status |
|---|---|---|---|
| VPN | `tailscale` | ✅ [core] | **MISSING from packages.x86_64** |
| Ofimática | `libreoffice-fresh` | ✅ [extra] | **MISSING** |
| Correo | `thunderbird` | ✅ [extra] | **MISSING** |
| Multimedia | `vlc` | ✅ [extra] | **MISSING** |
| PDF | `okular` | ✅ [extra] | **MISSING** |
| Calculadora | `qalculate-gtk` | ✅ [extra] | **MISSING** |
| Credenciales | `keepassxc` | ✅ [extra] | **MISSING** |
| Navegador | `chromium` | ✅ [extra] | **MISSING** |
| Gestor archivos | `yazi` | ✅ [extra] | **MISSING** |
| Virtualización | `virtualbox` | ✅ [extra] | **MISSING** |
| Contenedores | `docker` | ✅ [extra] | **MISSING** |
| Contenedores | `docker-compose` | ✅ [extra] | **MISSING** |
| Contenedores | `lazydocker` | ✅ [extra] | **MISSING** |
| Gaming | `steam` | ✅ [multilib] | **MISSING + needs multilib** |
| Gaming | `wine` | ✅ [multilib] | **MISSING + needs multilib** |
| Música | `spotify` | ❌ AUR only | **AUR** |
| Notas | `obsidian` | ❌ AUR only | **AUR** |
| Nube | `megasync` | ❌ AUR only | **AUR** |
| Bluetooth | `bluetui` | ❌ AUR only | **AUR** |
| WiFi | `impala` | ❌ AUR | **AUR / possibly crates.io** |

## pacman.conf Analysis

- **[multilib] is COMMENTED OUT** (lines 93-94)
- Must uncomment `[multilib]` and `Include = /etc/pacman.d/mirrorlist` for Steam/wine
- Current repos: `[core]` and `[extra]` only

## AUR Handling

- **No AUR scripts exist** — no `aur-packages.sh`, no AUR helper in packages.x86_64
- archiso does NOT support AUR packages natively (only pacman repos)
- AUR packages MUST be installed post-boot via a dedicated script

## Gap Summary

**13 packages** can be added directly to packages.x86_64 (official repos).
**2 packages** require multilib enabled (Steam, wine).
**5 packages** are AUR-only and need a separate install script.
**1 critical config change**: uncomment [multilib] in pacman.conf.