package cobrax

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestWalkCommands(t *testing.T) {
	// Create a test command tree:
	// root
	// ├── sub1
	// │   ├── sub1sub1
	// │   └── sub1sub2
	// └── sub2
	root := &cobra.Command{Use: "root"}
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	sub1sub1 := &cobra.Command{Use: "sub1sub1"}
	sub1sub2 := &cobra.Command{Use: "sub1sub2"}

	root.AddCommand(sub1, sub2)
	sub1.AddCommand(sub1sub1, sub1sub2)

	// Walk the commands and collect names
	var commandNames []string
	for cmd := range WalkCommands(root) {
		commandNames = append(commandNames, cmd.Use)
	}

	// Should visit commands in depth-first order
	expected := []string{"root", "sub1", "sub1sub1", "sub1sub2", "sub2"}
	assert.Equal(t, expected, commandNames)
}

func TestWalkCommandsSingleCommand(t *testing.T) {
	root := &cobra.Command{Use: "root"}

	var commandNames []string
	for cmd := range WalkCommands(root) {
		commandNames = append(commandNames, cmd.Use)
	}

	expected := []string{"root"}
	assert.Equal(t, expected, commandNames)
}

func TestWalkCommandsEarlyExit(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	sub1 := &cobra.Command{Use: "sub1"}
	sub2 := &cobra.Command{Use: "sub2"}
	root.AddCommand(sub1, sub2)

	var commandNames []string
	for cmd := range WalkCommands(root) {
		commandNames = append(commandNames, cmd.Use)
		if cmd.Use == "sub1" {
			break // Test early exit
		}
	}

	expected := []string{"root", "sub1"}
	assert.Equal(t, expected, commandNames)
}