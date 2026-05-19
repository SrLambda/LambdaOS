from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.containers import Container, Horizontal, Vertical
from textual.widgets import Header, Footer, Switch, ListView, ListItem, Select, Label

from .config_manager import ConfigManager

THEMES = [
    ("Catppuccin", "catppuccin"),
    ("Gruvbox", "gruvbox"),
    ("Tokyonight", "tokyonight"),
    ("Nord", "nord"),
    ("OneDark", "onedark"),
]


class OsTuiConfigurator(App):
    CSS_PATH = "style.tcss"

    TITLE = "LambdaOS \u2014 System Preferences"

    BINDINGS = [
        Binding("ctrl+s", "save", "Guardar"),
        Binding("q", "quit", "Salir"),
    ]

    def __init__(self):
        super().__init__()
        self.config = ConfigManager()
        self.settings = self.config.load_tui_settings()

    def compose(self) -> ComposeResult:
        current_theme = self.config.get_theme() or "catppuccin"

        yield Header()
        with Container(id="app-container"):
            yield ListView(
                ListItem(Label("Neovim")),
                ListItem(Label("Qtile")),
                id="sidebar",
            )
            with Vertical(id="content"):
                with Vertical(id="neovim-content"):
                    yield Label("Neovim Configuration", id="content-title")
                    yield Select(
                        options=THEMES,
                        value=current_theme,
                        prompt="Theme",
                        id="theme-select",
                    )
                    with Horizontal(classes="switch-row"):
                        yield Label("enable_lsp")
                        yield Switch(value=self.settings.get("enable_lsp", True), id="switch_lsp")
                    with Horizontal(classes="switch-row"):
                        yield Label("enable_copilot")
                        yield Switch(value=self.settings.get("enable_copilot", True), id="switch_copilot")
                    with Horizontal(classes="switch-row"):
                        yield Label("enable_neotree")
                        yield Switch(value=self.settings.get("enable_neotree", True), id="switch_neotree")
                yield Label("Qtile Configuration \u2014 Coming Soon", id="qtile-content")
        yield Footer()

    def on_mount(self) -> None:
        sidebar = self.query_one("#sidebar", ListView)
        sidebar.index = 0
        self._show_neovim()

    def on_list_view_selected(self, event: ListView.Selected) -> None:
        if event.list_view.index == 0:
            self._show_neovim()
        else:
            self._show_qtile()

    def _show_neovim(self) -> None:
        neovim_content = self.query_one("#neovim-content")
        qtile_content = self.query_one("#qtile-content")
        neovim_content.visible = True
        qtile_content.visible = False

    def _show_qtile(self) -> None:
        neovim_content = self.query_one("#neovim-content")
        qtile_content = self.query_one("#qtile-content")
        neovim_content.visible = False
        qtile_content.visible = True

    def on_switch_changed(self, event: Switch.Changed) -> None:
        switch = event.switch
        if switch.id == "switch_lsp":
            self.settings["enable_lsp"] = event.value
        elif switch.id == "switch_copilot":
            self.settings["enable_copilot"] = event.value
        elif switch.id == "switch_neotree":
            self.settings["enable_neotree"] = event.value

    def on_select_changed(self, event: Select.Changed) -> None:
        if event.select.id == "theme-select" and event.value != Select.BLANK:
            self.config.set_theme(str(event.value))
            self.config.save_os_theme(str(event.value))

    def action_save(self) -> None:
        self.config.save_tui_settings(self.settings)
        self.notify("Configuración guardada")
