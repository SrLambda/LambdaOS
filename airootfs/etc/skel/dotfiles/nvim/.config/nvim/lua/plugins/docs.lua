vim.opt.spelllang = "es,en"
vim.opt.spell = false

local spell_group = vim.api.nvim_create_augroup("docs_spell", { clear = true })
vim.api.nvim_create_autocmd("FileType", {
  group = spell_group,
  pattern = { "markdown", "tex", "latex", "text", "gitcommit", "rst" },
  callback = function()
    vim.opt_local.spell = true
  end,
})

return {
  {
    "iamcco/markdown-preview.nvim",
    cmd = { "MarkdownPreviewToggle", "MarkdownPreview", "MarkdownPreviewStop" },
    ft = { "markdown" },
    build = function()
      vim.fn["mkdp#util#install"]()
    end,
    config = function()
      vim.g.mkdp_auto_close = 0
      vim.g.mkdp_refresh_slow = 0
      vim.g.mkdp_open_to_the_world = 0
      vim.g.mkdp_browser = ""
      vim.g.mkdp_echo_preview_url = 1
      vim.g.mkdp_images_path = vim.fn.expand("~/.config/nvim/.markdown_images")
    end,
  },

  {
    "lervag/vimtex",
    ft = { "tex", "latex" },
    lazy = true,
    init = function()
      vim.g.vimtex_view_method = "zathura"
      vim.g.vimtex_compiler_method = "latexmk"
      vim.g.tex_flavor = "latex"
    end,
  },

  {
    "jbyuki/nabla.nvim",
    ft = { "tex", "latex", "markdown" },
    event = "VeryLazy",
    config = function()
      vim.keymap.set("n", "<leader>lp", "<cmd>lua require('nabla').popup()<CR>", { desc = "Preview LaTeX eq" })
    end,
  },

  {
    "MeanderingProgrammer/render-markdown.nvim",
    ft = { "markdown" },
    dependencies = { "nvim-treesitter/nvim-treesitter", "nvim-tree/nvim-web-devicons" },
    config = function()
      require("render-markdown").setup({})
    end,
  },

  {
    "aklt/plantuml-syntax",
    ft = { "plantuml" },
  },

  {
    "scrooloose/vim-slumlord",
    ft = { "plantuml" },
    cmd = { "Make" },
  },
}
