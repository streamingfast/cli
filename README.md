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
```

<details>
<summary><b>Comparison vs plain &cobra.Command construction</b></summary>

<table>
<tr><th>Before</th><th>After</th></tr>
<tr>
<td>
<pre>
var protogenCmd = &cobra.Command{
	Use:   "protogen [<manifest>]",
	Short: "Generate Rust bindings from a package",
	Long: cli.Dedent(`
		Generate Rust bindings from a package. The manifest is optional as it will try to find a file named
		'substreams.yaml' in current working directory if nothing entered. You may enter a directory that contains a 'substreams.yaml'
		file in place of '<manifest_file>', or a link to a remote .spkg file, using urls gs://, http(s)://, ipfs://, etc.'.
	`),
	RunE:         runProtogen,
	Args:         cobra.RangeArgs(0, 1),
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(protogenCmd)

	flags := protogenCmd.Flags()
	flags.StringP("output-path", "o", "src/pb", cli.FlagDescription(`
		Directory to output generated .rs files, if the received <package> argument is a local Substreams manifest file
		(e.g. a local file ending with .yaml), the output path will be made relative to it
	`))
	flags.StringArrayP("exclude-paths", "x", []string{}, "Exclude specific files or directories, for example \"proto/a/a.proto\" or \"proto/a\"")
	flags.Bool("generate-mod-rs", true, cli.FlagDescription(`
		Generate the protobuf 'mod.rs' file alongside the rust bindings. Include '--generate-mod-rs=false' If you wish to disable this generation.
		If there is a present 'buf.gen.yaml', consult https://github.com/neoeinstein/protoc-gen-prost/blob/main/protoc-gen-prost-crate/README.md to add 'mod.rs' generation functionality.
	`))
	flags.Bool("show-generated-buf-gen", false, "Whether to show the generated buf.gen.yaml file or not")
}
</pre>
</td>
<td>
<pre>
var protogenCmd = Command(
	runProtogen,
	"protogen [<manifest>]",
	"Generate Rust bindings from a package",
	Description(`
		Generate Rust bindings from a package. The manifest is optional as it will try to find a file named
		'substreams.yaml' in current working directory if nothing entered. You may enter a directory that contains a 'substreams.yaml'
		file in place of '<manifest_file>', or a link to a remote .spkg file, using urls gs://, http(s)://, ipfs://, etc.'.
	`),
	RangeArgs(0, 1),
	Flags(func(flags *pflag.FlagSet) {
		flags.StringP("output-path", "o", "src/pb", FlagDescription(`
			Directory to output generated .rs files, if the received <package> argument is a local Substreams manifest file
			(e.g. a local file ending with .yaml), the output path will be made relative to it
		`))
		flags.StringArrayP("exclude-paths", "x", []string{}, "Exclude specific files or directories, for example \"proto/a/a.proto\" or \"proto/a\"")
		flags.Bool("generate-mod-rs", true, FlagDescription(`
			Generate the protobuf 'mod.rs' file alongside the rust bindings. Include '--generate-mod-rs=false' If you wish to disable this generation.
			If there is a present 'buf.gen.yaml', consult https://github.com/neoeinstein/protoc-gen-prost/blob/main/protoc-gen-prost-crate/README.md to add 'mod.rs' generation functionality.
		`))
		flags.Bool("show-generated-buf-gen", false, "Whether to show the generated buf.gen.yaml file or not")
	}),
)
</pre>
</td>
</table>
</details>

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
