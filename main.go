package main

import (
	"encoding/json"
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

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	logFilePath := fmt.Sprintf("%s/%s", exPath, "log.txt")
	// open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("Logging to custom file:", exPath)

	handler = protocol.Handler{
		Initialize:          initialize,
		Initialized:         initialized,
		Shutdown:            shutdown,
		SetTrace:            setTrace,
		TextDocumentDidOpen: docDidOpen,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func docDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	fmt.Println(params)

	out, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("docDidOpen", string(out))
	diagnostics := make([]protocol.Diagnostic, 0)
	diagnostics = append(diagnostics, protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 1,
			},
			End: protocol.Position{
				Line:      0,
				Character: 2,
			},
		},

		Message: "Fix this!",
	})

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})

	return nil
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()

	log.Println("initialize", params)

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
	log.Println("SHUTDOWN")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
