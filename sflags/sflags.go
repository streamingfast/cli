package sflags

import (
	"github.com/spf13/cobra"
)

// Generate the flags based on Go code in this project directly, this however
// creates a chicken & egg problem if there is compilation error within the project
// but to fix them we must re-generate it.
//go:generate go run ./generator flags_generated.go sflags

// FlagDefined returns `true` if the flag is defined on the `cmd` (or one of its
// parent) and `false` otherwise.
func FlagDefined(cmd *cobra.Command, name string) bool {
	return cmd.Flags().Lookup(name) != nil
}
