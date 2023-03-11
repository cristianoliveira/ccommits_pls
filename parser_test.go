package main

import (
	"fmt"
	"strings"
	"testing"
)

func printValues(values *ConvCommit) {
	fmt.Println("results",
		values.CommitMessage.Type,
		values.CommitMessage.Modifier,
		values.CommitMessage.Colon,
		values.CommitMessage.Message,
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
				expectedErrorMessage: "unexpected token \"<EOF>\" (expected <colon> <whitespace> <message> <newline>+)",
			},
			{
				name:                 "it fails when missing <colon> ':'",
				commitMessage:        "feat my message",
				expectedErrorMessage: "unexpected token \" \"",
			},
			{
				name:                 "it fails when missing space between <colon> and <message>",
				commitMessage:        "feat:my message\n",
				expectedErrorMessage: "unexpected token \"my message\"",
			},
			{
				name: "it fails when missing space between <colon> and Message (multiple lines)",
				commitMessage: `feat:without space

			# Please enter the commit message for your changes. Lines starting
			# with '#' will be ignored, and an empty message aborts the commit.
			#
			# On branch main"b
			# Your branch is up to date with 'origin/main'.`,
				expectedErrorMessage: "unexpected token \"without space\" (expected <whitespace> <message> <newline>+)",
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
			expectedCommitMessage CommitMessage
		}{
			{
				name:          "it returns a parsed Conventional Commit",
				commitMessage: "feat: some foo bar\n",
				expectedCommitMessage: CommitMessage{
					Type:       "feat",
					Modifier:   "",
					Colon:      ":",
					Whitespace: " ",
					Message:    "some foo bar",
					Newline:    "\n",
				},
			},
			{
				name:          "it accepts modifiers for type",
				commitMessage: "feat!: some foo bar\n\n# Please",
				expectedCommitMessage: CommitMessage{
					Type:       "feat",
					Modifier:   "!",
					Colon:      ":",
					Whitespace: " ",
					Message:    "some foo bar",
					Newline:    "\n\n",
				},
			},
			{
				name:          "it accepts scope for type",
				commitMessage: "feat(foobarbaz)!: some foo bar\n\n",
				expectedCommitMessage: CommitMessage{
					Type:       "feat",
					Scope:      "foobarbaz",
					Modifier:   "!",
					Colon:      ":",
					Whitespace: " ",
					Message:    "some foo bar",
					Newline:    "\n\n",
				},
			},
		}

		for _, tcase := range testCases {
			td.Run(tcase.name, func(tt *testing.T) {
				values, err := ConvetionalCommitParse(fileName, tcase.commitMessage)

				if err != nil {
					tt.Fatal("ERROR: Should not have failed.\n", err.Error(), "\n")
					return
				}

				if values.CommitMessage == nil {
					tt.Fatal("Empty values.", values)
				}

				commitMessage := values.CommitMessage

				fmt.Println(commitMessage)

				if commitMessage.String() != tcase.expectedCommitMessage.String() {
					tt.Fatal(
						"CommitMessage doesn't match.",
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
