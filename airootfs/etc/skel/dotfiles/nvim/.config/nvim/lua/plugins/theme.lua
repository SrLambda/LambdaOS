return {
  {
    "catppuccin/nvim",
    name = "catppuccin",
    lazy = false,
    priority = 1000,
    config = function()
      local theme = vim.g.nvim_theme or "catppuccin"
      require("catppuccin").setup({
        flavour = "mocha",
        transparent_background = false,
        term_colors = true,
        integrations = {
          alpha = true,
          bufferline = true,
          gitsigns = true,
          indent_blankline = { enabled = true },
          lsp_trouble = true,
          mason = true,
          navic = { enabled = true },
          neotree = true,
          noice = true,
          notify = true,
          telescope = { enabled = true },
          treesitter = true,
          which_key = true,
        },
      })
      vim.cmd.colorscheme(theme)
    end,
  },
}
