package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Simple Conventional Commit v1.0.0 specification parser
// https://www.conventionalcommits.org/en/v1.0.0/#specification
var (
	convCommitLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`CommitType`, `(feat|fix)`},
		{`CommitTypeModifier`, `!`},
		{`Separator`, `:\s`},
		{`Message`, `.*[^\n]`},
		{`Description`, `[^#].*`},
		{"comment", `[#][^\n]*`},
	})

	conventionalCommitParser = participle.MustBuild[ConvCommit](
		participle.Lexer(convCommitLexer),
	)
)

type ConvCommit struct {
	Pos           lexer.Position
	EndPos        lexer.Position
	CommitMessage *CommitMessage `@@"\n"`
	Description   []*Description `@@*`
	Comments      []*Comment     `@@*`
}

func (cc *ConvCommit) String() string {
	values, err := json.Marshal(cc)
	if err != nil {
		fmt.Println("Error while stringifying values: ", err.Error())
		return ""
	}
	return string(values)
}

type Description struct {
	Value string `@Description`
}

type Comment struct {
	Value string `@comment`
}

type CommitMessage struct {
	Pos       lexer.Position
	EndPos    lexer.Position
	Type      string `@CommitType`
	Modifier  string `@CommitTypeModifier?`
	Separator string `@Separator`
	Message   string `@Message`
}

func (c *CommitMessage) String() string {
	values, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error while stringifying values: ", err.Error())
		return ""
	}
	return string(values)
}

type ViolationError struct {
	FileName string
	Err      error
}

func (r *ViolationError) Error() string {
	return r.Err.Error()
}

type ViolationPosition struct {
	Row int
	Col int
}

func (ep ViolationPosition) String() string {
	values, err := json.Marshal(ep)
	if err != nil {
		fmt.Println("Error while stringifying values: ", err.Error())
		return ""
	}
	return string(values)
}

func (ep *ViolationError) Position() ViolationPosition {
	errorMessage := ep.Err.Error()

	re := regexp.MustCompile(ep.FileName + `:(\d):(\d):`)
	match := re.FindStringSubmatch(errorMessage)

	row, err := strconv.Atoi(match[1])
	if err != nil {
		fmt.Println("Error while finding row position", errorMessage)
	}

	col, err := strconv.Atoi(match[2])
	if err != nil {
		fmt.Println("Error while finding col position", errorMessage)
	}

	return ViolationPosition{
		Row: row,
		Col: col,
	}
}

func NewViolationError(filename string, origError error) *ViolationError {
	return &ViolationError{
		FileName: filename,
		Err:      origError,
	}
}

func ConvetionalCommitParse(file, message string) (*ConvCommit, error) {
	values, err := conventionalCommitParser.ParseString(file, message)

	if err != nil {
		return values, NewViolationError(file, err)
	}

	return values, nil
}
