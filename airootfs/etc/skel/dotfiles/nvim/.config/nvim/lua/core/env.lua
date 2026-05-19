local tui_bridge = require("core.tui_bridge")

local theme = nil

local config_base = os.getenv("OS_CONFIG_DIR")
if config_base and config_base ~= "" then
  local theme_path = config_base .. "/os_theme.json"
  local f = io.open(theme_path, "r")
  if f then
    local content = f:read("*a")
    f:close()
    local ok, decoded = pcall(vim.json.decode, content)
    if ok and decoded and decoded.theme then
      theme = decoded.theme
    end
  end
end

if not theme or theme == "" then
  theme = os.getenv("NVIM_THEME")
end
if not theme or theme == "" then
  theme = "catppuccin"
end

vim.g.nvim_theme = theme
vim.g.tui_flags = tui_bridge.get_flags()
