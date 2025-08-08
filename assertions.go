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

// Ensure checks the condition and if it is false, it will call `cli.Quit` with the message
// effectively exiting the program with code 1.
func Ensure(condition bool, message string, args ...any) {
	if !condition {
		Quit(message, args...)
	}
}

// NoError checks if the error is nil, and if it is not, it will call `cli.Quit` with the message
// effectively exiting the program with code 1.
func NoError(err error, message string, args ...any) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

// Quit prints the message and exits the program with code 1.
// If OnAssertionFailure is set, it will be called with the message instead of printing it
// to stdout.
//
// If you want to exit with a different code, use `cli.Exit(code)` instead.
func Quit(message string, args ...any) {
	if OnAssertionFailure != nil {
		OnAssertionFailure(fmt.Sprintf(message, args...))
	} else {
		fmt.Printf(message+"\n", args...)
	}

	Exit(1)
}
