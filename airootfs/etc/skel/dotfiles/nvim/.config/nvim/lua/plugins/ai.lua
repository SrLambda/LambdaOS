local tui_flags = vim.g.tui_flags or {}
local enable_copilot = tui_flags.enable_copilot == true

return {
  {
    "zbirenbaum/copilot.lua",
    cmd = "Copilot",
    event = "InsertEnter",
    enabled = enable_copilot,
    config = function()
      require("copilot").setup({
        suggestion = {
          enabled = true,
          auto_trigger = true,
          keymap = {
            accept = "<Tab>",
            accept_word = "<C-Right>",
            accept_line = "<C-Down>",
            next = "<M-]>",
            prev = "<M-[>",
            dismiss = "<C-e>",
          },
        },
        panel = { enabled = false },
        filetypes = {
          markdown = true,
          help = false,
        },
      })
    end,
  },

  {
    "CopilotC-Nvim/CopilotChat.nvim",
    cmd = { "CopilotChat" },
    event = "VeryLazy",
    enabled = enable_copilot,
    dependencies = {
      { "zbirenbaum/copilot.lua" },
      { "nvim-lua/plenary.nvim" },
    },
    config = function()
      require("CopilotChat").setup({
        debug = false,
      })
    end,
  },
}
