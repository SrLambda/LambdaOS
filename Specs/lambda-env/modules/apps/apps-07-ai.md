# lambda-env: AI Agent Config (apps-07)

## Intent

Módulo TUI para configurar OpenCode y agentes IA: modelo, provider, API keys, reglas de comportamiento, proyectos.

## Scope

- Configurar provider de IA (OpenAI, Anthropic, local Ollama)
- Configurar API keys
- Seleccionar modelo por defecto
- Configurar reglas del agente (system prompt, context)
- Gestionar proyectos: directorios con config específica
- Configurar MCP servers

## Requirements

1. Leer/escribir config de OpenCode
2. Soportar múltiples providers con selección
3. API keys almacenadas de forma segura (no en texto plano si es posible)
4. Modelos: lista de modelos disponibles por provider
5. Proyectos: cada proyecto puede tener su propio agent config
6. Backend: OpenCode config files, variables de entorno

## Scenarios

### Escenario 1: Configurar OpenAI como provider
- Usuario abre AI Config → Provider
- Selecciona "OpenAI"
- Ingresa API key
- Selecciona modelo: "gpt-4o"
- Aplica

### Escenario 2: Configurar proyecto con reglas custom
- Abre AI Config → Projects → "Add project"
- Path: `~/Projects/LambdaOS`
- System prompt: "You are a Linux kernel developer..."
- Modelo: "claude-sonnet-4-20250514"
- Aplica

## Technical Notes

- OpenCode config: `~/.config/opencode/` o `opencode.json` en proyecto
- API keys: considerar `pass` o `keyring` para almacenamiento seguro
- Ollama como provider local (agregar a packages.x86_64 si se quiere)
- MCP servers: configurar en `opencode.json`

## Dependencies

- `core/01-hub-plugin-system`
- `core/02-settings-schema`
- opencode (distribuir via npm/cargo/empaquetado propio)
