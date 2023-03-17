package main

import (
	"fmt"
	"strings"
	"testing"
)

func printValues(values *ConvCommit) {
	fmt.Println("results",
		values.CommitTitle.Type,
		values.CommitTitle.Modifier,
		values.CommitTitle.Colon,
		values.CommitTitle.CommitDescription,
	)
}

func formatFailureString(expected, result string) string {
	return fmt.Sprintf("\n Expected: \"%s\" \n Result:   \"%s\" \n", expected, result)
}

func TestConventionalCommitsParser(t *testing.T) {
	fileName := "test"
	t.Run("wrong syntax", func(td *testing.T) {
		testCases := []struct {
			name                 string
			commitMessage        string
			expectedErrorMessage string
		}{
			{
				name:                 "it fails when missing type",
				commitMessage:        "foobar",
				expectedErrorMessage: "unexpected token \"foobar\"",
			},
			{
				name:                 "it fails when missing message",
				commitMessage:        "feat",
				expectedErrorMessage: "unexpected token \"<EOF>\" (expected <colon> <whitespace> <description>)",
			},
			{
				name:                 "it fails when missing <colon> ':'",
				commitMessage:        "feat my message",
				expectedErrorMessage: "unexpected token \" \"",
			},
			{
				name:                 "it fails when missing space between <colon> and <commitdescription>",
				commitMessage:        "feat:my message\n",
				expectedErrorMessage: "unexpected token \"my message\" (expected <whitespace>",
			},
			{
				name: "it fails when missing space between <colon> and CommitDescription (multiple lines)",
				commitMessage: `feat:without space

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main"b
# Your branch is up to date with 'origin/main'.`,
				expectedErrorMessage: "unexpected token \"without space\" (expected <whitespace> <description>)",
			},
		}

		for _, tcase := range testCases {
			td.Run(tcase.name, func(tt *testing.T) {
				_, err := ConvetionalCommitParse(fileName, tcase.commitMessage)

				if err == nil {
					tt.Fatal("Should have failed.\n")
					return
				}

				errorMessage := err.Error()
				if !strings.Contains(errorMessage, tcase.expectedErrorMessage) {
					tt.Fatal(formatFailureString(tcase.expectedErrorMessage, errorMessage))
					return
				}
			})
		}

		td.Run("contains the exact position where it failed", func(tt *testing.T) {
			_, err := ConvetionalCommitParse(fileName, "feat!:foobar")
			expectedErrorPosition := ViolationPosition{
				Row: 1,
				Col: 7, // Expected an space got "f"
			}

			if err == nil {
				tt.Fatal("Should have failed.\n")
				return
			}

			le, ok := err.(*ViolationError)
			if !ok {
				tt.Fatal("Error while parsing ErrorWithPosition.", err)
			}

			if le.Position().String() != expectedErrorPosition.String() {
				tt.Fatal(formatFailureString(le.Position().String(), expectedErrorPosition.String()))
			}
		})
	})

	t.Run("correct syntax", func(td *testing.T) {
		testCases := []struct {
			name                  string
			commitMessage         string
			expectedCommitMessage CommitTitle
		}{
			{
				name:          "it returns a parsed Conventional Commit",
				commitMessage: "feat: some foo bar\n#Please...\n",
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Modifier:          "",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "some foo bar",
				},
			},
			{
				name:          "it accepts modifiers for type",
				commitMessage: "feat!: some foo bar\n\n# Please ...\n",
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Modifier:          "!",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "some foo bar",
				},
			},
			{
				name:          "it accepts scope for type",
				commitMessage: "feat(foobarbaz)!: some foo bar\n\n# Please ...\n",
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Scope:             "foobarbaz",
					Modifier:          "!",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "some foo bar",
				},
			},
			{
				name: "it accepts scope for type",
				commitMessage: `feat!: foobar

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
`,
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Scope:             "",
					Modifier:          "!",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "foobar",
				},
			},
			{
				name: "it accepts git diff",
				commitMessage: `feat!: foobar

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
diff --git a/parser.go b/parser.go
index a4128be..3fcc654 100644
--- a/parser.go
+++ b/parser.go
@@ -32,6 +32,7 @@ var (
{"Whitespace", "\s"},
{"Colon", ":"},
{"Comment", "[#][^\n]*"},
+		{"GitDiff", "^diff.*\n"},

// Keywords
{"CommitType", "(feat|fix|chore|ci|docs|refactor|test)"},
@@ -51,6 +52,7 @@ type ConvCommit struct {
CommitTitle *CommitTitle "@@"
CommitDetails   []*CommitDetails "@@*"
Comments      []*Comment     "@@*"
+	Diff          []*string      "@GitDiff?"
}

func (cc *ConvCommit) String() string {
`,
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Scope:             "",
					Modifier:          "!",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "foobar",
				},
			},
			{
				name: "it accepts a commit body after commit title",
				commitMessage: `feat!: foobar

This is completly optional but possible

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
diff --git a/parser.go b/parser.go
index a4128be..3fcc654 100644
--- a/parser.go
+++ b/parser.go
@@ -32,6 +32,7 @@ var (
{"Whitespace", "\s"},
{"Colon", ":"},
{"Comment", "[#][^\n]*"},
+		{"GitDiff", "^diff.*\n"},

// Keywords
{"CommitType", "(feat|fix|chore|ci|docs|refactor|test)"},
@@ -51,6 +52,7 @@ type ConvCommit struct {
CommitTitle *CommitTitle "@@"
CommitDetails   []*CommitDetails "@@*"
Comments      []*Comment     "@@*"
+	Diff          []*string      "@GitDiff?"
}

func (cc *ConvCommit) String() string {
`,
				expectedCommitMessage: CommitTitle{
					Type:              "feat",
					Scope:             "",
					Modifier:          "!",
					Colon:             ":",
					Whitespace:        " ",
					CommitDescription: "foobar",
				},
			},
		}

		for _, tcase := range testCases {
			td.Run(tcase.name, func(tt *testing.T) {
				values, err := ConvetionalCommitParse(fileName, tcase.commitMessage)

				if err != nil {
					fmt.Println("parsed:", values)
					tt.Fatal("ERROR: Should not have failed.\n", err.Error(), "\n")
					return
				}

				if values.CommitTitle == nil {
					tt.Fatal("Empty values.", values)
				}

				commitMessage := values.CommitTitle

				fmt.Println(commitMessage)

				if commitMessage.String() != tcase.expectedCommitMessage.String() {
					tt.Fatal(
						"CommitTitle doesn't match.",
						formatFailureString(
							tcase.expectedCommitMessage.String(),
							commitMessage.String(),
						),
					)
				}
			})
		}
	})
}
