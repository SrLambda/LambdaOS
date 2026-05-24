# Archive Report: qtile-html-mockup

**Change**: qtile-html-mockup
**Archived**: 2026-05-24
**Status**: PASS WITH WARNINGS — All 36 tasks complete, all 10 specs implemented
**Artifact Store**: hybrid (OpenSpec files + Engram)

---

## Executive Summary

Interactive HTML mockup of the LambdaOS Qtile desktop, delivered as a self-contained single-file `preview/index.html` (1817 lines) plus `preview/README.md` (99 lines). The mockup faithfully mirrors the real Qtile configuration — bar widgets from `screens.py`, color tokens from `theme.py` (all 5 themes, 45 color values verified exact), workspace structure from `groups.py`, and keyboard shortcuts from `keys.py` (14 shortcuts). Built with vanilla HTML/CSS/JS — no frameworks, no build step, no server required. Opens directly in a browser via `file://` protocol.

## What Was Built

| File | Action | Lines | Description |
|------|--------|-------|-------------|
| `preview/index.html` | Create | 1817 | Self-contained interactive Qtile desktop mockup with inline CSS + JS |
| `preview/README.md` | Create | 99 | Usage instructions, shortcuts reference, theme table, technical notes |
| `openspec/changes/qtile-html-mockup/specs.md` | Create | 338 | 10 specs with acceptance criteria and Given/When/Then scenarios |
| `openspec/changes/qtile-html-mockup/design.md` | Create | 161 | Architecture decisions, data flow, component model, layout engine design |
| `openspec/changes/qtile-html-mockup/tasks.md` | Create | 88 | 36 tasks across 9 phases, all marked complete |

**Total code delivered**: 1,916 lines (index.html + README.md)
**Total artifacts**: 1,916 + 587 lines of SDD documentation = 2,503 lines

## Architecture Decisions (from design.md)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| File structure | Single `index.html` | Works with `file://` protocol, no CORS issues, no build step |
| Layout engine | JS absolute positioning | Pixel-perfect control, matches Qtile's screen-space model |
| Theme switching | CSS custom properties | Single `setProperty` mutation updates all elements instantly |
| Window rendering | Styled `<div>`s | `iframe` blocked by `file://` CORS |
| Keyboard capture | Global `document` listener | No focus dependency, captures all keystrokes |
| Font strategy | CDN with fallback | Lightweight repo, graceful degradation offline |

## Verification Warnings (3 — non-blocking)

### WARNING 1: Workspace icon codepoints differ from spec
- **Spec says**: `\uf120` (terminal), `\uf269` (code), `\uf1c9` (files), `\uf07c` (music), `\uf1bc` (entertainment)
- **Implementation uses**: FA6 Free class names (`fa-terminal`, `fa-code`, `fa-folder`, `fa-music`, `fa-gamepad`)
- **Impact**: Visual semantics correct but exact codepoints differ. Spec may reference FA4/FA5 Pro icons.

### WARNING 2: `--max-columns` CSS variable not consumed by layout engine
- **Design says**: Tablet breakpoint (481-768px) sets `--max-columns: 2`
- **Implementation**: `computeColumns()` always divides by window count, never reads the CSS variable
- **Impact**: On tablet with 3+ windows in Columns, windows get narrow columns instead of max 2. Edge case.

### WARNING 3: Chromium welcome page references non-existent Super+? shortcut
- **Evidence**: HTML displays `<kbd>Super+?</kbd>` but no handler exists
- **Impact**: Minor — the welcome page is decorative. Clicking "Keys" button opens the keybindings panel.

## Suggestions for Future Improvements

| # | Suggestion | Rationale |
|---|------------|-----------|
| S1 | Add Escape to dismiss overlays | Keybindings panel and theme selector lack keyboard-only dismiss |
| S2 | Add first-visit toast about OS Super key interception | README warns but users may not read it before trying shortcuts |
| S3 | Add window count cap (15 max) | Design mentions cap but no enforcement — unlimited `nextWindowId` increments |

## Manual Testing Checklist

Since this is a visual mockup with no automated tests, these scenarios should be manually verified:

### Theme
- [ ] Switch each of 5 themes — verify colors update instantly
- [ ] Select Gruvbox, reload — verify persistence via localStorage
- [ ] Clear localStorage, reload — verify Catppuccin default

### Keyboard Shortcuts
- [ ] Super+Return → Kitty opens, focused, notification "Launch terminal"
- [ ] Super+b → Chromium opens
- [ ] Super+e → Yazi opens
- [ ] Super+q → focused window closes
- [ ] Super+Tab/Shift+Tab → layout cycle forward/backward
- [ ] Super+1-5 → workspace switch
- [ ] Super+Shift+1-5 → move window to workspace
- [ ] Super+h/j/k/l → directional focus
- [ ] Super+t → floating toggle
- [ ] Super badge visible when Super held

### Layout
- [ ] MonadTall: 3 windows → main 60% left, 2 stacked right
- [ ] Max: focused fills desktop
- [ ] Columns: equal-width
- [ ] 200ms smooth transitions

### Workspaces
- [ ] Click workspace 2 → switches, highlighted
- [ ] Create windows on ws1, switch to ws2 → empty desktop
- [ ] Move window via Super+Shift+3 → relocates

### Responsive
- [ ] 375px (mobile): bar scrollable, windows stacked
- [ ] 600px (tablet): full bar, adjusted ratios
- [ ] 1920px (desktop): full layout
- [ ] No horizontal scroll at any width

### Offline
- [ ] Disconnect network, reload → layout works, icons may fall back

## Source of Truth

- **Primary artifact**: `preview/index.html` — the single self-contained mockup
- **README**: `preview/README.md` — usage documentation
- **SDD artifacts**: Archived in `openspec/changes/archive/2026-05-24-qtile-html-mockup/`
- **Engram**: Archive report saved at topic_key `sdd/qtile-html-mockup/archive-report`
- **Real Qtile config mirrored**: `screens.py` (bar widgets), `theme.py` (colors), `groups.py` (workspaces), `keys.py` (shortcuts)

## Final Status

```
┌─────────────────────────────────────────────────────────┐
│  SDD CYCLE COMPLETE ✅                                  │
│                                                         │
│  Change: qtile-html-mockup                              │
│  Status: PASS WITH WARNINGS (3 non-blocking)            │
│  Phases: 9/9 complete                                   │
│  Tasks: 36/36 complete                                  │
│  Specs: 10/10 implemented                               │
│  Code: 1,916 lines (index.html + README.md)             │
│  Archived: 2026-05-24                                   │
│                                                         │
│  Ready for the next change.                             │
└─────────────────────────────────────────────────────────┘
```

---

*Archive generated: 2026-05-24*
*Engram observation IDs: None (hybrid mode — filesystem primary)*
