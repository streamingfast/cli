package cobrax

import (
	"github.com/spf13/cobra"
)

// PrependPersistentPreRunE prepends (before existing) a PersistentPreRunE hook to the command but preserves existing hooks
// if they exists.
//
// If you don't use it, cobra takes a last writer wins approach applied to both PersistentPreRun and
// PersistentPreRunE.
func PrependPersistentPreRunE(cmd *cobra.Command, callback func(cmd *cobra.Command, args []string) error) {
	existingE := cmd.PersistentPreRunE
	existing := cmd.PersistentPreRun

	if existingE == nil && existing == nil {
		cmd.PersistentPreRunE = callback
		return
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := callback(cmd, args); err != nil {
			return err
		}

		if err := existingE(cmd, args); err != nil {
			return err
		}

		if existing != nil {
			existing(cmd, args)
		}

		return nil
	}
}

// AppendPersistentPreRunE appends a PersistentPreRunE hook to the command at the end and preserves existing hooks
// if they exists.
//
// If you don't use it, cobra takes a last writer wins approach applied to both PersistentPreRun and
// PersistentPreRunE.
func AppendPersistentPreRunE(cmd *cobra.Command, callback func(cmd *cobra.Command, args []string) error) {
	existingE := cmd.PersistentPreRunE
	existing := cmd.PersistentPreRun

	if existingE == nil && existing == nil {
		cmd.PersistentPreRunE = callback
		return
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := callback(cmd, args); err != nil {
			return err
		}

		if err := existingE(cmd, args); err != nil {
			return err
		}

		if existing != nil {
			existing(cmd, args)
		}

		return nil
	}
}
