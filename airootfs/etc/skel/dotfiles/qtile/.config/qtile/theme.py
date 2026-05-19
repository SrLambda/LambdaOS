import json
import os
from pathlib import Path

FALLBACK_THEME = "catppuccin"

THEMES = {
    "catppuccin": {
        "bg": "#1e1e2e",
        "fg": "#cdd6f4",
        "primary": "#cba6f7",
        "secondary": "#89b4fa",
        "accent": "#f5c2e7",
        "urgent": "#f38ba8",
        "inactive": "#45475a",
        "bar_bg": "#181825",
        "bar_fg": "#cdd6f4",
    },
    "gruvbox": {
        "bg": "#282828",
        "fg": "#ebdbb2",
        "primary": "#d79921",
        "secondary": "#458588",
        "accent": "#b16286",
        "urgent": "#cc241d",
        "inactive": "#3c3836",
        "bar_bg": "#1d2021",
        "bar_fg": "#ebdbb2",
    },
    "tokyonight": {
        "bg": "#1a1b26",
        "fg": "#c0caf5",
        "primary": "#bb9af7",
        "secondary": "#7aa2f7",
        "accent": "#ff9e64",
        "urgent": "#f7768e",
        "inactive": "#3b4261",
        "bar_bg": "#16161e",
        "bar_fg": "#c0caf5",
    },
    "nord": {
        "bg": "#2e3440",
        "fg": "#d8dee9",
        "primary": "#81a1c1",
        "secondary": "#88c0d0",
        "accent": "#b48ead",
        "urgent": "#bf616a",
        "inactive": "#434c5e",
        "bar_bg": "#3b4252",
        "bar_fg": "#d8dee9",
    },
    "onedark": {
        "bg": "#282c34",
        "fg": "#abb2bf",
        "primary": "#c678dd",
        "secondary": "#61afef",
        "accent": "#e5c07b",
        "urgent": "#e06c75",
        "inactive": "#3e4452",
        "bar_bg": "#21252b",
        "bar_fg": "#abb2bf",
    },
}


def load_theme():
    config_base = os.environ.get("OS_CONFIG_DIR", str(Path.home()))
    theme_path = Path(config_base) / "os_theme.json"

    try:
        with open(theme_path, "r") as f:
            data = json.load(f)
            theme_name = data.get("theme", FALLBACK_THEME)
    except (FileNotFoundError, json.JSONDecodeError, KeyError):
        theme_name = FALLBACK_THEME

    return THEMES.get(theme_name, THEMES[FALLBACK_THEME])
