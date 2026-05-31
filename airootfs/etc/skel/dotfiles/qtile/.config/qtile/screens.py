from libqtile import bar
from libqtile.config import Screen
from libqtile.widget import (
    Battery,
    Clock,
    CurrentLayout,
    GroupBox,
    Systray,
    Volume,
    WindowName,
)
from theme import load_theme

colors = load_theme()


def create_bar():
    return bar.Bar(
        [
            GroupBox(
                font="Monoid Nerd Font",
                fontsize=14,
                active=colors["fg"],
                inactive=colors["inactive"],
                highlight_color=colors["primary"],
                highlight_method="block",
                this_current_screen_border=colors["primary"],
                borderwidth=2,
                padding=6,
            ),
            CurrentLayout(
                font="Monoid Nerd Font",
                foreground=colors["accent"],
            ),
            WindowName(
                font="Monoid Nerd Font",
                fontsize=12,
                foreground=colors["fg"],
            ),
            Systray(),
            Volume(
                font="Monoid Nerd Font",
                foreground=colors["secondary"],
            ),
            Battery(
                font="Monoid Nerd Font",
                foreground=colors["secondary"],
                charge_char="\uf0e7",
                discharge_char="\uf240",
            ),
            Clock(
                format="%a %d/%m %H:%M",
                font="Monoid Nerd Font",
                foreground=colors["fg"],
            ),
        ],
        32,
        background=colors["bar_bg"],
        opacity=0.95,
    )


screens = [Screen(top=create_bar())]
