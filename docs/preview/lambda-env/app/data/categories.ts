import { ICONS } from './icon-map';

export interface ModuleDef {
  id: string;
  name: string;
  description: string;
  icon: string;
  fallbackIcon?: string;
  rootRequired?: boolean;
  implemented?: boolean;
}

export interface CategoryDef {
  id: string;
  name: string;
  icon: string;
  fallbackIcon?: string;
  description: string;
  modules: ModuleDef[];
}

export const CATEGORIES: CategoryDef[] = [
  {
    id: 'sistema',
    name: 'Sistema',
    icon: ICONS.categories.system.nerd,
    fallbackIcon: ICONS.categories.system.fallback,
    description: 'Configuración base del sistema operativo',
    modules: [
      { id: 'pantalla',       name: 'Pantalla',        icon: ICONS.modules.display.nerd,        fallbackIcon: ICONS.modules.display.fallback,        description: 'Resolución, múltiples monitores, brillo, perfiles',          implemented: true },
      { id: 'audio',          name: 'Audio',            icon: ICONS.modules.audio.nerd,          fallbackIcon: ICONS.modules.audio.fallback,          description: 'Volumen, dispositivos de entrada y salida, perfiles',        implemented: true },
      { id: 'red',            name: 'Red',              icon: ICONS.modules.network.nerd,        fallbackIcon: ICONS.modules.network.fallback,        description: 'WiFi, Ethernet, VPN, firewall, DNS',                        implemented: true },
      { id: 'bluetooth',      name: 'Bluetooth',        icon: ICONS.modules.bluetooth.nerd,      fallbackIcon: ICONS.modules.bluetooth.fallback,      description: 'Dispositivos emparejados, visibilidad, conexión',           implemented: true },
      { id: 'energia',        name: 'Energía',          icon: ICONS.modules.power.nerd,           fallbackIcon: ICONS.modules.power.fallback,           description: 'Batería, suspensión, rendimiento, acciones del sistema',    implemented: true },
      { id: 'teclado',        name: 'Teclado',          icon: ICONS.modules.keyboard.nerd,      fallbackIcon: ICONS.modules.keyboard.fallback,      description: 'Layout, variante, velocidad de repetición, atajos',         implemented: true },
      { id: 'fecha',          name: 'Fecha y Hora',     icon: '◷', fallbackIcon: '◷', description: 'Zona horaria, NTP, formato de fecha',                      implemented: true },
      { id: 'usuarios',       name: 'Usuarios',         icon: '⊙', fallbackIcon: '⊙', description: 'Cuentas, grupos, contraseñas, permisos',         rootRequired: true, implemented: true },
      { id: 'apariencia',     name: 'Apariencia',       icon: ICONS.modules.appearance.nerd,    fallbackIcon: ICONS.modules.appearance.fallback,    description: 'Tema, colores, fuente, wallpaper, opacidad',               implemented: true },
      { id: 'defaults',       name: 'Apps por Defecto', icon: ICONS.modules.defaults.nerd,        fallbackIcon: ICONS.modules.defaults.fallback,        description: 'Navegador, editor, terminal, reproductor multimedia',       implemented: true },
      { id: 'autostart',      name: 'Autostart',        icon: '▸', fallbackIcon: '▸', description: 'Programas y servicios al iniciar sesión',                   implemented: true },
      { id: 'servicios',      name: 'Servicios',        icon: '◎', fallbackIcon: '◎', description: 'Daemons del sistema, start/stop/restart',       rootRequired: true, implemented: true },
      { id: 'actualizaciones',name: 'Actualizaciones',  icon: '↻', fallbackIcon: '↻', description: 'Paquetes disponibles, actualizaciones automáticas',        implemented: true },
      { id: 'seguridad',      name: 'Seguridad',        icon: ICONS.modules.security.nerd,       fallbackIcon: ICONS.modules.security.fallback,       description: 'Firewall, SSH, GPG, LUKS, políticas',          rootRequired: true, implemented: true },
      { id: 'fuentes',        name: 'Fuentes',          icon: 'A', fallbackIcon: 'A', description: 'Fuentes instaladas, previsualización, instalación',         implemented: true },
      { id: 'notificaciones', name: 'Notificaciones',   icon: '◻', fallbackIcon: '◻', description: 'Posición, timeout, por aplicación, DND',                   implemented: true },
    ],
  },
  {
    id: 'aplicaciones',
    name: 'Aplicaciones',
    icon: ICONS.categories.apps.nerd,
    fallbackIcon: ICONS.categories.apps.fallback,
    description: 'Herramientas y apps del entorno de usuario',
    modules: [
      { id: 'neovim',     name: 'Neovim',       icon: ICONS.modules.neovim.nerd,      fallbackIcon: ICONS.modules.neovim.fallback,      description: 'Plugins, keymaps, colorscheme, LSP',              implemented: true },
      { id: 'qtile',      name: 'Qtile',         icon: ICONS.modules.qtile.nerd,       fallbackIcon: ICONS.modules.qtile.fallback,       description: 'Layouts, workspaces, keybindings, bars',          implemented: true },
      { id: 'screenshot', name: 'Screenshot',    icon: '◧', fallbackIcon: '◧', description: 'Área, formato, destino, atajo de teclado',        implemented: false },
      { id: 'grabacion',  name: 'Grabación',     icon: '⏺', fallbackIcon: '⏺', description: 'Fuente de video/audio, formato, directorio',      implemented: false },
      { id: 'terminal',   name: 'Terminal',      icon: '◰', fallbackIcon: '◰', description: 'Emulador, shell, fuente, colores, scrollback',    implemented: false },
      { id: 'filemanager',name: 'File Manager',  icon: '◱', fallbackIcon: '◱', description: 'Vista por defecto, bookmarks, editor de texto',   implemented: false },
      { id: 'ai',         name: 'AI',            icon: '◆', fallbackIcon: '◆', description: 'Integración con modelos, clave API, proveedor',   implemented: false },
    ],
  },
  {
    id: 'operaciones',
    name: 'Operaciones',
    icon: ICONS.categories.ops.nerd,
    fallbackIcon: ICONS.categories.ops.fallback,
    description: 'Mantenimiento y monitoreo del sistema',
    modules: [
      { id: 'monitor',      name: 'Monitor',      icon: ICONS.modules['hardware-dashboard'].nerd, fallbackIcon: ICONS.modules['hardware-dashboard'].fallback, description: 'CPU, RAM, disco, red en tiempo real',           implemented: true },
      { id: 'almacenamiento',name: 'Almacenamiento',icon: ICONS.modules.storage.nerd, fallbackIcon: ICONS.modules.storage.fallback, description: 'Discos, particiones, uso, montar/desmontar',  implemented: true },
      { id: 'logs',         name: 'Logs',          icon: ICONS.modules.logs.nerd,        fallbackIcon: ICONS.modules.logs.fallback,        description: 'Journal del sistema, filtros, exportar',       implemented: true },
      { id: 'backup',       name: 'Backup',        icon: '⊙', fallbackIcon: '⊙', description: 'Snapshots, rsync, destino, programación',      implemented: false },
      { id: 'dotfiles',     name: 'Dotfiles',      icon: ICONS.modules.dotfiles.nerd,  fallbackIcon: ICONS.modules.dotfiles.fallback,  description: 'Sincronización, perfil activo, stow',          implemented: true },
    ],
  },
  {
    id: 'setup',
    name: 'Setup',
    icon: '◇',
    fallbackIcon: '◇',
    description: 'Configuración inicial y perfiles',
    modules: [
      { id: 'wizard',     name: 'Asistente',  icon: '◑', fallbackIcon: '◑', description: 'Configuración guiada inicial del entorno',          implemented: false },
      { id: 'instalador', name: 'Instalador', icon: '↓', fallbackIcon: '↓', description: 'Paquetes, AUR, flatpaks, snippets',                implemented: false },
      { id: 'perfiles',   name: 'Perfiles',   icon: '◐', fallbackIcon: '◐', description: 'Perfiles de configuración, exportar/importar',      implemented: false },
    ],
  },
];
