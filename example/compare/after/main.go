package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
)

// Injected at build time
var version = ""

var zlog, _ = logging.RootLogger("project", "github.com/acme/project")

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

func generateImgE(*cobra.Command, []string) error { return nil }
func compareE(*cobra.Command, []string) error     { return nil }
