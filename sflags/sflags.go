package sflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

// MustGetViperKey returns the viper key associated with the flag `name` on the
// `cmd` (or one of its parent). It panics if the flag is not defined.
//
// You must have used `cli.ConfigureViperForCommand(cmd, prefix)` or `cli.ConfigureViper(prefix)`
// for this call to not panic.
func MustGetViperKey(cmd *cobra.Command, name string) string {
	flag := cmd.Flags().Lookup(name)
	if flag == nil {
		panic("flag not defined")
	}

	return MustGetViperKeyFromFlag(flag)
}

// MustGetViperKey returns the viper key associated with the flag.
//
// You must have used `cli.ConfigureViperForCommand(cmd, prefix)` or `cli.ConfigureViper(prefix)`
// for this call to not panic.
func MustGetViperKeyFromFlag(flag *pflag.Flag) string {
	viperKey, ok := getReboundKey(flag)
	if !ok {
		panic("flag not rebound, are you sure you have called 'cli.ConfigureViperForCommand(cmd, prefix)' or used 'cli.ConfigureViper(prefix)' when building up your CLI?")
	}

	return viperKey
}
