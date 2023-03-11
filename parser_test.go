package main

import (
	"fmt"
	"strings"
	"testing"
)

func printValues(values *ConvCommit) {
	if values.CommitMessage == nil {
		fmt.Println("Commit message parser got no matches")
		return
	}

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
				expectedErrorMessage: "test:1:5: unexpected token \"<EOF>\" (expected <separator> Message)",
			},
			{
				name:                 "it fails when missing <separator> ':'",
				commitMessage:        "feat my message",
				expectedErrorMessage: "test:1:5: unexpected token \" my message\" (expected <separator> Message)",
			},
			{
				name:                 "it fails when missing space between <separator> and Message",
				commitMessage:        "feat:my message",
				expectedErrorMessage: "test:1:5: unexpected token \":my message\" (expected <separator> Message)",
			},
		}

		for _, tcase := range testCases {
			td.Run(tcase.name, func(tt *testing.T) {
				values, err := ConvetionalCommitParse(fileName, tcase.commitMessage)

				if err == nil {
					tt.Fatal("Should have failed.")
					return
				}

				errorMessage := err.Error()
				if !strings.Contains(errorMessage, tcase.expectedErrorMessage) {
					tt.Fatal(formatFailureString(tcase.expectedErrorMessage, errorMessage))
					return
				}

				printValues(values)
			})
		}
	})

	t.Run("correct syntax", func(td *testing.T) {
		testCases := []struct {
			name                  string
			commitMessage         string
			expectedCommitMessage CommitMessage
		}{
			{
				name:          "it returns a parsed Conventional Commit",
				commitMessage: "feat: some foo bar",
				expectedCommitMessage: CommitMessage{
					Type:      "feat",
					Modifier:  "",
					Separator: ": ",
					Message: Message{
						Value: "some foo bar",
					},
				},
			},
			{
				name:          "it accepts modifiers for type",
				commitMessage: "feat!: some foo bar",
				expectedCommitMessage: CommitMessage{
					Type:      "feat",
					Modifier:  "!",
					Separator: ": ",
					Message: Message{
						Value: "some foo bar",
					},
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
							commitMessage.String(),
							tcase.expectedCommitMessage.String(),
						),
					)
				}
			})
		}
	})
}
