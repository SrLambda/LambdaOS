import os
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

# Expresión regular robusta: Busca $, # o % ignorando códigos ANSI (colores) entre el símbolo y el espacio
PROMPT = r"archiso[^\r\n]*[\$#%](?:\x1b\[[0-9;]*[a-zA-Z])*\s*"

def _run(cmd, **kwargs):
    return subprocess.run(cmd, capture_output=True, text=True, **kwargs)

def _extract_from_iso(iso_path, pattern):
    os.makedirs(EXTRACT_DIR, exist_ok=True)
    seven_z = shutil.which("7z")
    if seven_z is None:
        return None
    result = _run([seven_z, "x", "-aoa", f"-o{EXTRACT_DIR}", iso_path, pattern, "-y"])
    if result.returncode == 0:
        extracted = EXTRACT_DIR / pattern
        if extracted.is_file():
            return str(extracted)
    return None

def _get_iso_uuid(iso_path):
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
    candidates = sorted(out_dir.glob("lambda-os-*-x86_64.iso"))
    if not candidates:
        pytest.skip(f"No ISO matching lambda-os-*-x86_64.iso found in {out_dir}")
    return str(candidates[-1])

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
    child.sendline(username)
    idx = child.expect([r"Password:", PROMPT, pexpect.TIMEOUT, pexpect.EOF], timeout=timeout)
    if idx == 0:
        child.sendline(password)
        idx2 = child.expect([PROMPT, pexpect.TIMEOUT, pexpect.EOF], timeout=timeout)
        return idx2 == 0
    elif idx == 1:
        return True
    return False

@pytest.fixture(scope="session")
def qemu_booted(qemu_binary, kernel_path, initrd_path, iso_path, iso_uuid):
    append = (
        f"archisobasedir=arch archisosearchuuid={iso_uuid} "
        f"console=ttyS0,115200"
    )

    cmd = [
        qemu_binary,
        "-M", "pc",
        "-m", "2G",
        "-display", "none",
        "-serial", "stdio",             # <--- LA MAGIA ESTÁ AQUÍ (Universal Serial)
        "-device", "virtio-rng-pci",    # <--- Mantiene el arranque ultra rápido
        "-kernel", kernel_path,
        "-initrd", initrd_path,
        "-append", append,
        "-cdrom", iso_path,
    ]


    os.makedirs(DEBUG_LOG.parent, exist_ok=True)

    child = pexpect.spawn(
        cmd[0], args=cmd[1:], encoding="utf-8", timeout=TIMEOUT_PEXPECT,
        logfile=open(str(DEBUG_LOG), "w"),
    )

    try:
        idx = child.expect(
            [
                r"liveuser login:",              # 0
                r"login:",                       # 1
                PROMPT,                          # 2 (Entró directo por autologin)
                pexpect.TIMEOUT,
                pexpect.EOF,
            ],
            timeout=TIMEOUT_BOOT,
        )

        if idx in (0, 1):
            login_ok = _try_login(child, "root", "")
            if login_ok:
                child.sendline("useradd -m liveuser 2>/dev/null; passwd -d liveuser 2>/dev/null")
                child.expect([PROMPT, pexpect.TIMEOUT], timeout=30)
                child.sendline("su - liveuser")
                child.expect([PROMPT, pexpect.TIMEOUT, pexpect.EOF], timeout=30)
            else:
                _dump_last_output(child, "Fallo al loguear root")
        elif idx >= 3:
            _dump_last_output(child, "Timeout esperando login prompt")

        time.sleep(1)
        yield child
    finally:
        if child.isalive():
            try:
                child.sendline("sudo poweroff")
                child.expect(pexpect.EOF, timeout=15)
            except Exception:
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
    pytest.fail(f"QEMU did not reach expected state: {phase_name}.\nLog: {DEBUG_LOG}")

@pytest.fixture(scope="session")
def qemu_logged_in(qemu_booted):
    child = qemu_booted
    child.sendline("cd ~/dotfiles && stow */ 2>&1")
    child.expect([PROMPT, pexpect.TIMEOUT, pexpect.EOF], timeout=TIMEOUT_STOW)
    return child
