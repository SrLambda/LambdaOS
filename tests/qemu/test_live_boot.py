"""Integration tests for the LambdaOS (omarchy.iso) live boot in QEMU.

These tests use pexpect to drive qemu-system-x86_64 in headless/serial mode
(-display none + -serial stdio) with pure software emulation (no KVM) for
CI/CD compatibility. They verify:
  1. The ISO boots past the Syslinux bootloader and reaches the liveuser
     autologin shell prompt.
  2. GNU Stow deploys dotfile symlinks correctly.
  3. The Neovim init.lua exists under the stowed config.
"""

import pexpect
import pytest

from tests.qemu.conftest import TIMEOUT_CMD


class TestLiveBoot:
    """Tests that require a single QEMU session (session-scoped fixtures)."""

    def test_iso_boots_to_shell_prompt(self, qemu_booted):
        """Boot ISO, skip bootloader, and verify liveuser autologin shell.

        The ``qemu_booted`` fixture handles: detect Syslinux menu -> Enter
        to skip countdown -> wait for [liveuser@archiso ~]$ prompt.
        This test only asserts QEMU is still alive after the process.
        """
        assert qemu_booted.isalive(), "QEMU process is not alive"

    def test_liveuser_stow_symlinks_correct(self, qemu_logged_in):
        """Verify GNU Stow deployed dotfile symlinks pointing to ~/dotfiles/.

        Checks that:
        - ``~/.config/nvim`` points into ``~/dotfiles/nvim/.config/nvim``
        - ``~/.config/qtile`` points into ``~/dotfiles/qtile/.config/qtile``
        """
        child = qemu_logged_in

        child.sendline("readlink -f ~/.config/nvim")
        child.expect(r"[\$#%>] .*", timeout=TIMEOUT_CMD)
        nvim_target = child.before.strip() if child.before else ""

        assert nvim_target, "readlink for ~/.config/nvim returned empty output"
        assert "dotfiles/nvim/.config/nvim" in nvim_target, (
            f"Expected ~/.config/nvim -> .../dotfiles/nvim/.config/nvim\n"
            f"Got: {nvim_target}"
        )

        child.sendline("readlink -f ~/.config/qtile")
        child.expect(r"[\$#%>] .*", timeout=TIMEOUT_CMD)
        qtile_target = child.before.strip() if child.before else ""

        assert qtile_target, "readlink for ~/.config/qtile returned empty output"
        assert "dotfiles/qtile/.config/qtile" in qtile_target, (
            f"Expected ~/.config/qtile -> .../dotfiles/qtile/.config/qtile\n"
            f"Got: {qtile_target}"
        )

    def test_neovim_init_lua_exists(self, qemu_logged_in):
        """Verify ~/.config/nvim/init.lua is present and readable."""
        child = qemu_logged_in

        child.sendline("ls -la ~/.config/nvim/init.lua")
        child.expect(r"[\$#%>] .*", timeout=TIMEOUT_CMD)
        output = child.before.strip() if child.before else ""

        assert "init.lua" in output, (
            f"~/.config/nvim/init.lua not found via ls.\n"
            f"Command output:\n{output}"
        )

        child.sendline("test -f ~/.config/nvim/init.lua && echo EXISTS || echo MISSING")
        child.expect(r"[\$#%>] .*", timeout=TIMEOUT_CMD)
        output2 = child.before.strip() if child.before else ""

        assert "EXISTS" in output2, (
            f"~/.config/nvim/init.lua not found via test -f.\n"
            f"Command output:\n{output2}"
        )
