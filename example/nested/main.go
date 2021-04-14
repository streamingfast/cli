package main

import (
	"fmt"
	"os"

	"github.com/dfuse-io/cli"
	"github.com/dfuse-io/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	logging.TestingOverride()
}

func main() {
	cli.Run("runner", "Some random command runner with 2 sub-commands",
		Command(generateE,
			"generate",
			"Quick command summary, without a description",
		),
		Command(compareE,
			"compare <input_file>",
			"Quick command summary, with a description, the actual usage above is descriptive, you must handle the arguments manually",
			Description(`
				Description of the command, automatically de-indented by using first line identation,
				use 'go run ./example/nested compare --help to see it in action!
			`),
		),
	)
}

func generateE(cmd *cobra.Command, args []string) error {
	_, err := os.Getwd()
	cli.NoError(err, "unable to get working directory")

	fmt.Println("Generating something")
	return nil
}

func compareE(cmd *cobra.Command, args []string) error {
	shouldContinue, wasAnswered := cli.AskConfirmation(`Do you want to continue?`)
	if wasAnswered && shouldContinue {

	} else {
		fmt.Println("Not showing diff between files, run the following command to see it manually:")
	}

	return nil
}

var Run = cli.Run
var Command = cli.Command

type Description = cli.Description
