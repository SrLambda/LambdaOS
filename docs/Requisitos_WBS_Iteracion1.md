# Requisitos

|**WBS ID**|**Fase de Desarrollo**|**Tarea Específica**|**Archivo/Módulo Objetivo**|**Criterio de Aceptación (DoD)**|
|---|---|---|---|---|
|**1.0**|**Arquitectura Base**|1.1 Configurar gestor de paquetes (`lazy.nvim`)|`init.lua`, `lua/core/lazy.lua`|Gestor arranca sin errores al abrir `nvim`.|
|||1.2 Configurar opciones nativas y mapeos|`lua/core/options.lua`, `lua/core/keymaps.lua`|Tabs a 4 espacios, textwidth a 80, wrap desactivado.|
|||1.3 Implementar lectura de Qtile (Env Vars)|`lua/core/env.lua`|Función global que lee `$NVIM_THEME` y aplica fallbacks.|
|**2.0**|**Interfaz de Usuario (UI)**|2.1 Implementar Tema y Colores|`lua/plugins/theme.lua`|Catppuccin configurado, respondiendo a la variable de Qtile.|
|||2.2 Configurar Lualine, Alpha y Bufferline|`lua/plugins/ui.lua`|Barra de estado, pestañas superiores y pantalla de inicio funcionales.|
|||2.3 Configurar utilidades visuales|`lua/plugins/ui.lua`|`indent-blankline` y web-devicons instalados y visibles.|
|**3.0**|**Motor de Sintaxis**|3.1 Integrar y configurar Tree-sitter|`lua/plugins/treesitter.lua`|Resaltado sintáctico activo, `ensure_installed = "all"`. Parser "comment" para Batch.|
|**4.0**|**Navegación y Flujo**|4.1 Configurar Telescope, Neo-tree y Harpoon|`lua/plugins/navigation.lua`|Búsqueda difusa y explorador de archivos mapeados a atajos de teclado.|
|||4.2 Configurar Autopairs, Surround y Comment|`lua/plugins/editing.lua`|Cierre de llaves automático y atajos para comentarios activos.|
|**5.0**|**Herramientas de Apoyo**|5.1 Integrar Gitsigns y Toggleterm|`lua/plugins/tools.lua`|Integración de git en el gutter y terminal flotante funcional.|
|||5.2 Configurar Copilot|`lua/plugins/ai.lua`|`copilot.lua` inicializado y listo para autenticación.|
|**6.0**|**Core IDE: LSP y Auto-completado**|6.1 Instalar Mason y Mason-lspconfig|`lua/plugins/lsp.lua`|Binarios manejados centralmente.|
|||6.2 Configurar Servidores Base|`lua/plugins/lsp.lua`|Pyright, tsserver, clangd, html, css, lua_ls, bashls, docker, yaml, json, toml operativos.|
|||6.3 Configurar Servidores Complejos|`ftplugin/java.lua`, `ftplugin/rust.lua`|`nvim-jdtls` y `rustaceanvim` cargados solo en sus respectivos tipos de archivo.|
|**7.0**|**Core IDE: Formateo y Linters**|7.1 Configurar `conform.nvim`|`lua/plugins/formatting.lua`|Black, Prettier, google-java-format, etc., formateando al guardar (format-on-save).|
|||7.2 Configurar `nvim-lint`|`lua/plugins/linting.lua`|Flake8, eslint_d, shellcheck ejecutándose asíncronamente.|
|**8.0**|**Entornos Específicos**|8.1 Data Science (REPL)|`lua/plugins/data.lua`|`iron.nvim` o `molten.nvim` configurados para evaluación de bloques.|
|||8.2 Documentación (Markdown, LaTeX)|`lua/plugins/docs.lua`|Previsualizadores (markdown, mermaid, plantuml) y `vimtex` (con Zathura) funcionales. Corrector ortográfico en español activado (`spelllang=es`).|

# WBS

|**WBS ID**|**Fase de Desarrollo**|**Tarea Específica**|**Archivo/Módulo Objetivo**|**Criterio de Aceptación (DoD)**|
|---|---|---|---|---|
|**1.0**|**Arquitectura Base**|1.1 Configurar gestor de paquetes (`lazy.nvim`)|`init.lua`, `lua/core/lazy.lua`|Gestor arranca sin errores al abrir `nvim`.|
|||1.2 Configurar opciones nativas y mapeos|`lua/core/options.lua`, `lua/core/keymaps.lua`|Tabs a 4 espacios, textwidth a 80, wrap desactivado.|
|||1.3 Implementar lectura de Qtile (Env Vars)|`lua/core/env.lua`|Función global que lee `$NVIM_THEME` y aplica fallbacks.|
|**2.0**|**Interfaz de Usuario (UI)**|2.1 Implementar Tema y Colores|`lua/plugins/theme.lua`|Catppuccin configurado, respondiendo a la variable de Qtile.|
|||2.2 Configurar Lualine, Alpha y Bufferline|`lua/plugins/ui.lua`|Barra de estado, pestañas superiores y pantalla de inicio funcionales.|
|||2.3 Configurar utilidades visuales|`lua/plugins/ui.lua`|`indent-blankline` y web-devicons instalados y visibles.|
|**3.0**|**Motor de Sintaxis**|3.1 Integrar y configurar Tree-sitter|`lua/plugins/treesitter.lua`|Resaltado sintáctico activo, `ensure_installed = "all"`. Parser "comment" para Batch.|
|**4.0**|**Navegación y Flujo**|4.1 Configurar Telescope, Neo-tree y Harpoon|`lua/plugins/navigation.lua`|Búsqueda difusa y explorador de archivos mapeados a atajos de teclado.|
|||4.2 Configurar Autopairs, Surround y Comment|`lua/plugins/editing.lua`|Cierre de llaves automático y atajos para comentarios activos.|
|**5.0**|**Herramientas de Apoyo**|5.1 Integrar Gitsigns y Toggleterm|`lua/plugins/tools.lua`|Integración de git en el gutter y terminal flotante funcional.|
|||5.2 Configurar Copilot|`lua/plugins/ai.lua`|`copilot.lua` inicializado y listo para autenticación.|
|**6.0**|**Core IDE: LSP y Auto-completado**|6.1 Instalar Mason y Mason-lspconfig|`lua/plugins/lsp.lua`|Binarios manejados centralmente.|
|||6.2 Configurar Servidores Base|`lua/plugins/lsp.lua`|Pyright, tsserver, clangd, html, css, lua_ls, bashls, docker, yaml, json, toml operativos.|
|||6.3 Configurar Servidores Complejos|`ftplugin/java.lua`, `ftplugin/rust.lua`|`nvim-jdtls` y `rustaceanvim` cargados solo en sus respectivos tipos de archivo.|
|**7.0**|**Core IDE: Formateo y Linters**|7.1 Configurar `conform.nvim`|`lua/plugins/formatting.lua`|Black, Prettier, google-java-format, etc., formateando al guardar (format-on-save).|
|||7.2 Configurar `nvim-lint`|`lua/plugins/linting.lua`|Flake8, eslint_d, shellcheck ejecutándose asíncronamente.|
|**8.0**|**Entornos Específicos**|8.1 Data Science (REPL)|`lua/plugins/data.lua`|`iron.nvim` o `molten.nvim` configurados para evaluación de bloques.|
|||8.2 Documentación (Markdown, LaTeX)|`lua/plugins/docs.lua`|Previsualizadores (markdown, mermaid, plantuml) y `vimtex` (con Zathura) funcionales. Corrector ortográfico en español activado (`spelllang=es`).|