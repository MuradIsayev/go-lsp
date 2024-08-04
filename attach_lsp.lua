local vimClient = vim.lsp.start_client({
	name = "go-lsp",
	cmd = { "/Users/asmarisayeva/Documents/Projects/go-lsp/main" },
})

if not vimClient then
	vim.notify("Failed to start custom-go-lsp")
	return
end

vim.api.nvim_create_autocmd("FileType", {
	pattern = "markdown",
	callback = function()
		vim.lsp.buf_attach_client(0, vimClient)
	end,
})
