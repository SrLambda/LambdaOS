local opt = vim.opt

opt.tabstop = 4
opt.shiftwidth = 4
opt.softtabstop = 4
opt.expandtab = true
opt.smartindent = true

opt.textwidth = 80
opt.wrap = false

opt.number = true
opt.relativenumber = true
opt.cursorline = true
opt.signcolumn = "yes"

opt.mouse = "a"

opt.clipboard = "unnamedplus"

opt.swapfile = false
opt.backup = false
opt.undofile = true
opt.undodir = vim.fn.stdpath("data") .. "/undo"

opt.ignorecase = true
opt.smartcase = true

opt.termguicolors = true

opt.scrolloff = 8
opt.sidescrolloff = 8

opt.updatetime = 250
opt.timeoutlen = 300

opt.splitright = true
opt.splitbelow = true

opt.inccommand = "split"

opt.hlsearch = false
opt.incsearch = true

opt.list = true
opt.listchars = { tab = "» ", trail = "·", nbsp = "␣" }

vim.g.mapleader = " "
vim.g.maplocalleader = " "
