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
		values.CommitMessage.Separator,
		values.CommitMessage.Message,
	)
}

func formatFailureString(expected, result string) string {
	return fmt.Sprintf("\n Expected: \"%s\" \n Result: \"%s\" \n", expected, result)
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
				expectedErrorMessage: "test:1:1: unexpected token \"foobar\"",
			},
			{
				name:                 "it fails when missing message",
				commitMessage:        "feat",
				expectedErrorMessage: "test:1:5: unexpected token \"<EOF>\" (expected <separator> <message>)",
			},
			{
				name:                 "it fails when missing <separator> ':'",
				commitMessage:        "feat my message",
				expectedErrorMessage: "test:1:5: unexpected token \" my message\" (expected <separator> <message>)",
			},
			{
				name:                 "it fails when missing space between <separator> and <message>",
				commitMessage:        "feat:my message\n",
				expectedErrorMessage: "test:1:5: unexpected token \":my message\" (expected <separator> <message>)",
			},
			{
				name: "it fails when missing space between <separator> and Message (multiple lines)",
				commitMessage: `feat:without space

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main"b
# Your branch is up to date with 'origin/main'.`,
				expectedErrorMessage: "test:1:5: unexpected token \":without space\" (expected <separator> <message>)",
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
				Col: 6, // Expected an space got "f"
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
				commitMessage: "feat: some foo bar\n\n",
				expectedCommitMessage: CommitMessage{
					Type:      "feat",
					Modifier:  "",
					Separator: ": ",
					Message:   "some foo bar",
				},
			},
			{
				name:          "it accepts modifiers for type",
				commitMessage: "feat!: some foo bar\n\n",
				expectedCommitMessage: CommitMessage{
					Type:      "feat",
					Modifier:  "!",
					Separator: ": ",
					Message:   "some foo bar",
				},
			},
			{
				name:          "it accepts scope for type",
				commitMessage: "feat(foobarbaz)!: some foo bar\n\n",
				expectedCommitMessage: CommitMessage{
					Type:      "feat",
					Scope:     "foobarbaz",
					Modifier:  "!",
					Separator: ": ",
					Message:   "some foo bar",
				},
			},
		}

		for _, tcase := range testCases {
			td.Run(tcase.name, func(tt *testing.T) {
				values, err := ConvetionalCommitParse(fileName, tcase.commitMessage)
				commitMessage := values.CommitMessage

				if err != nil {
					tt.Fatal("Should not have failed.", err.Error())
					return
				}

				fmt.Println(commitMessage)

				if commitMessage.String() != tcase.expectedCommitMessage.String() {
					tt.Fatal(
						"CommitMessage.Type doesn't match.",
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
