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

## Wave 3: TUI Interface + System Modules ✅ COMPLETADO

**Estado**: Completado (junio 2026)
**Duración real**: ~10 días
**Resultado**: 7 módulos funcionales, TUI interactiva con sub-modelos Bubble Tea

### Lo que se construyó

| Categoría | Módulo | Acciones |
|-----------|--------|----------|
| system | **appearance** | set-theme (dark/light/nord/catppuccin), set-wallpaper, set-font-size |
| system | **audio** | set-volume, set-mute, set-sink |
| system | **keyboard** | set-layout (us), set-variant |
| system | **defaults** | set-browser, set-terminal, set-editor, set-file-manager, apply |
| apps | **neovim** | toggle-lsp, toggle-copilot, toggle-neotree, set-theme, apply |
| apps | **qtile** | set-terminal, set-browser, set-file-manager, reload |
| ops | **dotfiles** | stow, unstow, backup |

### TUI interface construida
- **3 vistas jerárquicas**: Categories → Modules → ModuleDetail
- **Sub-modelos Bubble Tea**: cada vista es un `tea.Model` independiente
- **5 tipos de widgets**: toggle, select, text, confirm, execute
- **Help overlay**: `?` muestra teclas disponibles
- **Status bar**: contexto + nombre de módulo + estado
- **Confirm dialogs**: para acciones destructivas (unstow, backup)

### Settings schema (v1.1.0)
17 categorías definidas en `settings.json` — solo 7 tienen módulo implementado.

### Lo que quedó PENDIENTE de la Wave 3 original
Los siguientes specs del Track B original NO se implementaron y se redistribuyen en Waves 4-6:
`system-01-screen`, `system-03-network`, `system-04-bluetooth`, `system-05-power`, `system-11-autostart`, `system-12-services`, `system-13-updates`, `system-14-security`, `system-15-fonts`, `system-16-notifications`

**Dependencias**: Wave 2 completa.

---

## Wave 4: Hardware Management

**Duración estimada**: 7-10 días
**Specs**: 5

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 12 | `display-module` | Detectar monitores, cambiar resolución, refresh rate, perfiles → xrandr/wlr-randr | system |
| 13 | `power-module` | Sleep timeout, lid close action, battery status → systemd-logind + upower | system |
| 14 | `keyboard-enhanced` | Ampliar layouts (es, la, de, fr), variants, compose key, options → setxkbmap + localectl | system |
| 15 | `audio-enhanced` | Sink/source selection dinámica, perfiles de audio, volume per-app | system |
| 16 | `hardware-dashboard` | Vista resumen de hardware: CPU, RAM, disk, temp, battery — solo lectura, actualizable | system |

**Criterio de salida**: TUI gestiona pantalla (resolución, monitores), energía (sleep, batería), teclado (10+ layouts), audio (selección de sinks). Dashboard muestra estado del hardware en tiempo real.

**Dependencias**: Wave 3 (hub + settings schema + TUI interface).

---

## Wave 5: Connectivity

**Duración estimada**: 7-10 días
**Specs**: 5

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 17 | `network-module` | WiFi scan/connect/disconnect, Ethernet status, IP info → NetworkManager/nmcli | system |
| 18 | `bluetooth-module` | Scan, pair, trust, connect, disconnect devices → bluez/bluetoothctl | system |
| 19 | `known-networks` | Lista de redes conocidas, forget, auto-connect toggle → nmcli connection | system |
| 20 | `vpn-stubs` | Placeholder para VPN (WireGuard/OpenVPN) — detectar configuraciones existentes | system |
| 21 | `connection-status` | Widget en status bar: WiFi/BT íconos con estado (connected/disconnected/scanning) | tui |

**Criterio de salida**: TUI gestiona WiFi y Bluetooth completamente desde la terminal. Conectarse a una red o pair-ear un dispositivo sin salir de la TUI.

**Dependencias**: Wave 3 (hub + TUI interface).

---

## Wave 6: System Management

**Duración estimada**: 7-10 días
**Specs**: 7

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 22 | `services-module` | Listar, enable, disable, start, stop systemd units → systemctl | system |
| 23 | `autostart-module` | Gestionar apps que arrancan con sesión → XDG autostart .desktop files | system |
| 24 | `updates-module` | Check updates disponibles, trigger pacman -Syu, mostrar changelog → pacman + checkupdates | system |
| 25 | `security-module` | Firewall (ufw enable/disable, allow/deny ports), sudo timeout, screen lock timeout | system |
| 26 | `fonts-module` | Listar fuentes instaladas, cambiar monospace/sans/serif defaults → fc-list + fontconfig | system |
| 27 | `notifications-module` | Do not disturb toggle, timeout, per-app rules → dunst config | system |
| 28 | `system-health` | Dashboard extendido: servicios running, temp, disk usage, últimas actualizaciones | ops |

**Criterio de salida**: TUI gestiona systemd, autostart, paquetes, firewall, fuentes y notificaciones. Dashboard de salud del sistema funcional.

**Dependencias**: Wave 3 (hub + TUI interface). Puede correr en paralelo con Waves 4 y 5.

---

## Wave 7: TUI UX Enhancement — Cross-cutting

**Duración estimada**: 10-14 días
**Specs**: 7

Esta wave no agrega módulos nuevos. Mejora la experiencia de uso de TODOS los módulos existentes.

| # | Spec | Qué valida | Impacto |
|---|---|---|---|
| 29 | `global-search` | `Ctrl+F` busca a través de TODOS los settings de todos los módulos. Resultados navegables con Enter para ir directo al setting. | 🔴 Crítico |
| 30 | `breadcrumbs` | Ruta de navegación visible y clickeable: `Sistema → Teclado → Layout`. Cada segmento es navegable. | 🟡 Alto |
| 31 | `restore-defaults` | "Reset to default" por módulo y global. Diálogo de confirmación. Defaults definidos en el schema. | 🟡 Alto |
| 32 | `import-export` | Exportar `settings.json` completo a un archivo. Importar desde archivo con validación + merge. Soporte para profiles. | 🟡 Alto |
| 33 | `real-time-preview` | Cambios de tema, fuente, wallpaper se aplican en tiempo real sin "Apply" explícito. Acciones no destructivas son instantáneas. | 🟢 Medio |
| 34 | `theme-sync` | `use_global_theme` conectado: cambiar tema en appearance → neovim y qtile reflejan el cambio automáticamente. Resuelve el open question de Wave 3. | 🟡 Alto |
| 35 | `settings-diff` | Antes de aplicar cambios, mostrar diff de qué cambió. Opción de undo por sesión. | 🟢 Medio |

**Criterio de salida**: La TUI se siente como una app de settings profesional. Buscar "wifi" te lleva al módulo de red. Los breadcrumbs te dicen dónde estás. Restaurar defaults es un comando. Exportás tu config y la importás en otra máquina.

**Dependencias**: Waves 3-6 completas (necesita todos los módulos para search y theme-sync).

---

## Wave 8: Apps & Operations

**Duración estimada**: 7-10 días
**Specs**: 9

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 36 | `apps-screenshot` | Configurar Flameshot: tecla, formato, destino → flameshot config | apps |
| 37 | `apps-recording` | Configurar OBS: escenas, fuentes, calidad → obs-websocket | apps |
| 38 | `apps-terminal` | Configurar Kitty: fuente, tema, opacidad → kitty.conf | apps |
| 39 | `apps-filemanager` | Configurar Yazi: tema, preview, plugins → yazi.toml | apps |
| 40 | `apps-ai` | Configurar OpenCode: modelo, API key, contexto → opencode.json | apps |
| 41 | `ops-monitor` | htop-like dentro de la TUI: procesos, CPU, memoria, red | ops |
| 42 | `ops-storage` | Discos, particiones, espacio libre → lsblk + df | ops |
| 43 | `ops-logs` | Viewer de journalctl con filtros (service, priority, time) | ops |
| 44 | `ops-backup` | Snapshots BTRFS + snapper: list, create, rollback | ops |

**Criterio de salida**: Todas las apps del ecosistema LambdaOS son configurables desde la TUI. Monitor, storage, logs y backup funcionales.

**Dependencias**: Wave 4 (OBS), Wave 5 (OpenCode), Wave 3 (hub).

---

## Wave 9: Setup, Users & Regional

**Duración estimada**: 10-14 días
**Specs**: 9

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 45 | `setup-wizard` | Wizard de primer boot: idioma, teclado, timezone, usuario, tema, apps default | setup |
| 46 | `setup-profiles` | Perfiles predefinidos: dev (gcc, go, rust, docker), gaming (steam, lutris, wine), rescue (herramientas de sistema) | setup |
| 47 | `users-module` | Crear/eliminar usuarios, grupos, cambiar contraseña, auto-login → useradd + passwd | system |
| 48 | `datetime-module` | Timezone, NTP toggle, formato 12/24h → timedatectl | system |
| 49 | `regional-module` | Locale, idioma del sistema, formato de números/moneda → localectl + locale.conf | system |
| 50 | `accessibility-basic` | Alto contraste, texto grande, sticky keys → gsettings + config files | system |
| 51 | `installer-calamares` | Integración Calamares: scaffolding, módulos, branding LambdaOS, launcher desde TUI | installer |
| 52 | `installer-disk` | Particionado guiado: automático, manual simple, LUKS opcional | installer |
| 53 | `system-advanced-dashboard` | Vista integrada de Users + Datetime + Regional + Accessibility | system |

**Criterio de salida**: Primer boot wizard funcional. Usuario creado con wizard. ISO instalable en disco con Calamares. Timezone y locale configurables desde TUI.

**Dependencias**: Waves 3-6 (system modules), Wave 8 (apps), Wave 0 (CI).

---

## Wave 10: Advanced & Release

**Duración estimada**: 10-14 días
**Specs**: 9

| # | Spec | Qué valida | Tipo |
|---|---|---|---|
| 54 | `printers-module` | Detectar impresoras, añadir, cola de trabajos → CUPS + system-config-printer | system |
| 55 | `online-accounts` | Conectar Google, Nextcloud (stubs con OAuth2 placeholder) → GNOME Online Accounts o similar | system |
| 56 | `sharing-module` | Compartir pantalla (VNC), archivos (Samba stubs) | system |
| 57 | `accessibility-advanced` | Lector de pantalla (Orca), teclado en pantalla, opciones de contraste avanzadas | system |
| 58 | `polish-sysctl` | Optimizaciones sysctl (swappiness, cache pressure, network buffers) | polish |
| 59 | `polish-services` | Servicios por defecto auditados y optimizados (mask/unmask innecesarios) | polish |
| 60 | `polish-mkinitcpio` | Initramfs optimizado (hook systemd, compresión zstd, early KMS) | polish |
| 61 | `cd-automation` | Push tag → CI buildea ISO, corre smoke + feature + install tests, crea GitHub Release | ci-cd |
| 62 | `public-demo` | Demo interactiva en GitHub Pages: TUI navegable via terminal emulator WASM | demo |

**Criterio de salida**: Distro completa. Push tag → release automático. Demo pública navegable. Sistema optimizado.
**Versión**: `v1.0.0` — Release público.

**Dependencias**: Todo lo anterior. Esta es la wave de salida.

---

## Resumen Actualizado

| Wave | Estado | Días | Specs | Entregable clave |
|---|---|---|---|---|
| **0** | ✅ | 2-3 | 4 | CI buildea ISO |
| **1** | ✅ | 3-5 | 3 | Hub + settings schema + repo pacman |
| **2** | ✅ | 5-7 | 4 | Neovim, Qtile, Dotfiles configurables |
| **3** | ✅ | 10 | 7 módulos | TUI interactiva, 7 módulos system/apps/ops |
| **4** | 🔲 | 7-10 | 5 | Display, power, keyboard+, audio+, dashboard |
| **5** | 🔲 | 7-10 | 5 | WiFi, Bluetooth, known networks, VPN stubs |
| **6** | 🔲 | 7-10 | 7 | Services, autostart, updates, security, fonts, notifications |
| **7** | 🔲 | 10-14 | 7 | Search global, breadcrumbs, restore defaults, import/export, preview, theme sync |
| **8** | 🔲 | 7-10 | 9 | Apps (screenshot, recording, terminal, fm, AI) + Ops (monitor, storage, logs, backup) |
| **9** | 🔲 | 10-14 | 9 | Wizard, perfiles, users, datetime, regional, accessibility, Calamares |
| **10** | 🔲 | 10-14 | 9 | Printers, online accounts, sharing, polish, CD auto, demo, v1.0.0 |
| **Total** | | **80-109 días** | **62 specs** | **LambdaOS v1.0.0** |

---

## Mapa de Dependencias entre Waves

```
Wave 0 (CI + pipeline) ✅
  │
  └─→ Wave 1 (core TUI: hub + settings + repo) ✅
        │
        └─→ Wave 2 (módulos: neovim, qtile, dotfiles) ✅
              │
              └─→ Wave 3 (TUI interface + 7 system modules) ✅
                    │
                    ├─→ Wave 4 (hardware: display, power, keyboard, audio)
                    │     │
                    │     └─→ Wave 8 (apps + ops)
                    │
                    ├─→ Wave 5 (connectivity: network, bluetooth)
                    │     │
                    │     └─→ Wave 8 (apps + ops)
                    │
                    └─→ Wave 6 (system: services, autostart, updates, security, fonts, notifications)
                          │
                          ├─→ Wave 7 (UX: search, breadcrumbs, restore, import/export, preview, sync)
                          │     │
                          │     └─→ Wave 9 (setup, users, regional, installer)
                          │           │
                          │           └─→ Wave 10 (advanced, polish, CD, demo → v1.0.0)
                          │
                          └─→ Wave 9 (setup, users, regional, installer)
```

**Nota**: Waves 4, 5 y 6 pueden desarrollarse en paralelo (no dependen entre sí, solo de Wave 3).

---

## Reglas de Progreso

1. No avanzar a la siguiente wave sin que la ISO de la wave actual buildea y bootea en QEMU.
2. Cada spec dentro de una wave es un commit independiente (o PR encadenado si excede 400 líneas).
3. Si una wave se bloquea por una decisión arquitectónica, documentar la decisión y continuar.
4. Las waves 4, 5 y 6 pueden desarrollarse en paralelo entre sí (solo dependen de Wave 3).
5. La Wave 7 (UX) requiere Waves 4-6 completas para indexar todos los módulos en el search global.
6. A partir de Wave 4, cada módulo nuevo usa el patrón establecido en Wave 3: `manifest.json` + binario Go en `internal/modules/<name>/`.
7. Las especificaciones vivas están en `openspec/specs/`. Las specs de planning histórico están en `docs/specs/lambda-env/`.
