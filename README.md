# dfuse CLI Library

This is a quick and simple library containing low-level primitive to create quick and hackish CLI applications without all the fuzz required for a full production app. This library is opiniated and aims at doing CLI tool rapidly.

**Note** This library is experimental and the API could change without notice.

## Example

The folder [./example](./example) contains example usage of the library. You can run them easily with a `go run ...`

 * Flat - `go run ./example/flat`
 * Nested - `go run ./example/nested`

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
```