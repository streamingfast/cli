package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/streamingfast/shutter"
	"go.uber.org/zap"
)

// Application is a simple object that can be used to manage the lifecycle of an application. You
// create the application with `NewApplication` and then you can use it to supervise and start
// child processes by using `SuperviseAndStart`.
//
// You then wait for the application to terminate by calling `WaitForTermination`. This call is blocking
// and register a signal handler to be notified of SIGINT and SIGTERM.
type Application struct {
	appCtx  context.Context
	shutter *shutter.Shutter
}

func NewApplication(ctx context.Context) *Application {
	shutter := shutter.New()

	appCtx, cancelApp := context.WithCancel(ctx)
	shutter.OnTerminating(func(_ error) {
		cancelApp()
	})

	return &Application{
		appCtx:  appCtx,
		shutter: shutter,
	}
}

func (a *Application) Context() context.Context {
	return a.appCtx
}

// Shutter interface over the `*shutter.Shutter` struct.
type Shutter interface {
	OnTerminated(f func(error))
	OnTerminating(f func(error))
	Shutdown(error)
}

// Runnable contracts is to be blocking and to return only when the task is done.
type Runnable interface {
	Run()
}

// RunnableError contracts is to be blocking and to return only when the task is done. We assume
// the error happens while bootstrapping the task.
type RunnableError interface {
	Run() error
}

// RunnableContext contracts is to be blocking and to return only when the task is done. The context
// must be **used** only for the bootstrap period of task, long running task should be tied to a
// `*shutter.Shutter` instance.
type RunnableContext interface {
	Run(ctx context.Context)
}

// RunnableContext contracts is to be blocking and to return only when the task is done. The context
// must be **used** only for the bootstrap period of task, long running task should be tied to a
// `*shutter.Shutter` instance. We assume the error happens while bootstrapping the task.
type RunnableContextError interface {
	Run(ctx context.Context) error
}

// Supervise the received child shutter, mainly, this ensures that on child's termination,
// the application is also terminated with the error that caused the child to terminate.
//
// If the application shuts down before the child, the child is also terminated but
// gracefully (it does **not** receive the error that caused the application to terminate).
//
// The child termination is always performed before the application fully complete, unless
// the gracecul shutdown delay has expired.
func (a *Application) Supervise(child Shutter) {
	child.OnTerminated(a.shutter.Shutdown)
	a.shutter.OnTerminating(func(_ error) {
		child.Shutdown(nil)
	})
}

// SuperviseAndStart calls [Supervise] and then starts the child in a goroutine. The received
// child must implement one of [Runnable], [RunnableContext], [RunnableError] or [RunnableContextError] to
// be able to be started correctly.
//
// The child is started in a goroutine and tied to the application lifecycle because we also
// called [Supervise]. Later the call to `WaitForTermination` will wait for the application to
// terminate which will also terminates and wait for all child.
func (a *Application) SuperviseAndStart(child Shutter) {
	a.Supervise(child)

	switch v := child.(type) {
	case Runnable:
		go v.Run()
	case RunnableContext:
		go v.Run(a.appCtx)
	case RunnableError:
		go func() {
			err := v.Run()
			if err != nil {
				child.Shutdown(err)
			}
		}()
	case RunnableContextError:
		go func() {
			err := v.Run(a.appCtx)
			if err != nil {
				child.Shutdown(err)
			}
		}()

	default:
		panic(fmt.Errorf("unsupported child type %T, must implement one of cli.Runnable, cli.RunnableContext, cli.RunnableError or cli.RunnableContextError", child))
	}
}

// WaitForTermination waits for the application to terminate. This first setup the signal handler and
// then wait for either the signal handler to be notified or the application to be terminating.
//
// On application terminating, all child registered with [Supervise] are also terminated. We then wait for
// all child to gracefully terminate. If the graceful shutdown delay is reached, we force the termination
// of the application right now.
//
// Doing Ctrl-C 4 times or more will lead to a force quit of the whole process by calling `os.Exit(1)`, this
// is performed by the signal handler code and is does **not** respect the graceful shutdown delay in this case.
func (a *Application) WaitForTermination(logger *zap.Logger, unreadyPeriodDelay, gracefulShutdownDelay time.Duration) error {
	// On any exit path, we synchronize the logger one last time
	defer func() {
		logger.Sync()
	}()

	signalHandler, isSignaled, _ := SetupSignalHandler(unreadyPeriodDelay, logger)
	select {
	case <-signalHandler:
		go a.shutter.Shutdown(nil)
		break
	case <-a.shutter.Terminating():
		logger.Info("run terminating", zap.Bool("from_signal", isSignaled.Load()), zap.Bool("with_error", a.shutter.Err() != nil))
		break
	}

	logger.Info("waiting for run termination")
	select {
	case <-a.shutter.Terminated():
	case <-time.After(gracefulShutdownDelay):
		logger.Warn("application did not terminate within graceful period of " + gracefulShutdownDelay.String() + ", forcing termination")
	}

	if err := a.shutter.Err(); err != nil {
		return err
	}

	logger.Info("run terminated gracefully")
	return nil
}
