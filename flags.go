package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Generate the flags based on Go code in this project directly, this however
// creates a chicken & egg problem if there is compilation error within the project
// but to fix them we must re-generate it.
//go:generate go run ./generate_flags flags_generated.go cli

// ConfigureViperForCommand sets env prefix to 'prefix', automatic env to check in env
// for any flags coming from anywhere (flag, config, default, etc.) as well as
// scoping flags to the command it's defined in for global acces.
func ConfigureViperForCommand(root *cobra.Command, prefix string) {
	viper.SetEnvPrefix(strings.ToUpper(prefix))
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	recurseCommands(root, nil)
}

func recurseCommands(root *cobra.Command, segments []string) {
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	zlog.Debug("re-binding flags", zap.String("cmd", root.Name()), zap.String("prefix", segmentPrefix))
	defer func() {
		zlog.Debug("reboung flags terminated", zap.String("cmd", root.Name()))
	}()

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		viper.BindPFlag(newVar, f)

		zlog.Debug("binding persistent flag", zap.String("actual", f.Name), zap.String("rebind", newVar))
	})

	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + f.Name
		viper.BindPFlag(newVar, f)

		zlog.Debug("binding flag", zap.String("actual", f.Name), zap.String("rebind", newVar))
	})

	for _, cmd := range root.Commands() {
		recurseCommands(cmd, append(segments, cmd.Name()))
	}
}
