# StreamingFast CLI Library
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/dgrpc)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Quick and opinionated library aiming to create CLI application rapidly. The library contains CLI primitives (around Cobra/Viper) as well as low-level primitives to ease the creation of developer scripts.

**Note** This library is experimental and the API could change without notice.

## Example

The folder [./example](./example) contains example usage of the library. You can run them easily them, open a terminal and navigate in the `example` folder, then `go run ...`

 * Flat - `go run ./flat`
 * Nested - `go run ./nested`

### Sample Boilerplate (copy/paste ready)

```golang
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
				Description of the command, automatically de-indented by using first line indentation,
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
	fmt.Println("Generating something from", cli.WorkingDirectory())
	return errors.New("showing on error look like (if you used `go run`, the 'exit status 1' is printed by it, compiled binary will not print it)")
}

func compareE(cmd *cobra.Command, args []string) error {
	shouldContinue, wasAnswered := cli.PromptConfirm(`Do you want to continue?`)
	if wasAnswered && shouldContinue {
		fmt.Println("Showing diff between files")
	} else {
		fmt.Println("Not showing diff between files, run the following command to see it manually:")
	}

	return nil
}
```

### Comparison between `cli.Command` and `&cobra.Command` manual construction

The best way to see the difference is by opening the before/after in your browser:

- [Before (`&cobra.Command` construction)](./example/compare/before/main.go)
- [After (`cli.Command` construction)](./example/compare/after/main.go)

And enjoy the feeling.

> **Note** The `cli` library is a wrapper around `cobra.Command`, so at the end you still deal with `*cobra.Command`.

## Contributing

**Issues and PR in this repo related strictly to the cli library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this codebase.

This codebase uses unit tests extensively, please write and run tests.

## License

[Apache 2.0](LICENSE)
