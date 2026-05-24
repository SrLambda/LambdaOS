from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
PACMAN_CONF = PROJECT_ROOT / "pacman.conf"


class TestPacmanConfig:
    """Tests for pacman.conf repository configuration."""

    @pytest.fixture(autouse=True)
    def pacman_conf_content(self):
        """Read pacman.conf once per test class."""
        assert PACMAN_CONF.exists(), f"pacman.conf not found at {PACMAN_CONF}"
        return PACMAN_CONF.read_text()

    def test_multilib_section_uncommented(self, pacman_conf_content):
        """The [multilib] section header is not prefixed with #."""
        lines = pacman_conf_content.splitlines()
        multilib_headers = [line for line in lines if line.strip().startswith("[multilib]")]
        assert len(multilib_headers) >= 1, "[multilib] section not found in pacman.conf"
        assert multilib_headers[0].strip() == "[multilib]", (
            f"[multilib] is commented out: {multilib_headers[0]}"
        )

    def test_multilib_include_uncommented(self, pacman_conf_content):
        """The Include line under [multilib] is not commented out."""
        lines = pacman_conf_content.splitlines()
        in_multilib = False
        include_line = None
        for line in lines:
            stripped = line.strip()
            if stripped == "[multilib]":
                in_multilib = True
                continue
            if in_multilib and stripped.startswith("["):
                break
            if in_multilib and stripped.startswith("Include"):
                include_line = stripped
                break
        assert include_line is not None, "No Include line found under [multilib]"
        assert not include_line.startswith("#"), (
            f"Include line under [multilib] is commented out: {include_line}"
        )
        assert "mirrorlist" in include_line, (
            f"Include line does not reference mirrorlist: {include_line}"
        )

    def test_core_section_uncommented(self, pacman_conf_content):
        """[core] section remains uncommented."""
        assert "\n[core]\n" in pacman_conf_content, "[core] section not found or commented"

    def test_extra_section_uncommented(self, pacman_conf_content):
        """[extra] section remains uncommented."""
        assert "\n[extra]\n" in pacman_conf_content, "[extra] section not found or commented"

    def test_repository_order(self, pacman_conf_content):
        """Repositories appear in correct order: core, extra, multilib."""
        core_idx = pacman_conf_content.find("\n[core]\n")
        extra_idx = pacman_conf_content.find("\n[extra]\n")
        multilib_idx = pacman_conf_content.find("\n[multilib]\n")

        assert core_idx != -1, "[core] not found"
        assert extra_idx != -1, "[extra] not found"
        assert multilib_idx != -1, "[multilib] not found"

        assert core_idx < extra_idx < multilib_idx, (
            f"Repository order incorrect: core@{core_idx}, extra@{extra_idx}, multilib@{multilib_idx}"
        )
