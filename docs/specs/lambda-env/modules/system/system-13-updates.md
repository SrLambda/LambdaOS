# lambda-env: Updates Manager (system-13)

## Intent

Módulo TUI para gestionar actualizaciones del sistema: paquetes oficiales de Arch + paquetes AUR.

## Scope

- Ver paquetes con actualizaciones disponibles
- Actualizar sistema (pacman)
- Actualizar paquetes AUR (yay/paru)
- Ver historial de actualizaciones
- Configurar frecuencia de chequeo

## Requirements

1. Check for updates: `pacman -Qu` + `yay/paru -Qua`
2. Mostrar lista de paquetes actualizables con versión actual → nueva
3. Actualizar todo o seleccionar paquetes individuales
4. Mostrar output de pacman/yay en tiempo real durante la actualización
5. Configurar mirror con `reflector` (ya incluido)

## Scenarios

### Escenario 1: Actualizar sistema
- Usuario abre Updates → "Check for updates"
- Ve: "42 packages to update (38 official, 4 AUR)"
- Selecciona "Update all" → ve progreso en tiempo real
- Confirmación al finalizar

### Escenario 2: Solo actualizar AUR
- Ve lista mixta
- Filtra por "AUR only"
- Actualiza solo los 4 paquetes AUR

## Technical Notes

- `pacman -Syu` para oficial
- `yay -Syu --noconfirm` o `paru -Syu` para AUR
- `reflector --latest 5 --sort rate --save /etc/pacman.d/mirrorlist` para mirrors
- Capturar output de pacman para mostrar progreso
- Manejar conflictos de paquetes (pacnews)

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Requiere root para pacman
- yay o paru (post-boot AUR)
