package main

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Simple Conventional Commit v1.0.0 specification parser
// https://www.conventionalcommits.org/en/v1.0.0/#specification
var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`CommitType`, `(feat|fix)`},
		{`CommitTypeModifier`, `!`},
		{`Separator`, `:\s+`},
		{`Message`, `.*`},
		{"comment", `[#][^\n]*`},
	})

	conventionalCommitParser = participle.MustBuild[ConvCommit](
		participle.Lexer(iniLexer),
	)
)

type ConvCommit struct {
	CommitMessage *CommitMessage `@@`
	Comments      []*Comment     `@@*`
}

type Comment struct {
	Value string `@comment`
}

type CommitMessage struct {
	Type      string  `@CommitType`
	Modifier  string  `@CommitTypeModifier?`
	Separator string  `@Separator`
	Message   Message `@@`
}

func (c CommitMessage) String() string {
	values, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error while stringifying values: ", err.Error())
		return ""
	}
	return string(values)
}

type Message struct {
	Value string `@Message`
}

func ConvetionalCommitParse(file, message string) (*ConvCommit, error) {
	values, err := conventionalCommitParser.ParseString(file, message)
	return values, err
}
