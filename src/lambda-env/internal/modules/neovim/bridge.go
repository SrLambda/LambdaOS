package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const tuiBridgeTemplate = `local M = {}

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
`

func UpdateTuiBridgeLua(nvimConfigPath string) error {
	if nvimConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home dir: %w", err)
		}
		nvimConfigPath = filepath.Join(home, ".config", "nvim")
	}

	bridgePath := filepath.Join(nvimConfigPath, "lua", "core", "tui_bridge.lua")
	backupPath := bridgePath + ".bak"

	if _, err := os.Stat(bridgePath); err == nil {
		data, readErr := os.ReadFile(bridgePath)
		if readErr != nil {
			return fmt.Errorf("read existing tui_bridge.lua: %w", readErr)
		}
		if writeErr := os.WriteFile(backupPath, data, 0644); writeErr != nil {
			return fmt.Errorf("backup tui_bridge.lua: %w", writeErr)
		}
	}

	bridgeDir := filepath.Dir(bridgePath)
	if err := os.MkdirAll(bridgeDir, 0755); err != nil {
		return fmt.Errorf("create tui_bridge.lua directory: %w", err)
	}

	if err := os.WriteFile(bridgePath, []byte(tuiBridgeTemplate), 0644); err != nil {
		return fmt.Errorf("write tui_bridge.lua: %w", err)
	}

	return nil
}
