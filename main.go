package main

import (
	"bufio"
	"encoding/json"
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

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Failed to decode message: ", err)
			continue
		}
		handleMessage(logger, state, method, contents)
	}

}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse the content with initialize method: %s", err)
		}

		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

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
	}

}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("Failed to open log file")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
