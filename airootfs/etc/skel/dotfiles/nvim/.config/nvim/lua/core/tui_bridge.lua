local M = {}

M.defaults = {
  enable_lsp = true,
  enable_copilot = true,
  enable_neotree = true,
}

function M.parse_settings()
  local config_path = vim.fn.expand("~/.config/lambdaos/settings.json")
  local flags = vim.deepcopy(M.defaults)

  local file = io.open(config_path, "r")
  if not file then
    vim.notify("[tui_bridge] settings.json not found, using defaults", vim.log.levels.WARN)
    return flags
  end

  local content = file:read("*a")
  file:close()

  local ok, decoded = pcall(vim.json.decode, content)
  if not ok then
    vim.notify("[tui_bridge] Failed to parse settings.json: " .. tostring(decoded), vim.log.levels.ERROR)
    return flags
  end

  -- Extract neovim section
  if decoded.neovim then
    if decoded.neovim.enable_lsp ~= nil then flags.enable_lsp = decoded.neovim.enable_lsp end
    if decoded.neovim.enable_copilot ~= nil then flags.enable_copilot = decoded.neovim.enable_copilot end
    if decoded.neovim.enable_neotree ~= nil then flags.enable_neotree = decoded.neovim.enable_neotree end
  end

  return flags
end

function M.get_flags()
  return M.parse_settings()
end

return M
