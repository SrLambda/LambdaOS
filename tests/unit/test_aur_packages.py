import shutil
import subprocess
from pathlib import Path

import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent
SCRIPT_PATH = PROJECT_ROOT / "scripts" / "aur-packages.sh"
SHELLCHECK = shutil.which("shellcheck")

AUR_PACKAGES = ["spotify", "obsidian", "megasync", "bluetui", "impala"]


class TestAurPackagesScript:
    """Tests for scripts/aur-packages.sh."""

    @pytest.mark.skipif(not SHELLCHECK, reason="shellcheck not installed")
    def test_shellcheck_passes(self):
        """aur-packages.sh passes shellcheck with zero warnings/errors."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        result = subprocess.run(
            [SHELLCHECK, str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result.returncode == 0
        ), f"shellcheck found issues:\n{result.stdout}\n{result.stderr}"

    def test_aur_helper_detection_yay(self, tmp_path, monkeypatch):
        """Script detects yay and proceeds (exit 0)."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        fake_yay = tmp_path / "yay"
        fake_yay.write_text("#!/bin/bash\nexit 0\n")
        fake_yay.chmod(0o755)
        monkeypatch.setenv("PATH", str(tmp_path))
        result = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result.returncode == 0
        ), f"Expected exit 0 with yay available, got {result.returncode}. Stderr: {result.stderr}"

    def test_aur_helper_detection_paru(self, tmp_path, monkeypatch):
        """Script detects paru when yay is missing and proceeds (exit 0)."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        fake_paru = tmp_path / "paru"
        fake_paru.write_text("#!/bin/bash\nexit 0\n")
        fake_paru.chmod(0o755)
        monkeypatch.setenv("PATH", str(tmp_path))
        result = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result.returncode == 0
        ), f"Expected exit 0 with paru available, got {result.returncode}. Stderr: {result.stderr}"

    def test_aur_helper_missing_exits_nonzero(self, tmp_path, monkeypatch):
        """Script exits with non-zero when neither yay nor paru is available."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        monkeypatch.setenv("PATH", str(tmp_path))
        result = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result.returncode != 0
        ), f"Expected non-zero exit when no AUR helper, got {result.returncode}"
        assert (
            result.returncode == 1
        ), f"Expected exit code 1 for missing helper, got {result.returncode}"

    def test_uses_needed_flag(self):
        """Script invokes AUR helper with --needed flag."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        source = SCRIPT_PATH.read_text()
        assert "--needed" in source, "Script must use --needed flag"

    def test_per_package_error_handling(self, tmp_path, monkeypatch):
        """Script continues when one AUR package fails and exits 2."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        fake_yay = tmp_path / "yay"
        fake_yay.write_text(
            "#!/bin/bash\n"
            'for arg in "$@"; do\n'
            '  if [[ "$arg" == "spotify" ]]; then\n'
            "    exit 1\n"
            "  fi\n"
            "done\n"
            "exit 0\n"
        )
        fake_yay.chmod(0o755)
        monkeypatch.setenv("PATH", str(tmp_path))
        result = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result.returncode == 2
        ), f"Expected exit code 2 for partial failure, got {result.returncode}. Stderr: {result.stderr}"

    def test_idempotency_rerun_exits_zero(self, tmp_path, monkeypatch):
        """Re-running script with --needed succeeds (exit 0)."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        fake_yay = tmp_path / "yay"
        fake_yay.write_text("#!/bin/bash\nexit 0\n")
        fake_yay.chmod(0o755)
        monkeypatch.setenv("PATH", str(tmp_path))

        # First run
        result1 = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert result1.returncode == 0, f"First run failed: {result1.stderr}"

        # Second run (idempotency)
        result2 = subprocess.run(
            ["/usr/bin/bash", str(SCRIPT_PATH)],
            capture_output=True,
            text=True,
        )
        assert (
            result2.returncode == 0
        ), f"Expected exit 0 on re-run, got {result2.returncode}. Stderr: {result2.stderr}"

    def test_script_is_executable(self):
        """aur-packages.sh has executable permissions."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        assert (
            SCRIPT_PATH.stat().st_mode & 0o111
        ), f"Script is not executable: {SCRIPT_PATH}"

    def test_all_aur_packages_listed(self):
        """Script references all documented AUR packages."""
        assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"
        source = SCRIPT_PATH.read_text()
        for package in AUR_PACKAGES:
            assert package in source, f"AUR package '{package}' not found in script"
