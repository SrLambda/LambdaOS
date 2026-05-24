# Lambda OS
---

## 1. Sistema Base y Arranque

- **OS:** Arch Linux
- **Sistema de Archivos:** BTRFS
- **Display Manager (Login):** Ly
- **Gestor de Ventanas:** Qtile (Servidor X11).
- **Snapshots:** Snapper

## 2. Servicios en Segundo Plano (Back-End)

- **Servidor de Audio:** PipeWire (con `wireplumber` y `pipewire-pulse` para máxima compatibilidad).
- **Autenticación (Polkit):** lxqt-policykit-agent
- **Gestor de Credenciales:** KeePassXC (con D-Bus Activation para integración fluida).
- **Servidor de Notificaciones:** Dunst

## 3. Entorno de Usuario y Utilidades (Front-End)

- **Terminal:** Kitty + Zsh
- **Tipografía del Sistema:** Monoid Nerd Font
- **Lanzador de Aplicaciones:** Rofi
- **Gestor de Archivos:** Yazi
- **Gestor de Wifi:** Impala
- **Gestor de Bluetooth:** BlueTUI

## 4. Funcionalidades Básicas (incluido en la ISO)

- **Navegador:** Chromium
- **Correo electrónico:** Mozilla Thunderbird
- **Ofimática:** LibreOffice
- **PDF:** Okular
- **Multimedia:** VLC
- **Calculadora:** Qalculate!
- **Notas:** Obsidian

## 5. Ocio (incluido en la ISO)

- **Gaming:** Steam, Battle.net (aprovechando la estabilidad nativa de X11 para Proton y Wine).
  > Requiere repositorio `[multilib]` habilitado (ya activo por defecto).
- **Música:** Spotify + Cliamp

## 6. Conectividad entre Dispositivos (incluido en la ISO)

- **VPN:** Tailscale
- **Nube:** Mega

## 7. Entorno de Desarrollo (incluido en la ISO)

- **IDE:** Nvim
- **Lenguajes de programación:** Python, C/C++, Web (HTML, CSS, JS), Rust, Go, Kotlin, Arduino.
- **Dev Tools:** Git + LazyGit, Docker + Docker Compose + LazyDocker.
- **Virtualización:** VirtualBox

## 8. Paquetes AUR (instalación post-boot requerida)

Los siguientes paquetes no están disponibles en los repositorios oficiales de Arch Linux y deben instalarse mediante un **AUR helper** (`yay` o `paru`) después de arrancar el sistema:

| Paquete | Descripción |
|---------|-------------|
| `spotify` | Cliente de música en streaming |
| `obsidian` | Base de conocimientos y notas |
| `megasync` | Sincronización con la nube de Mega.nz |
| `bluetui` | Gestor de Bluetooth en TUI |
| `impala` | Gestor de WiFi en TUI |

### Instalación rápida (copy-paste)

1. Instala un AUR helper (elige uno):

```bash
# Opción A: yay
sudo pacman -S yay

# Opción B: paru
sudo pacman -S paru
```

2. Ejecuta el script de paquetes AUR:

```bash
./scripts/aur-packages.sh
```

El script detecta automáticamente `yay` o `paru`, instala cada paquete con `--needed` (idempotente) y continúa incluso si un paquete individual falla.
