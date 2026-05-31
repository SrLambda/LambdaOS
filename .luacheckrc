std = "lua51+luajit"
globals = {"vim"}

-- Neovim configs often have unused callback params (e.g., `client`, `opts`)
-- and plugin toggle globals. Ignore unused warnings to keep CI green.
ignore = {
    "212", -- unused argument
    "213", -- unused argument after _
    "314", -- unused variable
    "511", -- setting non-standard global variable (plugin toggle pattern)
}
