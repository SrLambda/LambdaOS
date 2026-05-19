-- Java configuration using nvim-jdtls
-- Loaded only for java filetype

local config = {
  cmd = {
    "jdtls",
    "-configuration",
    vim.fn.expand("~/.cache/jdtls/config"),
    "-data",
    vim.fn.expand("~/.cache/jdtls/workspace/" .. vim.fn.fnamemodify(vim.fn.getcwd(), ":t")),
  },
  root_dir = vim.fs.dirname(vim.fs.find({ "gradlew", ".git", "mvnw" }, { upward = true })[1]),
}

local jdtls_ok, jdtls = pcall(require, "jdtls")
if jdtls_ok then
  jdtls.start_or_attach(config)
end

local map = vim.keymap.set
local opts = { buffer = 0, silent = true }

map("n", "<leader>jo", "<cmd>JdtOrganizeImports<CR>", { desc = "Organize imports" })
map("n", "<leader>jc", "<cmd>JdtCompile<CR>", { desc = "Compile" })
map("n", "<leader>jt", "<cmd>JdtTestClass<CR>", { desc = "Test class" })
map("n", "<leader>jtm", "<cmd>JdtTestMethod<CR>", { desc = "Test method" })

vim.opt_local.shiftwidth = 4
vim.opt_local.tabstop = 4
vim.opt_local.expandtab = true
