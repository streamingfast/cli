package cli

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog, tracer = logging.PackageLogger("cli", "github.com/streamingfast/cli")

func SetLogger(newZlog *zap.Logger, newTracer logging.Tracer) {
	zlog = newZlog
	tracer = newTracer
}
