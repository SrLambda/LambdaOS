# lambda-env — Suite de Configuración TUI para LambdaOS

## Visión

`lambda-env` es la suite de aplicaciones TUI que funciona como centro de configuración de LambdaOS. Equivale a "GNOME Settings" + "KDE System Settings" + "GNOME/KDE Gear" pero completamente en terminal.

## Principios de Diseño

1. **Todo es un módulo**: cada funcionalidad es un módulo independiente descubrible por el hub
2. **Settings unificados**: un solo archivo JSON (`~/.config/lambdaos/settings.json`) que leen tanto la TUI como Qtile, Neovim y cualquier app que necesite configuración
3. **TUI-first**: la configuración se hace desde la terminal, pero los cambios afectan todo el sistema (X11, systemd, dotfiles)
4. **Idempotente**: aplicar la misma configuración dos veces no rompe nada
5. **Sin dependencias de GUI**: funciona en tty pura, SSH, o dentro de un terminal emulator

## Arquitectura

```
lambda-env (hub principal)
├── core/                    ← plugin-loader, settings schema, theme engine
├── modules/system/          ← screen, audio, network, bluetooth, power, etc.
├── modules/apps/            ← neovim, qtile, screenshot, recording, ai, etc.
├── modules/ops/             ← monitor, storage, logs, backup, dotfiles
└── modules/setup/           ← wizard, installer, profiles
```

## Módulos Inventariados

### Core (2)
- `01-hub-plugin-system` — Hub principal + descubrimiento de módulos
- `02-settings-schema` — Schema JSON unificado de configuración

### System (16)
- `system-01-screen` — Gestión de monitores (xrandr)
- `system-02-audio` — Volumen, dispositivos, perfiles de audio (PipeWire)
- `system-03-network` — WiFi, Ethernet, VPN
- `system-04-bluetooth` — Pairing y gestión de dispositivos
- `system-05-power` — Brillo, batería, suspensión
- `system-06-keyboard` — Layout, repeat rate, atajos del sistema
- `system-07-users` — Gestión de cuentas de usuario
- `system-08-datetime` — Zona horaria, NTP, formato de hora
- `system-09-appearance` — Tema GTK, iconos, cursor, wallpaper
- `system-10-defaults` — Aplicaciones por defecto (xdg-mime)
- `system-11-autostart` — Apps al inicio (systemd user units)
- `system-12-services` — Habilitar/deshabilitar servicios systemd
- `system-13-updates` — Actualizar sistema + AUR
- `system-14-security` — Firewall, SSH keys, GPG
- `system-15-fonts` — Gestión de fuentes del sistema
- `system-16-notifications` — Configuración de Dunst

### Apps (7)
- `apps-01-neovim` — Configurar LSP, plugins, toggles de Neovim
- `apps-02-qtile` — Configurar layouts, keybindings, bar, groups
- `apps-03-screenshot` — Configuración de Flameshot
- `apps-04-recording` — Configuración de OBS Studio
- `apps-05-terminal` — Configuración de Kitty
- `apps-06-filemanager` — Configuración de Yazi
- `apps-07-ai` — Configuración de OpenCode / agentes IA

### Ops (5)
- `ops-01-monitor` — Monitor de sistema integrado (htop/btop)
- `ops-02-storage` — Discos, particiones, montaje
- `ops-03-logs` — Viewer de journalctl
- `ops-04-backup` — Snapper, rsync, snapshots
- `ops-05-dotfiles` — GNU Stow, sync, perfiles de dotfiles

### Setup (3)
- `setup-01-wizard` — Primer boot wizard
- `setup-02-installer` — Instalar LambdaOS en disco
- `setup-03-profiles` — Perfiles de sistema (dev, gaming, rescue, minimal)

## Priorización

### v1.0 (Must-have)
- Core: hub + settings schema
- System: screen, audio, network, bluetooth, power, appearance, updates, services
- Apps: neovim, qtile, dotfiles
- Setup: wizard

### v1.1+ (Nice-to-have)
- Resto de módulos system, apps, ops
- Installer
- Profiles

## Tecnología

Framework TUI a definir (opciones: `textual` Python, `bubbletea` Go, `whiptail/dialog` bash puro). La decisión se toma en `01-hub-plugin-system`.
