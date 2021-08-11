# StreamingFast CLI Library
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/dgrpc)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Quick and opinionated library aiming to create CLI application rapidly. The library contains CLI primitives (around Cobra/Viper) as well as low-level primitives to ease the creation of developer scripts.

**Note** This library is experimental and the API could change without notice.

## Example

The folder [./example](./example) contains example usage of the library. You can run them easily them, open a terminal and navigate in the `example` folder, then `go run ...`

 * Flat - `go run ./flat`
 * Nested - `go run ./nested`

### Public Helpers

| Method| Description |
|-|-|
| `cli.NoError(err, "file not found %q", fileName)` | Exit the process with exit code 1 and prints `fmt.Printf("file not found %q: %w\n", fileName, err)` if `err != nil` |
| `cli.Ensure(x == 0, "x point should be 0, got %d", x)` | Exit the process with exit code 1 and prints `fmt.Printf("x point should be 0, got %d\n", x)` if condition received is `false` |
| `cli.Quit("current date %d is too far away", time.Now())` | Exit the process with exit code 1 and prints `fmt.Printf("x point should be 0, got %d\n", x)` |
| `cli.FileExists("./some/file.png")` | Returns `true` if the file received in argument exists, `false` otherwise |
| `cli.CopyFile("current date %d is too far away", time.Now())` | Exit the process with exit code 1 and prints `fmt.Printf("x point should be 0, got %d\n", x)` |
|-|-|

### Sample Boilerplate (copy/paste ready)

```golang
package main

import (
	"fmt"
	"os"

	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()
var tracer = logging.ApplicationLogger("nested", "github.com/acme/nested", &zlog)

func main() {
	Run("runner", "Some random command runner with 2 sub-commands",
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
			Example("runner", `
				compare relative_file.json
				compare /absolute/file.json
			`),
		),
	)
}

func generateE(cmd *cobra.Command, args []string) error {
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
```

## Contributing

**Issues and PR in this repo related strictly to the cli library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

This codebase uses unit tests extensively, please write and run tests.

## License

[Apache 2.0](LICENSE)
