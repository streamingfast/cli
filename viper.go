package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ReboundFlagAnnotation = "github.com/streamingfast/cli#rebound-key"

// ConfigureViperForCommand sets env prefix to 'prefix', automatic env to check in env
// for any flags coming from anywhere (flag, config, default, etc.) as well as
// scoping flags to the command it's defined in for global acces.
func ConfigureViperForCommand(root *cobra.Command, envPrefix string) {
	viper.SetEnvPrefix(strings.ToUpper(envPrefix))
	viper.AutomaticEnv()

	// For backward compatibility, we support access through "_" and through "." for now,
	// configuring the actual key delimiter use on the global viper instance is not
	// possible.
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)

	recurseCommands(root, nil)
}

func recurseCommands(root *cobra.Command, segments []string) {
	if tracer.Enabled() {
		zlog.Debug("re-binding flags", zap.String("cmd", root.Name()), zap.Strings("segments", segments))

		defer func() {
			zlog.Debug("reboung flags terminated", zap.String("cmd", root.Name()))
		}()
	}

	persistentSegments := append(segments, "global")
	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		rebindFlag("persistent", f, append(persistentSegments, f.Name))
	})

	root.LocalNonPersistentFlags().VisitAll(func(f *pflag.Flag) {
		rebindFlag("local", f, append(segments, f.Name))
	})

	for _, cmd := range root.Commands() {
		recurseCommands(cmd, append(segments, cmd.Name()))
	}
}

func rebindFlag(tag string, f *pflag.Flag, segments []string) {
	newVarDash := strings.Join(segments, "-")
	newVarDot := strings.Join(segments, ".")

	addAnnotation(f, ReboundFlagAnnotation, newVarDot)

	viper.BindPFlag(newVarDash, f)
	viper.BindPFlag(newVarDot, f)

	zlog.Debug("binding "+tag+" flag", zap.String("actual", f.Name), zap.String("rebind_to", newVarDot+" (dash accepted)"))
}

func addAnnotation(flag *pflag.Flag, key string, value string) {
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}

	flag.Annotations[key] = append(flag.Annotations[key], value)
}
