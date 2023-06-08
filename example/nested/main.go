package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
)

// Injected at build time
var version = ""

var zlog, tracer = logging.RootLogger("project", "github.com/acme/project")

func main() {
	logging.InstantiateLoggers()

	Run(
		"runner",
		"Some random command runner with 2 sub-commands",

		ConfigureVersion(version),
		ConfigureViper("PROJECT"),

		Group(
			"generate",
			"Quick group summary, without a description",
			Command(generateImgE,
				"image",
				"Quick command summary, without a description",
				Flags(func(flags *pflag.FlagSet) {
					flags.String("version", "", "A flag description")
				}),
			),
		),

		Command(compareE,
			"compare <input_file>",
			"Quick command summary, with a description, the actual usage above is descriptive, you must handle the arguments manually",
			Description(`
				Description of the command, automatically de-indented by using first line identation,
				use 'go run ./example/nested compare --help' to see it in action!
			`),
			ExamplePrefixed("runner", `
				compare relative_file.json
				compare /absolute/file.json
			`),
			ExactArgs(2),
		),

		OnCommandErrorLogAndExit(zlog),
	)
}

func generateImgE(cmd *cobra.Command, args []string) error {
	_, err := os.Getwd()
	NoError(err, "unable to get working directory")

	fmt.Println("Generating something")
	return nil
}

func compareE(cmd *cobra.Command, args []string) error {
	shouldContinue, wasAnswered := AskConfirmation(`Do you want to continue?`)
	if wasAnswered && shouldContinue {

	} else {
		fmt.Println("Not showing diff between files, run the following command to see it manually:")
	}

	return nil
}
