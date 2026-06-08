# Specs Lambda-env UI/UX

> **Propósito de este documento**: Especificación conceptual de la experiencia de usuario e interfaz visual de `lambda-env`. Está escrito para ser la fuente de verdad que un diseñador use al prototipar en Figma. No describe tecnología de implementación — describe lo que el usuario ve, siente y hace.

---

## 1. Visión General

### Qué es lambda-env

El centro de configuración unificado de LambdaOS. El equivalente conceptual a "Configuración del Sistema" en macOS, "Settings" en GNOME, o "System Settings" en KDE: un único lugar desde donde el usuario configura cada aspecto del sistema operativo. La diferencia es que lambda-env está diseñado para correr en una terminal, pero el diseño visual debe tratarse como una aplicación de escritorio completa — con jerarquía clara, feedback inmediato y una experiencia pulida.

### Filosofía de Diseño

1. **Jerarquía clara**: El usuario siempre sabe dónde está, qué puede hacer, y cómo volver atrás. Tres niveles de profundidad máximo.
2. **Inmediatez**: Cada acción tiene feedback visible al instante. Sin estados de carga invisibles. Sin silencios.
3. **Consistencia**: El mismo patrón de interacción se repite en cada módulo. Si aprendiste a usar uno, sabés usar todos.
4. **Seguridad**: Las acciones destructivas o irreversibles piden confirmación. Los cambios de display tienen cuenta regresiva con auto-revert.
5. **Accesibilidad**: Contraste suficiente, navegación sin mouse posible, sin dependencia de colores como único canal de información.
6. **Elegancia minimalista**: La interfaz no compite con el contenido. Los módulos son los protagonistas; el marco es invisible.

### Personas

| Persona | Necesidad principal | Expectativa de UI |
|---------|-------------------|-------------------|
| **Usuario final** | Configurar apariencia, audio, pantalla, energía sin tocar archivos de configuración | Descubribilidad, seguridad (deshacer), feedback claro |
| **Sysadmin** | Gestionar servicios, red, firewall, usuarios desde terminal | Velocidad, precisión, atajos |
| **Developer** | Activar/desactivar tooling, gestionar dotfiles, configurar editor | Operaciones por lote, previsibilidad |

---

## 2. Arquitectura de Información

### Jerarquía

```
Lambda-env
  ├── Categorías (4)
  │     ├── Sistema (16 módulos)
  │     ├── Aplicaciones (7 módulos)
  │     ├── Operaciones (5 módulos)
  │     └── Setup (3 módulos)
  │
  └── Módulo individual
        └── Acciones del módulo (widgets interactivos)
```

### Categorías

| Categoría | Módulos | Descripción |
|-----------|---------|-------------|
| **Sistema** | 16 | Configuración base del OS: pantalla, audio, red, bluetooth, energía, teclado, fecha/hora, usuarios, apariencia, defaults, autostart, servicios, updates, seguridad, fuentes, notificaciones |
| **Aplicaciones** | 7 | Herramientas de usuario: neovim, qtile, screenshot, grabación, terminal, file manager, AI |
| **Operaciones** | 5 | Mantenimiento: monitor, almacenamiento, logs, backup, dotfiles |
| **Setup** | 3 | Configuración inicial: wizard, instalador, perfiles |

### Principio de navegación

La navegación es un **árbol de 3 niveles con backtracking explícito**:

```
Categorías  →  Módulos  →  Detalle del módulo
   ↑              ↑              │
   └──── Esc ─────└──── Esc ─────┘
```

- Avanzar: seleccionar y confirmar (Enter / clic)
- Retroceder: Escape o botón "volver"
- Salir: desde el nivel raíz, Escape o botón "salir"

---

## 3. Layout y Composición Visual

### Pantalla principal — Vista de Categorías

```
┌──────────────────────────────────────────────────┐
│  🔧 LambdaOS Settings                            │
│                                                  │
│  ┌──────────────────────────────────────────────┐│
│  │ ▶ Sistema                         16 módulos ││
│  │   Aplicaciones                     7 módulos ││
│  │   Operaciones                      5 módulos ││
│  │   Setup                            3 módulos ││
│  └──────────────────────────────────────────────┘│
│                                                  │
│  ─────────────────────────────────────────────── │
│  Categorías                │ ? Ayuda  Q Salir    │
└──────────────────────────────────────────────────┘
```

**Zonas**:
1. **Header** (8-10% alto): Título de la aplicación con icono. Fijo en todas las vistas.
2. **Contenido** (82-86% alto): Área principal. Cambia según la vista activa.
3. **Barra de estado** (6-8% alto, fija abajo): Contexto actual, atajos disponibles, resultado de última acción.

### Vista de Módulos dentro de una Categoría

```
┌──────────────────────────────────────────────────┐
│  🔧 LambdaOS Settings                            │
│                                                  │
│  Sistema                                         │
│  ┌──────────────────────────────────────────────┐│
│  │ ▶ Pantalla                                    ││
│  │     Resolución, múltiples monitores, perfiles ││
│  │                                               ││
│  │   Audio                                       ││
│  │     Volumen, dispositivos, perfiles de sonido ││
│  │                                               ││
│  │   Red                                         ││
│  │     WiFi, VPN, información de conexión        ││
│  └──────────────────────────────────────────────┘│
│                                                  │
│  ─────────────────────────────────────────────── │
│  Sistema                   │ Esc Volver  ? Ayuda │
└──────────────────────────────────────────────────┘
```

Cada módulo se muestra como una **fila** con:
- Nombre del módulo (destacado)
- Descripción breve en una línea (texto secundario, más pequeño o atenuado)
- Indicador de selección (cursor, highlight, o borde izquierdo coloreado)

### Vista de Detalle de Módulo

```
┌──────────────────────────────────────────────────┐
│  🔧 LambdaOS Settings                            │
│                                                  │
│  Audio                                           │
│  ┌──────────────────────────────────────────────┐│
│  │                                               ││
│  │  Volumen         [═══════════●────]  75%      ││
│  │                                               ││
│  │  Silenciar        ● Activado                  ││
│  │                                               ││
│  │  Dispositivo     ▶ Auriculares USB            ││
│  │                                               ││
│  │  Perfil          ▶ Altavoces                  ││
│  │                                               ││
│  │  [ Cargar perfil Guardado ]                   ││
│  │                                               ││
│  └──────────────────────────────────────────────┘│
│                                                  │
│  ─────────────────────────────────────────────── │
│  Audio · Volumen aplicado  │ Esc Volver  ? Ayuda │
└──────────────────────────────────────────────────┘
```

Las acciones del módulo se presentan como **widgets en una lista vertical**, con:
- Etiqueta a la izquierda (nombre de la acción)
- Control interactivo a la derecha (toggle, slider, selector, campo de texto, botón)
- Separación visual clara entre acciones

---

## 4. Sistema de Diseño

### Paleta de Colores

La paleta debe funcionar tanto sobre fondo oscuro (default) como sobre fondo claro (tema alternativo futuro).

| Rol | Color | Uso |
|-----|-------|-----|
| **Primario / Acento** | Púrpura `#7D56F4` | Elemento seleccionado, cursor, indicador de foco, botón primario, borde activo |
| **Éxito / Activo** | Verde `#04B575` | Toggle encendido, confirmación, estado "conectado", barra de progreso completa |
| **Error / Destructivo** | Rojo `#FF4672` | Error, acción destructiva, validación fallida, indicador "modificado sin guardar" |
| **Advertencia** | Amarillo/Ámbar `#F4D03F` | Advertencia no bloqueante, estado "atención requerida", valor fuera de rango |
| **Texto primario** | Blanco `#FFFFFF` o Negro `#1A1A1A` | Texto principal, etiquetas, títulos (contraste máximo contra fondo) |
| **Texto secundario** | Gris medio `#888888` | Descripciones, hints, placeholder text, valores no seleccionados |
| **Superficie** | Gris oscuro `#1E1E2E` o Blanco `#FAFAFA` | Fondo de tarjetas, área de contenido si se usa panel |
| **Fondo** | Negro `#0D0D0D` o Blanco `#FFFFFF` | Fondo general de la aplicación |

**Reglas de contraste**:
- Texto primario sobre fondo: ratio ≥ 7:1 (AAA)
- Texto secundario sobre fondo: ratio ≥ 4.5:1 (AA)
- Acento sobre fondo oscuro: verificar ≥ 4.5:1 para texto sobre acento

### Tipografía

Al ser una aplicación de terminal, la tipografía está limitada a **monoespacio**. Sin embargo, para el prototipo en Figma:

- **Títulos y headers**: Monoespacio bold, 1-2px más grande que el cuerpo
- **Cuerpo y etiquetas**: Monoespacio regular
- **Texto secundario**: Monoespacio regular, mismo tamaño, color atenuado
- **Datos e inputs**: Monoespacio regular, posiblemente con subrayado o recuadro para indicar editabilidad

El tamaño base recomendado para prototipo: 14-16px.

### Espaciado y Proporciones

- **Padding interno de componentes**: 8-12px (1-2 celdas de terminal)
- **Separación entre acciones en detalle**: 4-8px (1 celda)
- **Margen del área de contenido**: 16-24px (2-3 celdas) desde los bordes
- **Header**: altura fija, ~40-48px
- **Barra de estado**: altura fija, ~28-32px
- **Ancho mínimo efectivo**: 640px (80 columnas) — por debajo, el contenido se trunca, no se reorganiza

### Iconografía

**Nerd Fonts como estándar visual.** Se usa [Nerd Fonts](https://www.nerdfonts.com/) (específicamente Monoid Nerd Font) como conjunto de iconos principal. Esto permite glifos ricos y semánticos que funcionan en cualquier terminal moderna. El paquete `nerd-fonts-monoid` es dependencia base de LambdaOS, instalado desde la ISO mínima.

Para el prototipo en Figma, usar iconos equivalentes del sistema de diseño.

| Concepto | Unicode básico (fallback) | Nerd Fonts (estándar) |
|----------|--------------------------|----------------------|
| Cursor / seleccionado | `▶` | `` (nf-cod-chevron_right) |
| Toggle encendido | `●` | `` (nf-fa-toggle_on) |
| Toggle apagado | `○` | `` (nf-fa-toggle_off) |
| Expandir / abrir | `▶` / `▼` | `` / `` (nf-cod-chevron) |
| Ayuda | `?` | `` (nf-fa-question_circle) |
| Cerrar / volver | `✕` | `` (nf-fa-times) |
| Advertencia | `⚠` | `` (nf-fa-warning) |
| Error / Crítico | `✕` | `` (nf-fa-exclamation_circle) |
| Éxito / Check | `✓` | `` (nf-fa-check) |
| Candado / root | `🔒` | `` (nf-fa-lock) |
| Cargando / Sync | `⟳` | `` (nf-fa-refresh) |
| Búsqueda | `/` | `` (nf-fa-search) |
| Configuración | `⚙` | `` (nf-fa-cog) |

**Módulos con iconos Nerd Fonts específicos:**

| Módulo | Unicode fallback | Nerd Fonts |
|--------|-----------------|------------|
| Pantalla | `▣` | `` (nf-fa-desktop) o `` (nf-fa-laptop) |
| Audio | `♪` | `` (nf-fa-volume_up) |
| Red | `◉` | `` (nf-fa-globe) o `` (nf-fa-wifi) |
| Bluetooth | `⊛` | `` (nf-fa-bluetooth) |
| Energía | `⚡` | `` (nf-fa-bolt) o `` (nf-fa-battery_full) |
| Teclado | `⌨` | `` (nf-fa-keyboard_o) |
| Usuarios | `⊙` | `` (nf-fa-user) |
| Apariencia | `◈` | `` (nf-fa-paint_brush) |
| Seguridad | `⛨` | `` (nf-fa-user_secret) o `` (nf-fa-shield) |
| Actualizaciones | `↻` | `` (nf-fa-refresh) |
| Servicios | `◎` | `` (nf-fa-cogs) |
| Fuentes | `A` | `` (nf-fa-font) |
| Notificaciones | `◻` | `` (nf-fa-bell) |
| Neovim | `✦` | `` (nf-dev-vim) |
| Qtile | `⊞` | `` (nf-fa-window_maximize) |
| Dotfiles | `◌` | `` (nf-fa-files_o) o `` (nf-fa-git) |
| Logs | `≡` | `` (nf-fa-file_text_o) |
| Almacenamiento | `◫` | `` (nf-fa-hdd_o) |
| Monitor | `▦` | `` (nf-fa-bar_chart) |

**Principio de fallback**: Si Nerd Fonts no está disponible, cada icono DEBE degradar a su equivalente Unicode básico. El sistema detecta la presencia de Nerd Fonts al inicio y selecciona el set de iconos correspondiente.

### Bordes y Separadores

- **Separador horizontal**: Línea fina (1px), color secundario atenuado, usado entre header y contenido, y entre contenido y status bar
- **Tarjetas/paneles**: Sin bordes (flat design) o borde sutil de 1px en color superficie. No usar bordes dobles ┌┐└┘.
- **Diálogos modales**: Borde redondeado (4-6px radio), sombra o overlay detrás para indicar modalidad

---

## 5. Catálogo de Componentes Visuales

Cada acción de un módulo se materializa como uno de estos componentes. El diseñador debe tratar cada uno como un componente reutilizable en Figma.

### 5.1 Toggle (Interruptor)

**Propósito**: Activar/desactivar una opción booleana.

**Visual**:
```
┌──────────────────────────────────────────────────┐
│  Silenciar                           ● Activado  │
└──────────────────────────────────────────────────┘
```

- Etiqueta a la izquierda
- Indicador de estado a la derecha: círculo relleno (verde) = on, círculo vacío (gris) = off
- Texto de estado junto al indicador: "Activado" / "Desactivado"

**Interacción**: Clic o Enter sobre la fila cambia el estado. Feedback inmediato: el círculo se rellena/vacía y el color cambia.

### 5.2 Slider (Deslizador)

**Propósito**: Ajustar un valor numérico en un rango (ej: volumen 0-100%).

**Visual**:
```
┌──────────────────────────────────────────────────┐
│  Volumen         [═══════════●───────]     75%   │
└──────────────────────────────────────────────────┘
```

- Etiqueta a la izquierda
- Barra de progreso con indicador de posición (●)
- Valor numérico a la derecha
- Porción "rellena" de la barra en color primario, porción "vacía" en color secundario/gris

**Interacción**: Flechas izquierda/derecha o clic en la barra para ajustar. El valor se actualiza en tiempo real mientras se arrastra.

### 5.3 Selector (Dropdown / Lista de Opciones)

**Propósito**: Elegir una opción entre varias predefinidas.

**Visual (estado colapsado)**:
```
┌──────────────────────────────────────────────────┐
│  Dispositivo de salida        ▶ Auriculares USB  │
└──────────────────────────────────────────────────┘
```

**Visual (estado expandido)**:
```
┌──────────────────────────────────────────────────┐
│  Dispositivo de salida        ▼ Auriculares USB  │
│                               ┌────────────────┐ │
│                               │ Altavoces       │ │
│                               │ Auriculares USB │ │
│                               │ Salida HDMI     │ │
│                               └────────────────┘ │
└──────────────────────────────────────────────────┘
```

- Etiqueta a la izquierda
- Valor seleccionado a la derecha
- Indicador de expansión (▶ colapsado, ▼ expandido)
- Lista desplegable: fondo de superficie, opción seleccionada resaltada con color primario, opciones no seleccionadas con texto secundario

**Interacción**: Clic o Enter expande. Flechas arriba/abajo o clic seleccionan. Enter o clic en opción confirma y colapsa. Escape colapsa sin cambiar.

### 5.4 Campo de Texto (Input)

**Propósito**: Ingresar texto libre o valores numéricos con validación.

**Visual (estado normal)**:
```
┌──────────────────────────────────────────────────┐
│  Nombre de red WiFi            │                │
└──────────────────────────────────────────────────┘
```

**Visual (estado enfocado/editando)**:
```
┌──────────────────────────────────────────────────┐
│  Nombre de red WiFi            │ MiRed_5GHz│     │
└──────────────────────────────────────────────────┘
```

**Visual (estado con error de validación)**:
```
┌──────────────────────────────────────────────────┐
│  Puerto                        │ 99999 │ ⚠ Fuera │
│                                  de rango (1-65535)│
└──────────────────────────────────────────────────┘
```

- Etiqueta a la izquierda
- Recuadro de input con borde sutil
- Placeholder en texto secundario cuando está vacío
- Error de validación en rojo, debajo o a la derecha del campo
- Cursor parpadeante dentro del campo cuando está enfocado

**Interacción**: Enter o clic enfoca el campo. Escape desenfoca (cancela edición). Enter con campo enfocado confirma el valor. Validación al confirmar, no en cada tecla.

### 5.5 Botón de Acción (Execute)

**Propósito**: Disparar una acción puntual (escanear, conectar, instalar, aplicar).

**Visual**:
```
┌──────────────────────────────────────────────────┐
│  [  Buscar redes WiFi  ]                         │
│  [  Conectar  ]                                  │
│  [  Guardar perfil  ]                            │
└──────────────────────────────────────────────────┘
```

- Recuadro con padding horizontal generoso (16-24px)
- Fondo primario (púrpura) para acción principal, fondo superficie para acción secundaria
- Texto centrado en blanco sobre fondo primario
- Estado hover: ligero brillo o borde más claro
- Estado presionado: escala 97% o color más oscuro

**Interacción**: Clic o Enter dispara la acción. Durante la ejecución, el botón muestra estado "cargando" (texto "Ejecutando..." con elipsis animada). Al completar, feedback en barra de estado.

### 5.6 Diálogo de Confirmación (Modal)

**Propósito**: Pedir confirmación antes de una acción destructiva o irreversible.

**Visual**:
```
┌──────────────────────────────────────────────────┐
│  🔧 LambdaOS Settings                            │
│                                                  │
│  Audio                                           │
│  ┌──────────────────────────────────────────────┐│
│  │  overlay semitransparente oscuro             ││
│  │                                               ││
│  │  ┌──────────────────────────────────┐        ││
│  │  │  ⚠ ¿Olvidar dispositivo?        │        ││
│  │  │                                  │        ││
│  │  │  Se eliminará "Auriculares USB"  │        ││
│  │  │  de los dispositivos conocidos.  │        ││
│  │  │                                  │        ││
│  │  │  [  Sí, olvidar  ]   Cancelar   │        ││
│  │  └──────────────────────────────────┘        ││
│  │                                               ││
│  └──────────────────────────────────────────────┘│
│                                                  │
│  ─────────────────────────────────────────────── │
│  Audio · Diálogo de confirmación                 │
└──────────────────────────────────────────────────┘
```

- Overlay semitransparente (negro 60% opacidad) sobre el contenido
- Tarjeta modal centrada, fondo superficie, bordes redondeados (6-8px)
- Título con icono de advertencia
- Descripción breve de la acción
- Botón primario (acción destructiva en rojo si aplica) + botón secundario (cancelar)
- Navegación con Tab entre botones, Escape = cancelar

### 5.7 Indicador de Progreso

**Propósito**: Mostrar avance de una operación larga (instalación, escaneo, actualización).

**Visual**:
```
┌──────────────────────────────────────────────────┐
│  Actualizando paquetes...                        │
│  [══════════════●───────────]  67%               │
│  Instalando: firefox 121.0.1-2                   │
└──────────────────────────────────────────────────┘
```

- Barra de progreso con indicador
- Porcentaje a la derecha
- Texto descriptivo debajo (qué se está haciendo ahora)
- Color: primario para la barra de progreso, verde para completado

---

## 6. Flujos de Usuario por Módulo

Cada módulo sigue el mismo patrón estructural, pero con particularidades visuales según su dominio.

### 6.1 Pantalla (Display)

**Flujo principal**: Ver outputs → Seleccionar output → Elegir resolución/modo → Aplicar → Confirmar en cuenta regresiva

**Particularidad visual**:
- Lista de outputs como **tarjetas**, cada una con: nombre del output, resolución actual, estado (conectado/desconectado), posición relativa
- Output desconectado: tarjeta atenuada (opacidad 50%)
- Al aplicar cambio de resolución: **cuenta regresiva de 10 segundos** en la barra de estado, con botón "Confirmar" y "Revertir". Si no se confirma, revierte automáticamente.

### 6.2 Audio

**Flujo principal**: Ver volumen → Ajustar slider → Cambiar dispositivo → (opcional) Guardar perfil

**Particularidad visual**:
- Slider de volumen como elemento principal destacado
- Selectores de dispositivo de entrada y salida
- Sección de "Perfiles guardados" como lista de tarjetas con nombre y miniatura de configuración
- Per-app volume (si disponible): lista de aplicaciones con mini-sliders individuales, colapsable

### 6.3 Red (Network)

**Flujo principal**: Ver estado → (si WiFi) Escanear → Seleccionar red → Ingresar contraseña → Conectar

**Particularidad visual**:
- Panel de información actual: IP, gateway, DNS, tipo de conexión (como tarjeta de resumen, no interactiva)
- Lista de redes WiFi: cada fila muestra SSID, intensidad de señal (barras), tipo de seguridad (icono candado), estado (conectado/guardado)
- Red conectada: resaltada en verde con indicador ●
- Flujo de conexión: seleccionar red → si tiene contraseña, el campo de texto aparece inline (no en popup) debajo de la red seleccionada → botón "Conectar"

### 6.4 Bluetooth

**Flujo principal**: Activar BT → Escanear dispositivos → Seleccionar → Emparejar/Conectar

**Particularidad visual**:
- Toggle maestro: Bluetooth ON/OFF como primer elemento, destacado
- Lista de dispositivos: nombre, tipo (icono: 🎧 audio, ⌨️ input, 📱 otro), estado (conectado/emparejado/no emparejado)
- Dispositivo conectado: fila en verde con ●
- Acciones contextuales: botones "Conectar", "Desconectar", "Olvidar" aparecen al seleccionar un dispositivo (no siempre visibles)

### 6.5 Energía (Power)

**Flujo principal**: Ver batería → Configurar timeouts → Elegir acción al cerrar tapa

**Particularidad visual**:
- Indicador de batería: barra horizontal con porcentaje e ícono de batería, verde/amarillo/rojo según nivel
- Campos numéricos para timeouts (suspensión, apagar pantalla) con unidad de medida visible ("minutos")
- Selector de "Acción al cerrar tapa": suspender / hibernar / nada / apagar pantalla

### 6.6 Teclado (Keyboard)

**Flujo principal**: Seleccionar layout → (opcional) Elegir variante → Cambio inmediato

**Particularidad visual**:
- Selector de layout principal con indicador de layout activo
- Vista previa del layout: representación visual del teclado con las teclas (simplificado, puede ser una tabla de referencia)
- Selector de variante aparece/desaparece según el layout seleccionado (dinámico)
- Sin botón "aplicar": el cambio es inmediato al seleccionar

### 6.7 Usuarios (Users)

**Flujo principal**: Ver lista de usuarios → Crear/Editar/Eliminar usuario

**Particularidad visual**:
- Lista de usuarios como tarjetas con: nombre, grupos, última sesión, estado
- Flujo de creación: botón "+" abre formulario inline con campos: nombre, contraseña, grupos (multi-select)
- Requiere autenticación root: indicador 🔒 visible en la barra de estado
- Eliminar usuario: diálogo de confirmación con advertencia roja

### 6.8 Apariencia (Appearance)

**Flujo principal**: Seleccionar tema → Elegir wallpaper → Ajustar tamaño de fuente

**Particularidad visual**:
- Selector de temas como **tarjetas con preview**: miniatura del tema (color primario + secundario), nombre, descripción
- Selector de wallpaper: grid de miniaturas o lista con nombres + preview pequeño
- Slider de tamaño de fuente con preview en vivo de texto de ejemplo

### 6.9 Defaults (Aplicaciones por defecto)

**Flujo principal**: Por cada categoría (navegador, editor, terminal, etc.) → Seleccionar aplicación

**Particularidad visual**:
- Lista agrupada por tipo de archivo/acción: "Navegador web", "Editor de texto", "Terminal", "Reproductor de video", etc.
- Cada fila: tipo de default + aplicación actual + selector para cambiar
- Diseño de formulario: etiquetas alineadas, consistente

### 6.10 Resto de módulos

| Módulo | Widgets principales | Particularidad visual |
|--------|-------------------|----------------------|
| **Fecha/Hora** | Selector de zona horaria con search, toggle NTP | Mapa de zonas horarias o lista con buscado. Preview de hora actual |
| **Autostart** | Lista de servicios con toggles | Similar a servicios pero solo user units. Toggle simple |
| **Servicios** | Lista con indicador de estado + start/stop/restart | Estado: ● verde (running), ○ gris (stopped), ⚠ amarillo (failed). Botones contextuales |
| **Actualizaciones** | Botón "Buscar" → Lista de paquetes → Botón "Actualizar" | Barra de progreso durante actualización. Lista con checkboxes para seleccionar paquetes |
| **Seguridad** | Toggle firewall, lista de reglas, gestión de claves SSH/GPG | Secciones colapsables: Firewall, SSH, GPG. Formulario para generar claves |
| **Fuentes** | Lista de fuentes con preview, botón instalar | Preview de cada fuente: "El veloz murciélago hindú..." en la fuente misma |
| **Notificaciones** | Sliders de timeout, selectores de posición, toggles por app | Representación visual de la posición de notificación en pantalla |

---

## 7. Estados y Feedback

### 7.1 Barra de Estado

Elemento persistente en la parte inferior. Cambia según el contexto:

**Formato general**:
```
[Contexto] · [Mensaje de feedback]  │  [Atajos]
```

**Estados posibles**:

| Estado | Color | Ejemplo |
|--------|-------|---------|
| Normal / Contexto | Primario (púrpura) | `Sistema · 16 módulos` |
| Éxito | Verde | `Audio · Volumen aplicado: 75%` |
| Error | Rojo | `Red · Error: No se pudo conectar a "MiWiFi"` |
| Advertencia | Amarillo | `Audio · Per-app volume requiere PipeWire` |
| Cargando | Primario con elipsis | `Red · Escaneando redes...` |
| Requiere root | Rojo + 🔒 | `Usuarios · Se requiere autenticación de administrador` |
| Modificado (sin guardar) | Rojo + `*` | `Teclado · Cambios sin guardar *` |

### 7.2 Estados de Componentes

**Toggle**:
- Normal: indicador y texto de estado visibles
- Hover/foco: ligero brillo o borde en la fila
- Presionado: transición visual del indicador (animación de llenado/vaciado)
- Deshabilitado: opacidad reducida (50%), no interactivo

**Slider**:
- Normal: barra + indicador + valor
- Arrastrando/focus: indicador ligeramente más grande y brillante
- Deshabilitado: barra gris uniforme, sin indicador

**Selector**:
- Colapsado: indicador ▶ + valor actual
- Expandido: indicador ▼ + lista desplegable con highlight en seleccionado
- Cargando opciones dinámicas: texto "Cargando opciones..." en la lista

**Campo de texto**:
- Normal: borde sutil, placeholder visible si está vacío
- Enfocado: borde en color primario, cursor visible
- Error: borde en rojo, mensaje de error debajo
- Validación exitosa: borde en verde por un instante, luego vuelve a normal

**Botón**:
- Normal: fondo primario, texto blanco
- Hover: fondo ligeramente más claro
- Presionado: escala 97%
- Cargando: texto "Ejecutando..." (elipsis animada)
- Completado: transformación breve a verde con checkmark, luego vuelve
- Deshabilitado: gris, sin sombra, no interactivo

### 7.3 Indicador "Modificado"

Cuando un módulo tiene cambios no aplicados/guardados:
- Asterisco `*` rojo junto al nombre del módulo en la barra de estado
- En la vista de módulos, un indicador `*` junto al nombre del módulo modificado
- Al intentar salir sin guardar: diálogo de confirmación "Hay cambios sin guardar. ¿Salir de todas formas?"

---

## 8. Principios de Interacción

### Navegación

- **Jerarquía plana**: máximo 3 niveles. El usuario nunca se pierde.
- **Backtracking universal**: Escape siempre vuelve atrás. En el nivel raíz, Escape pregunta si salir.
- **Estado preservado**: volver a un módulo muestra la misma posición de scroll y selección que cuando se salió.
- **Sin modales anidados**: solo un diálogo de confirmación a la vez.

### Teclado y Mouse

El diseño debe contemplar ambos:

| Acción | Teclado | Mouse |
|--------|---------|-------|
| Navegar entre items | ↑ ↓ o Tab | Clic directo |
| Seleccionar / Activar | Enter | Clic |
| Volver atrás | Escape | Botón "Volver" en UI |
| Toggle | Space o Enter | Clic en la fila |
| Cambiar valor slider | ← → | Arrastrar indicador |
| Abrir selector | Enter o Space | Clic en el selector |
| Ayuda | ? | Botón "?" en barra de estado |
| Buscar (futuro) | / | Clic en campo de búsqueda |
| Salir | q | Botón "Salir" |

### Seguridad

- **Acciones destructivas**: siempre requieren confirmación (capa modal)
- **Cambios de display**: cuenta regresiva de 10s con auto-revert si no se confirma (previene pantalla negra)
- **Operaciones como root**: indicador visual claro (🔒 rojo) y posiblemente diálogo de contraseña
- **Salir con cambios**: advertencia antes de descartar cambios no guardados

### Accesibilidad

- **Contraste**: todos los textos cumplen WCAG AA (4.5:1) como mínimo
- **No solo color**: los estados usan texto + color + posición (el toggle dice "Activado"/"Desactivado", no solo cambia de color)
- **Foco visible**: siempre hay un indicador claro de qué elemento tiene el foco
- **Navegación secuencial**: el orden de Tab es lógico (de arriba a abajo, de izquierda a derecha)
- **Sin animaciones bloqueantes**: el usuario puede seguir interactuando durante transiciones

---

## 9. Apéndice: Resumen Visual por Módulo

| # | Módulo | Widget principal | Widgets secundarios | Root | Complejidad visual |
|---|--------|-----------------|-------------------|------|-------------------|
| 01 | Pantalla | Selector de outputs (tarjetas) | Selector modo, Toggle primario, Botón guardar/load | No | Alta (múltiples outputs, perfiles) |
| 02 | Audio | Slider de volumen | Toggle mute, Selector sink/source, Selector perfil, Mini-sliders por app | No | Alta (muchos controles) |
| 03 | Red | Lista de redes WiFi | Toggle WiFi, Campo texto password, Panel de info | No | Media (flujo de conexión) |
| 04 | Bluetooth | Lista de dispositivos | Toggle BT, Botones contextuales | No | Media |
| 05 | Energía | Indicador de batería | Campos numéricos timeout, Selector acción tapa | Parcial | Baja |
| 06 | Teclado | Selector de layout | Selector de variante, Preview visual | No | Baja |
| 07 | Usuarios | Lista de usuarios (tarjetas) | Formulario crear/editar, Diálogo eliminar | Sí | Alta (formularios) |
| 08 | Fecha/Hora | Selector de zona horaria | Toggle NTP | Parcial | Media (buscador en selector) |
| 09 | Apariencia | Grid de temas (tarjetas con preview) | Selector wallpaper, Slider fuente | No | Alta (visual, previews) |
| 10 | Defaults | Lista de selectores agrupados | — | No | Baja |
| 11 | Autostart | Lista de servicios (toggles) | — | No | Baja |
| 12 | Servicios | Lista con estado | Botones start/stop/restart contextuales | Parcial | Media |
| 13 | Actualizaciones | Lista de paquetes (checkboxes) | Botón buscar, Botón actualizar, Barra progreso | Parcial | Alta (durante actualización) |
| 14 | Seguridad | Secciones colapsables | Toggle firewall, Lista reglas, Formulario claves | Parcial | Alta (múltiples secciones) |
| 15 | Fuentes | Lista con preview tipográfico | Botón instalar | No | Media (previews) |
| 16 | Notificaciones | Sliders + Selectores | Toggles por app, Representación visual posición | No | Media |

---

*Documento de especificación conceptual de UI/UX para prototipado. No describe implementación — describe experiencia.*
