# Factibilidad v2 — Prototipo lambda-env → Go + Bubbletea

> Actualización tras nuevos módulos: Neovim, Qtile, Storage, Logs, Dotfiles.  
> Estado: **22/31 módulos implementados (71%)**.

---

## 1. Qué Cambió

| Métrica | v1 (anterior) | v2 (ahora) |
|---------|--------------|------------|
| Módulos implementados | 17 | **22** |
| Categoría Sistema | 16/16 ✅ | 16/16 ✅ |
| Categoría Apps | 0/7 | **2/7** (Neovim, Qtile) |
| Categoría Ops | 1/5 | **4/5** (Monitor, Storage, Logs, Dotfiles) |
| Categoría Setup | 0/3 | 0/3 |
| Componentes TUI usados | 7 | 7 (todos) |
| Nuevos patrones visuales | — | Preview de código, miniaturas de layouts, follow en vivo |

### Nuevos módulos y sus particularidades

| Módulo | Lo interesante | Complejidad visual |
|--------|---------------|-------------------|
| **Neovim** | Preview de colorscheme con código Lua sintetizado en los colores del tema. Healthcheck animado secuencial. 12 plugins categorizados, 8 servidores LSP. | **Muy alta** — la preview card es el componente visual más rico del prototipo |
| **Qtile** | Bar preview dinámica (top/bottom). **Miniaturas visuales de 6 layouts** (MonadTall, Bsp, Floating, etc.) renderizadas con CSS grid. Selector de apps por defecto. | **Muy alta** — las miniaturas de layouts son imposibles en terminal |
| **Storage** | Tarjetas de discos con tipo/health/temp. Tabla de particiones con **TUIProgress por fila**. Modal de confirmación para desmontar. | **Alta** — primer uso real de TUIProgress en contexto de datos |
| **Logs** | **Follow en vivo** simulado (nuevas líneas cada 2s). Filtros combinables (unidad + prioridad + tiempo + texto). Colores por severidad (8 niveles). Contador de errores/warnings en vivo. | **Media** — el follow es nativo en Bubbletea, los filtros son estado |
| **Dotfiles** | Tabla de paquetes con stow/unstow. Modal de resolución de conflictos. **Progress bar animada** durante backup. Output de git status simulado. | **Alta** — progress bar + modales + tabla |

---

## 2. Matriz de Factibilidad Actualizada

Features nuevas respecto a v1 marcadas con 🆕.

| Feature | Prototipo | Go | Complejidad | Notas |
|---------|-----------|-----|------------|-------|
| **🆕 Preview de colorscheme con código Lua** | ✅ | **IMPOSSIBLE** | — | No hay parser de Lua ni sintaxis coloreada arbitraria en terminal |
| **🆕 Miniaturas visuales de layouts Qtile** | ✅ | **IMPOSSIBLE** | — | CSS grids coloreados. Máximo: diagramas ASCII. |
| **🆕 Bar preview dinámica (top/bottom)** | ✅ | **DOABLE** | Baja | Render condicional con lipgloss, ancho variable |
| **🆕 Healthcheck animado secuencial** | ✅ | **NATIVE** | Baja | `tea.Tick` + append a slice |
| **🆕 Tabla de particiones con progress bars** | ✅ | **HACKY** | Alta | Progress bar por fila es posible pero el grid CSS requiere `JoinHorizontal` |
| **🆕 Tarjetas de discos con health/temp** | ✅ | **DOABLE** | Media | Cards horizontales con lipgloss borders |
| **🆕 Follow de logs en vivo** | ✅ | **NATIVE** | Baja | `tea.Tick` cada 2s + append |
| **🆕 Filtros múltiples en logs** | ✅ | **NATIVE** | Baja | Estado local + filter function |
| **🆕 Colores por severidad (8 niveles)** | ✅ | **NATIVE** | — | Lipgloss colors |
| **🆕 Contador errores/warnings en vivo** | ✅ | **NATIVE** | Baja | Computado del estado filtrado |
| **🆕 Dotfiles: tabla stow/unstow** | ✅ | **HACKY** | Alta | Grid CSS → `JoinHorizontal` o lista vertical |
| **🆕 Dotfiles: resolución de conflictos** | ✅ | **NATIVE** | — | Modal `Confirm` ya existe |
| **🆕 Dotfiles: progress bar de backup** | ✅ | **NATIVE** | Baja | `TUIProgress` + `tea.Tick` |
| **🆕 Dotfiles: git status output** | ✅ | **NATIVE** | — | Texto coloreado |
| **🆕 Secciones con sub-headers categorizados** | ✅ | **DOABLE** | Baja | Render extra en TUISection |
| **🆕 Botón disabled** | ✅ | **NATIVE** | Baja | Lipgloss `Faint(true)` |

### Resumen acumulado (v1 + v2)

| Clasificación | Cantidad | Qué significa |
|--------------|----------|--------------|
| **NATIVE** | 12 | Ya existe o es trivial |
| **DOABLE** | 14 | Alcanzable con esfuerzo medio |
| **HACKY** | 4 | Posible con workarounds incómodos (tablas, grids) |
| **IMPOSSIBLE** | 8 | Principalmente efectos visuales CSS y previews ricos |

---

## 3. Lo que el prototipo v2 revela sobre el diseño

### 3.1 Patrones que sí funcionan en terminal

Estos módulos nuevos demuestran patrones que **fortalecen** la factibilidad:

- **Neovim healthcheck**: Animación secuencial de checks → `tea.Sequence` de comandos. Muy Bubbletea-friendly.
- **Logs follow**: `tea.Tick` + append. Es exactamente lo que Bubbletea hace bien.
- **Filtros múltiples**: Estado local inmutable + filter puro. Es Go idiomático.
- **Storage progress bars**: `TUIProgress` usado en contexto de datos reales. Valida el componente.
- **Dotfiles conflict resolution**: `TUIModal` danger variant. Ya existe en Go como `Confirm`.
- **Estados loading en botones**: "Ejecutando...", "Sincronizando...", "Creando backup...". El prototipo muestra que el patrón escala bien a múltiples módulos.

### 3.2 Patrones que NO funcionan en terminal

Estos son los verdaderos blockers:

| Feature | Por qué no | Qué hacer en su lugar |
|---------|-----------|----------------------|
| **Preview de código Lua coloreado** | Sin parser de sintaxis ni colores arbitrarios por theme | Mostrar nombre del colorscheme + swatch de colores (3-4 cuadraditos con el color de fondo) |
| **Miniaturas de layouts Qtile** | CSS grids. No hay forma de hacer miniaturas proporcionales en terminal | Diagrama ASCII fijo por layout o descripción textual. Ej: `MonadTall: [60%|40%]` |
| **Tablas con CSS grid** (Storage, Dotfiles) | Terminal es monoespacio, no hay columnas proporcionales | `lipgloss.JoinHorizontal` con anchos fijos, o lista vertical con info en 2 líneas |
| **Tarjetas flex-wrap** (Storage disks) | Sin flex-wrap en terminal | Lista horizontal con ancho fijo por tarjeta |

---

## 4. Cobertura vs Implementación Go Real

Comparando lo que el prototipo muestra contra lo que realmente corre en `src/lambda-env/`:

| Módulo | Backend Go | TUI Go | Prototipo UI | Gap |
|--------|-----------|--------|-------------|-----|
| Display | ✅ `display/main.go` | Básico | Rico (cards, countdown, slider) | **Alto** — falta toda la capa visual |
| Audio | ✅ `audio/main.go` | Básico | Rico (per-app, perfiles, EQ, secciones) | **Alto** |
| Keyboard | ✅ `keyboard/main.go` | Básico | Medio (selectores) | **Medio** |
| Appearance | ✅ `appearance/main.go` | Básico | Medio | **Medio** |
| Power | ✅ `power/main.go` | Básico | Medio (sliders, selectores) | **Medio** |
| Defaults | ✅ `defaults/main.go` | Básico | Medio (9 selectores) | **Medio** |
| Hardware Dashboard | ✅ `hardware-dashboard/main.go` | Básico | Medio (Monitor) | **Medio** |
| Neovim | ✅ `neovim/main.go` | Básico | Muy rico (preview, health, plugins) | **Alto** — preview imposible |
| Qtile | ✅ `qtile/main.go` | Básico | Muy rico (bar preview, layouts) | **Alto** — miniaturas imposibles |
| Dotfiles | ✅ `dotfiles/main.go` | Básico | Rico (tabla, conflictos, backups) | **Alto** |
| **Network** | ❌ No existe | — | Rico (signal bars, conexión inline) | **Total** — backend + TUI desde cero |
| **Bluetooth** | ❌ No existe | — | Medio | **Total** |
| **Security** | ❌ No existe | — | Rico (nivel dinámico, SSH, firewall) | **Total** |
| **Services** | ❌ Parcial | — | Medio | **Total** |
| **Updates** | ❌ No existe | — | Medio | **Total** |
| **Storage** | ✅ `almacenamiento/main.go` | — | Rico (tabla, progress) | **Alto** |
| **Logs** | ✅ `logs/main.go` | — | Rico (follow, filtros) | **Alto** |
| **DateTime** | ❌ | — | Simple | **Total** |
| **Users** | ❌ | — | Medio (avatares, modales) | **Total** |
| **Autostart** | ❌ | — | Simple | **Total** |
| **Fonts** | ❌ | — | Medio (preview tipográfico) | **Total** — preview imposible |
| **Notifications** | ❌ | — | Medio | **Total** |

### Resumen de gaps

- **10 módulos** ya tienen backend Go pero TUI mínima → la brecha es puramente visual
- **11 módulos** no tienen backend Go → requieren binario nuevo + TUI
- **4 módulos** sin backend pero con UI simple → buenos candidatos para empezar

---

## 5. Orden de Implementación Actualizado

### Fase 1 — TUI Foundation (lo que más valor aporta)

| # | Qué | Esfuerzo | Impacto |
|---|-----|----------|---------|
| 1 | **Header bar** (breadcrumb, reloj, hint búsqueda) | Bajo | 🔥 Alto — presente en todo el prototipo |
| 2 | **Search overlay** (`/`) | Medio | 🔥 Alto |
| 3 | **Estados loading/done** en botones | Bajo | 🔥 Alto — el prototipo los usa en TODOS los módulos |
| 4 | **TUISection colapsable** | Bajo | 🔥 Alto — usado en 18/22 módulos |
| 5 | **TUISlider** | Medio | Alto — Audio, Display, Power, Neovim, Qtile, Notifications |
| 6 | **TUIProgress** | Bajo | Medio — Storage, Dotfiles, Updates |
| 7 | **Indicador modificado** | Bajo | Medio — ya existe parcialmente |

### Fase 2 — Módulos con backend existente (enriquecer TUI)

| # | Módulo | Qué agregar | Esfuerzo |
|---|--------|------------|----------|
| 8 | **Display** | Monitor cards, countdown timer | Medio |
| 9 | **Audio** | Per-app sliders, perfiles guardados, secciones | Medio |
| 10 | **Storage** | Tabla particiones + progress bars | Alto (grid) |
| 11 | **Logs** | Follow en vivo, filtros, colores severidad | Bajo |
| 12 | **Dotfiles** | Tabla stow/unstow, conflictos, backups | Alto (grid) |
| 13 | **Neovim** | Healthcheck, plugins, swatch de colorscheme | Medio |
| 14 | **Qtile** | Bar preview, selectores, layouts (texto, no miniaturas) | Medio |

### Fase 3 — Módulos nuevos (backend + TUI)

| # | Módulo | Backend | TUI | Prioridad |
|---|--------|---------|-----|-----------|
| 15 | **Network** | Alto | Medio | Alta — feature core |
| 16 | **Bluetooth** | Medio | Bajo | Media |
| 17 | **Security** | Medio | Medio | Media |
| 18 | **Services** | Bajo (systemctl wrapper) | Bajo | Alta |
| 19 | **Updates** | Alto (pacman wrapper) | Bajo | Baja |
| 20 | **DateTime** | Bajo (timedatectl) | Bajo | Baja |
| 21 | **Users** | Medio | Bajo | Media |
| 22 | **Autostart** | Bajo | Bajo | Baja |
| 23 | **Fonts** | Bajo (fc-list) | Bajo | Baja |
| 24 | **Notificaciones** | Bajo (dunstctl) | Bajo | Baja |

---

## 6. Nerd Fonts: El Game-Changer Visual

### Por qué Nerd Fonts cambia la factibilidad

El spec original restringía a "Unicode-only". Con **Nerd Fonts como estándar** (Monoid Nerd Font incluido en la ISO base de LambdaOS), varios problemas que marqué como IMPOSSIBLE o HACKY se vuelven DOABLE o NATIVE:

| Feature | Con Unicode básico | Con Nerd Fonts |
|---------|-------------------|----------------|
| Preview de colorscheme Neovim | IMPOSSIBLE (sin código coloreado) | **DOABLE** — swatch `███` + icono `` + nombre |
| Miniaturas de layouts Qtile | IMPOSSIBLE (CSS grids) | **DOABLE** — icono `` + diagrama ASCII + label |
| Preview tipográfico Fonts | IMPOSSIBLE (múltiples tamaños) | **DOABLE** — icono `` + bold/italic + nombre real |
| Signal bars WiFi | DOABLE (▁▂▃▄▅) | **NATIVE** — `` + barras `█` con color semántico |
| Estados loading | DOABLE (texto) | **NATIVE** — `` spin animado con `tea.Tick` |
| Toggle on/off | NATIVE (`●`/`○`) | **NATIVE** — ``/`` con color, visualmente más rico |
| Severidad en logs | NATIVE (colores) | **NATIVE** — `` warning, `` error, `` info |
| Batería | DOABLE | **NATIVE** — ``→`` con 5 niveles |
| Módulos en lista | NATIVE (Unicode) | **NATIVE** — iconos semánticos por módulo (ver tabla abajo) |

### Tabla de iconos por módulo

| Módulo | Unicode fallback | Nerd Fonts |
|--------|-----------------|------------|
| Pantalla | `▣` | `` (desktop) o `` (laptop) |
| Audio | `♪` | `` (volume_up) |
| Red | `◉` | `` (wifi) |
| Bluetooth | `⊛` | `` (bluetooth) |
| Energía | `⚡` | `` (bolt) |
| Teclado | `⌨` | `` (keyboard) |
| Usuarios | `⊙` | `` (user) |
| Apariencia | `◈` | `` (paint_brush) |
| Seguridad | `⛨` | `` (user_secret) |
| Actualizaciones | `↻` | `` (refresh) |
| Servicios | `◎` | `` (cogs) |
| Fuentes | `A` | `` (font) |
| Notificaciones | `◻` | `` (bell) |
| Neovim | `✦` | `` (dev-vim) |
| Qtile | `⊞` | `` (window_maximize) |
| Dotfiles | `◌` | `` (git) |
| Logs | `≡` | `` (file_text) |
| Almacenamiento | `◫` | `` (hdd) |
| Monitor | `▦` | `` (bar_chart) |

### Principio de fallback

Si Nerd Fonts no está presente (raw tty, SSH sin forwarding, terminal sin soporte), el sistema detecta la ausencia y degrada automáticamente al set Unicode básico. Esto se implementa como un mapa de iconos cargado al inicio según la capacidad de la terminal.

---

## 7. Veredicto Final (Actualizado con Nerd Fonts)

| Pregunta | Respuesta |
|----------|-----------|
| ¿Es factible? | **Sí, ~88% del prototipo** (subió del 78%) |
| ¿Qué sigue siendo IMPOSSIBLE? | Efectos CRT (viñeta, scanlines), CSS grids multi-columna, sombras CSS. Efectos puramente cosméticos. |
| ¿Qué sigue siendo HACKY? | Tablas con CSS grid (Storage, Dotfiles, Services) — requieren `lipgloss.JoinHorizontal` con ancho fijo. |
| ¿Qué módulo era difícil y ahora es viable? | **Neovim** (swatch + icono), **Qtile** (icono + ASCII), **Fonts** (icono + nombre real). |
| ¿Vale la pena Nerd Fonts? | **Sí, rotundamente.** Sube la factibilidad un 10%, elimina 3 IMPOSSIBLE del tablero, y unifica el lenguaje visual a costo de UNA dependencia de paquete (`nerd-fonts-monoid`). |

### Las 5 cosas que más valor aportan implementar AHORA

1. **Nerd Fonts foundation** → detección + fallback + set de iconos. El piso sobre el que se construye todo lo demás.
2. **Header bar + search** → unifica la experiencia, cierra la deuda UX #1
3. **Estados loading/done** → feedback inmediato que hoy no existe en ningún lado, enriquecido con iconos Nerd Fonts
4. **TUISection colapsable** → usada en 18/22 módulos, sin esto el detail view es un paredón
5. **Logs module** → es el más fácil (follow + filtros = nativo) y demuestra el patrón correcto con iconos de severidad Nerd Fonts

---

*Análisis actualizado comparando `docs/preview/lambda-env/` (v2, 22 módulos) contra `src/lambda-env/` (implementación Go). Nerd Fonts incorporado como estándar visual.*
