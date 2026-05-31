import os
import re
import shutil
import subprocess
import time
from pathlib import Path

import pexpect
import pytest

PROJECT_ROOT = Path(__file__).resolve().parent.parent.parent

TIMEOUT_BOOT = 300
TIMEOUT_STOW = 120
TIMEOUT_PEXPECT = 120
TIMEOUT_CMD = 30

DEBUG_LOG = Path("/tmp/opencode/qemu_boot_debug.log")
EXTRACT_DIR = Path("/tmp/opencode/iso_extract")


def _run(cmd, **kwargs):
    """Shortcut for subprocess.run with capture."""
    return subprocess.run(cmd, capture_output=True, text=True, **kwargs)


def _extract_from_iso(iso_path, pattern):
    """Extract a file from ISO using 7z. Returns path string or None."""
    os.makedirs(EXTRACT_DIR, exist_ok=True)
    seven_z = shutil.which("7z")
    if seven_z is None:
        return None
    result = _run(
        [seven_z, "x", "-aoa", f"-o{EXTRACT_DIR}", iso_path, pattern, "-y"],
    )
    if result.returncode != 0:
        return None
    extracted = EXTRACT_DIR / pattern
    if extracted.is_file():
        return str(extracted)
    return None


def _get_iso_uuid(iso_path):
    """Get the ISO9660 filesystem UUID via blkid."""
    result = _run(["blkid", "-s", "UUID", "-o", "value", iso_path])
    if result.returncode == 0 and result.stdout.strip():
        return result.stdout.strip()
    return None


def pytest_configure(config):
    config.addinivalue_line("markers", "qemu: tests that require QEMU to run")


@pytest.fixture(scope="session")
def project_root():
    return PROJECT_ROOT


@pytest.fixture(scope="session")
def qemu_binary():
    binary = shutil.which("qemu-system-x86_64")
    if binary is None:
        pytest.skip("qemu-system-x86_64 not found in PATH")
    return binary


@pytest.fixture(scope="session")
def iso_path(project_root):
    out_dir = project_root / "out"
    if not out_dir.is_dir():
        pytest.skip(f"Output directory not found: {out_dir}")
    candidates = sorted(out_dir.glob("LambdaOS-*-x86_64.iso"))
    if not candidates:
        pytest.skip(f"No ISO matching LambdaOS-*-x86_64.iso found in {out_dir}")
    iso = candidates[-1]
    return str(iso)


@pytest.fixture(scope="session")
def kernel_path(iso_path):
    path = _extract_from_iso(iso_path, "arch/boot/x86_64/vmlinuz-linux")
    if path is None:
        pytest.skip("Could not extract vmlinuz-linux from ISO")
    return path


@pytest.fixture(scope="session")
def initrd_path(iso_path):
    path = _extract_from_iso(iso_path, "arch/boot/x86_64/initramfs-linux.img")
    if path is None:
        pytest.skip("Could not extract initramfs-linux.img from ISO")
    return path


@pytest.fixture(scope="session")
def iso_uuid(iso_path):
    uuid = _get_iso_uuid(iso_path)
    if uuid is None:
        pytest.skip("Could not determine ISO filesystem UUID")
    return uuid


def _try_login(child, username, password, timeout=90):
    """Try to log in via getty; return True if we get a shell prompt."""
    child.sendline(username)
    
    # Archiso root usually drops straight to a shell without asking for a password.
    # We expect EITHER a password prompt OR a successful shell prompt.
    idx = child.expect(
        [
            r"Password:",                        # 0: Requires password
            r"\[" + username + r"@.*\][\$#%>]",  # 1: Shell prompt (bracketed)
            username + r"@.*[\$#%>]",            # 2: Shell prompt (no brackets)
            r"[\$#%>] ",                         # 3: Bare prompt (e.g. root's '# ')
            pexpect.TIMEOUT,                     # 4
            pexpect.EOF,                         # 5
        ],
        timeout=30,
    )
    
    if idx == 0:
        # Prompted for password
        child.sendline(password)
        idx2 = child.expect(
            [
                r"\[" + username + r"@.*\][\$#%>]",
                username + r"@.*[\$#%>]",
                r"[\$#%>] ",
                pexpect.TIMEOUT,
                pexpect.EOF,
            ],
            timeout=timeout,
        )
        return idx2 < 3
    elif idx in (1, 2, 3):
        # Dropped straight to shell (no password needed)
        return True
    else:
        return False

@pytest.fixture(scope="session")
def qemu_booted(qemu_binary, kernel_path, initrd_path, iso_path, iso_uuid):
    """Boot QEMU with direct kernel boot to get serial console output.

    Since the ISO's syslinux config does not include console=ttyS0, we
    bypass the bootloader entirely using -kernel, -initrd, and -append.
    The ISO filesystem is still exposed via -cdrom so archiso's init
    system can find and mount it via archisosearchuuid.

    We use -device isa-serial + -chardev explicitly (instead of -serial)
    to ensure the ISA serial port is registered in the guest's ACPI/PNP
    tables so the kernel detects it and creates /dev/ttyS0.
    """
    append = (
        f"archisobasedir=arch archisosearchuuid={iso_uuid} "
        f"console=ttyS0,115200 8250.nr_uarts=1"
    )

    cmd = [
        qemu_binary,
        "-M", "pc",
        "-m", "2G",
        "-display", "none",
        "-device", "isa-serial,chardev=serial0",
        "-chardev", "stdio,id=serial0",
        "-device", "virtio-rng-pci",     # <--- AÑADE ESTA LÍNEA AQUÍ
        "-kernel", kernel_path,
        "-initrd", initrd_path,
        "-append", append,
        "-cdrom", iso_path,
    ]

    child = pexpect.spawn(
        cmd[0], args=cmd[1:], encoding="utf-8", timeout=TIMEOUT_PEXPECT,
        logfile=open(str(DEBUG_LOG), "w"),
    )

    try:
        idx = child.expect(
            [
                r"\[liveuser@.*\][\$#%>]",      # 0: autologin shell (bracketed)
                r"liveuser@.*[\$#%>]",           # 1: autologin shell (no brackets)
                r"liveuser@\S+",                 # 2: catch-all liveuser
                r"[\$#%>] ",                     # 3: bare prompt
                r"liveuser login:",              # 4: missing autologin (liveuser)
                r"login:",                       # 5: generic login prompt
                pexpect.TIMEOUT,                 # 6
                pexpect.EOF,                     # 7
            ],
            timeout=TIMEOUT_BOOT,
        )

        if idx in (4, 5):
            # No autologin on serial — log in manually.
            # Liveuser may not exist (created at runtime, sometimes fails).
            # Fall back to root (empty password) and su to liveuser.
            login_ok = _try_login(child, "root", "")
            if login_ok:
                child.sendline("id liveuser 2>/dev/null || (useradd -m liveuser && passwd -d liveuser)")
                child.expect([r"[\$#%>] ", pexpect.TIMEOUT], timeout=30)
                child.sendline("su - liveuser")
                child.expect(
                    [
                        r"\[liveuser@.*\][\$#%>]",
                        r"liveuser@.*[\$#%>]",
                        r"liveuser@\S+",
                        r"[\$#%>] ",
                        pexpect.TIMEOUT,
                        pexpect.EOF,
                    ],
                    timeout=30,
                )
            else:
                _dump_last_output(child, "root manual login prompt")
        elif idx >= 6:
            _dump_last_output(child, "autologin shell prompt")

        time.sleep(1)
        child.expect([r"[\$#%>] ", pexpect.TIMEOUT], timeout=5)

        yield child
    finally:
        if child.isalive():
            try:
                child.sendline("sudo poweroff")
                child.expect(pexpect.EOF, timeout=15)
            except (pexpect.TIMEOUT, pexpect.EOF, OSError):
                pass
            child.close(force=True)


def _dump_last_output(child, phase_name):
    before = getattr(child, "before", None) or ""
    after = getattr(child, "after", None) or ""
    snippet = before[-3000:] if before else "(empty)"
    trace = (
        f"\n--- QEMU output (last 3000 chars) at phase '{phase_name}' ---\n"
        f"{snippet}\n"
        f"--- after: {after!r} ---\n"
        f"Full log: {DEBUG_LOG}\n"
    )
    with open(str(DEBUG_LOG), "a") as f:
        f.write(trace)
    pytest.fail(
        f"QEMU did not reach expected state: {phase_name}.\n"
        f"Last output (truncated to 3000 chars):\n{snippet}\n"
        f"Full debug log: {DEBUG_LOG}"
    )


@pytest.fixture(scope="session")
def qemu_logged_in(qemu_booted):
    child = qemu_booted

    child.sendline("cd ~/dotfiles && stow */ 2>&1")
    child.expect([r"[\$#%>] ", pexpect.TIMEOUT, pexpect.EOF], timeout=TIMEOUT_STOW)

    return child
