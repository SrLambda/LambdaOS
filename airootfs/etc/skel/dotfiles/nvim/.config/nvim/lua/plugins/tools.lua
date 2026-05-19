return {
  {
    "lewis6991/gitsigns.nvim",
    event = { "BufReadPost", "BufNewFile" },
    config = function()
      require("gitsigns").setup({
        signs = {
          add = { text = "│" },
          change = { text = "│" },
          delete = { text = "_" },
          topdelete = { text = "‾" },
          changedelete = { text = "~" },
          untracked = { text = "┆" },
        },
        current_line_blame = true,
        current_line_blame_opts = {
          delay = 500,
        },
        on_attach = function(bufnr)
          local gitsigns = require("gitsigns")
          local map = vim.keymap.set

          map("n", "<leader>gb", gitsigns.toggle_current_line_blame, { buffer = bufnr, desc = "Toggle line blame" })
          map("n", "<leader>gd", gitsigns.diffthis, { buffer = bufnr, desc = "Diff this" })
          map("n", "<leader>gh", gitsigns.preview_hunk, { buffer = bufnr, desc = "Preview hunk" })
          map("n", "<leader>gr", gitsigns.reset_hunk, { buffer = bufnr, desc = "Reset hunk" })
          map("n", "<leader>gs", gitsigns.stage_hunk, { buffer = bufnr, desc = "Stage hunk" })
          map("v", "<leader>gs", function() gitsigns.stage_hunk({ vim.fn.line("."), vim.fn.line("v") }) end,
            { buffer = bufnr, desc = "Stage selected hunk" })
          map("n", "<leader>gu", gitsigns.undo_stage_hunk, { buffer = bufnr, desc = "Undo stage hunk" })
          map("n", "<leader>gD", gitsigns.reset_buffer, { buffer = bufnr, desc = "Reset buffer" })
          map("n", "]h", gitsigns.next_hunk, { buffer = bufnr, desc = "Next hunk" })
          map("n", "[h", gitsigns.prev_hunk, { buffer = bufnr, desc = "Prev hunk" })
        end,
      })
    end,
  },

  {
    "akinsho/toggleterm.nvim",
    version = "*",
    event = "VeryLazy",
    config = function()
      require("toggleterm").setup({
        size = 20,
        open_mapping = [[<C-\>]],
        hide_numbers = true,
        shade_terminals = true,
        shading_factor = 2,
        start_in_insert = true,
        insert_mappings = true,
        persist_size = true,
        direction = "float",
        float_opts = {
          border = "curved",
        },
      })

      local Terminal = require("toggleterm.terminal").Terminal
      local lazygit = Terminal:new({
        cmd = "lazygit",
        dir = "git_dir",
        direction = "float",
        float_opts = { border = "curved" },
        on_open = function(term)
          vim.cmd("startinsert!")
          vim.api.nvim_buf_set_keymap(term.bufnr, "t", "<Esc>", "<C-\\><C-n>", { noremap = true, silent = true })
        end,
      })

      function _lazygit_toggle()
        lazygit:toggle()
      end

      vim.keymap.set("n", "<leader>gg", "<cmd>lua _lazygit_toggle()<CR>", { desc = "Lazygit" })
      vim.keymap.set("n", "<leader>tf", "<cmd>ToggleTerm direction=float<CR>", { desc = "Float terminal" })
      vim.keymap.set("n", "<leader>th", "<cmd>ToggleTerm direction=horizontal<CR>", { desc = "Horizontal terminal" })
      vim.keymap.set("n", "<leader>tv", "<cmd>ToggleTerm direction=vertical<CR>", { desc = "Vertical terminal" })
    end,
  },
}
