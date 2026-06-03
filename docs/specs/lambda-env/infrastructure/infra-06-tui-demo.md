# lambda-env: TUI Interactive Demo (GitHub Pages)

## Intent / Intención

**EN**: Create an interactive public demo of the LambdaOS TUI that runs in a browser and is deployed to GitHub Pages. This demo combines all TUI view prototypes into a navigable, interactive experience that showcases the full configuration hub.

**ES**: Crear una demo pública interactiva de la TUI de LambdaOS que corre en un browser y se despliega en GitHub Pages. Esta demo combina todos los prototipos de vistas TUI en una experiencia navegable e interactiva que muestra el hub de configuración completo.

## Requirements / Requisitos

1. **EN**: The demo MUST be deployed to GitHub Pages automatically on vX.0 tag push.
   **ES**: La demo DEBE desplegarse en GitHub Pages automáticamente al hacer push de un tag vX.0.

2. **EN**: The demo MUST combine all prototypes from `src/lambda-env/prototypes/` into a single navigable application.
   **ES**: La demo DEBE combinar todos los prototipos de `src/lambda-env/prototypes/` en una sola aplicación navegable.

3. **EN**: The demo MUST use a lightweight framework (Alpine.js, HTMX, or vanilla JS with router) — no heavy SPA frameworks.
   **ES**: La demo DEBE usar un framework liviano (Alpine.js, HTMX, o vanilla JS con router) — sin frameworks SPA pesados.

4. **EN**: The demo MUST preserve the terminal aesthetic — monospace fonts, dark background, terminal-like colors and layout.
   **ES**: La demo DEBE preservar la estética terminal — fuentes monospace, fondo oscuro, colores y layout tipo terminal.

5. **EN**: Users MUST be able to navigate between modules, open settings, and simulate configuration changes.
   **ES**: Los usuarios DEBEN poder navegar entre módulos, abrir configuraciones, y simular cambios de configuración.

6. **EN**: The demo MUST be a single `index.html` entry point with all assets using relative paths (GitHub Pages compatible).
   **ES**: La demo DEBE ser un único punto de entrada `index.html` con todos los assets usando rutas relativas (compatible con GitHub Pages).

## Technical Notes / Notas Técnicas

**EN**: The demo reuses the layout and CSS from the Phase 1 prototypes (`src/lambda-env/prototypes/`). Instead of rebuilding from scratch, the demo adds interactivity (navigation, state simulation) on top of the existing prototype files. This ensures the demo always matches what the TUI actually looks like.

**ES**: La demo reutiliza el layout y CSS de los prototipos de Fase 1 (`src/lambda-env/prototypes/`). En lugar de reconstruir desde cero, la demo agrega interactividad (navegación, simulación de estado) sobre los archivos de prototipo existentes. Esto asegura que la demo siempre coincida con lo que la TUI realmente se ve.

### Directory Structure

```
src/lambda-env/demo/
├── index.html              ← Entry point + navigation shell
├── css/
│   └── terminal.css        ← Shared terminal aesthetic (from prototypes)
├── js/
│   ├── router.js           ← Simple view router
│   └── state.js            ← Simulated settings state
└── views/                  ← Copied/adapted from prototypes/
    ├── hub-menu.html
    ├── system-screen.html
    ├── audio-config.html
    └── ...
```

### CI/CD Integration

The demo deploys via a separate GitHub Actions job in `cd.yml`:

```yaml
deploy-demo:
  if: startsWith(github.ref, 'refs/tags/v')
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Deploy to GitHub Pages
      uses: peaceiris/actions-gh-pages@v4
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./src/lambda-env/demo
```
