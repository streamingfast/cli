package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

// Injected at build time
var version = ""

var zlog, _ = logging.RootLogger("project", "github.com/acme/project")

var rootCmd = &cobra.Command{
	Use:     "runner",
	Short:   "Some random command runner with 2 sub-commands",
	Version: version,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Quick group summary, without a description",
}

var generateImgCmd = &cobra.Command{
	Use:   "image",
	Short: "Quick command summary, without a description",
	RunE:  generateImgE,
}

var compareCmd = &cobra.Command{
	Use:   "compare <input_file>",
	Short: "Quick command summary, with a description, the actual usage above is descriptive, you must handle the arguments manually",
	Long: "Description of the command, automatically de-indented by using first line identation\n" +
		"use 'go run ./example/nested compare --help' to see it in action!",
	Example: "runner compare relative_file.json\n" +
		"runner compare /absolute/file.json",
	Args: cobra.ExactArgs(2),
	RunE: compareE,
}

func before() {
	logging.InstantiateLoggers()

	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateImgCmd)
	rootCmd.AddCommand(compareCmd)

	generateImgCmd.Flags().String("version", "", "A flag description")

	cli.ConfigureViperForCommand(rootCmd, "PROJECT")

	// It's not possible to configure the version of a command as easily requires a bit of code to extract
	// the commit/date from the binary, it add ~30 lines of code. Some other googdies from cli command
	// building are lost also like default silence usage on error which is annoying to set up right.

	logAndExit := func(code int, err error) {
		zlog.Error("unable to execute command", zap.Error(err))
		zlog.Sync()
		os.Exit(code)
	}

	cli.OnAssertionFailure = func(message string) {
		logAndExit(1, fmt.Errorf(message))
	}

	if err := rootCmd.Execute(); err != nil {
		logAndExit(1, err)
	}
}

func generateImgE(*cobra.Command, []string) error { return nil }
func compareE(*cobra.Command, []string) error     { return nil }
