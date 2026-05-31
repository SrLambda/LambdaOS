# lambda-env: First Boot Wizard (setup-01)

## Intent

Módulo TUI que se ejecuta automáticamente en el primer boot de LambdaOS (live o instalada) para guiar al usuario en la configuración inicial.

## Scope

- Bienvenida + explicación de LambdaOS
- Selección de idioma del sistema
- Selección de layout de teclado
- Configuración de WiFi (si no hay Ethernet)
- Creación de usuario (si es instalación nueva)
- Selección de tema visual
- Configuración de zona horaria
- Resumen y aplicar

## Requirements

1. Ejecutar automáticamente en primer boot (detectar via flag file `~/.config/lambdaos/first-boot-done`)
2. Flujo paso a paso con navegación simple (flechas + Enter)
3. Cada paso valida la entrada antes de continuar
4. Al finalizar, crea el flag file para no ejecutarse de nuevo
5. Todos los settings se escriben en `settings.json`
6. Opción de "Skip" para usuarios avanzados

## Scenarios

### Escenario 1: Primer boot en live ISO
- Usuario bootea LambdaOS por primera vez
- Se abre el wizard automáticamente
- Paso 1: Bienvenida → Next
- Paso 2: Idioma → "Español (Latinoamérica)" → Next
- Paso 3: Keyboard → "es" / "latam" → Next
- Paso 4: WiFi → escanea → selecciona red → ingresa password → Next
- Paso 5: Theme → selecciona "catppuccin" → Next
- Paso 6: Timezone → "America/Argentina/Buenos_Aires" → Next
- Paso 7: Resumen → "Apply"
- Wizard cierra, settings aplicados, flag creado

### Escenario 2: Usuario avanzado saltea wizard
- Presiona "Skip wizard" en la bienvenida
- Flag se crea igual para no preguntar de nuevo
- Usuario configura manualmente con lambda-env

## Technical Notes

- Flag file: `~/.config/lambdaos/first-boot-done`
- Ejecutar desde `.zlogin` o systemd user service `lambdaos-wizard.service`
- Cada paso escribe en `settings.json` (que se crea en el proceso)
- WiFi: usar `iwctl` para conectar
- Keyboard: `localectl set-x11-keymap`
- Theme: escribir en `settings.json` → Qtile lee al iniciar
- Timezone: `timedatectl set-timezone`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- Todos los módulos que configura (keyboard, network, appearance, datetime)
