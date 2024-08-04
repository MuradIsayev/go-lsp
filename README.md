# `go-lsp`

### Implemented _Language Server Protocol (LSP)_ in Go, providing features like **autocompletion**, **hover**, **code actions** and **diagnostics**. The server handles RPC communication between the text editor (in my case, Neovim) and language services. A Neovim integration script is included, allowing the LSP to attach to markdown files for testing (ideally it would be a programming language with a real compiler to analyse). Additionally, Logs (log.txt) capture all activities such as connections, method calls, and changes. 
