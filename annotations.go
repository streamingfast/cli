package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var annotationOnError = "on-error"

func cliAnnotationKey(in string) string {
	return fmt.Sprintf("github.com/streamingfast/cli-%s", in)
}

var cmdToAnnotations = map[*cobra.Command]map[string]any{}

func setCommandAnnotation(cmd *cobra.Command, key string, value any) {
	cmdAnnotations := cmdToAnnotations[cmd]
	if cmdAnnotations == nil {
		cmdAnnotations = make(map[string]any)
		cmdToAnnotations[cmd] = cmdAnnotations
	}

	cmdAnnotations[cliAnnotationKey(key)] = value
}

func getCommandAnnotation(cmd *cobra.Command, key string) (any, bool) {
	cmdAnnotations, knownCmd := cmdToAnnotations[cmd]
	if knownCmd {
		value, found := cmdAnnotations[cliAnnotationKey(key)]
		return value, found
	}

	return nil, false
}
