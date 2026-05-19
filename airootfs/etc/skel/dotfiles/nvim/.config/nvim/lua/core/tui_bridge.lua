local M = {}

M.defaults = {
  enable_lsp = true,
  enable_copilot = false,
  enable_neotree = true,
}

function M.parse_tui_settings()
  local config_path = vim.fn.expand("~/.config/nvim/tui_settings.json")
  local flags = vim.deepcopy(M.defaults)

  local file = io.open(config_path, "r")
  if not file then
    vim.notify("[tui_bridge] tui_settings.json not found, using defaults", vim.log.levels.WARN)
    return flags
  end

  local content = file:read("*a")
  file:close()

  local ok, decoded = pcall(vim.json.decode, content)
  if not ok then
    vim.notify("[tui_bridge] Failed to parse tui_settings.json: " .. tostring(decoded), vim.log.levels.ERROR)
    return flags
  end

  for k, v in pairs(decoded) do
    if flags[k] ~= nil then
      flags[k] = v
    end
  end

  return flags
end

function M.get_flags()
  return M.parse_tui_settings()
end

return M
