# Qtile Interface HTML Mockup — Specs

> Change: `qtile-html-mockup`
> Created: 2026-05-24
> Purpose: Interactive HTML mockup of the LambdaOS Qtile desktop for development reference and user showcase

---

## Spec 1: Estructura del Desktop

### Requirement
The mockup SHALL render a faithful HTML/CSS representation of the LambdaOS Qtile desktop layout, including the top bar, desktop area, and window management zones.

### Acceptance Criteria
- [ ] A top bar spanning the full width of the viewport, 32px height
- [ ] A desktop area filling the remaining viewport below the bar
- [ ] A solid background color matching the active theme's `bg` color
- [ ] The layout is responsive: adapts to any viewport width >= 320px

### Scenarios

**Given** a user opens `preview/index.html` in a browser
**When** the page loads
**Then** the desktop renders with a top bar and a colored background area

**Given** the viewport is resized to a narrow width (e.g., 480px)
**When** the page reflows
**Then** the bar remains visible and the desktop area scales proportionally

---

## Spec 2: Top Bar Components

### Requirement
The top bar SHALL contain the same widgets as the actual Qtile configuration in `screens.py`, rendered in the same order and with equivalent visual styling.

### Bar Widget Order (left to right)
1. **GroupBox** — 5 workspace indicators with FontAwesome icon equivalents
2. **CurrentLayout** — shows active layout name/icon
3. **WindowName** — shows focused window title
4. **Systray** — placeholder icon area
5. **Volume** — volume indicator with icon
6. **Battery** — battery indicator with icon
7. **Clock** — live clock showing `%a %d/%m %H:%M` format

### Acceptance Criteria
- [ ] GroupBox shows 5 workspace indicators with icons matching the Qtile groups:
  - `1` → terminal icon (`\uf120`)
  - `2` → code icon (`\uf269`)
  - `3` → files icon (`\uf1c9`)
  - `4` → music icon (`\uf07c`)
  - `5` → entertainment icon (`\uf1bc`)
- [ ] Active workspace is visually highlighted with `highlight_method="block"` style
- [ ] CurrentLayout displays the active layout name (e.g., "MonadTall", "Max", "Columns")
- [ ] WindowName displays the title of the focused simulated window
- [ ] Systray shows placeholder icons (network, volume, battery)
- [ ] Volume shows a speaker icon with a numeric value
- [ ] Battery shows a battery icon with percentage
- [ ] Clock updates every second in the format `ddd dd/mm HH:MM`
- [ ] Bar background uses `bar_bg` color from the active theme
- [ ] Bar opacity is visually set to ~0.95 (semi-transparent effect)

### Scenarios

**Given** the mockup is rendered with the Catppuccin theme
**When** the user views the top bar
**Then** the bar background is `#181825`, text is `#cdd6f4`, and the active workspace has a `#cba6f7` highlight block

**Given** the clock widget is visible
**When** 60 seconds pass
**Then** the displayed time has updated to reflect the current time

---

## Spec 3: Theme System

### Requirement
The mockup SHALL support all 5 themes from the Qtile `theme.py` configuration, with an interactive selector that persists the user's choice.

### Supported Themes
| Theme | bg | fg | primary | secondary | accent | urgent | inactive | bar_bg |
|-------|-----|-----|---------|-----------|--------|--------|----------|--------|
| catppuccin | #1e1e2e | #cdd6f4 | #cba6f7 | #89b4fa | #f5c2e7 | #f38ba8 | #45475a | #181825 |
| gruvbox | #282828 | #ebdbb2 | #d79921 | #458588 | #b16286 | #cc241d | #3c3836 | #1d2021 |
| tokyonight | #1a1b26 | #c0caf5 | #bb9af7 | #7aa2f7 | #ff9e64 | #f7768e | #3b4261 | #16161e |
| nord | #2e3440 | #d8dee9 | #81a1c1 | #88c0d0 | #b48ead | #bf616a | #434c5e | #3b4252 |
| onedark | #282c34 | #abb2bf | #c678dd | #61afef | #e5c07b | #e06c75 | #3e4452 | #21252b |

### Acceptance Criteria
- [ ] A theme selector UI is visible (dropdown or button group)
- [ ] Selecting a theme immediately updates ALL colors across the entire mockup
- [ ] The selected theme is persisted in `localStorage` under key `lambdaos-qtile-theme`
- [ ] On page reload, the previously selected theme is restored from `localStorage`
- [ ] If no theme is stored, Catppuccin is the default (matching `FALLBACK_THEME` in `theme.py`)
- [ ] Window borders use `primary` for focused and `inactive` for unfocused windows
- [ ] Bar text uses `bar_fg`, bar background uses `bar_bg`

### Scenarios

**Given** no theme is stored in localStorage
**When** the page loads for the first time
**Then** Catppuccin theme colors are applied by default

**Given** the user selects "Gruvbox" from the theme selector
**When** the selection is made
**Then** all colors update to Gruvbox values AND `localStorage.setItem("lambdaos-qtile-theme", "gruvbox")` is called

**Given** the user previously selected "Nord"
**When** the page is reloaded
**Then** Nord theme colors are applied automatically

---

## Spec 4: Simulated Windows

### Requirement
The desktop area SHALL display simulated application windows that represent the actual apps used in LambdaOS, styled to resemble their real appearance.

### Simulated Applications
1. **Kitty Terminal** — dark background with monospace font, showing a shell prompt with LambdaOS branding
2. **Chromium Browser** — browser chrome with address bar, tabs, and a sample webpage content area
3. **Yazi File Manager** — TUI-style file manager with file list, preview pane, and status bar

### Acceptance Criteria
- [ ] Each simulated window has a title bar with the app name and close/minimize buttons
- [ ] Focused window has `primary` color border; unfocused windows have `inactive` color border
- [ ] Clicking a window brings it to front and marks it as focused
- [ ] Kitty terminal shows:
  - Dark background (`bg` color)
  - Monospace font (Monoid Nerd Font or fallback)
  - A shell prompt like `lambda@lambdaos ~ %`
  - At least 3 lines of simulated terminal output
- [ ] Chromium shows:
  - Tab bar with at least one tab
  - Address bar with a sample URL
  - Content area with placeholder content (e.g., "LambdaOS — Welcome")
- [ ] Yazi shows:
  - TUI-style layout with file/folder list on the left
  - Preview pane on the right showing selected file content
  - Status bar at bottom with current path and mode indicator

### Scenarios

**Given** three windows are visible on the desktop
**When** the user clicks the Yazi window
**Then** Yazi's border changes to `primary` color, Kitty and Chromium borders change to `inactive`, and WindowName in the bar updates to "Yazi"

**Given** the Kitty window is focused
**When** the theme changes from Catppuccin to Tokyo Night
**Then** Kitty's background, text, and border colors all update to Tokyo Night values

---

## Spec 5: Layout Simulation

### Requirement
The mockup SHALL support switching between the 3 Qtile layouts (MonadTall, Max, Columns), visually rearranging the simulated windows to match each layout's behavior.

### Layout Behaviors
- **MonadTall**: One main window on the left (60% width), remaining windows stacked on the right
- **Max**: Single focused window fills the entire desktop area
- **Columns**: Windows arranged in equal-width columns side by side

### Acceptance Criteria
- [ ] CurrentLayout widget in the bar reflects the active layout
- [ ] Switching layouts repositions all visible windows according to the layout rules
- [ ] Window proportions are visually accurate for each layout
- [ ] Layout transitions are smooth (CSS transition of ~200ms)

### Scenarios

**Given** the mockup is in MonadTall layout with 3 windows visible
**When** the user switches to Max layout
**Then** the focused window expands to fill the entire desktop area and other windows are hidden behind it

**Given** the mockup is in Max layout
**When** the user switches to Columns layout
**Then** all 3 windows are displayed side by side with equal widths

---

## Spec 6: Keyboard Shortcut Simulation

### Requirement
The mockup SHALL respond to keyboard shortcuts that mirror the actual Qtile keybindings from `keys.py`, providing visual feedback for each action.

### Supported Shortcuts (all using Super/Mod4 as the modifier)
| Shortcut | Action |
|----------|--------|
| `Super+Return` | Open a new Kitty terminal window |
| `Super+b` | Open a new Chromium window |
| `Super+e` | Open a new Yazi file manager window |
| `Super+q` | Close the focused window |
| `Super+Tab` | Cycle to next layout |
| `Super+Shift+Tab` | Cycle to previous layout |
| `Super+1` through `Super+5` | Switch to workspace 1-5 |
| `Super+Shift+1` through `Super+Shift+5` | Move focused window to workspace 1-5 |
| `Super+h/j/k/l` | Focus window left/down/up/right |
| `Super+t` | Toggle floating mode on focused window |

### Acceptance Criteria
- [ ] A visual indicator shows when the Super key is held (e.g., "SUPER" badge in corner)
- [ ] Pressing a supported shortcut executes the corresponding action
- [ ] An on-screen notification shows the shortcut description (matching the `desc` field from `keys.py`)
- [ ] Notifications auto-dismiss after 2 seconds
- [ ] Shortcuts work regardless of which element has focus (global key listener)

### Scenarios

**Given** the user presses `Super+Tab`
**When** the current layout is MonadTall
**Then** the layout changes to Max, the CurrentLayout widget updates, and a notification "Next layout" appears briefly

**Given** the user presses `Super+q` with Yazi focused
**When** the shortcut is executed
**Then** the Yazi window is removed from the desktop and focus shifts to the next available window

---

## Spec 7: Workspace System

### Requirement
The mockup SHALL implement a workspace system with 5 workspaces, where windows can be assigned to and viewed per workspace.

### Acceptance Criteria
- [ ] Clicking a workspace indicator in the GroupBox switches to that workspace
- [ ] Each workspace maintains its own set of visible windows independently
- [ ] The active workspace is highlighted with `primary` color block
- [ ] Empty workspaces show only the desktop background with no windows
- [ ] `Super+1` through `Super+5` keyboard shortcuts switch workspaces

### Scenarios

**Given** workspace 1 has a Kitty window and workspace 2 is empty
**When** the user clicks workspace 2 in the GroupBox
**Then** the desktop shows only the background (no windows) and workspace 2 is highlighted

**Given** a Yazi window exists on workspace 3
**When** the user presses `Super+Shift+1`
**Then** the Yazi window moves to workspace 1 and disappears from workspace 3

---

## Spec 8: Keybindings Documentation Section

### Requirement
The mockup SHALL include a dedicated section documenting all Qtile keybindings, accessible from the mockup interface.

### Acceptance Criteria
- [ ] A "Keybindings" button or link is visible in the mockup (e.g., in a corner or as a toggleable panel)
- [ ] The keybindings section displays all shortcuts from `keys.py` in a readable table format
- [ ] Shortcuts are grouped by category:
  - **Launchers**: terminal, rofi, browser, file manager
  - **Window Management**: kill, fullscreen, floating, focus, move, resize
  - **Layouts**: next/previous layout
  - **Workspaces**: switch, move window
- [ ] The keybindings section is styled consistently with the active theme
- [ ] The section can be toggled open/closed without leaving the mockup

### Scenarios

**Given** the user clicks the "Keybindings" button
**When** the keybindings panel opens
**Then** all Qtile shortcuts are displayed in categorized tables with the current theme's colors

**Given** the keybindings panel is open
**When** the user changes the theme
**Then** the keybindings panel colors update to match the new theme

---

## Spec 9: Responsive Design

### Requirement
The mockup SHALL be fully responsive and usable on viewports from 320px (mobile) to 2560px+ (ultrawide).

### Acceptance Criteria
- [ ] At 320px-480px (mobile): bar widgets collapse gracefully or become scrollable; windows stack vertically
- [ ] At 481px-768px (tablet): bar is fully visible; windows adapt to available space
- [ ] At 769px+ (desktop): full layout as specified
- [ ] No horizontal scrolling at any viewport width
- [ ] Touch-friendly: workspace indicators and theme selector are tappable on mobile

### Scenarios

**Given** the viewport is 375px wide (iPhone SE)
**When** the mockup is viewed
**Then** the bar is visible, windows stack vertically, and all interactive elements are tappable

**Given** the viewport is 1920px wide (standard desktop)
**When** the mockup is viewed
**Then** the full desktop layout renders as specified with proper proportions

---

## Spec 10: File Structure and Delivery

### Requirement
The mockup SHALL be delivered as a self-contained HTML file (or minimal file set) in a `preview/` directory at the repository root.

### Acceptance Criteria
- [ ] Entry point is `preview/index.html`
- [ ] CSS and JavaScript are either inline or in separate files within `preview/`
- [ ] No external build step required — opens directly in a browser
- [ ] FontAwesome or equivalent icon font is loaded via CDN for workspace icons and widget icons
- [ ] Monoid Nerd Font is loaded via CDN or falls back gracefully to a monospace font
- [ ] The `preview/` directory contains only the mockup files (no build artifacts)
- [ ] A `preview/README.md` documents how to open and use the mockup

### Scenarios

**Given** a user clones the LambdaOS repository
**When** they open `preview/index.html` in any modern browser
**Then** the mockup loads and functions without errors

**Given** the user has no internet connection
**When** they open `preview/index.html`
**Then** the mockup still renders (icons may show fallback characters, fonts fall back to system monospace)

---

## Technical Constraints

- **No frameworks required**: Vanilla HTML, CSS, and JavaScript only
- **Single-file preference**: If possible, keep everything in one `index.html` for simplicity
- **CDN dependencies allowed**: FontAwesome for icons, Google Fonts or CDN for Monoid Nerd Font
- **Browser support**: Modern browsers (Chrome, Firefox, Safari, Edge) — no IE support needed
- **No server required**: Must work with `file://` protocol

---

## Out of Scope

- Actual Qtile functionality (this is a visual mockup, not a WM)
- Backend integration with the real OS
- Multi-monitor simulation
- Custom widget configuration beyond what's in the existing Qtile config
- Animation beyond CSS transitions for layout changes
