return {
  {
    "stevearc/conform.nvim",
    event = { "BufReadPost", "BufNewFile" },
    cmd = { "ConformInfo" },
    config = function()
      local conform = require("conform")

      conform.setup({
        formatters_by_ft = {
          lua = { "stylua" },
          python = { "isort", "black" },
          javascript = { "prettierd", "prettier", stop_after_first = true },
          typescript = { "prettierd", "prettier", stop_after_first = true },
          javascriptreact = { "prettierd", "prettier", stop_after_first = true },
          typescriptreact = { "prettierd", "prettier", stop_after_first = true },
          json = { "prettierd", "prettier", stop_after_first = true },
          jsonc = { "prettierd", "prettier", stop_after_first = true },
          yaml = { "prettierd", "prettier", stop_after_first = true },
          markdown = { "prettierd", "prettier", stop_after_first = true },
          html = { "prettierd", "prettier", stop_after_first = true },
          css = { "prettierd", "prettier", stop_after_first = true },
          scss = { "prettierd", "prettier", stop_after_first = true },
          rust = { "rustfmt" },
          go = { "gofumpt", "goimports" },
          c = { "clang-format" },
          cpp = { "clang-format" },
          java = { "google-java-format" },
          sh = { "shfmt" },
          bash = { "shfmt" },
          zsh = { "shfmt" },
          toml = { "taplo" },
          sql = { "sql_formatter" },
        },
        format_on_save = function(bufnr)
          if vim.g.disable_autoformat or vim.b[bufnr].disable_autoformat then
            return
          end
          return { timeout_ms = 5000, lsp_fallback = true }
        end,
        default_format_opts = {
          lsp_fallback = true,
        },
      })

      vim.keymap.set({ "n", "v" }, "<leader>cf", function()
        conform.format({ lsp_fallback = true, timeout_ms = 5000 })
      end, { desc = "Format" })

      vim.api.nvim_create_user_command("FormatDisable", function(args)
        if args.bang then
          vim.b.disable_autoformat = true
        else
          vim.g.disable_autoformat = true
        end
      end, { desc = "Disable auto-format-on-save", bang = true })

      vim.api.nvim_create_user_command("FormatEnable", function()
        vim.g.disable_autoformat = nil
        vim.b.disable_autoformat = nil
      end, { desc = "Enable auto-format-on-save" })
    end,
  },
}
