return {
  {
    "Vigemus/iron.nvim",
    event = "VeryLazy",
    config = function()
      local iron = require("iron.core")
      iron.setup({
        config = {
          scratch_repl = true,
          repl_definition = {
            python = {
              command = { "python" },
              format = require("iron.fts.common").bracketed_paste,
            },
            sh = {
              command = { "bash" },
            },
            lua = {
              command = { "lua" },
            },
            r = {
              command = { "R" },
            },
            julia = {
              command = { "julia" },
            },
          },
          repl_open_cmd = require("iron.view").split.vertical.botright(0.4),
        },
        keymaps = {
          send_motion = "<leader>sc",
          visual_send = "<leader>sc",
          send_file = "<leader>sf",
          send_line = "<leader>sl",
          send_paragraph = "<leader>sp",
          send_until_cursor = "<leader>su",
          send_mark = "<leader>sm",
          cr = "<leader>s<CR>",
          interrupt = "<leader>si",
          exit = "<leader>sq",
          clear = "<leader>sx",
        },
        ignore_blank_lines = true,
      })
    end,
  },

  {
    "GCBallesteros/jupytext.nvim",
    event = "VeryLazy",
    config = true,
  },
}
