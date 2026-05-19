# Arquitectura Modular del Sistema

## 1. Filosofía y Objetivo
El objetivo de esta distribución Arch Linux personalizada es ser **100% modular**. Cualquier usuario que instale o ejecute esta ISO debe poder elegir qué partes de la configuración (Dotfiles) desea utilizar y cuáles no, sin romper el sistema. Se utilizará **GNU Stow** como gestor de enlaces simbólicos.

## 2. Estructura de Directorios (GNU Stow y Lua)
Todas las configuraciones de usuario vivirán como "paquetes" independientes dentro de la carpeta plantilla (`skel`) del LiveCD (`airootfs/etc/skel/dotfiles/`).

**Estructura esperada para la Iteración 1 (Neovim):**

airootfs/etc/skel/dotfiles/
├── nvim/                   
│   └── .config/
│       └── nvim/
│           ├── init.lua           # Entry point de Neovim
│           ├── tui_settings.json  # Archivo puente generado por la TUI
│           ├── lua/
│           │   ├── core/          # WBS 1.0 (lazy.lua, options.lua, keymaps.lua, env.lua)
│           │   ├── plugins/       # WBS 2.0 al 8.0 (theme, ui, lsp, formatting, etc.)
│           │   └── utils/
│           └── ftplugin/          # WBS 6.3 (java.lua, rust.lua)
├── qtile/                  
│   └── ...
└── tui_os/                 
    └── ...