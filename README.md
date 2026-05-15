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

## 4. Funcionalidades Básicas

- **Navegador:** Chromium
- **Correo electrónico:** Mozilla Thunderbird
- **Ofimática:** LibreOffice
- **PDF:** Okular
- **Multimedia:** VLC
- **Calculadora:** Qalculate!
- **Notas:** Obsidian

## 5. Ocio

- **Gaming:** Steam, Battle.net (aprovechando la estabilidad nativa de X11 para Proton y Wine).
- **Música:** Spotify + Cliamp

## 6. Conectividad entre Dispositivos

- **VPN:** Tailscale
- **Nube:** Mega

## 7. Entorno de Desarrollo

- **IDE:** Nvim
- **Lenguajes de programación:** Python, C/C++, Web (HTML, CSS, JS), Rust, Go, Kotlin, Arduino.
- **Dev Tools:** Git + LazyGit, Docker + Docker Compose + LazyDocker.
- **Virtualización:** VirtualBox
