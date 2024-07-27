package analysis

import (
	"fmt"
	"strings"

	"github.com/MuradIsayev/go-lsp/lsp"
)

type State struct {
	// map of file names and their contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File %s, Total Characters %d, Hover Position: Line=%d Character=%d", uri, len(document), position.Line, position.Character),
		},
	}

}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}

	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "fmt") {
			actions = append(actions, lsp.CodeAction{
				Title: "Remove fmt",
				Edit: &lsp.WorkspaceEdit{
					Changes: map[string][]lsp.TextEdit{
						uri: {
							{
								Range: lsp.Range{
									Start: lsp.Position{
										Line:      row,
										Character: strings.Index(line, "fmt") - 1,
									},
									End: lsp.Position{
										Line:      row,
										Character: strings.Index(line, "fmt") + len("fmt") + 1,
									},
								},
								NewText: "",
							},
						},
					},
				},
			})

			// another action to comment out the line
			actions = append(actions, lsp.CodeAction{
				Title: "Comment out the line",
				Edit: &lsp.WorkspaceEdit{
					Changes: map[string][]lsp.TextEdit{
						uri: {
							{
								Range: lsp.Range{
									Start: lsp.Position{
										Line:      row,
										Character: 0,
									},
									End: lsp.Position{
										Line:      row,
										Character: len(line),
									},
								},
								NewText: "// " + line,
							},
						},
					},
				},
			})
		}
	}

	return lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}
