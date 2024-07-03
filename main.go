package main

import (
	"bufio"
	"log"
	"os"

	"github.com/MuradIsayev/go-lsp/rpc"
)

func main() {
	logger := GetLogger("log.txt")
	logger.Println("Starting LSP server")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Failed to decode message: ", err)
			continue
		}
		handleMessage(logger, method, contents)
	}

}

func handleMessage(logger *log.Logger, method string, msg []byte) {
	logger.Printf("Received message with method: %s", method)
}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("Failed to open log file")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
