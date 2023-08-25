package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// SetupSignalHandler registers a signal handler for SIGINT and SIGTERM and returns a channel to receive
// determines if the service has been signaled to shutdown gracefully as well as 2 atomic booleans.
//
// The first one `hasBeenSignaled` determines if the service has been signaled already, this could be true
// before the channel is notified if `unreadyPeriodDelay > 0`. The usefulness of this is to change the
// readiness of your app based on the value in `hasBeenSignaled`. This way, your app is unready as soon
// as it has been signaled to shutdown which helps a lot in Kubernetes and other orchestration systems. See
// the `unreadyPeriodDelay` parameter details below for more information.
//
// The second one `waitedFullDelay` determines if the service has waited the full delay period before notifying
// the channel. This will be `true` unless `unreadyPeriodDelay > 0` and that `Ctrl-C` was pressed 4 times or more
// which forces a kill `os.Exit(1)`!
//
// The parameter `unreadyPeriodDelay` influences when you will be notified of the signal through the channel. If you
// set it to 0, you will be notified immediately. If you set it to 5 seconds, you will be notified of the signal
// after 5 seconds.
//
// The usefulness of this is to give your app some time to become unready before it is killed. This is important
// for some orchestration systems like Kubernetes where you can use this delay to mark your app as unready by
// checking the value of `hasBeenSignaled` while your app continues to work properly. This will give Kubernetes
// `unreadyPeriodDelay` time to stop sending traffic and removing from the active list of endpoints.
//
// This should usually be used for HTTP/gRPC servers, for other types of apps, you can set this to 0 to avoid
// this needlessly waiting period.
func SetupSignalHandler(unreadyPeriodDelay time.Duration, logger *zap.Logger) (receiveOutgoingSignals <-chan os.Signal, hasBeenSignaled, waitedFullDelay *atomic.Bool) {
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	seen := 0
	outgoingSignals := make(chan os.Signal, 10)

	receiveOutgoingSignals = outgoingSignals
	hasBeenSignaled = atomic.NewBool(false)
	waitedFullDelay = atomic.NewBool(false)

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
					hasBeenSignaled.Store(true)

					if unreadyPeriodDelay <= 0 {
						logger.Info("received termination signal and no unready period delay configured, exiting now")
						waitedFullDelay.Store(true)
						outgoingSignals <- s
						break
					}

					logger.Info(
						fmt.Sprintf("received termination signal, waiting for unready period delay %s before notifying listner (Ctrl+C again 3 times to force kill!)", unreadyPeriodDelay),
						zap.Stringer("signal", s),
					)

					go time.AfterFunc(unreadyPeriodDelay, func() {
						waitedFullDelay.Store(true)
						outgoingSignals <- s
					})
					break
				}

				logger.Info("received termination signal twice, shutting down now", zap.Stringer("signal", s))
				outgoingSignals <- s
			}
		}
	}()

	return outgoingSignals, hasBeenSignaled, waitedFullDelay
}
