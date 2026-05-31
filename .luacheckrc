std = "lua51+luajit"
globals = {"vim", "_lazygit_toggle"}

-- Neovim configs often have unused callback params (e.g., client, opts)
-- Suppress unused warnings to keep CI green.
unused = false
