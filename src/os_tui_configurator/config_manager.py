import json
import os
from pathlib import Path


class ConfigManager:
    DEFAULT_SETTINGS = {
        "enable_lsp": True,
        "enable_copilot": True,
        "enable_neotree": True,
    }

    THEME_FILENAME = ".nvim_theme"
    OS_THEME_FILENAME = "os_theme.json"

    def _get_config_base(self) -> Path:
        base = os.environ.get("OS_CONFIG_DIR", str(Path.home()))
        return Path(base)

    def _get_settings_path(self) -> Path:
        return self._get_config_base() / "nvim" / ".config" / "nvim" / "tui_settings.json"

    def _get_theme_path(self) -> Path:
        return self._get_config_base() / "nvim" / ".config" / "nvim" / self.THEME_FILENAME

    def load_tui_settings(self) -> dict:
        path = self._get_settings_path()
        if not path.exists():
            return dict(self.DEFAULT_SETTINGS)
        with open(path, "r") as f:
            return json.load(f)

    def save_tui_settings(self, settings: dict) -> None:
        path = self._get_settings_path()
        path.parent.mkdir(parents=True, exist_ok=True)
        with open(path, "w") as f:
            json.dump(settings, f, indent=2)

    def set_theme(self, theme: str) -> None:
        path = self._get_theme_path()
        path.parent.mkdir(parents=True, exist_ok=True)
        with open(path, "w") as f:
            f.write(f"NVIM_THEME={theme}\n")

    def get_theme(self) -> str | None:
        path = self._get_theme_path()
        if not path.exists():
            return None
        content = path.read_text().strip()
        if "=" in content:
            return content.split("=", 1)[1].strip()
        return content

    def _get_os_theme_path(self) -> Path:
        return self._get_config_base() / self.OS_THEME_FILENAME

    def save_os_theme(self, theme: str) -> None:
        path = self._get_os_theme_path()
        path.parent.mkdir(parents=True, exist_ok=True)
        theme_data = {
            "theme": theme,
            "nvim_theme": theme,
            "qtile_theme": theme,
        }
        with open(path, "w") as f:
            json.dump(theme_data, f, indent=2)

    def load_os_theme(self) -> dict | None:
        path = self._get_os_theme_path()
        if not path.exists():
            return None
        with open(path, "r") as f:
            return json.load(f)
