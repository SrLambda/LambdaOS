# Especificaciones Técnicas: TUI "Preferencias del Sistema"

## 1. Visión General
Se desarrollará una interfaz gráfica en terminal (TUI) usando Python y **Textual** para gestionar las configuraciones modulares del sistema operativo, actuando como un panel de control unificado.

## 2. Stack Tecnológico
- **Lenguaje:** Python (>= 3.11)
- **Framework TUI:** Textual (`textual>=0.47.0`)
- **Interoperabilidad:** La TUI modificará archivos puente (JSON/Variables de Entorno) que las aplicaciones (como Neovim y Qtile) consumirán, evitando que la TUI tenga que reescribir código Lua o Python.

## 3. Diseño de la Interfaz (UI/UX)
- **Header & Footer:** Título de la app y atajos globales (Q: Salir, S: Guardar).
- **Sidebar (Izquierda):** Menú de navegación (Módulo Neovim, Módulo Qtile).
- **Content Area (Derecha):** Formularios interactivos y *switches*.

## 4. Alcance: Iteración 1 (Integración con Neovim)
Basado en el WBS de Neovim, la TUI interactuará de la siguiente manera:
- **Gestión del Tema (WBS 1.3 / 2.1):** 
  - La TUI ofrecerá un selector de temas (Catppuccin, Gruvbox, etc.).
  - Al guardar, la TUI actualizará la variable de entorno `$NVIM_THEME` en el archivo de configuración base del usuario (ej. `~/.profile` o el entorno de Qtile) para que el módulo `lua/core/env.lua` de Neovim lo recoja en su próximo inicio.
- **Gestión de Plugins (Toggles):**
  - La TUI escribirá un archivo `~/.config/nvim/tui_settings.json`.
  - Este JSON contendrá banderas: `{"enable_lsp": true, "enable_copilot": false, "enable_neotree": true}`.
  - El gestor `lazy.nvim` leerá este JSON en `init.lua` para establecer la propiedad `enabled = true/false` en los plugins correspondientes (WBS 4.1, 5.2, 6.0).

## 5. Reglas de Desarrollo para Agentes IA (`dev_tui`)
- Utilizar CSS de Textual (`.tcss`) para el diseño. Usar *Reactive Attributes* para el estado.
- **Manejo de Errores:** Si los dotfiles no existen en `~/dotfiles/nvim`, mostrar un modal de error.
- **Testing Primero (TDD):** Usar `textual.testing` para asegurar que el botón "Guardar Tema" efectivamente escribe la variable `$NVIM_THEME` correctamente.