# Tasks: Qtile Interface HTML Mockup

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 1000–1200 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 (HTML+CSS) → PR 2 (JS functionality) → PR 3 (Polish+README) |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | HTML skeleton + all CSS (themes, responsive) | PR 1 | Base = main tracker branch. Visual-only, fully verifiable in browser (~400 lines) |
| 2 | JS state + layout engine + windows + keyboard | PR 2 | Base = PR 1 branch. Adds interactivity to PR 1's structure (~500 lines) |
| 3 | Polish — floating mode, transitions, README | PR 3 | Base = PR 2 branch. Minor additions + docs (~100 lines) |

## Phase 1: HTML Structure

- [x] 1.1 Create `preview/index.html` with `<header id="bar">`, `<main id="desktop">`, theme selector, keybindings panel, notification container
- [x] 1.2 Build GroupBox: 5 workspace indicators with FontAwesome icons (terminal, code, files, music, entertainment)
- [x] 1.3 Build bar widgets row: CurrentLayout, WindowName, Systray placeholders, Volume, Battery, Clock
- [x] 1.4 Build simulated app divs: Kitty terminal with prompt, Chromium with tabs/address bar, Yazi with TUI file list + preview

## Phase 2: CSS — Base + Themes

- [x] 2.1 Define CSS custom properties: `:root[data-theme="catppuccin"]` block with all 9 tokens (verified exact match with theme.py)
- [x] 2.2 Add remaining 4 theme blocks: gruvbox, tokyonight, nord, onedark (exact colors from theme.py)
- [x] 2.3 Style top bar: 32px height, semi-transparent `--bar-bg`, `--bar-fg` text, flex layout for widgets (verified complete from Phase 1)
- [x] 2.4 Style windows: title bar, app content areas, focused (`--primary`) and unfocused (`--inactive`) borders (verified complete from Phase 1)
- [x] 2.5 Style theme selector dropdown and keybindings panel with theme-aware colors (verified complete from Phase 1)
- [x] 2.6 Style notification toasts with `--primary` background, auto-dismiss animation (verified complete from Phase 1)

## Phase 3: CSS — Responsive

- [x] 3.1 Add `@media (max-width: 480px)` block: minimal bar (clock + icons), stacked windows, fullscreen keybindings panel
- [x] 3.2 Add `@media (max-width: 768px)` block: full bar widgets, MonadTall main at 50%, max 2 columns
- [x] 3.3 Add `@media (min-width: 769px)` block: full desktop layout with proper proportions
- [x] 3.4 Add `transition: left/top/width/height 200ms` on windows for smooth layout switching (verified present from Phase 1)

## Phase 4: JS — State Management + Theme Engine

- [x] 4.1 Create `state` object: theme, layout, activeWorkspace, workspaces map (1–5), focusedWindowId, superHeld, nextWindowId
- [x] 4.2 Implement theme engine: `setAttribute('data-theme')` on `:root`, `localStorage` get/set with Catppuccin default
- [x] 4.3 Implement live clock widget: `setInterval` updating Clock element in `%a %d/%m %H:%M` format
- [x] 4.4 Implement theme selector onChange handler and localStorage restore on page load

## Phase 5: JS — Layout Engine

- [x] 5.1 Implement `computeMonadTall(windows)`: main at 60% left, stack remaining on right
- [x] 5.2 Implement `computeMax(focusedWindow)`: single window fills desktop, others hidden
- [x] 5.3 Implement `computeColumns(windows)`: equal-width side-by-side
- [x] 5.4 Implement `applyLayout()`: calls compute function, sets `left/top/width/height` on each window div

## Phase 6: JS — Window Management

- [x] 6.1 Implement `createWindow(type)`: create window div, push to active workspace, auto-focus
- [x] 6.2 Implement `destroyWindow(id)`: remove from workspace, focus next, handle empty workspace
- [x] 6.3 Implement click-to-front: `mousedown` handler brings clicked window to front + sets focus
- [x] 6.4 Implement workspace switching: move windows between workspace arrays, re-render on switch

## Phase 7: JS — Keyboard Handler

- [x] 7.1 Implement global `keydown`/`keyup` listeners on `document`, track Super/Meta key state
- [x] 7.2 Build `SHORTCUT_MAP`: map all 14 shortcuts to action functions (open apps, close, cycle layouts, switch workspaces, move windows, toggle float)
- [x] 7.3 Implement Super badge UI: show/hide "SUPER" indicator in corner
- [x] 7.4 Implement notification system: show toast with shortcut description, auto-dismiss after 2s

## Phase 8: JS — Keybindings Panel

- [x] 8.1 Build categorized shortcuts data structure (Launchers, Window Mgmt, Layouts, Workspaces)
- [x] 8.2 Implement toggle handler: show/hide `#keybindings` panel on button click
- [x] 8.3 Render categorized table rows using the active theme's colors

## Phase 9: Polish + README

- [x] 9.1 Implement floating mode visual: Super+t applies translateY(-4px) + box-shadow + floating badge
- [x] 9.2 Ensure all 200ms CSS transitions are applied and smooth
- [x] 9.3 Create `preview/README.md` with usage instructions, shortcuts list, theme list, no-server note
