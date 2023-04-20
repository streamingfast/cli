package cli

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// SetupSignalHandler handle termination signals, it returns
//
// this is a graceful delay to allow residual traffic sent by the load balancer to be processed
// without returning 500. Once the delay has passed then the service can be shutdown
func SetupSignalHandler(gracefulShutdownDelay time.Duration, logger *zap.Logger) (receiveOutgoingSignals <-chan os.Signal, hasBeenSignaled, isGraceful *atomic.Bool) {
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	seen := 0
	outgoingSignals := make(chan os.Signal, 10)

	receiveOutgoingSignals = outgoingSignals
	hasBeenSignaled = atomic.NewBool(false)
	isGraceful = atomic.NewBool(false)

	go func() {
		for {
			s := <-signals
			switch s {
			case syscall.SIGTERM, syscall.SIGINT:
				seen++

				if seen > 3 {
					logger.Info("received termination signal 3 times, forcing kill")
					logger.Sync()

					Exit(1)
				}

				if !hasBeenSignaled.Load() {
					if gracefulShutdownDelay <= 0 {
						logger.Info("received termination signal and no graceful shutdown delay configured, exiting now")
						isGraceful.Store(true)
						outgoingSignals <- s
						break
					}

					logger.Info("received termination signal (Ctrl+C multiple times to force kill)", zap.Stringer("signal", s))
					hasBeenSignaled.Store(true)

					go time.AfterFunc(gracefulShutdownDelay, func() {
						isGraceful.Store(true)
						outgoingSignals <- s
					})
					break
				}

				logger.Info("received termination signal twice, shutting down now", zap.Stringer("signal", s))
				outgoingSignals <- s
			}
		}
	}()

	return outgoingSignals, hasBeenSignaled, isGraceful
}
