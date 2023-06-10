package main

import (
	"fmt"

	"github.com/spf13/cobra"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
)

// Injected at build time
var version = ""

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

		ConfigureVersion(version),
		ConfigureViper("PROJECT"),
		OnCommandErrorLogAndExit(zlog),
	)
}

func run(cmd *cobra.Command, args []string) error {
	zlog.Info("Executing")

	zlog.Debug("Will be displayed if DLOG=debug (equivalent to DLOG='.*=debug')")
	if tracer.Enabled() {
		zlog.Debug("Will be displayed if DLOG=trace (equivalent to DLOG='.*=trace')")
	}

	return fmt.Errorf("testing")
}
