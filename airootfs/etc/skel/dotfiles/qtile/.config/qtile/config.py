import os
import subprocess
from pathlib import Path

from groups import groups
from keys import keys
from libqtile import hook, qtile
from libqtile.config import Key, Screen
from libqtile.layout import Columns, Max, MonadTall
from theme import load_theme

colors = load_theme()

layouts = [
    MonadTall(
        border_width=2,
        border_focus=colors["primary"],
        border_normal=colors["inactive"],
        margin=4,
    ),
    Max(
        border_width=2,
        border_focus=colors["primary"],
        border_normal=colors["inactive"],
        margin=4,
    ),
    Columns(
        border_width=2,
        border_focus=colors["primary"],
        border_normal=colors["inactive"],
        margin=4,
    ),
]

widget_defaults = dict(
    font="Monoid Nerd Font",
    fontsize=12,
    padding=3,
    foreground=colors["fg"],
)

extension_defaults = widget_defaults.copy()

floating_layout = None


@hook.subscribe.startup_once
def autostart():
    home = Path.home()

    dotfiles_dir = home / "dotfiles"
    if dotfiles_dir.is_dir():
        subprocess.Popen(["stow", "*/"], cwd=dotfiles_dir)

    subprocess.Popen(["xsetroot", "-solid", colors["bg"]])

    subprocess.Popen(["picom", "--experimental-backends"])

    subprocess.Popen(["flameshot"])
