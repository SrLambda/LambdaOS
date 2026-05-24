# LambdaOS — Specs de Mejora

> Documento de requerimientos para la próxima iteración de LambdaOS.
> Generado: 2026-05-24 | Estado actual: 6 fases completadas, 30/30 tests unitarios pasando, 0/3 E2E pasando.

---

## Spec 1: Corregir Paquetes Faltantes (CRÍTICO)

### Problema
El `README.md` describe funcionalidades que **no están en `packages.x86_64`**:
- Gaming: Steam, Battle.net → requiere `multilib` habilitado
- VPN: Tailscale → paquete `tailscale` ausente
- Nube: Mega → `megasync` ausente (AUR)
- Ofimática: LibreOffice, Thunderbird, VLC, Okular → ausentes
- Notas: Obsidian → ausente (AUR)
- Música: Spotify → `spotify` ausente (AUR)
- Bluetooth: BlueTUI → no existe en repos oficiales
- WiFi: Impala → no existe en repos oficiales

### Requerimiento
1. Habilitar repositorio `[multilib]` en `pacman.conf`.
2. Agregar paquetes faltantes que existen en repos oficiales a `packages.x86_64`.
3. Documentar qué paquetes requieren AUR y crear un script `aur-packages.sh` separado.
4. Eliminar del README o marcar como "pendiente" lo que no se pueda incluir en esta iteración.

### Criterios de Aceptación
- [ ] `pacman.conf` tiene `[multilib]` descomentado.
- [ ] `packages.x86_64` incluye: `tailscale`, `libreoffice-fresh`, `thunderbird`, `vlc`, `okular`, `qalculate-gtk`, `obsidian` (si existe en repos), `spotify` (si existe en repos).
- [ ] Script `aur-packages.sh` documenta paquetes AUR con instrucciones de instalación.
- [ ] README actualizado para reflejar qué está incluido en la ISO vs qué requiere instalación post-boot.

---

## Spec 2: Limpiar Artefactos Huérfanos (CRÍTICO)

### Problema
- `backup_test/` es una copia muerta de tests QEMU antiguos.
- `.atl/` no está en `.gitignore` y genera ruido en `git status`.
- `tests/.venv/` existe como directorio con dependencias instaladas.

### Requerimiento
1. Eliminar `backup_test/` o documentar por qué existe.
2. Agregar `.atl/` a `.gitignore`.
3. Agregar `tests/.venv/` a `.gitignore` (ya cubierto por `.venv/` en root, pero verificar).

### Criterios de Aceptación
- [ ] `git status` no muestra `.atl/` como untracked.
- [ ] `backup_test/` eliminado o justificado con README interno.
- [ ] `.gitignore` cubre todos los directorios virtuales y de cache.

---

## Spec 3: Formalizar Paquete Python del TUI

### Problema
- `main.py` usa `from src.os_tui_configurator.app import ...` que depende del cwd.
- No hay `pyproject.toml` ni `setup.py`.
- `requirements-dev.txt` y `src/requirements.txt` están duplicados.

### Requerimiento
1. Crear `pyproject.toml` con:
   - Metadata del proyecto (nombre, versión, descripción).
   - Dependencias declaradas (`textual`, `pexpect`, `pytest`).
   - Entry point CLI: `lambdaos-tui = "os_tui_configurator.main:main"`.
   - Configuración de pytest (asyncio mode, test paths).
2. Unificar dependencias en un solo archivo.
3. Corregir imports para que funcionen independientemente del cwd.

### Criterios de Aceptación
- [ ] `pyproject.toml` existe con metadata completa.
- [ ] `pip install -e .` instala el paquete en modo editable.
- [ ] `lambdaos-tui` ejecuta la app desde cualquier directorio.
- [ ] `pytest` funciona sin necesidad de estar en el root del proyecto.
- [ ] Un solo archivo de dependencias (eliminado el duplicado).

---

## Spec 4: Arreglar Tests QEMU E2E

### Problema
Los 3 tests E2E nunca pasan. El fixture `qemu_booted`:
- Intenta login como `root` con password vacío.
- Crea `liveuser` on-the-fly con `useradd` + `passwd -d`.
- Hace `su - liveuser` dentro del mismo shell.
- Es frágil y dependiente del timing de boot.

### Requerimiento
1. Simplificar el flujo de autenticación: usar autologin directo como `liveuser` (ya configurado en `getty@tty1`).
2. Mejorar la detección del prompt con regex más robustos.
3. Agregar test de diagnóstico que solo verifique boot sin assertions complejas.
4. Documentar requisitos para correr tests E2E (KVM, RAM, ISO pre-compilada).

### Criterios de Aceptación
- [ ] `test_iso_boots_to_shell_prompt` pasa consistentemente con KVM.
- [ ] `test_liveuser_stow_symlinks_correct` verifica symlinks correctamente.
- [ ] `test_neovim_init_lua_exists` encuentra init.lua.
- [ ] Fixture `qemu_booted` no necesita crear usuarios on-the-fly.
- [ ] Documentación clara de cómo correr tests E2E.

---

## Spec 5: Completar Configuración de Qtile

### Problema
- `floating_layout = None` en `config.py` deshabilita floating windows.
- El sidebar de la TUI muestra "Qtile Configuration — Coming Soon".
- No hay tests de integración entre TUI y Qtile (cambio de tema).

### Requerimiento
1. Configurar `floating_layout` con layouts útiles (defaults de Qtile o custom).
2. Implementar panel de configuración Qtile en la TUI (al menos tema y barra).
3. Agregar test que verifique que `os_theme.json` escrito por la TUI es leído correctamente por Qtile.

### Criterios de Aceptación
- [ ] `floating_layout` configurado con al menos un layout funcional.
- [ ] TUI permite cambiar tema de Qtile (UI implementada, aunque sea básica).
- [ ] Test de integración TUI -> os_theme.json -> Qtile theme.py.

---

## Spec 6: Implementar CI/CD Básico

### Problema
No hay integración continua. Los tests solo corren localmente.

### Requerimiento
1. Crear `.github/workflows/ci.yml` que:
   - Ejecute tests unitarios en cada push/PR.
   - Verifique sintaxis de archivos Lua y Python.
   - Verifique que `packages.x86_64` no tenga duplicados.
2. (Opcional) Build de ISO en nightly o en tag.

### Criterios de Aceptación
- [ ] Workflow CI ejecuta `pytest tests/unit/ -v` en Ubuntu runner.
- [ ] Workflow falla si hay syntax errors en Python o Lua.
- [ ] Badge de CI en README.md.

---

## Spec 7: Tests para tui_bridge.lua

### Problema
El módulo `tui_bridge.lua` es el puente crítico entre la TUI Python y Neovim, pero no tiene tests.

### Requerimiento
1. Crear tests unitarios para `tui_bridge.lua` usando busted o plenary.nvim.
2. Verificar que:
   - Parsea correctamente `tui_settings.json`.
   - Usa defaults cuando el archivo no existe.
   - Maneja JSON malformado sin crashear.
   - Expone flags correctos en `vim.g.tui_flags`.

### Criterios de Aceptación
- [ ] Tests para `tui_bridge.lua` existen y pasan.
- [ ] Test de JSON malformado verifica fallback a defaults.
- [ ] Test de archivo ausente verifica defaults.

---

## Spec 8: Mejorar `run_qemu.sh`

### Problema
- Usa `blkid` y escribe a `/tmp/iso_extract` sin manejar permisos.
- No verifica que la ISO existe antes de ejecutar.
- Hardcodea `-M pc` (no usa Q35, más moderno).
- No limpia archivos temporales después.

### Requerimiento
1. Agregar verificación de existencia de ISO al inicio.
2. Usar `mktemp -d` para directorio temporal en vez de path fijo.
3. Agregar trap para limpiar temporales al salir.
4. Usar `-M q35` en vez de `-M pc`.
5. Agregar opción `-cpu host` para mejor performance con KVM.

### Criterios de Aceptación
- [ ] Script falla con mensaje claro si la ISO no existe.
- [ ] Directorio temporal se crea y limpia automáticamente.
- [ ] QEMU usa máquina q35 y CPU host.

---

## Spec 9: Agregar Snapper y BTRFS

### Problema
El README menciona Snapper y BTRFS pero:
- No hay configuración de Snapper en `airootfs/`.
- `packages.x86_64` no incluye `snapper` ni `snap-pac`.
- No hay hooks de pacman para snapshots automáticos.

### Requerimiento
1. Agregar `snapper` y `snap-pac` a `packages.x86_64`.
2. Crear configuración de Snapper en `airootfs/etc/snapper/configs/root`.
3. Crear hook de pacman para snapshots pre/post transacción.
4. Configurar timeline de snapshots automáticos.

### Criterios de Aceptación
- [ ] `snapper` y `snap-pac` en `packages.x86_64`.
- [ ] Configuración de Snapper en `airootfs/etc/snapper/`.
- [ ] Hook de pacman para snapshots automáticos.
- [ ] Test que verifique configuración de Snapper.

---

## Spec 10: Documentación de Arquitectura Actualizada

### Problema
- `docs/Arquitectura_Modular.md` está incompleto (solo 25 líneas).
- `docs/ESTADO_ACTUAL.md` tiene información de bugs corregidos pero no el estado real actual.
- No hay diagrama de flujo de boot.

### Requerimiento
1. Actualizar `Arquitectura_Modular.md` con estructura completa actual.
2. Crear `docs/BOOT_FLOW.md` con diagrama del proceso de boot.
3. Agregar `docs/TESTING.md` con instrucciones para correr cada tipo de test.
4. Crear `docs/CONTRIBUTING.md` con guía para nuevos contribuidores.

### Criterios de Aceptación
- [ ] `Arquitectura_Modular.md` refleja la estructura real de `airootfs/`.
- [ ] `BOOT_FLOW.md` existe con diagrama de boot (BIOS/UEFI -> GRUB/Syslinux -> kernel -> initramfs -> systemd -> ly -> qtile).
- [ ] `TESTING.md` documenta tests unitarios y E2E con comandos exactos.
- [ ] `CONTRIBUTING.md` existe con guía de estilo y proceso de PR.

---

## Priorización Sugerida

| Prioridad | Spec | Impacto | Esfuerzo |
|-----------|------|---------|----------|
| P0 | Spec 1: Paquetes faltantes | Alto | Medio |
| P0 | Spec 2: Limpiar artefactos | Bajo | Bajo |
| P1 | Spec 4: Tests E2E | Alto | Alto |
| P1 | Spec 3: Paquete Python TUI | Medio | Medio |
| P2 | Spec 6: CI/CD | Medio | Bajo |
| P2 | Spec 8: Mejorar run_qemu.sh | Medio | Bajo |
| P3 | Spec 5: Qtile completo | Medio | Medio |
| P3 | Spec 7: Tests tui_bridge | Medio | Bajo |
| P3 | Spec 9: Snapper/BTRFS | Alto | Medio |
| P4 | Spec 10: Documentación | Medio | Medio |
