# lambda-env: Hub Principal + Plugin System

## Intent

Crear el entry point `lambda-env` (o `lambdaos-tui`) que funciona como hub principal de la suite. Debe descubrir, cargar y ejecutar módulos de forma dinámica.

## Scope

- Entry point CLI: `lambda-env` abre el menú principal
- Sistema de plugins: escanea `~/.local/share/lambda-env/modules/` y `/usr/share/lambda-env/modules/`
- Cada módulo es un script ejecutable o un paquete Python/Go con un `manifest.json`
- Navegación consistente entre módulos (teclas de navegación unificadas)
- Theme engine: colores consistentes en todos los módulos

## Requirements

1. El hub debe listar todos los módulos disponibles agrupados por categoría (System, Apps, Ops, Setup)
2. Cada módulo debe declarar: nombre, descripción, categoría, dependencias, si requiere root
3. El hub debe validar que las dependencias del módulo están instaladas antes de ejecutarlo
4. La navegación debe ser consistente: flechas para moverse, Enter para seleccionar, Esc para volver, q para salir
5. Debe funcionar en tty pura (sin X11)

## Scenarios

### Escenario 1: Usuario abre lambda-env por primera vez
- Se muestra el menú principal con categorías
- El usuario navega con flechas y Enter
- Selecciona "System" → ve lista de módulos de sistema
- Selecciona "Screen" → se carga el módulo de pantalla

### Escenario 2: Módulo con dependencias faltantes
- El usuario selecciona "Recording" (OBS)
- El hub detecta que OBS no está instalado
- Muestra: "Este módulo requiere: obs-studio. Instalar ahora? [y/N]"
- Si el usuario acepta, instala y luego abre el módulo

### Escenario 3: Módulo que requiere root
- El usuario selecciona "Services"
- El hub detecta que requiere privilegios de root
- Ejecuta el módulo con `sudo` o pide contraseña

## Technical Notes

- Framework a definir: `textual` (Python), `bubbletea` (Go), o `whiptail/dialog` (bash)
- El manifest de cada módulo: `{ "name", "description", "category", "requires_root", "dependencies" }`
- Los módulos se comunican con el hub via exit codes y stdout
- El settings unificado (`settings.json`) es leído por el hub y pasado a los módulos como contexto

## Dependencies

- Ninguno (es el módulo base)
