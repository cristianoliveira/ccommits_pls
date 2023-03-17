package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func RemoveParentesis(types ...string) participle.Option {
	if len(types) == 0 {
		types = []string{"String"}
	}
	return participle.Map(func(t lexer.Token) (lexer.Token, error) {
		value := strings.TrimPrefix(strings.TrimSuffix(t.Value, ")"), "(")

		t.Value = value
		return t, nil
	}, types...)
}

// Simple Conventional Commit v1.0.0 specification parser
// https://www.conventionalcommits.org/en/v1.0.0/#specification
var (
	convCommitLexer = lexer.MustSimple([]lexer.SimpleRule{
		// Fixed grammar
		{"Tab", `\t`},
		{"Newline", `\n`},
		{"Whitespace", `[\s\t]+`},
		{"Colon", `:`},
		{"Comment", `^#[^\n]*`},

		// Keywords
		{"CommitType", `(feat|fix|chore|ci|docs|refactor|test)`},
		{"CommitScope", `\(.*\)`},
		{"CommitTypeModifier", `!`},
		{"Details", `[^\n]*`},
	})

	conventionalCommitParser = participle.MustBuild[ConvCommit](
		participle.Lexer(convCommitLexer),
		RemoveParentesis("CommitScope"),
	)
)

type ConvCommit struct {
	CommitTitle   *CommitTitle   `@@`
	CommitDetails *CommitDetails `@@`
	Comments      []*Comment     `@@*`
	Rest          []*Rest        `@@*`
}

func (cc *ConvCommit) String() string {
	values, err := json.Marshal(cc)
	if err != nil {
		fmt.Println("Error while stringifying values: ", err.Error())
		return ""
	}
	return string(values)
}

type CommitDetails struct {
	Value string `~"#"@Newline* @Details? @Newline*`
}

type Rest struct {
	Value string `(?!"#") @Details? @Newline`
}

type Comment struct {
	Value string "@Comment @Newline"
}

type CommitTitle struct {
	Type       string "@CommitType"
	Scope      string "@CommitScope?"
	Modifier   string "@CommitTypeModifier?"
	Colon      string "@Colon"
	Whitespace string "@Whitespace"
	Details    string `@Details`
}

func (c *CommitTitle) String() string {
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

	if len(match) < 3 {
		return ViolationPosition{
			Row: 0,
			Col: 0,
		}
	}

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
	// It ignores git diff when parsing.
	gitDiffStartPoint := strings.Index(message, "diff --git")
	if gitDiffStartPoint <= 0 {
		gitDiffStartPoint = len(message)
	}

	values, err := conventionalCommitParser.ParseString(file, message[0:gitDiffStartPoint])

	if err != nil {
		return values, NewViolationError(file, err)
	}

	return values, nil
}
