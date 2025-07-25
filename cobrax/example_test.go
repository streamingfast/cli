package cobrax_test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli/cobrax"
)

func ExampleWalkCommands() {
	// Create a simple command tree
	root := &cobra.Command{Use: "myapp"}
	start := &cobra.Command{Use: "start"}
	stop := &cobra.Command{Use: "stop"}
	config := &cobra.Command{Use: "config"}
	set := &cobra.Command{Use: "set"}
	get := &cobra.Command{Use: "get"}

	root.AddCommand(start, stop, config)
	config.AddCommand(set, get)

	// Walk through all commands
	for cmd := range cobrax.WalkCommands(root) {
		fmt.Println(cmd.Use)
	}

	// Output:
	// myapp
	// config
	// get
	// set
	// start
	// stop
}