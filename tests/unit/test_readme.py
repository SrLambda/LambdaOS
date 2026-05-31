from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
README_PATH = PROJECT_ROOT / "README.md"

AUR_PACKAGES = ["spotify", "obsidian", "megasync", "bluetui", "impala"]
OFFICIAL_PACKAGES = [
    "chromium",
    "docker",
    "docker-compose",
    "keepassxc",
    "lazydocker",
    "libreoffice-fresh",
    "okular",
    "qalculate-gtk",
    "steam",
    "tailscale",
    "thunderbird",
    "virtualbox",
    "vlc",
    "wine",
    "yazi",
]


class TestReadmeAccuracy:
    """Tests that README.md accurately distinguishes ISO-included vs AUR packages."""

    @pytest.fixture(autouse=True)
    def readme_content(self):
        """Read README.md once per test class."""
        assert README_PATH.exists(), f"README.md not found at {README_PATH}"
        return README_PATH.read_text()

    def test_has_iso_included_section(self, readme_content):
        """README has a section or label for ISO-included packages."""
        content_lower = readme_content.lower()
        assert (
            "incluido en la iso" in content_lower or "included in iso" in content_lower
        ), "README must label packages as included in the ISO"

    def test_has_aur_post_install_section(self, readme_content):
        """README has a section for AUR post-install packages."""
        assert (
            "aur" in readme_content.lower() and "post" in readme_content.lower()
        ), "README must have a post-install AUR section"

    def test_references_aur_script(self, readme_content):
        """README references scripts/aur-packages.sh."""
        assert (
            "scripts/aur-packages.sh" in readme_content
        ), "README must reference scripts/aur-packages.sh"

    def test_all_aur_packages_labeled(self, readme_content):
        """Each documented AUR package is labeled as post-install."""
        for package in AUR_PACKAGES:
            assert (
                package.lower() in readme_content.lower()
            ), f"AUR package '{package}' not mentioned in README"

    def test_has_copy_paste_instructions(self, readme_content):
        """README provides copy-pasteable instructions for AUR helper and script."""
        assert (
            "yay" in readme_content.lower() or "paru" in readme_content.lower()
        ), "README must mention yay or paru for AUR helper installation"
        assert (
            "./scripts/aur-packages.sh" in readme_content
            or "bash scripts/aur-packages.sh" in readme_content
        ), "README must provide copy-pasteable script invocation"

    def test_multilib_packages_labeled(self, readme_content):
        """Steam and Wine are labeled as requiring multilib."""
        content_lower = readme_content.lower()
        assert "steam" in content_lower, "README must mention Steam"
        assert "wine" in content_lower, "README must mention Wine"
        # They should be in the ISO-included section, not AUR section
        assert "multilib" in content_lower, "README must mention multilib requirement"
