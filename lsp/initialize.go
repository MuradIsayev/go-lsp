package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync   int            `json:"textDocumentSync"`
	HoverProvider      bool           `json:"hoverProvider"`
	DefinitionProvider bool           `json:"definitionProvider"`
	CodeActionProvider bool           `json:"codeActionProvider"`
	CompletionProvider map[string]any `json:"completionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				// tells IDE (client) to always send the whole document/file
				TextDocumentSync: 1,
				// tells IDE (client) that we support hover provider
				HoverProvider: true,
				// tells IDE (client) that we support definition provider
				DefinitionProvider: true,
				// tells IDE (client) that we support code action provider
				CodeActionProvider: true,
				// tells IDE (client) that we support completion provider
				CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "go-lsp",
				Version: "0.0.0.1-test",
			},
		},
	}
}
