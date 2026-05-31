std = "lua51+luajit"
globals = {"vim"}

-- Neovim configs often have unused callback params (e.g., client, opts)
-- and plugin toggle globals. Suppress unused warnings.
unused = false
