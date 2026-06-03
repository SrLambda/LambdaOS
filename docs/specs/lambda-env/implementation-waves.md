# lambda-env — Plan de Implementación por Waves

## Visión

`lambda-env` es la suite de configuración TUI de LambdaOS. Equivale a "GNOME Settings" + "KDE System Settings" pero completamente en terminal.

## Principios de Implementación

1. **Cada wave genera una ISO testeable** — no hay waves "solo código". Cada wave produce una ISO que buildea y bootea.
2. **Cambios atómicos** — cada spec es un cambio independiente con su propio commit.
3. **Monorepo hasta v1.0** — todo en este repo.
4. **Framework agnóstico** — la TUI no está atada a un lenguaje. El hub descubre y ejecuta módulos como CLI tools independientes.
5. **Proyecto personal** — waves pequeñas (2-9 specs) para mantener en ventanas de tiempo limitadas.

## Estrategia de CI/CD progresiva

| Nivel | Waves | Qué hace |
|---|---|---|
| **Smoke test** | 0+ | ISO buildea y bootea en QEMU sin crash |
| **Feature test** | 2+ | ISO buildea, bootea, y verifica que la feature de la wave funciona |
| **Install test** | 8+ | ISO buildea, instala en disco virtual, y el sistema instalado bootea |
| **CD automático** | 9 | Push tag → release automático con ISO |

---

## Wave 0: Pipeline + ISO mínima funcional

**Duración estimada**: 2-3 días
**Specs**: 4

| # | Spec | Qué valida |
|---|---|---|
| 1 | `ci-01-ci-workflow` | CI lintea + buildea ISO en cada push |
| 2 | `polish-04-release-tag` | Tags generan releases con versionado semántico |
| 3 | `pkg-01-flameshot` | Flameshot instalado + keybinding Qtile (Mod+Shift+S) |
| 4 | `branding-06-iso-name` | ISO con nombre LambdaOS correcto |

**Criterio de salida**: Push a main → CI buildea ISO → ISO bootea en QEMU → Flameshot disponible. ISO se llama `LambdaOS-v0.0.1-x86_64.iso`.

**Dependencias**: Ninguna.

---

## Wave 1: Decisiones arquitectónicas de la TUI

**Duración estimada**: 3-5 días
**Specs**: 3

| # | Spec | Qué valida |
|---|---|---|
| 5 | `core/02-settings-schema` | `settings.json` existe y se lee correctamente |
| 6 | `core/01-hub-plugin-system` | Framework elegido, hub binario, descubre módulos |
| 7 | `infra-01-repo-pacman-setup` | Repo pacman local configurado |

**Criterio de salida**: `lambda-env` ejecuta desde terminal. Muestra menú. El contrato de módulos está definido. Framework decidido (Go/Python agnóstico/bash).

**Dependencias**: Wave 0 (CI funcionando).

---

## Wave 2: Primeros módulos funcionales

**Duración estimada**: 5-7 días
**Specs**: 4

| # | Spec | Qué valida |
|---|---|---|
| 8 | `apps-01-neovim` | TUI togglea LSP/Copilot → nvim respeta |
| 9 | `apps-02-qtile` | TUI cambia terminal default → Qtile respeta |
| 10 | `ops-05-dotfiles` | TUI stow/unstow → dotfiles aplicados |
| 11 | `infra-02-repo-package-tui` | TUI instalable como paquete pacman |

**Criterio de salida**: `lambda-env` → Neovim → toggle LSP → nvim abre sin LSP. Stow/unstow desde TUI. `pacman -S lambdaos-tui` funciona.

**Dependencias**: Wave 1 (hub + settings schema).

---

## Wave 3: TUI Interface + System Modules

**Duración estimada**: 7-10 días
**Specs**: 18 (2 tracks paralelos)

### Track A: TUI Interface Development (prioridad alta)
Actualmente la TUI solo muestra categorías + módulos y ejecuta. Necesita vistas interactivas:

| # | Spec | Qué valida |
|---|---|---|
| 12 | `tui-01-interactive-views` | Forms con inputs, toggles visuales, listas con estado |
| 13 | `tui-02-sub-navigation` | Sub-menus por módulo (Neovim → Toggles, Theme, Plugins) |
| 14 | `tui-03-status-bar` | Status bar muestra estado actual de settings, no solo errores |
| 15 | `tui-04-confirm-dialogs` | Confirmación antes de acciones destructivas (unstow, backup) |
| 16 | `tui-05-help-overlay` | Help overlay con teclas disponibles por vista |

### Track B: System Modules
| # | Spec | Qué valida |
|---|---|---|
| 17 | `system-09-appearance` | Tema global → sincroniza con neovim + qtile |
| 18 | `system-01-screen` | Detectar monitores, cambiar resolución → xrandr |
| 19 | `system-02-audio` | Volumen, sink default, mute → pipewire |
| 20 | `system-03-network` | WiFi/Ethernet desde TUI → NetworkManager |
| 21 | `system-04-bluetooth` | Bluetooth pairing → bluetoothctl |
| 22 | `system-05-power` | Power profiles, suspend → systemd-logind |
| 23 | `system-06-keyboard` | Layout, variant, options → setxkbmap |
| 24 | `system-10-defaults` | Apps por defecto (terminal, browser, file manager) |
| 25 | `system-11-autostart` | Servicios al inicio (picom, flameshot, etc.) |
| 26 | `system-12-services` | Enable/disable systemd services |
| 27 | `system-13-updates` | Updates disponibles, trigger upgrade |
| 28 | `system-14-security` | Firewall, sudo rules, SSH keys |
| 29 | `system-15-fonts` | Fuentes instaladas, font size global |
| 30 | `system-16-notifications` | Notification daemon, rules |

### Wave 3 Open Questions (from Wave 2)
- Neovim module: gestionar imports de plugins/ más allá de los 3 toggles
- Qtile module: parameterizar keys.py más allá de terminal
- Migración os_theme.json → settings.json (tema sincronizado appearance → neovim → qtile)

**Criterio de salida**:
- TUI tiene vistas interactivas con toggles, inputs, listas con estado visual
- Al menos 4 módulos system funcionales con vistas dedicadas
- Tema global sincronizado entre appearance → neovim → qtile
- lambda-env muestra UI rica, no solo launcher de módulos

**Dependencias**: Wave 2 completa + bug fix de package main.

---

## Wave 4: Hardware esencial

**Duración estimada**: 5-7 días
**Specs**: 4

| # | Spec | Qué valida |
|---|---|---|
| 16 | `system-01-screen` | Detectar monitores, cambiar resolución |
| 17 | `system-02-audio` | Volumen, dispositivos de audio |
| 18 | `system-05-power` | Brillo, batería |
| 19 | `pkg-02-obs-studio` | OBS instalado |

**Criterio de salida**: TUI gestiona pantalla, audio y energía. OBS disponible.

**Dependencias**: Wave 3 (system modules base + TUI interface).

---

## Wave 5: Conectividad

**Duración estimada**: 5-7 días
**Specs**: 5

| # | Spec | Qué valida |
|---|---|---|
| 20 | `system-03-network` | WiFi/Ethernet desde TUI |
| 21 | `system-04-bluetooth` | Bluetooth pairing desde TUI |
| 22 | `system-06-keyboard` | Layout de teclado |
| 23 | `pkg-03-davinci-resolve` | Resolve en script AUR |
| 24 | `pkg-04-opencode` | OpenCode disponible |

**Criterio de salida**: TUI gestiona toda la conectividad. Todos los paquetes de la Fase 1 listos.

**Dependencias**: Wave 3 (system modules base + TUI interface).

---

## Wave 6: Sistema + Docs

**Duración estimada**: 5-7 días
**Specs**: 7

| # | Spec | Qué valida |
|---|---|---|
| 25 | `system-12-services` | Habilitar/deshabilitar servicios |
| 26 | `system-13-updates` | Ver e instalar updates |
| 27 | `system-10-defaults` | Apps por defecto |
| 28 | `system-11-autostart` | Apps al inicio |
| 29 | `infra-03-repo-package-configs` | Configs empaquetadas |
| 30 | `infra-04-docs-local` | Docs en localhost:8080 |
| 31 | `infra-05-docs-content` | Contenido de docs |

**Criterio de salida**: TUI gestiona servicios, updates, defaults y autostart. Docs accesibles desde browser.

**Dependencias**: Wave 3 (system modules base + TUI interface).

---

## Wave 7: Apps + Ops

**Duración estimada**: 5-7 días
**Specs**: 9

| # | Spec | Qué valida |
|---|---|---|
| 32 | `apps-03-screenshot` | Config Flameshot desde TUI |
| 33 | `apps-04-recording` | Config OBS desde TUI |
| 34 | `apps-05-terminal` | Config Kitty desde TUI |
| 35 | `apps-06-filemanager` | Config Yazi desde TUI |
| 36 | `apps-07-ai` | Config OpenCode desde TUI |
| 37 | `ops-01-monitor` | Monitor de sistema |
| 38 | `ops-02-storage` | Discos y particiones |
| 39 | `ops-03-logs` | Logs del sistema |
| 40 | `ops-04-backup` | Snapshots y backups |

**Criterio de salida**: Todas las apps configurables desde TUI. Monitor, storage, logs y backup funcionales.

**Dependencias**: Wave 4 (OBS), Wave 5 (OpenCode), Wave 1 (hub).

---

## Wave 8: Setup + Installer

**Duración estimada**: 7-10 días
**Specs**: 8

| # | Spec | Qué valida |
|---|---|---|
| 41 | `setup-01-wizard` | Primer boot wizard |
| 42 | `setup-03-profiles` | Perfiles dev/gaming/rescue |
| 43 | `installer-01-calamares-scaffold` | Calamares instalado |
| 44 | `installer-02-calamares-modules` | Instalación funcional |
| 45 | `installer-03-calamares-branding` | Calamares con branding |
| 46 | `installer-04-launcher` | Lanzable desde TUI |
| 47 | `system-07-users` | Gestión de usuarios |
| 48 | `system-08-datetime` | Fecha, hora, timezone |

**Criterio de salida**: Se puede instalar LambdaOS en disco. Wizard de primer boot funcional.

**Dependencias**: Wave 8 (installer), Wave 9 (branding para Calamares).

---

## Wave 9: Polish + CD + Demo Pública

**Duración estimada**: 5-7 días
**Specs**: 9

| # | Spec | Qué valida |
|---|---|---|
| 49 | `polish-01-sysctl-tweaks` | Optimizaciones sysctl |
| 50 | `polish-02-services-default` | Servicios por defecto |
| 51 | `polish-03-mkinitcpio-custom` | Initramfs optimizado |
| 52 | `system-14-security` | Firewall, SSH, GPG |
| 53 | `system-15-fonts` | Gestión de fuentes |
| 54 | `system-16-notifications` | Config Dunst |
| 55 | `ci-02-cd-workflow` | Releases automáticos con tags |
| 56 | `ci-03-nightly-builds` | Builds nocturnos |
| 57 | `infra-06-tui-demo` | Demo interactiva de la TUI en GitHub Pages |

**Criterio de salida**: Push tag → ISO se buildea y publica automáticamente. Sistema optimizado. Demo interactiva en GitHub Pages navegable con todos los módulos.

**Dependencias**: Wave 8 (installer), Wave 0 (CI base), Wave 9 (branding para la demo).

---

## Resumen

| Wave | Días | Specs | Entregable clave |
|---|---|---|---|
| **0** | 2-3 | 4 | CI buildea ISO, Flameshot, nombre correcto |
| **1** | 3-5 | 3 | Framework decidido, hub abre, repo pacman |
| **2** | 5-7 | 4 | TUI controla Neovim, Qtile, dotfiles |
| **3** | 7-10 | 18 | TUI interactiva + 14 módulos system + tema sincronizado |
| **4** | 5-7 | 4 | Pantalla, audio, energía, OBS |
| **5** | 5-7 | 5 | Red, BT, teclado, todos los paquetes |
| **6** | 5-7 | 7 | Servicios, updates, docs |
| **7** | 5-7 | 9 | Todas las apps + ops |
| **8** | 7-10 | 8 | Wizard + Installer |
| **9** | 5-7 | 9 | Polish + CD automático + Demo pública |
| **Total** | **52-72 días** | **66 specs** | **Distro v1.0** |

---

## Mapa de Dependencias entre Waves

```
Wave 0 (CI + pipeline)
  │
  ├─→ Wave 1 (core TUI: hub + settings + repo)
  │     │
  │     ├─→ Wave 2 (módulos: neovim, qtile, dotfiles, empaquetado)
  │     │     │
  │     │     └─→ Wave 3 (TUI interface + system modules + tema sincronizado)
  │     │           │
  │     │           ├─→ Wave 4 (hardware: screen, audio, power)
  │     │           │     │
  │     │           │     └─→ Wave 7 (apps + ops: screenshot, recording, monitor, etc.)
  │     │           │
  │     │           ├─→ Wave 5 (conectividad: network, BT, keyboard)
  │     │           │     │
  │     │           │     └─→ Wave 7 (apps + ops)
  │     │           │
  │     │           └─→ Wave 6 (sistema: services, updates, docs)
  │     │                 │
  │     │                 └─→ Wave 8 (installer: Calamares + wizard)
  │     │
  │     └─→ Wave 8 (installer: Calamares + wizard)
  │
  └─→ Wave 9 (polish + CD + demo pública: depende de todo lo anterior)
```

## Reglas de Progreso

1. No avanzar a la siguiente wave sin que la ISO de la wave actual buildea y bootea en QEMU.
2. Cada spec dentro de una wave es un commit independiente.
3. Si una wave se bloquea por una decisión arquitectónica (ej: framework de la TUI en Wave 1), documentar la decisión y continuar.
4. Las waves 4, 5 y 6 pueden desarrollarse en paralelo entre sí (no dependen unas de otras, solo de Wave 1).
