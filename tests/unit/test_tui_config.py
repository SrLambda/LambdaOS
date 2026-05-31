import json
import os

import pytest

from src.os_tui_configurator.app import OsTuiConfigurator
from src.os_tui_configurator.config_manager import ConfigManager


@pytest.fixture
def tmp_config_dir(tmp_path):
    """Crea estructura tests/tmp/nvim/.config/nvim/ para pruebas aisladas."""
    config_dir = tmp_path / "nvim" / ".config" / "nvim"
    config_dir.mkdir(parents=True)

    settings = {"enable_lsp": True, "enable_copilot": False, "enable_neotree": True}
    settings_file = config_dir / "tui_settings.json"
    settings_file.write_text(json.dumps(settings))

    old = os.environ.get("OS_CONFIG_DIR")
    os.environ["OS_CONFIG_DIR"] = str(tmp_path)
    yield tmp_path
    if old is not None:
        os.environ["OS_CONFIG_DIR"] = old
    else:
        os.environ.pop("OS_CONFIG_DIR", None)


class TestConfigManager:
    """Tests unitarios para el ConfigManager."""

    def test_load_tui_settings_returns_dict(self, tmp_config_dir):
        """Carga tui_settings.json y devuelve un diccionario con las banderas."""
        cm = ConfigManager()
        settings = cm.load_tui_settings()

        assert isinstance(settings, dict)
        assert settings == {
            "enable_lsp": True,
            "enable_copilot": False,
            "enable_neotree": True,
        }

    def test_load_tui_settings_missing_file_returns_defaults(self, tmp_config_dir):
        """Si el JSON no existe, devuelve defaults (todo true)."""
        settings_file = (
            tmp_config_dir / "nvim" / ".config" / "nvim" / "tui_settings.json"
        )
        settings_file.unlink()

        cm = ConfigManager()
        settings = cm.load_tui_settings()

        assert settings == {
            "enable_lsp": True,
            "enable_copilot": True,
            "enable_neotree": True,
        }

    def test_save_tui_settings_writes_correct_values(self, tmp_config_dir):
        """Guarda flags y verifica que el archivo se escribio correctamente."""
        cm = ConfigManager()
        new_settings = {
            "enable_lsp": False,
            "enable_copilot": True,
            "enable_neotree": False,
        }
        cm.save_tui_settings(new_settings)

        settings_file = (
            tmp_config_dir / "nvim" / ".config" / "nvim" / "tui_settings.json"
        )
        saved = json.loads(settings_file.read_text())

        assert saved == new_settings

    def test_save_and_reload_preserves_data(self, tmp_config_dir):
        """Round-trip: guarda, recarga, verifica que los datos son identicos."""
        cm = ConfigManager()
        new_settings = {
            "enable_lsp": False,
            "enable_copilot": True,
            "enable_neotree": False,
        }
        cm.save_tui_settings(new_settings)

        reloaded = cm.load_tui_settings()

        assert reloaded == new_settings

    def test_set_theme_writes_env_var(self, tmp_config_dir):
        """set_theme('gruvbox') debe escribir NVIM_THEME=gruvbox en archivo de env."""
        cm = ConfigManager()
        cm.set_theme("gruvbox")

        theme_path = tmp_config_dir / "nvim" / ".config" / "nvim" / ".nvim_theme"
        assert theme_path.exists()
        content = theme_path.read_text().strip()
        assert content == "NVIM_THEME=gruvbox"

    def test_get_theme_returns_none_when_no_file(self, tmp_config_dir):
        """get_theme() devuelve None si el archivo de tema no existe."""
        cm = ConfigManager()
        result = cm.get_theme()
        assert result is None

    def test_get_theme_returns_stored_theme(self, tmp_config_dir):
        """get_theme() devuelve el tema guardado previamente."""
        cm = ConfigManager()
        cm.set_theme("tokyonight")
        assert cm.get_theme() == "tokyonight"

    def test_respects_os_config_dir_env(self, tmp_path, monkeypatch):
        """Con OS_CONFIG_DIR apuntando a path custom, las operaciones usan ese path."""
        monkeypatch.setenv("OS_CONFIG_DIR", str(tmp_path))
        cm = ConfigManager()

        expected_path = tmp_path / "nvim" / ".config" / "nvim" / "tui_settings.json"
        assert cm._get_settings_path() == expected_path

    def test_save_os_theme_creates_file(self, tmp_config_dir):
        """save_os_theme escribe os_theme.json con la estructura correcta."""
        cm = ConfigManager()
        cm.save_os_theme("gruvbox")
        theme_path = tmp_config_dir / "os_theme.json"
        assert theme_path.exists()
        data = json.loads(theme_path.read_text())
        assert data == {
            "theme": "gruvbox",
            "nvim_theme": "gruvbox",
            "qtile_theme": "gruvbox",
        }

    def test_load_os_theme_returns_none_when_missing(self, tmp_config_dir):
        """load_os_theme devuelve None si el archivo no existe."""
        cm = ConfigManager()
        assert cm.load_os_theme() is None

    def test_load_os_theme_returns_saved_data(self, tmp_config_dir):
        """load_os_theme recupera los datos guardados."""
        cm = ConfigManager()
        cm.save_os_theme("nord")
        data = cm.load_os_theme()
        assert data["theme"] == "nord"


class TestOsTuiConfiguratorApp:
    """Tests de interfaz para la app Textual."""

    @pytest.mark.asyncio
    async def test_app_mounts_header(self, tmp_config_dir):
        """La app tiene un Header visible."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            header = pilot.app.query_one("Header")
            assert header.visible

    @pytest.mark.asyncio
    async def test_app_mounts_footer_with_shortcuts(self, tmp_config_dir):
        """El Footer muestra atajos de teclado (Q=Salir, Ctrl+S=Guardar)."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            footer = pilot.app.query_one("Footer")
            assert footer.visible

    @pytest.mark.asyncio
    async def test_app_mounts_sidebar(self, tmp_config_dir):
        """El Sidebar contiene items de navegacion (Neovim, Qtile)."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            sidebar = pilot.app.query_one("#sidebar")
            assert sidebar is not None
            assert sidebar.visible

    @pytest.mark.asyncio
    async def test_app_has_three_switches(self, tmp_config_dir):
        """El area de contenido contiene los tres switches de configuracion."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            switches = pilot.app.query("Switch")
            assert len(switches) == 3

    @pytest.mark.asyncio
    async def test_toggle_lsp_switch_updates_state(self, tmp_config_dir):
        """Al togglear el switch enable_lsp, el estado cambia en app.settings."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            initial = app.settings["enable_lsp"]
            await pilot.click("#switch_lsp")
            assert app.settings["enable_lsp"] != initial

    @pytest.mark.asyncio
    async def test_switches_have_correct_initial_values(self, tmp_config_dir):
        """Los switches reflejan los valores iniciales cargados del JSON."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            switch_lsp = pilot.app.query_one("#switch_lsp")
            switch_copilot = pilot.app.query_one("#switch_copilot")
            switch_neotree = pilot.app.query_one("#switch_neotree")

            assert switch_lsp.value is True
            assert switch_copilot.value is False
            assert switch_neotree.value is True

    @pytest.mark.asyncio
    async def test_save_action_writes_to_file(self, tmp_config_dir):
        """Al presionar Ctrl+S, se persiste el estado actual en el JSON."""
        app = OsTuiConfigurator()
        async with app.run_test() as pilot:
            await pilot.click("#switch_copilot")
            await pilot.press("ctrl+s")

            settings_file = (
                tmp_config_dir / "nvim" / ".config" / "nvim" / "tui_settings.json"
            )
            saved = json.loads(settings_file.read_text())

            assert saved == app.settings
