package main

import (
	"fmt"
	"log"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func AnalizeContent(content string) ([]protocol.Diagnostic, error) {
	diagnostics := make([]protocol.Diagnostic, 0)

	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	// Git doesn't accept empty messages for commits, so no need to validate.
	if fmt.Sprintf("(%s)", lines[0]) == "()" {
		return diagnostics, nil
	}

	_, err := ConvetionalCommitParse("COMMIT_MESSAGE", content)
	if err != nil {
		errMessage := err.Error()
		log.Println("CCommits: Error", errMessage)

		violationError, ok := err.(*ViolationError)
		if !ok {
			log.Println("CCommits: Error while casting violation error")
		}

		col := violationError.Position().Col
		startCol := uint32(col)
		endCol := uint32(col + 1)

		errorDiagnostic := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: startCol,
				},
				End: protocol.Position{
					Line:      0,
					Character: endCol,
				},
			},

			Message: fmt.Sprintf("Violation: %s", errMessage), // err.Error(),
		}

		diagnostics = append(diagnostics, errorDiagnostic)
	}

	return diagnostics, nil
}
