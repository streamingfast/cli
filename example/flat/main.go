package main

import (
	"fmt"

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
	cli.Run(
		"flat", "A flat command",
		Execute(run),
		Description(`
			Description of the command, automatically de-indented by using first line identation,
			use 'runner generate --help to see it in action!
		`),
	)
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("Executed!")
	return nil
}

var Run = cli.Run
var Command = cli.Command

type Execute = cli.Execute
type Description = cli.Description
