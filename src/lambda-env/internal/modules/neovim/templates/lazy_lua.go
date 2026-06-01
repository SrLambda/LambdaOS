package templates

const LazyLuaTemplate = `local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({"git", "clone", "--filter=blob:none", "https://github.com/folke/lazy.nvim.git", "--branch=stable", lazypath})
end
vim.opt.rtp:prepend(lazypath)

require("lazy").setup({
  spec = {
    { import = "plugins" },
{{if .EnableLSP}}    { import = "plugins.lsp" },{{end}}
{{if .EnableCopilot}}    { import = "plugins.ai" },{{end}}
{{if .EnableNeotree}}    { "nvim-neo-tree/neo-tree.nvim", branch = "v3.x", config = function() require("neo-tree").setup({ close_if_last_window = true, window = { position = "left", width = 30 } }) end },{{end}}
  },
  defaults = { lazy = false, version = false },
  install = { colorscheme = { "{{.Theme}}" } },
  checker = { enabled = true, notify = false },
  change_detection = { notify = false },
})
`
