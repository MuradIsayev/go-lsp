package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/MuradIsayev/go-lsp/analysis"
	"github.com/MuradIsayev/go-lsp/lsp"
	"github.com/MuradIsayev/go-lsp/rpc"
)

func main() {
	logger := GetLogger("log.txt")
	logger.Println("Starting LSP server")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Failed to decode message: ", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}

}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse the content with initialize method: %s", err)
		}

		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)

		writeResponse(writer, msg)

		logger.Print("Sent the reply")

	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocomentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Couldn't parse the content with textDocument/didOpen method: %s", err)
			return
		}

		logger.Printf("Opened: %s", notification.Params.TextDocument.URI)

		// We got a document, and we have opened it
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
	case "textDocument/didChange":
		var notification lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Couldn't parse the content with textDocument/didChange method: %s", err)
			return
		}

		logger.Printf("Changed: %s %s", notification.Params.TextDocument.TextDocumentIdentifier.URI, notification.Params.ContentChanges)

		for _, change := range notification.Params.ContentChanges {
			state.UpdateDocument(notification.Params.TextDocument.TextDocumentIdentifier.URI, change.Text)
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse the content with textDocument/hover method: %s", err)
			return
		}

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse the content with textDocument/definition method: %s", err)
			return
		}

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse the content with textDocument/codeAction method: %s", err)
			return
		}

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)

		writeResponse(writer, response)
	}

}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("Failed to open log file")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
