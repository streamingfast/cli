package main

import (
	"fmt"

	"github.com/spf13/cobra"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
)

var zlog, tracer = logging.RootLogger("project", "github.com/acme/project")

func main() {
	logging.InstantiateLoggers()

	Run(
		"flat <arg> [<optional_arg]",
		"A flat command short description",
		Execute(run),
		MinimumNArgs(1),
		Description(`
			Description of the command, automatically de-indented by using first line identation,
			use 'runner generate --help to see it in action!
		`),
		Example(`
			flat <value>
		`),
	)
}

func run(cmd *cobra.Command, args []string) error {
	zlog.Info("Executed")

	zlog.Debug("Will be displayed if DEBUG=true, DEBUG=* or DEBUG=flat (specific logger) is specified (or with TRACE env)")
	if tracer.Enabled() {
		zlog.Debug("Will be displayed if TRACE=true, TRACE=* or TRACE=flat (specific logger) is specified")
	}

	return fmt.Errorf("testing")
}
