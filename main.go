package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "ConventionalCommits"

var version string = "0.0.1"
var handler protocol.Handler

var inMemoryDocument []byte

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	log.Println("CCommits log folder: ", exPath)
	logFilePath := fmt.Sprintf("%s/%s", exPath, "log.txt")
	log.Println("CCommits initial log file: ", logFilePath)
	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	handler = protocol.Handler{
		Initialize:            initialize,
		Initialized:           initialized,
		Shutdown:              shutdown,
		SetTrace:              setTrace,
		TextDocumentDidOpen:   docDidOpen,
		TextDocumentDidChange: docDidChange,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func panicOnError(err error) {
	if err != nil {
		log.Println("CCommits fatal error:", err)
	}
}

func docDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	content := string(inMemoryDocument)

	for _, change := range params.ContentChanges {
		if change_, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			startIndex, endIndex := change_.Range.IndexesIn(content)
			content = content[:startIndex] + change_.Text + content[endIndex:]
		} else if change_, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			content = change_.Text
		}
	}

	inMemoryDocument = []byte(content)

	log.Println("CCommits in memory file content", content)

	diagnostics, err := AnalizeContent(content)
	if err != nil {
		return err
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})

	return nil
}

func docDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	log.Println("docDidOpen", params)

	// Limit the size of analized content. Usually a commit contains all changes
	commitMsgContent := params.TextDocument.Text[0:600]

	log.Println("commitMsgContent", commitMsgContent)

	inMemoryDocument = []byte(params.TextDocument.Text)
	return nil
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	log.Println("CCommits shutdown.")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
