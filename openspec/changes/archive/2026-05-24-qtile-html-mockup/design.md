# Design: Qtile Interface HTML Mockup

## Technical Approach

Single self-contained `preview/index.html` with embedded CSS and JS — no build step, framework, or server required. Mirrors the real Qtile `screens.py` bar widget order, `theme.py` color system, `groups.py` workspace icons, and `keys.py` shortcut map. CSS custom properties enable zero-lag theme switching. Vanilla JS manages 5-workspace state, 3-layout engine, and global keyboard dispatch.

## Architecture Decisions

| Decision | Options | Tradeoff | Choice |
|----------|---------|----------|--------|
| File structure | Single `index.html` vs. split `html`/`css`/`js` | Split: cleaner dev but breaks `file://` for modular JS (CORS). Single: ~1500 lines but works everywhere | **Single `index.html`** — matches spec 10 single-file preference |
| Layout engine | CSS Grid vs. JS absolute positioning with calculation | Grid: simpler code, but columns/monadtall resize complexity. JS abs-pos: pixel-perfect, matches Qtile's screen-space model | **JS absolute positioning** — recalculates on layout switch, resize, window changes |
| Theme switching | Inline style injection vs. CSS custom properties on `:root` | Inline: many DOM queries. Custom props: single `setProperty` per token, all elements react | **CSS custom properties** — set 9 tokens on `:root`, zero DOM walks |
| Window rendering | `<div>` with border + titlebar vs. `<iframe>` for real apps | iframe: real content but blocked by `file://` CORS. Div: pure styling, always works | **Styled `<div>`** — Kitty/Chromium/Neovim as divs with CSS-simulated content |
| Keyboard capture | `keydown` on `<body>` vs. per-element handlers | Per-element: misses unfocused targets. Body-level: captures always but needs Super-key tracking | **`document.addEventListener('keydown')`** — global, no focus dependency |
| Font strategy | Self-host fonts in repo vs. CDN links | Self-host: offline works but ~10MB in repo. CDN: lightweight, spec allows | **CDN** — FontAwesome via jsDelivr, Monoid Nerd Font via cdn-fonts. Fallback to system monospace on offline |

## Data Flow

```
                   ┌──────────────┐
  localStorage     │  index.html  │     CDN (FA fonts + Monoid)
  theme pref ─────→│  on load     │←────  <link> tags
                   └──────┬───────┘
                          │
         ┌────────────────┼────────────────────┐
         ▼                ▼                     ▼
  ┌──────────┐    ┌──────────────┐    ┌────────────────┐
  │ Theme    │    │ State Manager│    │ Keyboard       │
  │ Engine   │    │ (workspaces, │    │ Dispatcher     │
  │ (CSS     │    │  layout,     │    │ (Super+key     │
  │  props)  │    │  windows,    │    │  → actions)    │
  └────┬─────┘    │  focus)      │    └───────┬────────┘
       │          └──────┬───────┘            │
       │                 │                    │
       ▼                 ▼                    ▼
  ┌────────────┐  ┌──────────────┐  ┌────────────────┐
  │ DOM update │  │ Layout Engine│  │ Notification   │
  │ (all color │  │ (window      │  │ System (toast  │
  │  elements) │  │  positions)  │  │  auto-dismiss) │
  └────────────┘  └──────────────┘  └────────────────┘
```

## Component Model

| Component | DOM Root | State Coupling | Renders |
|-----------|----------|----------------|---------|
| **Bar** | `<header id="bar">` | theme, activeWs, layout, focusedWindow | GroupBox (5 icons + highlight), CurrentLayout, WindowName, SystemTray placeholder, Volume, Battery, Clock |
| **Desktop** | `<main id="desktop">` | activeWs, layout | Window divs positioned by layout engine |
| **Window** | `<div class="window">` (per instance) | app type (kitty/chromium/neovim), focused, floating | Title bar + app-specific content inner div |
| **ThemeSelector** | `<select id="theme-select">` | theme | 5 `<option>` elements, onChange dispatches theme switch |
| **KeybindingsPanel** | `<aside id="keybindings">` | theme (colors only) | Categorized table from keys.py data, togglable via button |
| **NotificationSystem** | `<div id="notifications">` | queue of {message, timeout} | Animated toast divs, auto-remove after 2s |

## State Management (vanilla JS object)

```javascript
const state = {
  theme: 'catppuccin',           // loaded from localStorage on init
  layout: 'monadtall',           // 'monadtall' | 'max' | 'columns'
  activeWorkspace: 1,            // 1-5
  workspaces: {                  // keyed by 1-5, each holding window array
    1: [], 2: [], 3: [], 4: [], 5: []
  },
  focusedWindowId: null,         // id of window with primary border
  superHeld: false,              // tracks Super key state
  nextWindowId: 0,               // monotonic counter for window IDs
};
// Theme persistence: localStorage key 'lambdaos-qtile-theme'
```

## Theme System Design

```css
:root[data-theme="catppuccin"] {
  --bg: #1e1e2e; --fg: #cdd6f4; --primary: #cba6f7;
  --secondary: #89b4fa; --accent: #f5c2e7; --urgent: #f38ba8;
  --inactive: #45475a; --bar-bg: #181825; --bar-fg: #cdd6f4;
}
/* 4 more :root[data-theme="..."] blocks */
```

Theme switch: `document.documentElement.setAttribute('data-theme', name)` — one DOM mutation touches every `var(--*)` reference. `localStorage.setItem` on change, `getItem` on load.

## Layout Engine Design

Windows are absolutely positioned children of `#desktop`. The engine is a pure function `computeLayout(state, viewport)` → `Map<windowId, {x, y, w, h}>`:

| Layout | Calculation |
|--------|-------------|
| **MonadTall** | Main window at `{x:0, y:0, w:vw*0.6, h:vh-32}`. Remaining windows stacked on right: `{x:vw*0.6, y:row*h/n, w:vw*0.4, h:h/n}` |
| **Max** | Focused window at `{x:0, y:0, w:vw, h:vh-32}`. All others hidden. |
| **Columns** | N windows get `{x:i*vw/n, y:0, w:vw/n, h:vh-32}` for i=0..n-1 |

Layout transitions use `transition: left 200ms, top 200ms, width 200ms, height 200ms` on window elements.

## Keyboard Handler Design

```javascript
document.addEventListener('keydown', (e) => {
  if (e.key === 'Meta') { state.superHeld = true; showSuperBadge(); return; }
  if (!state.superHeld) return;
  e.preventDefault();
  const action = SHORTCUT_MAP[e.key]; // { 'Return': openKitty, 'b': openChromium, ... }
  if (action) action();
});
document.addEventListener('keyup', (e) => {
  if (e.key === 'Meta') { state.superHeld = false; hideSuperBadge(); }
});
```

`SHORTCUT_MAP` is keyed by `event.key` (not `event.code`). Special handling for `Tab`/`Shift+Tab` via `e.shiftKey` check. Numeric keys `1`-`5` dispatch to workspace switch or move based on `e.shiftKey`.

## Responsive Strategy

| Breakpoint | Behavior |
|------------|----------|
| 320-480px (mobile) | Bar: clock + workspace icons only, scrollable if needed. Desktop: windows stack vertically. Tappable workspace switches. Floated Keybindings panel covers full screen |
| 481-768px (tablet) | Full bar widgets, windows adapt proportionally. MonadTall main reduced to 50%. Columns max 2 per row |
| 769px+ (desktop) | Full layout as designed. All widgets visible |

Implemented via `@media (max-width: 768px)` and `@media (max-width: 480px)` blocks. No horizontal scroll at any width.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `preview/index.html` | Create | Single self-contained mockup: ~400 lines CSS custom props + theme blocks, ~500 lines JS (state, layout engine, keyboard, notifications), ~300 lines HTML (bar, desktop, keybindings panel). Total ~1200 lines |
| `preview/README.md` | Create | How to open, shortcuts reference, theme list, no-server note |

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Manual — Theme | All 5 themes apply correctly | Visual inspection: switch each theme, verify bar bg, text, workspace highlight, window borders match theme.py values |
| Manual — Layout | MonadTall/Max/Columns reposition windows | Open 3 windows, cycle layouts via Super+Tab, verify proportions match spec |
| Manual — Shortcuts | All 14 shortcuts execute | Systematic walkthrough of each shortcut, verify action + notification appears |
| Manual — Persistence | Theme survives reload | Select Gruvbox, reload page, verify Gruvbox applied |
| Manual — Responsive | Mobile/tablet/desktop breakpoints | Resize browser to each range, verify no horizontal scroll, all elements functional |
| Manual — Offline | Works without CDN | Disconnect network, reload — verify layout works, icons fall back to text |

No automated tests — this is a visual HTML mockup. Manual verification per the scenarios in specs.md.

## Migration / Rollout

No migration required. New `preview/` directory is additive to the repository. No existing files modified. Rollback: delete `preview/` directory.

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| CDN downtime breaks FontAwesome icons | Low — icons become fallback text, layout intact | FontAwesome fallback: `font-family: "Font Awesome 6 Free", sans-serif`. Spec confirms offline fallback is acceptable |
| Super key captured by OS before reaching JS | Medium — some browsers/OS intercept Meta key | `keydown` event still fires; show warning in README about OS-level Super key bindings (e.g., GNOME overview) |
| Window limit with many keyboards opens degrades UX | Low — realistic max ~10 windows | Cap at 15 windows, show notice if exceeded |
| `file://` protocol blocks localStorage on some browsers (Safari) | Medium — theme persistence lost | Detect localStorage availability on init, disable persistence gracefully with console warning |

## Open Questions (Resolved)

- [x] **Super+e follows keys.py**: `Super+e` opens Yazi (file manager), not Neovim. The mockup mirrors the real Qtile config. Neovim is still simulated as a window type but launched via a separate shortcut or pre-placed on a workspace.
- [x] **Floating mode visual**: Super+t applies `transform: translateY(-4px)` + `box-shadow: 0 8px 24px rgba(0,0,0,0.4)` to the focused window, plus a subtle "floating" badge in the title bar.
- [x] **Monoid Nerd Font CDN**: If no reliable CDN exists, system monospace is the safe fallback. The mockup will attempt CDN load first, then fall back gracefully. Already acceptable per spec.
