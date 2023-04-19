package cli

import (
	"fmt"
	"os"
)

// OnQuit is a global variable that can be overriden to control how the
// program should exit when cli module wants to quit.
//
// Implementation must enforce an hard-stop on the goroutine by doing
// either a 'panic' or an 'os.Exit(1)'
//
// The message can be "" in which case it should not be printed/logged.
var OnQuit = func(message string) {
	if message != "" {
		fmt.Println(message)
	}

	os.Exit(1)
}

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
	OnQuit(fmt.Sprintf(message, args...))
}
