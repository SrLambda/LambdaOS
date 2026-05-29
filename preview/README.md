# LambdaOS — Qtile Desktop Mockup

An interactive HTML/CSS/JS preview of the LambdaOS Qtile desktop interface.

## Quick Start

Open `index.html` in any modern browser. No server required.

```bash
# macOS
open index.html

# Linux
xdg-open index.html

# Windows
start index.html
```

## Features

- **6 Themes**: LambdaOS (default), Catppuccin, Gruvbox, Tokyo Night, Nord, One Dark — switch via the dropdown in the top-right corner
- **3 Layouts**: MonadTall, Max, Columns — cycle with Super+Tab
- **5 Workspaces**: Click the workspace icons in the bar or use Super+1-5
- **Keyboard Shortcuts**: Full Qtile keybinding simulation (see Keybindings section)
- **Responsive Design**: Works on desktop, tablet, and mobile viewports
- **Theme Persistence**: Your theme choice is saved in localStorage

## Keybindings

All shortcuts use the Super (⌘/Win) key as the modifier.

### Launchers
| Shortcut | Action |
|----------|--------|
| Super+Return | Open Kitty terminal |
| Super+b | Open Chromium browser |
| Super+e | Open Yazi file manager |

### Window Management
| Shortcut | Action |
|----------|--------|
| Super+q | Close focused window |
| Super+f | Toggle fullscreen |
| Super+t | Toggle floating mode |

### Navigation
| Shortcut | Action |
|----------|--------|
| Super+h/j/k/l | Focus left/down/up/right |

### Layouts
| Shortcut | Action |
|----------|--------|
| Super+Tab | Next layout |
| Super+Shift+Tab | Previous layout |

### Workspaces
| Shortcut | Action |
|----------|--------|
| Super+1-5 | Switch to workspace |
| Super+Shift+1-5 | Move window to workspace |

## Themes

| Theme | Background | Primary | Accent |
|-------|-----------|---------|--------|
| **LambdaOS** | `#0B0F19` | `#6D40FF` | `#2E00C7` |
| Catppuccin | `#1e1e2e` | `#cba6f7` | `#f5c2e7` |
| Gruvbox | `#282828` | `#d79921` | `#b16286` |
| Tokyo Night | `#1a1b26` | `#bb9af7` | `#ff9e64` |
| Nord | `#2e3440` | `#81a1c1` | `#b48ead` |
| One Dark | `#282c34` | `#c678dd` | `#e5c07b` |

## Technical Notes

- **No build step**: Pure HTML, CSS, and JavaScript
- **Works offline**: After first load (CDN fonts may need initial connection)
- **CDN dependencies**: FontAwesome 6.5 for icons, Monoid Nerd Font for typography
- **Browser support**: Chrome, Firefox, Safari, Edge (modern versions)
- **localStorage**: Used for theme persistence (not available in Safari with file:// protocol)

## Known Limitations

- The Super key may be intercepted by your OS (e.g., GNOME overview, Windows Start menu)
- This is a visual mockup only — it does not replace Qtile
- Window management is simulated, not a real WM

## Project Structure

```
preview/
├── index.html    # Self-contained desktop mockup
├── logo.tsx      # LambdaOS brand: logo SVG, color palette, neofetch/boot mockups
└── README.md     # This file
```

## Brand

The LambdaOS brand uses a geometric lowercase lambda (λ) as its symbol:

- **Logo**: SVG with gradient from `#6D40FF` (Neon Purple) to `#2E00C7` (Deep Lambda)
- **ASCII Art**: Terminal-friendly lambda shape using `/`, `\`, `S`, `O`, and `·` characters
- **Full palette**: See `logo.tsx` for the complete brand system including gradient demos and boot screen mockup

## Related

- [LambdaOS README](../README.md)
- [Qtile Config](../airootfs/etc/skel/dotfiles/qtile/.config/qtile/config.py)
