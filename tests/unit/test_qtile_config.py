import ast
from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
QTILE_DIR = PROJECT_ROOT / "airootfs" / "etc" / "skel" / "dotfiles" / "qtile" / ".config" / "qtile"

QTILE_FILES = ["config.py", "theme.py", "keys.py", "groups.py", "screens.py"]


class TestQtileConfigSyntax:
    """Pruebas de validacion sintactica para la configuracion de Qtile."""

    @pytest.mark.parametrize("filename", QTILE_FILES)
    def test_python_syntax_valid(self, filename):
        """Cada archivo .py de Qtile tiene sintaxis Python valida."""
        filepath = QTILE_DIR / filename
        assert filepath.exists(), f"{filename} not found at {filepath}"

        source = filepath.read_text()
        try:
            ast.parse(source)
        except SyntaxError as e:
            pytest.fail(f"{filename} has a syntax error: {e}")

    def test_config_imports_resolve(self, tmp_path, monkeypatch):
        """config.py hace imports locales a modulos existentes y sin errores."""
        monkeypatch.setenv("OS_CONFIG_DIR", str(tmp_path))

        config_path = QTILE_DIR / "config.py"
        source = config_path.read_text()
        tree = ast.parse(source)

        local_imports = []
        for node in ast.walk(tree):
            if isinstance(node, ast.ImportFrom) and node.level == 1:
                local_imports.append(node.module)

        for module in local_imports:
            module_path = QTILE_DIR / f"{module}.py"
            assert module_path.exists(), f"Local import '{module}' not found at {module_path}"
            module_source = module_path.read_text()
            try:
                ast.parse(module_source)
            except SyntaxError as e:
                pytest.fail(f"Imported module '{module}' has a syntax error: {e}")

    def test_theme_py_has_try_except_fallback(self):
        """theme.py usa try/except para manejar la ausencia de os_theme.json."""
        theme_path = QTILE_DIR / "theme.py"
        source = theme_path.read_text()
        tree = ast.parse(source)

        has_try = any(isinstance(node, ast.Try) for node in ast.walk(tree))
        assert has_try, "theme.py debe usar try/except al cargar os_theme.json"

    def test_keys_py_defines_mod4(self):
        """keys.py define la tecla modificadora mod4 (Super)."""
        source = (QTILE_DIR / "keys.py").read_text()
        assert "mod4" in source, "keys.py debe definir mod4 como modificador principal"

    def test_theme_py_defines_five_themes(self):
        """theme.py contiene al menos 5 paletas de colores."""
        source = (QTILE_DIR / "theme.py").read_text()

        expected_themes = ["catppuccin", "gruvbox", "tokyonight", "nord", "onedark"]
        for theme in expected_themes:
            assert theme in source.lower(), f"theme.py debe contener el tema '{theme}'"

    def test_groups_py_defines_five_groups(self):
        """groups.py define 5 workspaces."""
        source = (QTILE_DIR / "groups.py").read_text()
        tree = ast.parse(source)

        group_count = 0
        for node in ast.walk(tree):
            if isinstance(node, ast.Call):
                if hasattr(node.func, "id") and node.func.id == "Group":
                    group_count += 1

        assert group_count >= 5, f"groups.py debe definir al menos 5 Group, encontrados {group_count}"

    def test_screens_py_creates_bar(self):
        """screens.py crea una barra con widgets."""
        source = (QTILE_DIR / "screens.py").read_text()
        assert "bar.Bar" in source, "screens.py debe crear una bar.Bar"

    def test_all_files_are_utf8(self):
        """Todos los archivos de Qtile son UTF-8."""
        for filename in QTILE_FILES:
            filepath = QTILE_DIR / filename
            try:
                filepath.read_text(encoding="utf-8")
            except UnicodeDecodeError:
                pytest.fail(f"{filename} no es UTF-8 valido")
