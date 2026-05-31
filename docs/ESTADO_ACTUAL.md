# Estado Actual — LambdaOS

Documento generado el 2026-05-19. Resume el progreso del proyecto hasta la Fase 6 (build + test E2E).

---

## Fases Completadas

| Fase | Descripción | Agentes | Tests |
|------|-------------|---------|-------|
| **Fase 1** | Testing First — Pruebas QEMU con pytest + pexpect | qa_tester | `tests/qemu/test_live_boot.py` (3 tests) |
| **Fase 2** | Infraestructura Archiso — autologin getty TTY1, stow en skel | sysadmin_stow | — |
| **Fase 3** | Configuración Neovim — lazy.nvim, WBS Iteración 1 (22 archivos) | dev_nvim + sysadmin_stow | — |
| **Fase 4** | TUI "System Preferences" — Panel de Control en Textual | dev_tui + qa_tester | `tests/unit/test_tui_config.py` (18 tests) |
| **Fase 5** | Qtile modular + servicios systemd + tema unificado os_theme.json | dev_tui + sysadmin_stow + qa_tester | `tests/unit/test_qtile_config.py` (12 tests) |
| **Fase 6** | Script build_and_test.sh — compilación ISO y test E2E | sysadmin_stow | Pendiente de pasar |

**Total tests unitarios: 30/30 pasando** (18 TUI + 12 Qtile).
**Tests QEMU E2E: 0/3 pasando** (bloqueados por bug en curso).

---

## Ubicación de Archivos Clave

### TUI (`src/os_tui_configurator/`)
| Archivo | Función |
|---------|---------|
| `app.py` | `OsTuiConfigurator` — Header, Sidebar (ListView Neovim/Qtile), Content (Select tema + 3 Switches), Footer, Ctrl+S |
| `config_manager.py` | `ConfigManager` — I/O de `tui_settings.json`, `.nvim_theme`, `os_theme.json` usando `OS_CONFIG_DIR` |
| `style.tcss` | CSS externo de Textual (sidebar 24 cols, content 1fr, switch-rows) |
| `main.py` | Entry point: `python -m src.os_tui_configurator.main` |

### Dotfiles Stow — Neovim (`airootfs/etc/skel/dotfiles/nvim/.config/nvim/`)
| Archivo | Función |
|---------|---------|
| `init.lua` | Entry point: carga `tui_bridge` → `env` → `options` → `keymaps` → `lazy` |
| `tui_settings.json` | Banderas toggleables: `enable_lsp`, `enable_copilot`, `enable_neotree` |
| `lua/core/tui_bridge.lua` | Parsea `tui_settings.json` y expone `vim.g.tui_flags` |
| `lua/core/env.lua` | Lee `$NVIM_THEME` y `os_theme.json` (primario), fallback a `.nvim_theme` |
| `lua/core/lazy.lua` | Bootstrap de lazy.nvim + `import "plugins"` |
| `lua/core/options.lua` | tab=4, tw=80, wrap=false, number, mouse |
| `lua/core/keymaps.lua` | Leader=space |
| `lua/plugins/theme.lua` | Catppuccin configurable vía `$NVIM_THEME` |
| `lua/plugins/ui.lua` | Lualine, Alpha, Bufferline, indent-blankline |
| `lua/plugins/treesitter.lua` | ensure_installed=all |
| `lua/plugins/navigation.lua` | Telescope + Harpoon + Neo-tree (toggleable) |
| `lua/plugins/editing.lua` | Autopairs, Surround, Comment |
| `lua/plugins/tools.lua` | Gitsigns, Toggleterm |
| `lua/plugins/ai.lua` | Copilot (toggleable) |
| `lua/plugins/lsp.lua` | Mason + mason-lspconfig + nvim-cmp + 15 LSPs (toggleable) |
| `lua/plugins/formatting.lua` | conform.nvim format-on-save |
| `lua/plugins/linting.lua` | nvim-lint async |
| `lua/plugins/data.lua` | iron.nvim REPL |
| `lua/plugins/docs.lua` | Markdown preview, vimtex, spelllang=es |
| `ftplugin/java.lua` | nvim-jdtls (solo filetype java) |
| `ftplugin/rust.lua` | rustaceanvim (solo filetype rust) |

### Dotfiles Stow — Qtile (`airootfs/etc/skel/dotfiles/qtile/.config/qtile/`)
| Archivo | Función |
|---------|---------|
| `theme.py` | 5 temas (Catppuccin, Gruvbox, Tokyonight, Nord, OneDark). `load_theme()` lee `os_theme.json` con try/except y fallback |
| `keys.py` | 40+ keybindings. mod4/Super. Lanzadores: kitty, rofi, chromium, yazi |
| `groups.py` | 5 workspaces con labels Unicode |
| `screens.py` | Barra superior con GroupBox, WindowName, Clock, Systray, Volume, Battery |
| `config.py` | Layouts (MonadTall, Max, Columns) + hook `startup_once` (picom, xsetroot, stow) |

### Infraestructura Archiso (`airootfs/`)
| Ruta | Función |
|------|---------|
| `etc/systemd/system/getty@tty1.service.d/autologin.conf` | Autologin `liveuser` en TTY1 |
| `etc/systemd/system/display-manager.service` | Symlink → `/usr/lib/systemd/system/ly.service` |
| `etc/skel/.bash_profile` | Ejecuta `stow */` desde `~/dotfiles/` al iniciar bash |
| `etc/skel/.zprofile` | Igual para zsh |
| `etc/skel/.config/systemd/user/` | 5 unidades user (pipewire, pipewire-pulse, wireplumber, dunst, lxqt-policykit-agent) |
| `etc/skel/.config/systemd/user/default.target.wants/` | Symlinks de activación para cada unidad |

### Build & Test
| Archivo | Función |
|---------|---------|
| `build_and_test.sh` | Script: clean → pacman.conf check → mkarchiso → symlink ISO → pytest QEMU |
| `profiledef.sh` | Configuración mkarchiso: `iso_name="LambdaOS"`, bootmodes bios+uefi, squashfs |
| `packages.x86_64` | 172 paquetes (neovim, qtile, pipewire, ly, kitty, etc.) |
| `pacman.conf` | Configuración de pacman para el entorno archiso |
| `tests/unit/test_tui_config.py` | 18 tests (11 ConfigManager + 7 App UI) |
| `tests/unit/test_qtile_config.py` | 12 tests (sintaxis AST, imports, temas, barra) |
| `tests/qemu/test_live_boot.py` | 3 tests E2E (boot, stow symlinks, init.lua) |
| `tests/qemu/conftest.py` | Fixtures session-scoped: qemu_booted, qemu_logged_in |

---

## Error Actual (Bug Bloqueante) — CORREGIDO 2026-05-19

### Síntoma
El test `test_iso_boots_to_shell_prompt` en `tests/qemu/test_live_boot.py` falla por timeout. QEMU arranca la ISO pero la prueba nunca detecta el prompt esperado.

### Cambios aplicados en el fix

**1. `tests/qemu/conftest.py` — fixture `qemu_booted`:**
- **Fase 1 (bootloader)**: Se ampliaron los regex de 3 a 6 patrones: `boot:`, `ISOLINUX|SYSLINUX`, `Arch Linux install medium`, `Automatic boot in`, `Press \[Tab\]|Press Enter`, `Welcome to Arch Linux`.
- **Fase 2 (shell prompt)**: Se ampliaron los regex de 4 a 6 patrones, ahora soportan bash (`$`), zsh (`%`), root (`#`), con/sin brackets, con/sin hostname, con ANSI codes.
- **Timeout**: `TIMEOUT_BOOT` aumentado de 180s a 300s para emulación pura sin KVM.
- **Logging**: Se agregó `pexpect.spawn(logfile=)` escribiendo todo el output a `/tmp/opencode/qemu_boot_debug.log`.
- **Función `_dump_last_output()`**: Escribe los últimos 3000 chars del buffer al log y al mensaje de fallo para diagnóstico.
- Se agregó `time.sleep(1)` post-login para dejar que la shell se asiente.

**2. `airootfs/etc/systemd/system/getty@tty1.service.d/autologin.conf`:**
- Se cambió `--autologin liveuser -` por `--autologin liveuser %I`, usando el especificador de instancia de systemd en vez de `-` (stdin), mejor compatibilidad con `getty@.service`.

**3. `tests/qemu/test_live_boot.py`:**
- Todos los patrones `[\$#] ` se actualizaron a `[\$#%>] ` para soportar prompts zsh (`%`).

### Archivos modificados
- `tests/qemu/conftest.py` — fixture `qemu_booted` (líneas 54-127), nueva función `_dump_last_output` (línea 130-147), fixture `qemu_logged_in` (línea 150-162)
- `tests/qemu/test_live_boot.py` — 4 patrones de prompt actualizados
- `airootfs/etc/systemd/system/getty@tty1.service.d/autologin.conf` — `%I` en vez de `-`

---

## Comandos de Prueba

```bash
# Tests unitarios (30 tests, pasan todos)
python -m pytest tests/unit/ -v

# Tests QEMU E2E (requiere ISO compilada)
python -m pytest tests/qemu/test_live_boot.py -v

# Build + test completo (requiere sudo para mkarchiso)
./build_and_test.sh

# Previsualizar TUI
OS_CONFIG_DIR=/tmp/test_tui_config python -m src.os_tui_configurator.main
```
