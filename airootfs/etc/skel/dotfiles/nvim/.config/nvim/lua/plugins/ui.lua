return {
  {
    "nvim-lualine/lualine.nvim",
    event = "VeryLazy",
    config = function()
      require("lualine").setup({
        options = {
          theme = "catppuccin",
          component_separators = { left = "", right = "" },
          section_separators = { left = "", right = "" },
        },
        sections = {
          lualine_a = { "mode" },
          lualine_b = { "branch", "diff", "diagnostics" },
          lualine_c = { "filename" },
          lualine_x = { "encoding", "fileformat", "filetype" },
          lualine_y = { "progress" },
          lualine_z = { "location" },
        },
      })
    end,
  },

  {
    "goolord/alpha-nvim",
    event = "VimEnter",
    config = function()
      local alpha = require("alpha")
      local dashboard = require("alpha.themes.dashboard")
      dashboard.section.header.val = {
        "  ╻  ┏━┓ ┏┳┓ ┏┓   ┏┳┓ ┏━┓",
        "  ┃  ┣━┫ ┃┃┃ ┣┻┓   ┃ ┃ ┗━┓",
        "  ╹  ╹ ╹ ╹ ╹ ╹ ╹   ╹ ╹ ┗━┛",
      }
      dashboard.section.buttons.val = {
        dashboard.button("e", "  New file", "<cmd>ene<CR>"),
        dashboard.button("f", "  Find file", "<cmd>Telescope find_files<CR>"),
        dashboard.button("r", "  Recent files", "<cmd>Telescope oldfiles<CR>"),
        dashboard.button("g", "  Find text", "<cmd>Telescope live_grep<CR>"),
        dashboard.button("c", "  Config", "<cmd>e $MYVIMRC<CR>"),
        dashboard.button("q", "  Quit", "<cmd>qa<CR>"),
      }
      dashboard.section.footer.val = "LambdaOS Neovim"
      alpha.setup(dashboard.config)
    end,
  },

  {
    "akinsho/bufferline.nvim",
    event = "VeryLazy",
    dependencies = { "nvim-tree/nvim-web-devicons" },
    config = function()
      require("bufferline").setup({
        options = {
          mode = "buffers",
          numbers = "none",
          close_command = "bdelete! %d",
          right_mouse_command = "bdelete! %d",
          left_mouse_command = "buffer %d",
          middle_mouse_command = nil,
          indicator = { style = "underline" },
          separator_style = "thin",
          diagnostics = "nvim_lsp",
          offsets = {
            {
              filetype = "neo-tree",
              text = "File Explorer",
              highlight = "Directory",
              text_align = "left",
            },
          },
        },
      })
    end,
  },

  {
    "lukas-reineke/indent-blankline.nvim",
    event = "VeryLazy",
    main = "ibl",
    config = function()
      require("ibl").setup({
        scope = { enabled = true },
        indent = { char = "│" },
      })
    end,
  },

  {
    "nvim-tree/nvim-web-devicons",
    lazy = true,
  },
}
