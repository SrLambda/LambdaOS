from libqtile.config import Key
from libqtile.lazy import lazy

mod = "mod4"
terminal = "kitty"

keys = [
    Key([mod], "Return", lazy.spawn(terminal), desc="Launch terminal"),
    Key([mod], "r", lazy.spawn("rofi -show drun"), desc="Launch rofi"),
    Key([mod], "b", lazy.spawn("chromium"), desc="Launch browser"),
    Key([mod], "e", lazy.spawn("yazi"), desc="Launch file manager"),

    Key([mod], "q", lazy.window.kill(), desc="Kill focused window"),
    Key([mod, "shift"], "q", lazy.shutdown(), desc="Shutdown Qtile"),
    Key([mod], "f", lazy.window.toggle_fullscreen(), desc="Toggle fullscreen"),
    Key([mod], "t", lazy.window.toggle_floating(), desc="Toggle floating"),

    Key([mod], "Tab", lazy.next_layout(), desc="Next layout"),
    Key([mod, "shift"], "Tab", lazy.prev_layout(), desc="Previous layout"),

    Key([mod], "h", lazy.layout.left(), desc="Focus left"),
    Key([mod], "j", lazy.layout.down(), desc="Focus down"),
    Key([mod], "k", lazy.layout.up(), desc="Focus up"),
    Key([mod], "l", lazy.layout.right(), desc="Focus right"),

    Key([mod, "shift"], "h", lazy.layout.shuffle_left(), desc="Move window left"),
    Key([mod, "shift"], "j", lazy.layout.shuffle_down(), desc="Move window down"),
    Key([mod, "shift"], "k", lazy.layout.shuffle_up(), desc="Move window up"),
    Key([mod, "shift"], "l", lazy.layout.shuffle_right(), desc="Move window right"),

    Key([mod, "control"], "h", lazy.layout.grow_left(), desc="Grow left"),
    Key([mod, "control"], "j", lazy.layout.grow_down(), desc="Grow down"),
    Key([mod, "control"], "k", lazy.layout.grow_up(), desc="Grow up"),
    Key([mod, "control"], "l", lazy.layout.grow_right(), desc="Grow right"),

    Key([mod], "1", lazy.group["1"].toscreen(), desc="Switch to group 1"),
    Key([mod], "2", lazy.group["2"].toscreen(), desc="Switch to group 2"),
    Key([mod], "3", lazy.group["3"].toscreen(), desc="Switch to group 3"),
    Key([mod], "4", lazy.group["4"].toscreen(), desc="Switch to group 4"),
    Key([mod], "5", lazy.group["5"].toscreen(), desc="Switch to group 5"),
    Key([mod, "shift"], "1", lazy.window.togroup("1"), desc="Move window to group 1"),
    Key([mod, "shift"], "2", lazy.window.togroup("2"), desc="Move window to group 2"),
    Key([mod, "shift"], "3", lazy.window.togroup("3"), desc="Move window to group 3"),
    Key([mod, "shift"], "4", lazy.window.togroup("4"), desc="Move window to group 4"),
    Key([mod, "shift"], "5", lazy.window.togroup("5"), desc="Move window to group 5"),
]
