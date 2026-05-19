-- Rust configuration using rustaceanvim
-- Loaded only for rust filetype

local rust_ok, rust_tools = pcall(require, "rustaceanvim")
if rust_ok then
  vim.g.rustaceanvim = {
    tools = {
      inlay_hints = {
        auto = true,
        only_current_line = false,
        show_parameter_hints = true,
      },
    },
    server = {
      on_attach = function(client, bufnr)
        local map = vim.keymap.set
        local opts = { buffer = bufnr, silent = true }

        map("n", "<leader>rr", "<cmd>RustRunnables<CR>", { desc = "Rust runnables" })
        map("n", "<leader>rd", "<cmd>RustDebuggables<CR>", { desc = "Rust debuggables" })
        map("n", "<leader>rh", "<cmd>RustHoverActions<CR>", { desc = "Rust hover actions" })
        map("n", "<leader>ra", "<cmd>RustCodeAction<CR>", { desc = "Rust code action" })
        map("n", "<leader>re", "<cmd>RustExpandMacro<CR>", { desc = "Rust expand macro" })
        map("n", "<leader>rc", "<cmd>RustOpenCargo<CR>", { desc = "Open Cargo.toml" })
      end,
      settings = {
        ["rust-analyzer"] = {
          checkOnSave = {
            command = "clippy",
          },
          cargo = {
            allFeatures = true,
          },
        },
      },
    },
    dap = {},
  }
end

vim.opt_local.shiftwidth = 4
vim.opt_local.tabstop = 4
vim.opt_local.expandtab = true
