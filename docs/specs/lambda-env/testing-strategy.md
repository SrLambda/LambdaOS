# lambda-env: Testing Strategy by Wave

## Overview

This document defines how each wave is tested. Every wave produces a bootable ISO that is validated before the wave is considered complete.

**Este documento define cómo se testa cada wave. Cada wave produce una ISO booteable que se valida antes de considerar la wave completa.**

---

## CI/CD Testing Levels

| Level | Waves | What it does | Where it runs |
|---|---|---|---|
| **Build** | 0+ | ISO compila sin errores | GitHub Actions (CI) |
| **Smoke** | 0+ | ISO bootea en QEMU, llega a prompt, servicios clave corriendo | GitHub Actions (CI) con QEMU TCG |
| **Feature** | 2+ | Verifica que la feature de la wave funciona dentro de la VM | GitHub Actions (CI) con QEMU TCG + serial |
| **Install** | 8+ | Instala LambdaOS en disco virtual, bootea el sistema instalado | GitHub Actions (CI) con QEMU TCG |

### Recommendation: QEMU TCG on GitHub Actions

**EN**: GitHub Actions runners support KVM but it's not guaranteed (depends on runner type, nested virtualization availability). For reliability, we use QEMU TCG (software emulation) in CI. It's slower (~3x) but deterministic. For local development, use KVM (`-enable-kvm`) for fast iteration.

**ES**: Los runners de GitHub Actions soportan KVM pero no está garantizado (depende del tipo de runner, disponibilidad de virtualización anidada). Para confiabilidad, usamos QEMU TCG (emulación por software) en CI. Es más lento (~3x) pero determinístico. Para desarrollo local, usar KVM (`-enable-kvm`) para iteración rápida.

| Environment | Acceleration | Boot time | Use case |
|---|---|---|---|
| GitHub Actions CI | TCG (software) | ~3-5 min | Reliable, deterministic |
| Local dev | KVM (hardware) | ~30-60 sec | Fast iteration |
| Local no-KVM | TCG (software) | ~3-5 min | Fallback |

---

## Smoke Test (Wave 0+)

**Goal**: Verify the ISO is at least usable — it boots, reaches a login prompt, and key services are running.

**EN**: The smoke test verifies the ISO boots in QEMU, reaches the Ly display manager (or tty fallback), and key packages are installed. The test fails if after 3 minutes without any activity/output, nothing happens.

**ES**: El smoke test verifica que la ISO bootea en QEMU, llega al display manager Ly (o fallback a tty), y paquetes clave están instalados. El test falla si después de 3 minutos sin ninguna actividad/salida, no sucede nada.

### Test Flow

```
1. Boot ISO in QEMU (no KVM, TCG mode)
2. Wait up to 180 seconds for activity
   - Activity = serial output, VGA text change, or process start
3. Verify:
   a. Ly display manager is running (or tty prompt visible)
   b. flameshot --version returns successfully
   c. ISO name contains "LambdaOS"
4. If 180s timeout with no activity → FAIL
```

### QEMU Configuration

```bash
qemu-system-x86_64 \
  -m 2048 \
  -smp 2 \
  -drive file=LambdaOS-*.iso,format=raw,media=cdrom \
  -drive file=test-disk.qcow2,format=qcow2 \
  -serial stdio \
  -display none \
  -no-reboot \
  -timeout 300
```

### Activity Detection

**EN**: We monitor the QEMU serial console for output. If no output is received for 180 consecutive seconds, the test fails. This catches both boot failures and hangs.

**ES**: Monitoreamos la consola serial de QEMU por salida. Si no se recibe salida durante 180 segundos consecutivos, el test falla. Esto detecta tanto fallos de boot como cuelgues.

```python
class ActivityMonitor:
    def __init__(self, timeout=180):
        self.timeout = timeout
        self.last_activity = time.time()
        self.output_buffer = ""

    def feed(self, data):
        self.output_buffer += data
        self.last_activity = time.time()

    def is_alive(self):
        return (time.time() - self.last_activity) < self.timeout
```

### Verification Commands (via serial)

After boot reaches a prompt, the test sends commands via serial:

```bash
# Verify flameshot is installed
flameshot --version

# Verify ISO name
cat /etc/hostname

# Verify key services
systemctl is-active ly
systemctl is-active pipewire
```

---

## Feature Test (Wave 2+)

**Goal**: Verify that the feature introduced in the wave actually works inside the VM.

**EN**: Feature tests use the QEMU serial console to interact with the system. This is the most efficient approach because: (1) no network setup required, (2) works in TCG mode, (3) direct access to the system, (4) can test TUI modules by sending keystrokes.

**ES**: Los feature tests usan la consola serial de QEMU para interactuar con el sistema. Este es el enfoque más eficiente porque: (1) no requiere setup de red, (2) funciona en modo TCG, (3) acceso directo al sistema, (4) puede testear módulos TUI enviando keystrokes.

### Why Serial (not SSH)?

| Approach | Setup | TCG Compatible | TUI Testable | Network Required |
|---|---|---|---|---|
| **Serial** | None (built-in) | ✅ Yes | ✅ Yes (keystrokes) | ❌ No |
| SSH | Network config | ❌ Needs virtio-net | ❌ No TUI | ✅ Yes |
| Agent (QMP) | QMP socket | ✅ Yes | ❌ No TUI | ❌ No |

**Decision**: Serial console. It's the only approach that works without network, supports TCG, and can test TUI modules via keystroke injection.

### Feature Test Pattern

Each wave adds feature tests that verify its specific functionality:

```python
def test_wave_2_tui_opens():
    """Wave 2: lambda-env opens and shows menu"""
    vm = boot_iso()
    vm.serial.wait_for_prompt(timeout=180)
    vm.serial.send("lambda-env\n")
    vm.serial.wait_for_text("System", timeout=10)
    vm.serial.wait_for_text("Apps", timeout=10)
    vm.serial.wait_for_text("Ops", timeout=10)
    vm.serial.send("q\n")  # Quit
    assert vm.exit_code == 0

def test_wave_2_neovim_toggle():
    """Wave 2: TUI can toggle Neovim LSP"""
    vm = boot_iso()
    vm.serial.wait_for_prompt(timeout=180)
    vm.serial.send("lambda-env\n")
    vm.serial.wait_for_text("Neovim", timeout=10)
    vm.serial.send_key("ENTER")  # Select Neovim module
    vm.serial.wait_for_text("LSP", timeout=10)
    vm.serial.send_key("SPACE")  # Toggle
    vm.serial.wait_for_text("disabled", timeout=5)
    # Verify settings.json was updated
    vm.serial.send("cat ~/.config/lambdaos/settings.json | jq '.neovim.enable_lsp'\n")
    vm.serial.wait_for_text("false", timeout=5)
```

### Wave-by-Wave Feature Tests

| Wave | Feature Test | What it verifies |
|---|---|---|
| **0** | `test_wave_0_smoke` | ISO boots, flameshot installed, Ly running |
| **1** | `test_wave_1_hub_opens` | `lambda-env` opens, shows categories, settings.json exists |
| **2** | `test_wave_2_neovim_toggle` | TUI toggles Neovim LSP, settings.json updates |
| **2** | `test_wave_2_dotfiles_stow` | TUI stows kitty module, config appears in home |
| **3** | `test_wave_3_theme_change` | TUI changes theme, Qtile reloads with new colors |
| **4** | `test_wave_4_screen_detect` | TUI detects monitors via xrandr |
| **4** | `test_wave_4_audio_volume` | TUI reads/sets volume via wpctl |
| **5** | `test_wave_5_network_scan` | TUI scans WiFi networks via iwctl |
| **5** | `test_wave_5_bluetooth_status` | TUI reads Bluetooth status via bluetoothctl |
| **6** | `test_wave_6_services_list` | TUI lists systemd services |
| **6** | `test_wave_6_docs_accessible` | `curl localhost:8080` returns docs |
| **7** | `test_wave_7_monitor_runs` | TUI launches btop/htop |
| **8** | `test_wave_8_wizard_runs` | First boot wizard appears |
| **8** | `test_wave_8_calamares_available` | `sudo calamares` launches |
| **9** | `test_wave_9_sysctl_applied` | `sysctl net.ipv4.tcp_congestion_control` returns "bbr" |

---

## Install Test (Wave 8+)

**Goal**: Verify that LambdaOS can be installed to disk and the installed system boots.

**EN**: We recommend a headless installer script for CI testing rather than automating Calamares GUI. The script mimics what Calamares does: partition disk, format, copy files, install bootloader, configure fstab. This is deterministic and testable. Calamares is tested manually in QA.

**ES**: Recomendamos un script de instalación headless para testing en CI en lugar de automatizar la GUI de Calamares. El script imita lo que hace Calamares: particionar disco, formatear, copiar archivos, instalar bootloader, configurar fstab. Esto es determinístico y testeable. Calamares se testa manualmente en QA.

### Why Headless Script (not Calamares automation)?

| Approach | Deterministic | CI Compatible | Covers Real Flow | Maintenance |
|---|---|---|---|---|
| **Headless script** | ✅ Yes | ✅ Yes | ✅ Core flow | Low |
| Calamares auto | ❌ GUI-dependent | ❌ Needs X11 | ✅ Full flow | High |

**Decision**: Headless installer script for CI. Manual Calamares testing for QA.

### Install Test Flow

```
1. Boot LambdaOS ISO in QEMU
2. Run headless installer script:
   a. Partition virtual disk (EFI + BTRFS)
   b. Format partitions
   c. Copy root filesystem (arch-chroot style)
   d. Install bootloader (GRUB/systemd-boot)
   e. Configure fstab, hostname, users
3. Reboot into installed system
4. Verify:
   a. System boots from disk (not ISO)
   b. Login prompt appears
   c. User can login
   d. lambda-env runs
5. Timeout: 300 seconds for install + 180 seconds for boot
```

### Headless Installer Script

```bash
#!/usr/bin/env bash
# scripts/headless-install.sh
# Used by CI install test

set -euo pipefail

DISK="/dev/vda"
EFI_PART="${DISK}1"
ROOT_PART="${DISK}2"

# Partition
sgdisk --zap-all "$DISK"
sgdisk -n1:1M:+512M -t1:EF00 -c1:"EFI" "$DISK"
sgdisk -n2:0:0 -t2:8304 -c2:"Linux BTRFS" "$DISK"

# Format
mkfs.fat32 -F 32 "$EFI_PART"
mkfs.btrfs -f -L "LambdaOS" "$ROOT_PART"

# Mount
mount "$ROOT_PART" /mnt
btrfs subvolume create /mnt/@
btrfs subvolume create /mnt/@home
btrfs subvolume create /mnt/@log
btrfs subvolume create /mnt/@snapshots
umount /mnt

mount -o subvol=@ "$ROOT_PART" /mnt
mkdir -p /mnt/{boot,home,var/log,.snapshots}
mount "$EFI_PART" /mnt/boot
mount -o subvol=@home "$ROOT_PART" /mnt/home
mount -o subvol=@log "$ROOT_PART" /mnt/var/log
mount -o subvol=@snapshots "$ROOT_PART" /mnt/.snapshots

# Copy filesystem
cp -a /run/archiso/airootfs/* /mnt/

# Configure
arch-chroot /mnt /bin/bash -c '
  grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=LambdaOS
  grub-mkconfig -o /boot/grub/grub.cfg
  echo "lambdaos" > /etc/hostname
  mkinitcpio -P
'

echo "Installation complete"
```

---

## Test Infrastructure

### Base Test Framework

**EN**: We rescue and improve the existing `tests/qemu/test_live_boot.py`. The new framework adds: activity monitoring, serial console interaction, feature test patterns, and install test support.

**ES**: Rescatamos y mejoramos el `tests/qemu/test_live_boot.py` existente. El nuevo framework agrega: monitoreo de actividad, interacción con consola serial, patrones de feature test, y soporte de install test.

### Directory Structure

```
tests/
├── unit/                          ← Unit tests (existing)
│   └── ...
├── qemu/
│   ├── conftest.py                ← QEMU fixtures (VM boot, serial)
│   ├── test_live_boot.py          ← Smoke test (improved)
│   ├── test_features.py           ← Feature tests (wave 2+)
│   ├── test_install.py            ← Install test (wave 8+)
│   └── utils/
│       ├── vm.py                  ← QEMU VM wrapper
│       ├── serial.py              ← Serial console interaction
│       └── activity.py            ← Activity monitor (3min timeout)
└── __init__.py
```

### QEMU Fixture (conftest.py)

```python
import pytest
import subprocess
import os

@pytest.fixture(scope="module")
def iso_path():
    """Find the built ISO"""
    out_dir = os.path.join(os.getcwd(), "out")
    isos = [f for f in os.listdir(out_dir) if f.endswith(".iso")]
    assert len(isos) == 1, f"Expected 1 ISO, found {len(isos)}"
    return os.path.join(out_dir, isos[0])

@pytest.fixture
def vm(iso_path):
    """Boot ISO in QEMU, yield VM wrapper, cleanup after"""
    vm = QEMUVM(iso_path, memory=2048, cpus=2, kvm=False)
    vm.start()
    yield vm
    vm.stop()
```

### CI Integration

```yaml
# .github/workflows/ci.yml
jobs:
  smoke-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build ISO
        run: sudo mkarchiso -v -w work/ -o out/ .
      - name: Smoke test
        run: |
          pip install -r requirements-dev.txt
          pytest tests/qemu/test_live_boot.py -v --timeout=600

  feature-test:
    runs-on: ubuntu-latest
    needs: smoke-test
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - name: Build ISO
        run: sudo mkarchiso -v -w work/ -o out/ .
      - name: Feature tests
        run: |
          pip install -r requirements-dev.txt
          pytest tests/qemu/test_features.py -v --timeout=900
```

---

## 3-Minute Timeout Rule

**EN**: All QEMU tests fail if after 180 seconds (3 minutes) of consecutive inactivity, nothing happens. "Inactivity" means no serial output, no VGA text change, no new process start. This rule catches both boot failures and silent hangs.

**ES**: Todos los tests de QEMU fallan si después de 180 segundos (3 minutos) de inactividad consecutiva, no sucede nada. "Inactividad" significa: sin salida serial, sin cambio de texto VGA, sin nuevo proceso iniciado. Esta regla detecta tanto fallos de boot como cuelgues silenciosos.
