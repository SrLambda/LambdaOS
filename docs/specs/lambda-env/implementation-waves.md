# lambda-env — Plan de Implementación por Waves

> **Última actualización**: Junio 2026 — Agregada Wave 3.5 (Visual Design Foundation con Nerd Fonts).  
> **Prototipo de referencia**: `docs/preview/lambda-env/` (22/31 módulos, React + TypeScript).  
> **Spec UI/UX**: `docs/Specs_lambda-env_UI/UX.md`.

## Visión

`lambda-env` es la suite de configuración TUI de LambdaOS. Equivale a "GNOME Settings" + "KDE System Settings" pero completamente en terminal.

## Principios de Implementación

1. **Cada wave genera una ISO testeable** — no hay waves "solo código". Cada wave produce una ISO que buildea y bootea.
2. **Cambios atómicos** — cada spec es un cambio independiente con su propio commit.
3. **Monorepo hasta v1.0** — todo en este repo.
4. **Framework agnóstico** — la TUI no está atada a un lenguaje. El hub descubre y ejecuta módulos como CLI tools independientes.
5. **Nerd Fonts es el estándar visual** — `nerd-fonts-monoid` incluido en la ISO base. Cada icono tiene fallback Unicode.
6. **Proyecto personal** — waves pequeñas (2-9 specs) para mantener en ventanas de tiempo limitadas.

## Estrategia de CI/CD progresiva

| Nivel | Waves | Qué hace |
|---|---|---|
| **Smoke test** | 0+ | ISO buildea y bootea en QEMU sin crash |
| **Feature test** | 2+ | ISO buildea, bootea, y verifica que la feature de la wave funciona |
| **Install test** | 8+ | ISO buildea, instala en disco virtual, y el sistema instalado bootea |
| **CD automático** | 9 | Push tag → release automático con ISO |

---

## Wave 0: Pipeline + ISO mínima funcional ✅

**Duración estimada**: 2-3 días
**Specs**: 4

| # | Spec | Qué valida |
|---|---|---|
| 1 | `ci-01-ci-workflow` | CI lintea + buildea ISO en cada push |
| 2 | `polish-04-release-tag` | Tags generan releases con versionado semántico |
| 3 | `pkg-01-flameshot` | Flameshot instalado + keybinding Qtile (Mod+Shift+S) |
| 4 | `branding-06-iso-name` | ISO con nombre LambdaOS correcto |

**Criterio de salida**: Push a main → CI buildea ISO → ISO bootea en QEMU → Flameshot disponible.

---

## Wave 1: Decisiones arquitectónicas de la TUI ✅

**Duración estimada**: 3-5 días
**Specs**: 3

| # | Spec | Qué valida |
|---|---|---|
| 5 | `core/02-settings-schema` | `settings.json` existe y se lee correctamente |
| 6 | `core/01-hub-plugin-system` | Framework elegido, hub binario, descubre módulos |
| 7 | `infra-01-repo-pacman-setup` | Repo pacman local configurado |

**Criterio de salida**: `lambda-env` ejecuta desde terminal. Contrato de módulos definido.

---

## Wave 2: Primeros módulos funcionales ✅

**Duración estimada**: 5-7 días
**Specs**: 4

| # | Spec | Qué valida |
|---|---|---|
| 8 | `apps-01-neovim` | TUI togglea LSP/Copilot → nvim respeta |
| 9 | `apps-02-qtile` | TUI cambia terminal default → Qtile respeta |
| 10 | `ops-05-dotfiles` | TUI stow/unstow → dotfiles aplicados |
| 11 | `infra-02-repo-package-tui` | TUI instalable como paquete pacman |

---

## Wave 3: TUI Interface + System Modules ✅

**Estado**: Completado (junio 2026)
**Duración real**: ~10 días
**Resultado**: 7 módulos funcionales, TUI interactiva con sub-modelos Bubble Tea

### Lo que se construyó

| Categoría | Módulo | Acciones |
|-----------|--------|----------|
| system | **appearance** | set-theme, set-wallpaper, set-font-size |
| system | **audio** | set-volume, set-mute, set-sink |
| system | **keyboard** | set-layout, set-variant |
| system | **defaults** | set-browser/terminal/editor/file-manager, apply |
| apps | **neovim** | toggle-lsp, toggle-copilot, toggle-neotree, set-theme |
| apps | **qtile** | set-terminal, set-browser, set-file-manager, reload |
| ops | **dotfiles** | stow, unstow, backup |

### TUI construida
- 3 vistas jerárquicas: Categories → Modules → ModuleDetail
- Sub-modelos Bubble Tea
- 5 tipos de widgets: toggle, select, text, confirm, execute
- Help overlay (`?`), status bar, confirm dialogs

---

## Wave 3.5: Visual Design Foundation 🆕

**Duración estimada**: 3-5 días
**Specs**: 5

Esta wave no agrega módulos nuevos. Establece el **lenguaje visual** que todas las waves futuras usarán. Todo lo construido hasta Wave 3 se actualiza retroactivamente.

| # | Spec | Qué valida | Impacto |
|---|---|---|---|
| 3.5.1 | **nerd-fonts-foundation** | `nerd-fonts-monoid` instalado en ISO base. Mapa de iconos cargado al inicio. Detección de disponibilidad + fallback automático a Unicode. | 🔴 Fundacional |
| 3.5.2 | **module-icon-set** | Cada módulo recibe su icono Nerd Fonts según tabla de diseño. Vista de categorías y módulos actualizada con nuevos iconos. | 🟡 Alto |
| 3.5.3 | **widget-icons** | Toggle, loading, success, error, warning, search, confirm — todos los widgets migrados a iconos Nerd Fonts. | 🟡 Alto |
| 3.5.4 | **color-palette-refined** | Paleta de 5 colores ajustada para contraste WCAG AA. Status bar, toggles, y errores verificados contra el spec UI/UX. | 🟢 Medio |
| 3.5.5 | **prototype-icon-sync** | Los 22 módulos del prototipo React actualizados para usar Nerd Fonts. El prototipo y la implementación Go comparten el mismo set de iconos. | 🟢 Medio |

**Criterio de salida**: 
- `lambda-env` ejecuta en terminal con Nerd Fonts → todos los iconos renderizan correctamente.
- Terminal sin Nerd Fonts (raw tty, SSH básico) → fallback Unicode automático, sin glifos rotos.
- Contraste de colores verificado con tooling (WCAG AA).
- Prototipo React y código Go comparten el mismo `icon-map.json`.

**Dependencias**: Wave 3 completada.

---

## Wave 4: Hardware Management

**Duración estimada**: 7-10 días
**Specs**: 5

| # | Spec | Qué valida | Nerd Fonts |
|---|---|---|---|
| 12 | `display-module` | Detectar monitores, cambiar resolución, refresh rate, perfiles, **countdown 10s con auto-revert** | `` desktop + `` laptop |
| 13 | `power-module` | Sleep timeout, lid close action, battery status con **iconos de batería 5 niveles** (``→``) | `` bolt + ``→`` |
| 14 | `keyboard-enhanced` | Layouts (10+), variants, compose key, options → setxkbmap + localectl | `` keyboard |
| 15 | `audio-enhanced` | Sink/source dinámica, perfiles de audio, **per-app volume sliders** | `` volume |
| 16 | `hardware-dashboard` | Vista resumen hardware con **barras de progreso** (CPU, RAM, disco, temp) | `` chart |

**Criterio de salida**: TUI gestiona pantalla (countdown auto-revert), energía (iconos batería), audio (per-app volume). Dashboard con progress bars en tiempo real.

**Dependencias**: Wave 3.5 (Nerd Fonts) + Wave 3 (hub + TUI).

---

## Wave 5: Connectivity

**Duración estimada**: 7-10 días
**Specs**: 5

| # | Spec | Qué valida | Nerd Fonts |
|---|---|---|---|
| 17 | `network-module` | WiFi scan/connect, **signal bars visuales**, inline connection panel, Ethernet status | `` wifi + `` globe |
| 18 | `bluetooth-module` | Scan, pair, trust, **disconnect modal**, device list con estados | `` bluetooth |
| 19 | `known-networks` | Redes conocidas, forget, auto-connect toggle | `` check |
| 20 | `vpn-stubs` | Placeholder VPN (WireGuard/OpenVPN) | `` lock |
| 21 | `connection-status-bar` | Widget en status bar con íconos WiFi/BT + estado en tiempo real | `` + `` |

**Criterio de salida**: WiFi y Bluetooth 100% desde TUI. Conectarse a una red o pair-ear auriculares sin salir de la TUI.

**Dependencias**: Wave 3.5 + Wave 3.

---

## Wave 6: System Management

**Duración estimada**: 7-10 días
**Specs**: ??

| # | Spec | Qué valida | Nerd Fonts |
|---|---|---|---|
| 22 | `services-module` | Listar, enable/disable, start/stop systemd units. **Indicador de estado**: `●` running, `○` stopped, `` failed | `` cogs |
| 23 | `autostart-module` | apps que arrancan con sesión → XDG autostart .desktop files | `` cog |
| 24 | `updates-module` | Check updates, trigger pacman -Syu, **progress bar** durante actualización | `` refresh |
| 25 | `security-module` | Firewall toggle, reglas, SSH/GPG keys, **nivel de seguridad dinámico** (CRÍTICO→MÁXIMO) | `` shield |
| 26 | `fonts-module` | Listar fuentes, preview con `` + nombre real, instalar/desinstalar | `` font |
| 27 | `notifications-module` | DND toggle ``, timeout, posición, per-app rules | `` bell |

**Criterio de salida**: TUI gestiona systemd, autostart, paquetes, firewall, fuentes y notificaciones. Nivel de seguridad dinámico funcional.

**Dependencias**: Wave 3.5 + Wave 3. Puede correr en paralelo con Waves 4 y 5.

---

## Wave 7: UX Enhancement — Cross-cutting

**Duración estimada**: 10-14 días
**Specs**: 7

Esta wave no agrega módulos nuevos. Eleva la experiencia de TODOS los módulos existentes al nivel del prototipo.

| # | Spec | Qué valida | Impacto |
|---|---|---|---|
| 28 | **header-bar** | Header con **breadcrumb clickable** + **reloj** (`tea.Tick`) + **hint de búsqueda** (`/`). Visible en todas las vistas. | 🔴 Crítico |
| 29 | **global-search** | `Ctrl+F` o `/` busca en TODOS los settings. **Overlay** con resultados navegables. Enter va directo al setting. | 🔴 Crítico |
| 30 | **section-collapse** | `TUISection` colapsable con ``/``. Usado retroactivamente en todos los módulos existentes. | 🔴 Crítico |
| 31 | **loading-states** | Botones con estados `` (loading), `` (done), `` (error). Aplicado retroactivamente. | 🟡 Alto |
| 32 | **slider-widget** | `TUISlider` con `←/→`, track visual, min/max labels. Reemplaza inputs numéricos en audio, display, power, neovim, qtile. | 🟡 Alto |
| 33 | **restore-defaults** | "Reset to default" por módulo y global. Diálogo `` confirmación. Defaults del schema. | 🟡 Alto |
| 34 | **import-export** | Exportar/importar `settings.json` con validación + merge. Soporte para profiles. | 🟢 Medio |

**Criterio de salida**: La TUI se siente como una app profesional. Buscar "wifi" te lleva al módulo de red. Breadcrumbs te dicen dónde estás. Botones tienen feedback visual. Secciones colapsables organizan módulos complejos.

**Dependencias**: Waves 3.5 + 4 + 5 + 6 completas (necesita todos los módulos para search global).

---

## Wave 8: Apps & Operations

**Duración estimada**: 7-10 días
**Specs**: 9

| # | Spec | Qué valida | Nerd Fonts |
|---|---|---|---|
| 35 | `apps-screenshot` | Flameshot: tecla, formato, destino | `` camera |
| 36 | `apps-recording` | OBS: escenas, fuentes, calidad | `` film |
| 37 | `apps-terminal` | Kitty/Alacritty: fuente, tema, opacidad | `` terminal |
| 38 | `apps-filemanager` | Yazi: tema, preview, plugins | `` folder |
| 39 | `apps-ai` | OpenCode: modelo, API key, contexto | `` brain |
| 40 | `ops-monitor` | htop-like TUI: procesos, CPU, memoria, red | `` chart |
| 41 | `ops-storage` | Discos, particiones, **progress bars** por uso | `` hdd |
| 42 | `ops-logs` | **journalctl viewer con follow en vivo**, filtros, colores severidad 8 niveles, **iconos ``/``/``** | `` file_text |
| 43 | `ops-backup` | Snapshots BTRFS + snapper: list, create, rollback | `` archive |

**Criterio de salida**: Apps del ecosistema configurables desde TUI. Logs con follow en vivo. Storage con progress bars.

**Dependencias**: Wave 3.5 + Waves 3-6.

---

## Wave 9: Setup, Users & Regional

**Duración estimada**: 10-14 días
**Specs**: 9

| # | Spec | Qué valida | Nerd Fonts |
|---|---|---|---|
| 44 | `setup-wizard` | Wizard primer boot: idioma, teclado, timezone, usuario, tema | `` wizard |
| 45 | `setup-profiles` | Perfiles: dev (gcc, go, rust), gaming (steam, lutris), rescue | `` cubes |
| 46 | `users-module` | Crear/eliminar usuarios, grupos, contraseña, **avatares con iniciales** `[L]` | `` user |
| 47 | `datetime-module` | Timezone, NTP toggle, formato 12/24h, **reloj grande** en vista | `` clock |
| 48 | `regional-module` | Locale, idioma, formato números/moneda | `` globe |
| 49 | `accessibility-basic` | Alto contraste, texto grande, sticky keys | `` universal_access |
| 50 | `installer-calamares` | Integración Calamares: scaffolding, branding LambdaOS | `` download |
| 51 | `installer-disk` | Particionado guiado: automático, manual, LUKS opcional | `` hdd |
| 52 | `system-advanced-dashboard` | Vista integrada: Users + Datetime + Regional + Accessibility | `` cog |

**Criterio de salida**: Primer boot wizard funcional. ISO instalable con Calamares. Timezone/locale configurables.

**Dependencias**: Waves 3-7.

---

## Wave 10: Advanced & Release (v1.0.0)

**Duración estimada**: 10-14 días
**Specs**: 9

| # | Spec | Qué valida |
|---|---|---|
| 53 | `printers-module` | Detectar impresoras, añadir, cola → CUPS |
| 54 | `online-accounts` | Conectar Google, Nextcloud (OAuth2 stubs) |
| 55 | `sharing-module` | Compartir pantalla (VNC), archivos (Samba stubs) |
| 56 | `accessibility-advanced` | Lector pantalla (Orca), teclado en pantalla |
| 57 | `polish-sysctl` | swappiness, cache pressure, network buffers |
| 58 | `polish-services` | Servicios auditados y optimizados |
| 59 | `polish-mkinitcpio` | Initramfs optimizado (systemd, zstd, early KMS) |
| 60 | `cd-automation` | Push tag → CI build + test + release automático |
| 61 | `public-demo` | Demo interactiva en GitHub Pages (terminal WASM) |

**Criterio de salida**: **LambdaOS v1.0.0**. Release público.

---

## Resumen Actualizado

| Wave | Estado | Días | Specs | Entregable clave |
|---|---|---|---|---|
| **0** | ✅ | 2-3 | 4 | CI buildea ISO |
| **1** | ✅ | 3-5 | 3 | Hub + settings schema |
| **2** | ✅ | 5-7 | 4 | Neovim, Qtile, Dotfiles |
| **3** | ✅ | 10 | 7 módulos | TUI interactiva, 7 módulos |
| **3.5** 🆕 | 🔲 | 3-5 | 5 | Nerd Fonts foundation, iconos, paleta WCAG AA |
| **4** | 🔲 | 7-10 | 5 | Hardware (display countdown, per-app audio, dashboard) |
| **5** | 🔲 | 7-10 | 5 | Connectivity (WiFi signal bars, BT modals) |
| **6** | 🔲 | 7-10 | 6 | System (services status, security level, fonts preview) |
| **7** | 🔲 | 10-14 | 7 | UX (header, search, sections, loading, slider, undo) |
| **8** | 🔲 | 7-10 | 9 | Apps + Ops (logs follow, storage progress) |
| **9** | 🔲 | 10-14 | 9 | Setup, users, wizard, Calamares |
| **10** | 🔲 | 10-14 | 9 | Advanced, polish, CD, demo, **v1.0.0** |
| **Total** | | **83-112 días** | **61 specs** | **LambdaOS v1.0.0** |

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
                    └─→ Wave 3.5 (VISUAL DESIGN: Nerd Fonts + iconos + paleta) 🆕
                          │
                          ├─→ Wave 4 (hardware: display, power, keyboard, audio, dash)
                          │     │
                          │     └─→ Wave 8 (apps + ops)
                          │
                          ├─→ Wave 5 (connectivity: network, bluetooth)
                          │     │
                          │     └─→ Wave 8 (apps + ops)
                          │
                          └─→ Wave 6 (system: services, autostart, updates, security, fonts, notif)
                                │
                                ├─→ Wave 7 (UX: header, search, sections, loading, slider, undo)
                                │     │
                                │     └─→ Wave 9 (setup, users, regional, installer)
                                │           │
                                │           └─→ Wave 10 (advanced, polish, CD, demo → v1.0.0)
                                │
                                └─→ Wave 9 (setup, users, regional, installer)
```

**Notas**:
- Waves 4, 5 y 6 pueden desarrollarse en paralelo (solo dependen de Wave 3.5).
- Wave 3.5 es requisito para TODO lo que sigue. Ningún módulo nuevo se construye sin Nerd Fonts.
- Wave 7 (UX) requiere Waves 4-6 completas para aplicar los enhancements a todos los módulos.

---

## Reglas de Progreso

1. No avanzar a la siguiente wave sin que la ISO de la wave actual buildea y bootea en QEMU.
2. Cada spec dentro de una wave es un commit independiente (o PR encadenado si excede 400 líneas).
3. Si una wave se bloquea por una decisión arquitectónica, documentar la decisión y continuar.
4. Waves 4, 5 y 6 pueden desarrollarse en paralelo entre sí.
5. A partir de Wave 3.5, **todo módulo nuevo debe usar iconos Nerd Fonts con fallback Unicode**.
6. El prototipo React en `docs/preview/lambda-env/` debe mantenerse sincronizado con los cambios de diseño visual.
7. Las especificaciones vivas usan Engram como artifact store (sin archivos `openspec/` en el repo).
