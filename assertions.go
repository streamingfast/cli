package cli

import (
	"fmt"
)

// OnAssertionFailure is a global variable that can be overriden to control how the
// program should print/process when cli module assertion fail.
//
// The message can be "" in which case it should not be printed/logged.
//
// If your handler does not exit by itself, a call to `cli.Exit(1)` is performed
// after the handler has executed.
//
// If you exit yourself, you should use `cli.Exit(code)` so that exit handlers
// are called if any present.
var OnAssertionFailure func(message string)

func Ensure(condition bool, message string, args ...interface{}) {
	if !condition {
		Quit(message, args...)
	}
}

func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func Quit(message string, args ...interface{}) {
	if OnAssertionFailure != nil {
		OnAssertionFailure(fmt.Sprintf(message, args...))
	} else {
		fmt.Printf(message+"\n", args...)
	}

	Exit(1)
}
