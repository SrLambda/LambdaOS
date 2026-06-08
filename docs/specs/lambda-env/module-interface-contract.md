# lambda-env: Module Interface Contract

## Overview

This document defines the contract between the `lambda-env` hub and its modules. Every module must comply with this contract to be discovered and executed by the hub.

**Este documento define el contrato entre el hub de `lambda-env` y sus módulos. Todo módulo debe cumplir este contrato para ser descubierto y ejecutado por el hub.**

---

## Architecture Decision: Go as Primary Language

**EN**: The hub and all default modules are written in Go. Other languages (bash, Python) are ONLY permitted for adapter modules that wrap external service APIs where Go would add unnecessary complexity (e.g., `bluetoothctl` wrapper, `iwctl` wrapper). This keeps the binary self-contained, dependency-free, and fast.

**ES**: El hub y todos los módulos default están escritos en Go. Otros lenguajes (bash, Python) SOLO se permiten para módulos adaptadores que envuelven APIs de servicios externos donde Go agregaría complejidad innecesaria (ej: wrapper de `bluetoothctl`, wrapper de `iwctl`). Esto mantiene el binario autocontenido, sin dependencias y rápido.

---

## 1. Module Discovery

**EN**: The hub scans module directories at startup. No caching — the module list is small (max ~40 modules), scanning two directories takes <10ms. Caching adds complexity (stale cache, invalidation) for negligible gain.

**ES**: El hub escanea los directorios de módulos al inicio. Sin caché — la lista de módulos es pequeña (máx ~40), escanear dos directorios toma <10ms. El caché agrega complejidad (caché desactualizado, invalidación) con ganancia negligible.

### Discovery Paths

| Priority | Path | Purpose |
|---|---|---|
| 1 | `/usr/share/lambda-env/modules/` | System modules (installed via pacman) |
| 2 | `~/.local/share/lambda-env/modules/` | User modules (custom, overrides) |

User modules override system modules with the same name.

### Discovery Algorithm

```
1. Scan /usr/share/lambda-env/modules/ for directories containing manifest.json
2. Scan ~/.local/share/lambda-env/modules/ for directories containing manifest.json
3. Merge: user modules override system modules with same name
4. Sort by category, then name
5. Validate each manifest (required fields present, valid JSON)
6. Skip invalid modules with warning
```

---

## 2. Module Structure

Each module is a directory with this structure:

```
/usr/share/lambda-env/modules/<module-name>/
├── manifest.json          ← Required: module metadata
├── module                 ← Required: executable (Go binary, script with shebang)
├── README.md              ← Optional: documentation
└── config/                ← Optional: default config files
    └── <module>.json
```

### manifest.json

```json
{
  "name": "screen",
  "version": "0.1.0",
  "description": "Manage display configuration (xrandr)",
  "description_es": "Gestionar configuración de pantalla (xrandr)",
  "category": "system",
  "icon": "display",
  "requires_root": false,
  "dependencies": ["xorg-xrandr"],
  "min_hub_version": "0.1.0",
  "timeout": 30,
  "tags": ["display", "monitor", "resolution", "xrandr"],
  "author": "LambdaOS Team"
}
```

### Required Fields

| Field | Type | Description |
|---|---|---|
| `name` | string | Unique module identifier (lowercase, hyphens) |
| `version` | string | Semantic version of the module |
| `description` | string | Short description (English) |
| `description_es` | string | Short description (Spanish) |
| `category` | string | One of: `system`, `apps`, `ops`, `setup` |
| `requires_root` | boolean | Whether the module needs sudo |
| `dependencies` | string[] | Package names required |
| `min_hub_version` | string | Minimum hub version for API compatibility |

### Optional Fields

| Field | Type | Description |
|---|---|---|
| `icon` | string | Icon name for future GUI integration |
| `timeout` | int | Max execution time in seconds (default: 30) |
| `tags` | string[] | Search/filter keywords |
| `author` | string | Module author |

---

## 3. Communication Protocol: JSON over stdout/stderr

**EN**: Modules communicate with the hub via JSON on stdout. This was chosen over TUI control codes because: (1) JSON is parseable by any language, (2) easy to test with automated tools, (3) structured data survives pipe/redirection, (4) TUI rendering is the hub's responsibility, not the module's.

**ES**: Los módulos se comunican con el hub via JSON en stdout. Se eligió sobre códigos de control TUI porque: (1) JSON es parseable por cualquier lenguaje, (2) fácil de testear con herramientas automatizadas, (3) datos estructurados sobreviven pipes/redirección, (4) el renderizado TUI es responsabilidad del hub, no del módulo.

### Module Execution Flow

```
Hub                          Module
 │                              │
 ├── exec module run ──────────→│
 │                              │
 │ ←── JSON on stdout ──────────┤  (structured response)
 │ ←── text on stderr ──────────┤  (human-readable messages, errors)
 │ ←── exit code ───────────────┤  (0=success, 1=error, 2=warning)
 │                              │
```

### Module Input

The hub passes input to the module via environment variables and stdin:

```bash
# Environment variables
LAMBDA_ENV_ACTION="run"          # Action: run, validate, help
LAMBDA_ENV_SETTINGS="/home/user/.config/lambdaos/settings.json"
LAMBDA_ENV_HUB_VERSION="0.1.0"
LAMBDA_ENV_LOCALE="en_US"        # or es_AR

# Stdin (optional, for actions that need data)
echo '{"output": "HDMI-1", "mode": "1920x1080@60"}' | module run
```

### Module Output (stdout — JSON only)

#### Success Response

```json
{
  "status": "ok",
  "action": "run",
  "data": {
    "outputs": [
      {
        "name": "eDP-1",
        "connected": true,
        "current_mode": "1920x1080@60",
        "available_modes": ["1920x1080@60", "1920x1080@50", "1366x768@60"]
      },
      {
        "name": "HDMI-1",
        "connected": true,
        "current_mode": null,
        "available_modes": ["1920x1080@60", "2560x1440@60"]
      }
    ]
  }
}
```

#### Error Response

```json
{
  "status": "error",
  "action": "run",
  "code": "XRANDR_FAILED",
  "message": "Failed to query xrandr: command not found",
  "message_es": "Error al consultar xrandr: comando no encontrado",
  "details": {
    "exit_code": 127,
    "stderr": "xrandr: command not found"
  }
}
```

#### Warning Response

```json
{
  "status": "warning",
  "action": "run",
  "code": "DEPENDENCY_MISSING",
  "message": "xorg-xrandr is not installed. Install it to use this module.",
  "message_es": "xorg-xrandr no está instalado. Instálalo para usar este módulo.",
  "suggestion": "pacman -S xorg-xrandr"
}
```

### Module Output (stderr — human-readable text)

Modules write human-readable messages to stderr for logging and debugging:

```
[screen] Querying displays via xrandr...
[screen] Found 2 connected outputs: eDP-1, HDMI-1
[screen] Applying mode 1920x1080@60 to HDMI-1
```

### Exit Codes

| Code | Meaning | Hub Behavior |
|---|---|---|
| `0` | Success | Display data, return to menu |
| `1` | Error | Show error message, log to file, return to menu |
| `2` | Warning | Show warning, continue if user confirms |

---

## 4. Error Handling

**EN**: When a module fails (exit code != 0), the hub does BOTH: (1) shows the error message to the user in the TUI, AND (2) logs the full error to `/var/log/lambda-env/modules.log`. The user-facing message includes a reference to the log file.

**ES**: Cuando un módulo falla (exit code != 0), el hub hace AMBAS cosas: (1) muestra el mensaje de error al usuario en la TUI, Y (2) loguea el error completo a `/var/log/lambda-env/modules.log`. El mensaje visible al usuario incluye referencia al archivo de log.

### User-Facing Error

```
┌─────────────────────────────────────────────┐
│  Error in module: screen                    │
│                                             │
│  Failed to query xrandr: command not found  │
│                                             │
│  Full log: /var/log/lambda-env/modules.log  │
│                                             │
│              [ OK ]                         │
└─────────────────────────────────────────────┘
```

### Log File Format

```
2026-05-30T14:32:01Z [ERROR] module=screen action=run exit_code=1
  stdout: {"status":"error","code":"XRANDR_FAILED","message":"..."}
  stderr: [screen] Querying displays via xrandr...
          xrandr: command not found
  env: LAMBDA_ENV_ACTION=run, LAMBDA_ENV_LOCALE=en_US
```

---

## 5. Language Policy

**EN**: Default language for hub and modules is Go. Other languages are ONLY permitted for adapter modules that wrap external service CLI tools where Go would add unnecessary complexity.

**ES**: El lenguaje default para el hub y módulos es Go. Otros lenguajes SOLO se permiten para módulos adaptadores que envuelven herramientas CLI de servicios externos donde Go agregaría complejidad innecesaria.

### When to Use Non-Go

| Scenario | Language | Example |
|---|---|---|
| Wrapping `bluetoothctl` | bash | `modules/system/bluetooth/module` (bash script) |
| Wrapping `iwctl` | bash | `modules/system/network/module` (bash script) |
| Wrapping `pacman`/`yay` | bash | `modules/system/updates/module` (bash script) |
| Core logic, data processing | Go | `modules/system/screen/module` (Go binary) |
| Settings manipulation | Go | `modules/apps/neovim/module` (Go binary) |

### Adapter Module Pattern

When using bash/python for an adapter:

```bash
#!/usr/bin/env bash
# modules/system/bluetooth/module
# This is an adapter for bluetoothctl CLI

set -euo pipefail

action="${LAMBDA_ENV_ACTION:-run}"

case "$action" in
  run)
    # Execute bluetoothctl, parse output, emit JSON
    output=$(bluetoothctl list 2>&1) || {
        echo '{"status":"error","code":"BTCTL_FAILED","message":"bluetoothctl failed"}'
        exit 1
    }
    # Parse and emit JSON...
    ;;
  help)
    echo '{"status":"ok","action":"help","data":{"usage":"lambda-env bluetooth run"}}'
    ;;
esac
```

The adapter MUST still emit JSON on stdout and follow the error handling contract.

---

## 6. Settings Integration

Modules read settings from the unified schema at `~/.config/lambdaos/settings.json` (path provided via `LAMBDA_ENV_SETTINGS` env var).

**Rules**:
- The hub is the ONLY writer to `settings.json` (via atomic write: temp file + rename)
- Modules READ their section from settings.json
- Modules WRITE to settings.json by emitting a `settings_delta` in their response:

```json
{
  "status": "ok",
  "action": "run",
  "data": { ... },
  "settings_delta": {
    "display": {
      "active_profile": "home"
    }
  }
}
```

The hub merges the delta into settings.json atomically.

---

## 7. Module Lifecycle

```
1. Discovery: hub finds module directories with manifest.json
2. Validation: hub checks manifest fields, dependencies, min_hub_version
3. Dependency check: hub verifies dependencies are installed (pacman -Q)
4. Execution: hub runs module with env vars + optional stdin
5. Response: hub parses JSON from stdout, text from stderr
6. Settings merge: if settings_delta present, hub merges atomically
7. Display: hub renders data in TUI
8. Logging: hub logs execution to /var/log/lambda-env/modules.log
```
