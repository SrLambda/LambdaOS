# Ejecutar GNU Stow para desplegar dotfiles modulares
if [ -d "$HOME/dotfiles" ]; then
    cd "$HOME/dotfiles" && stow */ 2>/dev/null
fi
