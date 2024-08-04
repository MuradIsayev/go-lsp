# `go-lsp`

Implemented _Language Server Protocol (LSP)_ in Go, providing features like **autocompletion**, **hover**, **code actions** and **diagnostics**. The server handles RPC communication between the text editor (in my case, Neovim) and language services. A Neovim integration script is included, allowing the LSP to attach to markdown files for testing (ideally it would be a programming language with a real compiler to analyse). Additionally, Logs (log.txt) capture all activities such as connections, method calls, and changes. 

Adding the code snippet inside _attach_lsp.lua_ to the Neovim config will basically tell Neovim to attach the given command (In my case, an executable binary file) to the markdown files (it can be any type of file). This means that every time we open a markdown file in Neovim, our language server will be attached to that file/buffer and remain active until the file/buffer is closed.
