#!/usr/bin/env bash
set -euo pipefail

# dev-tui.sh — Build and run lambda-env TUI with all modules in dev mode
# Usage: ./scripts/dev-tui.sh [--build-only]

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LAMBDA_ENV="$REPO_ROOT/src/lambda-env"
MODULES_DIR="$HOME/.local/share/lambda-env/modules"
SETTINGS_FILE="$HOME/.config/lambdaos/settings.json"

echo "=== lambda-env dev TUI ==="
echo ""

# --- Build lambda-env binary ---
echo "[1/3] Building lambda-env..."
cd "$LAMBDA_ENV"
go build -o lambda-env ./cmd/lambda-env
echo "      OK ($(du -h lambda-env | cut -f1))"

# --- Build all 7 module binaries ---
echo "[2/3] Building modules..."
MODULES=(neovim qtile dotfiles keyboard appearance audio defaults)

for mod in "${MODULES[@]}"; do
    mkdir -p "$MODULES_DIR/$mod"
    cd "$LAMBDA_ENV/internal/modules/$mod"
    go build -o module . 2> /dev/null || {
        echo "      WARNING: $mod module build failed (may need dependencies)"
        continue
    }
    cp module "$MODULES_DIR/$mod/module"
    cp manifest.json "$MODULES_DIR/$mod/manifest.json"
    echo "      $mod OK"
done

# --- Create settings.json if missing ---
echo "[3/3] Settings..."
if [ -f "$SETTINGS_FILE" ]; then
    echo "      Using existing: $SETTINGS_FILE"
else
    mkdir -p "$(dirname "$SETTINGS_FILE")"
    cat > "$SETTINGS_FILE" << 'EOF'
{
  "version": "1.0.0",
  "appearance": {"theme": "dark", "font_size": 14, "opacity": 100, "wallpaper": ""},
  "display": {"active_profile": "default", "profiles": []},
  "audio": {"default_sink": "", "volume": 75, "muted": false},
  "network": {"wifi_enabled": true, "known_networks": []},
  "bluetooth": {"enabled": true, "paired_devices": []},
  "keyboard": {"layout": "us", "variant": "", "options": ""},
  "neovim": {
    "theme": "tokyonight", "font": "JetBrainsMono",
    "lines": 40, "columns": 120,
    "enable_lsp": true, "enable_copilot": true, "enable_neotree": true,
    "lsp_servers": ["gopls", "pyright"],
    "use_global_theme": true
  },
  "qtile": {
    "bar_position": "top", "bar_size": 24,
    "layouts": [], "terminal": "kitty", "browser": "firefox",
    "default_file_manager": "thunar",
    "groups": [{"name":"1"},{"name":"2"},{"name":"3"}],
    "color_scheme": "dracula",
    "use_global_theme": true
  },
  "services": {"enabled": []}
}
EOF
    echo "      Created: $SETTINGS_FILE"
fi

echo ""
echo "================================================"
echo "  lambda-env ready!"
echo "  Keys: ↑↓/kj navigate | Enter select | Esc back"
echo "  ? help overlay | q quit"
echo "================================================"

if [ "${1:-}" = "--build-only" ]; then
    echo ""
    echo "Build complete. To launch, run:"
    echo "  $LAMBDA_ENV/lambda-env"
    exit 0
fi

echo ""
exec "$LAMBDA_ENV/lambda-env"
