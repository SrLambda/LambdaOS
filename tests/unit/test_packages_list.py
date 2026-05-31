from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
PACKAGES_FILE = PROJECT_ROOT / "packages.x86_64"

# Packages that MUST remain in the base ISO (essential tools).
# Heavy/optional packages live in scripts/post-install/ profiles.
EXPECTED_PACKAGES = [
    "keepassxc",
    "lazydocker",
    "qalculate-gtk",
    "tailscale",
    "yazi",
]


class TestPackagesList:
    """Tests for packages.x86_64 content and ordering."""

    @pytest.fixture(autouse=True)
    def packages_content(self):
        """Read packages.x86_64 once per test class."""
        assert PACKAGES_FILE.exists(), f"packages.x86_64 not found at {PACKAGES_FILE}"
        return PACKAGES_FILE.read_text()

    @pytest.fixture(autouse=True)
    def packages_lines(self):
        """Read packages.x86_64 as list of lines."""
        assert PACKAGES_FILE.exists(), f"packages.x86_64 not found at {PACKAGES_FILE}"
        return PACKAGES_FILE.read_text().splitlines()

    @pytest.mark.parametrize("package", EXPECTED_PACKAGES)
    def test_package_present(self, package, packages_lines):
        """Each expected package appears exactly once in packages.x86_64."""
        found = [line for line in packages_lines if line.strip() == package]
        assert (
            len(found) == 1
        ), f"Expected '{package}' to appear exactly once, found {len(found)} times"

    def test_alphabetical_order(self, packages_lines):
        """packages.x86_64 is strictly alphabetically sorted."""
        # Filter out empty lines and comments (if any)
        clean_lines = [line.strip() for line in packages_lines if line.strip()]
        assert clean_lines == sorted(
            clean_lines
        ), "packages.x86_64 is not in alphabetical order"
