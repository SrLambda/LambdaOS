local tui_flags = vim.g.tui_flags or {}
local enable_neotree = tui_flags.enable_neotree ~= false

local plugins = {
  {
    "nvim-telescope/telescope.nvim",
    branch = "0.1.x",
    event = "VeryLazy",
    dependencies = {
      "nvim-lua/plenary.nvim",
      {
        "nvim-telescope/telescope-fzf-native.nvim",
        build = "make",
      },
    },
    config = function()
      local telescope = require("telescope")
      local actions = require("telescope.actions")
      telescope.setup({
        defaults = {
          mappings = {
            i = {
              ["<C-j>"] = actions.move_selection_next,
              ["<C-k>"] = actions.move_selection_previous,
              ["<C-q>"] = actions.send_selected_to_qflist + actions.open_qflist,
            },
          },
          file_ignore_patterns = { "node_modules", ".git/" },
        },
        pickers = {
          find_files = { hidden = true },
          live_grep = { additional_args = { "--hidden" } },
        },
      })
      telescope.load_extension("fzf")

      local map = vim.keymap.set
      map("n", "<leader>ff", "<cmd>Telescope find_files<CR>", { desc = "Find files" })
      map("n", "<leader>fg", "<cmd>Telescope live_grep<CR>", { desc = "Live grep" })
      map("n", "<leader>fb", "<cmd>Telescope buffers<CR>", { desc = "Find buffers" })
      map("n", "<leader>fh", "<cmd>Telescope help_tags<CR>", { desc = "Help tags" })
      map("n", "<leader>fr", "<cmd>Telescope oldfiles<CR>", { desc = "Recent files" })
      map("n", "<leader>fs", "<cmd>Telescope lsp_document_symbols<CR>", { desc = "LSP symbols" })
      map("n", "<leader>ft", "<cmd>Telescope treesitter<CR>", { desc = "Treesitter symbols" })
    end,
  },

  {
    "ThePrimeagen/harpoon.nvim",
    branch = "harpoon2",
    event = "VeryLazy",
    dependencies = { "nvim-lua/plenary.nvim" },
    config = function()
      local harpoon = require("harpoon")
      harpoon:setup()

      local map = vim.keymap.set
      map("n", "<leader>ha", function() harpoon:list():add() end, { desc = "Harpoon add" })
      map("n", "<leader>hm", function() harpoon.ui:toggle_quick_menu(harpoon:list()) end, { desc = "Harpoon menu" })
      map("n", "<leader>h1", function() harpoon:list():select(1) end, { desc = "Harpoon 1" })
      map("n", "<leader>h2", function() harpoon:list():select(2) end, { desc = "Harpoon 2" })
      map("n", "<leader>h3", function() harpoon:list():select(3) end, { desc = "Harpoon 3" })
      map("n", "<leader>h4", function() harpoon:list():select(4) end, { desc = "Harpoon 4" })
    end,
  },
}

if enable_neotree then
  table.insert(plugins, {
    "nvim-neo-tree/neo-tree.nvim",
    branch = "v3.x",
    event = "VeryLazy",
    dependencies = {
      "nvim-lua/plenary.nvim",
      "nvim-tree/nvim-web-devicons",
      "MunifTanjim/nui.nvim",
    },
    config = function()
      require("neo-tree").setup({
        close_if_last_window = true,
        window = {
          position = "left",
          width = 30,
        },
        filesystem = {
          filtered_items = {
            visible = true,
            hide_dotfiles = false,
            hide_gitignored = false,
          },
        },
      })

      vim.keymap.set("n", "<leader>e", "<cmd>Neotree toggle<CR>", { desc = "Toggle Neo-tree" })
      vim.keymap.set("n", "<leader>bf", "<cmd>Neotree reveal<CR>", { desc = "Buffer in Neo-tree" })
    end,
  })
end

return plugins
