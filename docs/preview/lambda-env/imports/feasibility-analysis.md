# Análisis de Factibilidad: Prototipo → Implementación Go + Bubbletea

> Evalúa la traducción del prototipo web `docs/preview/lambda-env/` a la implementación real en terminal con Go, Bubbletea y Lipgloss.

---

## 1. Resumen Ejecutivo

El prototipo React cubre ~95% del spec UI/UX con 7 componentes TUI y 17 módulos funcionales. La implementación real en Go (`src/lambda-env/`) ya tiene la arquitectura correcta (máquina de estados 3 niveles, sistema de plugins, 10 módulos binarios), pero su capa de presentación es **minimalista** comparada con el prototipo: sin header, sin search, sin secciones colapsables, sin estados de carga, y con widgets más básicos.

**Veredicto**: ~80% del prototipo es factible en terminal. Un ~10% es imposible (efectos CRT, múltiples tamaños de fuente, sombras). El ~10% restante requiere workarounds creativos pero es alcanzable.

---

## 2. Estado Actual de la Implementación Go

### Lo que ya existe

| Componente | Archivo Go | Estado |
|-----------|-----------|--------|
| Máquina de estados 3 niveles | `internal/tui/model.go` | ✅ Completo |
| Navegación por teclado (j/k, ↑↓, Esc, Enter) | `internal/tui/update.go` | ✅ Completo |
| Help overlay (`?`) | `internal/tui/components/help.go` | ✅ Completo |
| Status bar con colores semánticos | `internal/tui/components/statusbar.go` | ✅ Completo (5 tipos) |
| Toggle (●/○) | `internal/tui/components/toggle.go` | ✅ Completo |
| TextInput con validación | `internal/tui/components/textinput.go` | ✅ Completo |
| Confirm modal | `internal/tui/components/confirm.go` | ✅ Completo |
| Sistema de plugins (hub + binaries) | `pkg/module/` | ✅ Completo |
| Persistencia en JSON | `~/.config/lambdaos/settings.json` | ✅ Completo |

### Módulos implementados (10 binarios)

| Módulo | Binario | Acciones |
|--------|---------|----------|
| Display | `display/main.go` | run, set-mode, set-position, set-primary, save/load-profile |
| Audio | `audio/main.go` | run, set-volume, set-mute, set-sink, set-source, set-profile, set-app-volume |
| Keyboard | `keyboard/main.go` | run, set-layout, set-variant, set-compose, set-options |
| Appearance | `appearance/main.go` | run, set-theme, set-wallpaper, set-font-size |
| Power | `power/main.go` | run, set-screen-timeout, set-sleep-timeout, set-lid-close-action |
| Defaults | `defaults/main.go` | run, set-browser/terminal/editor/file-manager, apply |
| Hardware Dashboard | `hardware-dashboard/main.go` | run (CPU, RAM, disco, temps, batería, uptime) |
| Neovim | `neovim/main.go` | run, set-theme, set-font, set-lines/columns, enable/disable features |
| Qtile | `qtile/main.go` | run, set-bar-position, set-bar-size, set-terminal/browser, set-layouts |
| Dotfiles | `dotfiles/main.go` | run, backup, stow, list-conflicts |

### Lo que NO existe en Go

- **Header bar**: sin breadcrumb, sin reloj, sin hint de búsqueda
- **Search overlay** (`/`): inexistente; solo navegación por lista
- **Slider widget**: sin componente de rango numérico visual
- **Progress bar**: sin componente de barra de progreso
- **Secciones colapsables**: la vista de detalle es una lista plana de acciones
- **Estados de carga en botones**: sin "Ejecutando..." ni "✓ Listo"
- **Indicador de nivel de seguridad**: no implementado en el módulo
- **Countdown timer**: no implementado para cambios de resolución
- **Signal bars**: sin representación visual de intensidad WiFi
- **Per-app volume sliders**: el backend soporta `set-app-volume` pero la UI no lo renderiza
- **Preview tipográfico**: no aplicable en terminal (un solo tamaño de fuente)
- **Tarjetas de monitores en grid**: no hay grid CSS en terminal
- **Avatares con iniciales**: sin renderizado de avatar

---

## 3. Matriz de Factibilidad

Cada feature del prototipo clasificada según su viabilidad en terminal:

| Feature | Prototipo | Factibilidad Go | Complejidad | Notas |
|---------|-----------|----------------|------------|-------|
| **Navegación y Layout** |
| Estados 3 niveles (cat→mod→detail) | ✅ | **NATIVE** | — | Ya implementado |
| Breadcrumb en header | ✅ | **DOABLE** | Baja | Agregar al status bar o nuevo header |
| Search overlay (`/`) | ✅ | **DOABLE** | Media | Sub-modelo de búsqueda + overlay |
| Reloj en header | ✅ | **NATIVE** | Baja | `tea.Tick` cada segundo |
| Scroll en contenido | ✅ | **DOABLE** | Media | `bubbles/viewport` |
| **Widgets** |
| TUISlider (track+thumb, drag) | ✅ | **DOABLE** | Media | `←/→` en vez de drag; `bubbles/slider` existe |
| TUIToggle (●/○, texto on/off) | ✅ | **NATIVE** | — | Ya existe como `Toggle` |
| TUISelect (dropdown, teclado) | ✅ | **DOABLE** | Media | `bubbles/list` o custom; `◄►` ya funciona |
| TUIInput (focus, validación, error) | ✅ | **NATIVE** | — | Ya existe como `TextInput` |
| TUIButton (loading, done, variantes) | ✅ | **DOABLE** | Baja | Agregar estados al renderizado de acciones |
| TUIModal (overlay, backdrop, danger) | ✅ | **NATIVE** | — | Ya existe como `Confirm` |
| TUIProgress (barra+thumb, color) | ✅ | **NATIVE** | Baja | Lipgloss renderiza barras fácilmente |
| TUISection (colapsable, rootRequired) | ✅ | **DOABLE** | Baja | Agregar header de sección a detail view |
| **Features específicas por módulo** |
| Tarjetas de monitores (grid 3 columnas) | ✅ | **HACKY** | Alta | Sin grid; usar `lipgloss.JoinHorizontal` |
| Countdown 10s con auto-revert | ✅ | **DOABLE** | Media | `tea.Tick` + estado en detail view |
| Signal bars (5 barras WiFi) | ✅ | **NATIVE** | Baja | Unicode: ▁▂▃▄▅▆▇█ |
| Panel de conexión inline | ✅ | **DOABLE** | Media | Cursor anidado en detail view |
| Per-app volume sliders | ✅ | **DOABLE** | Media | Múltiples sliders en una vista |
| Nivel de seguridad dinámico | ✅ | **DOABLE** | Baja | Calculado desde estado de toggles |
| Preview tipográfico (múltiples tamaños) | ✅ | **IMPOSSIBLE** | — | Terminal usa un solo tamaño de fuente |
| Avatares con iniciales | ✅ | **NATIVE** | Baja | `[L]` con borde lipgloss |
| Modal de confirmación eliminar | ✅ | **NATIVE** | — | Ya existe como `Confirm` |
| **Efectos visuales** |
| Persistencia (localStorage) | ✅ | **NATIVE** | — | Go usa JSON (mejor) |
| Transiciones suaves (CSS) | ✅ | **DOABLE** | Media | Tick-based animation en Bubbletea |
| Viñeta CRT (radial gradient) | ✅ | **IMPOSSIBLE** | — | Terminal no soporta gradientes |
| Scanlines CRT | ✅ | **IMPOSSIBLE** | — | Sin control de píxeles |
| Flicker CRT (CSS keyframes) | ✅ | **HACKY** | Media | Alternar dim/bright en ticks |
| Brillo/sombra de botones | ✅ | **IMPOSSIBLE** | — | Sin `box-shadow` en terminal |
| Varios tamaños de fuente | ✅ | **IMPOSSIBLE** | — | Terminal tiene un solo tamaño |
| Opacidad (implemented: false) | ✅ | **DOABLE** | Baja | `Faint(true)` o color atenuado |
| Indicador modificado (*) | ✅ | **NATIVE** | — | Ya existe en StatusBar |

**Resumen**: 8 NATIVE, 12 DOABLE, 2 HACKY, 6 IMPOSSIBLE.

---

## 4. Brechas Críticas

Features del prototipo que **no pueden replicarse** en una terminal real:

### 4.1 Imposibles (sin workaround viable)

| Feature | Razón |
|---------|-------|
| **Viñeta CRT** | Requiere `radial-gradient` CSS. La terminal es una grilla de caracteres. |
| **Scanlines CRT** | Requiere `repeating-linear-gradient` a nivel píxel. |
| **Múltiples tamaños de fuente** | La terminal tiene un solo tamaño de fuente por sesión. Para jerarquía visual se usa bold, color, mayúsculas. |
| **Sombras y glows** | Sin `box-shadow`. Se usa intensidad de color y bordes. |

### 4.2 Hacks necesarios (posibles pero incómodos)

| Feature | Workaround |
|---------|-----------|
| **Grid de 3 columnas** | Usar `lipgloss.JoinHorizontal` para unir cards lado a lado. Frágil con texto largo, requiere ancho fijo. |
| **Flicker CRT** | Alternar colores dim/bright con `tea.Tick`. Efecto sutil, no justifica la complejidad. |
| **Drag de slider con mouse** | Bubbletea soporta `tea.MouseMsg` pero es limitado. Mejor usar `←/→`. |

### 4.3 Lo que el prototipo muestra y la terminal puede igualar

El 80% del valor del prototipo está en los **flujos de interacción**, no en los efectos visuales:

- Flujo de conexión WiFi con panel inline → **DOABLE**
- Countdown de display con auto-revert → **DOABLE**
- Perfiles de audio guardados → **DOABLE** (widgets estándar)
- Lista de usuarios con acciones contextuales → **DOABLE**
- Búsqueda con `/` → **DOABLE**
- Secciones colapsables → **DOABLE**
- Estados de carga y feedback → **DOABLE**

---

## 5. Orden de Implementación Recomendado

### Fase 1 — Fundación (mayor impacto, menor esfuerzo)

| # | Feature | Esfuerzo | Prototipo paridad |
|---|---------|----------|------------------|
| 1 | **Header bar** (título, breadcrumb, reloj, hint búsqueda) | Bajo | Parcial — sin search overlay aún |
| 2 | **Indicador modificado (*)** en breadcrumb | Bajo | ✅ Completo |
| 3 | **Estados loading/done en acciones** ("Ejecutando...", "✓ Listo") | Bajo | ✅ Completo |
| 4 | **Secciones colapsables** (`▼/▶` toggle en detail view) | Bajo | ✅ Completo |
| 5 | **Opacidad para no implementados** (`Faint(true)`) | Bajo | ✅ Completo |

### Fase 2 — Widgets (esfuerzo medio)

| # | Feature | Esfuerzo | Prototipo paridad |
|---|---------|----------|------------------|
| 6 | **Slider** (`←/→`, track+thumb, min/max) | Medio | ✅ Completo |
| 7 | **Progress bar** (Lipgloss rendering) | Bajo | ✅ Completo |
| 8 | **Search overlay** (sub-modelo + `/` trigger) | Medio | ✅ Completo |
| 9 | **Viewport scrolleable** (`bubbles/viewport`) | Medio | Parcial |

### Fase 3 — Features por módulo (esfuerzo medio-alto)

| # | Feature | Esfuerzo | Notas |
|---|---------|----------|-------|
| 10 | **Signal bars WiFi** (▁▂▃▄▅▆▇█) | Bajo | Unicode nativo |
| 11 | **Countdown timer display** | Medio | `tea.Tick` + estado |
| 12 | **Indicador nivel seguridad** | Bajo | Badge calculado |
| 13 | **Avatares con iniciales** | Bajo | `[L]` con borde |
| 14 | **Per-app volume sliders** | Medio | Múltiples widgets |
| 15 | **Tarjetas de monitores** | Alto | `JoinHorizontal`, hacky |
| 16 | **Panel conexión inline** | Medio | Cursor anidado |

### Fase 4 — Pulido (opcional)

| # | Feature | Esfuerzo | Valor |
|---|---------|----------|-------|
| 17 | Animaciones tick-based (transiciones color) | Medio | Bajo — cosmético |
| 18 | Flicker CRT opcional | Medio | Bajo — gimmick |

---

## 6. Mapeo Componente a Componente

| Prototipo (React/TS) | Go equivalente | Acción requerida |
|----------------------|---------------|-----------------|
| `App.tsx` (nav state, search, status) | `internal/tui/model.go` + `update.go` + `view.go` | Agregar header, clock, search state |
| `TUIWindow.tsx` (chrome) | `internal/tui/view.go` (nuevo header render) | Agregar renderizado de header bar |
| `TUIToggle.tsx` | `internal/tui/components/toggle.go` | ✅ Existe; agregar texto "Activado"/"Desactivado" |
| `TUIInput.tsx` | `internal/tui/components/textinput.go` | ✅ Existe |
| `TUISelect.tsx` | `views/detail.go` (inline) | Extraer a componente; agregar modo dropdown |
| `TUISlider.tsx` | — | **NUEVO** `internal/tui/components/slider.go` |
| `TUIButton.tsx` | `views/detail.go` (inline) | Extraer; agregar estados loading/done |
| `TUIModal.tsx` | `internal/tui/components/confirm.go` | ✅ Existe; agregar variante danger |
| `TUIProgress.tsx` | — | **NUEVO** `internal/tui/components/progress.go` |
| `TUISection.tsx` | — | **NUEVO** o integrar en `detail.go` |
| `categories.ts` (datos) | `pkg/module/manifest.go` + Hub | ✅ Discovery dinámico ya existe |
| `tokens.ts` (diseño) | `internal/tui/view.go` (estilos) | ✅ Colores ya definidos; agregar constantes |
| `crt.css` | — | **NO IMPLEMENTAR** — inviable en terminal |
| `StatusBar` (en App.tsx) | `internal/tui/components/statusbar.go` | ✅ Existe; agregar breadcrumb, clock |
| `DisplayModule.tsx` | `internal/modules/display/main.go` | Backend ✅; TUI necesita cards + countdown |
| `AudioModule.tsx` | `internal/modules/audio/main.go` | Backend ✅; TUI necesita sliders + perfiles |
| `NetworkModule.tsx` | — | **NUEVO módulo binario** + TUI rendering |
| `SecurityModule.tsx` | — | **NUEVO módulo binario** + TUI rendering |
| `SimpleModules.tsx` (DateTime, Users, etc.) | Parcial (vía settings schema) | Mayoría necesitan nuevos binarios |

---

## 7. Conclusión

**El prototipo es un 80% factible en terminal.** Las partes inviables son exclusivamente efectos visuales cosméticos (viñeta CRT, scanlines, sombras, múltiples tamaños de fuente) que no afectan la funcionalidad ni la experiencia de usuario real.

La arquitectura del prototipo (máquina de estados, sistema de componentes, flujos de interacción) es **directamente trasladable** a Bubbletea. De hecho, la implementación Go ya tiene los fundamentos correctos — solo necesita enriquecer su capa de presentación.

**Lo que más valor aportaría con menor esfuerzo**: header con breadcrumb y reloj, search overlay, estados de carga en acciones, y secciones colapsables. Con esas 4 cosas, la experiencia se acerca significativamente al prototipo.

**Lo que requiere decisiones de diseño**: el grid de monitores y el panel de conexión inline son los features más "web" del prototipo. En terminal requieren repensar el layout: ¿lista vertical con preview inline? ¿vista dividida? Son decisiones de UX que el spec no resuelve y que el prototipo asume resueltas por el paradigma web.

---

*Análisis generado comparando `docs/preview/lambda-env/` (prototipo React) contra `src/lambda-env/` (implementación Go + Bubbletea).*
