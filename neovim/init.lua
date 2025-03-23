local ui = require("ui")
vim.api.nvim_set_keymap(
	"n",
	"<leader>db",
	":lua require('ui').open_floating_window()<CR>",
	{ noremap = true, silent = true }
)
