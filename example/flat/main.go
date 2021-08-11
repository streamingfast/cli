package main

import (
	. "github.com/streamingfast/cli"
	"github.com/dfuse-io/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()
var tracer = logging.ApplicationLogger("nested", "github.com/acme/nested", &zlog)

func main() {
	Run(
		"flat", "A flat command",
		Execute(run),
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

	return nil
}
