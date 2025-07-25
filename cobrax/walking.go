package cobrax

import (
	"iter"

	"github.com/spf13/cobra"
)

// WalkCommands recursively walks through a Cobra command tree using Go's iterator pattern.
// It returns an iterator that yields each command in the tree, starting with the root command
// and then recursively walking through all subcommands in depth-first order.
//
// Usage:
//
//	for cmd := range WalkCommands(rootCmd) {
//		fmt.Println(cmd.Name())
//	}
func WalkCommands(root *cobra.Command) iter.Seq[*cobra.Command] {
	return func(yield func(*cobra.Command) bool) {
		walkCommandsRecursive(root, yield)
	}
}

// walkCommandsRecursive performs the actual recursive walking of the command tree.
// It returns false if the iteration should stop early (when yield returns false).
func walkCommandsRecursive(cmd *cobra.Command, yield func(*cobra.Command) bool) bool {
	// Yield the current command
	if !yield(cmd) {
		return false
	}

	// Recursively walk through all subcommands
	for _, subCmd := range cmd.Commands() {
		if !walkCommandsRecursive(subCmd, yield) {
			return false
		}
	}

	return true
}