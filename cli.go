package cli

import (
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandOption interface {
	Apply(cmd *cobra.Command)
}

type CommandOptionFunc func(cmd *cobra.Command)

func (f CommandOptionFunc) Apply(cmd *cobra.Command) {
	f(cmd)
}

type Flags func(flags *pflag.FlagSet)

func (f Flags) Apply(cmd *cobra.Command) {
	f(cmd.Flags())
}

type PersistentFlags func(flags *pflag.FlagSet)

func (f PersistentFlags) Apply(cmd *cobra.Command) {
	f(cmd.PersistentFlags())
}

// NoArgs returns an error if any args are included.
func NoArgs() Args {
	return Args(cobra.NoArgs)
}

// OnlyValidArgs returns an error if any args are not in the list of ValidArgs.
func OnlyValidArgs() Args {
	return Args(cobra.OnlyValidArgs)
}

// ArbitraryArgs never returns an error.
func ArbitraryArgs() Args {
	return Args(cobra.ArbitraryArgs)
}

// MinimumNArgs returns an error if there is not at least N args.
func MinimumNArgs(n int) Args {
	return Args(cobra.MinimumNArgs(n))
}

// MaximumNArgs returns an error if there are more than N args.
func MaximumNArgs(n int) Args {
	return Args(cobra.MaximumNArgs(n))
}

// ExactArgs returns an error if there are not exactly n args.
func ExactArgs(n int) Args {
	return Args(cobra.ExactArgs(n))
}

// ExactValidArgs returns an error if
// there are not exactly N positional args OR
// there are any positional args that are not in the `ValidArgs` field of `Command`
func ExactValidArgs(n int) Args {
	return Args(cobra.ExactValidArgs(n))
}

// RangeArgs returns an error if the number of args is not within the expected range.
func RangeArgs(min int, max int) Args {
	return Args(cobra.RangeArgs(min, max))
}

type Args cobra.PositionalArgs

func (a Args) Apply(cmd *cobra.Command) {
	cmd.Args = cobra.PositionalArgs(a)
}

func Description(value string) description {
	return description(strings.TrimSpace(dedent.Dedent(value)))
}

type description string

func (d description) Apply(cmd *cobra.Command) {
	cmd.Long = string(d)
}

func Group(usage, short string, opts ...CommandOption) CommandOption {
	return CommandOptionFunc(func(parent *cobra.Command) {
		parent.AddCommand(command(nil, usage, short, opts...))
	})
}

func Command(execute func(cmd *cobra.Command, args []string) error, usage, short string, opts ...CommandOption) CommandOption {
	return CommandOptionFunc(func(parent *cobra.Command) {
		parent.AddCommand(command(execute, usage, short, opts...))
	})
}

type BeforeAllHook func(cmd *cobra.Command)

func (f BeforeAllHook) Apply(cmd *cobra.Command) {
	f(cmd)
}

type AfterAllHook func(cmd *cobra.Command)

func (f AfterAllHook) Apply(cmd *cobra.Command) {
	f(cmd)
}

func Execute(f func(cmd *cobra.Command, args []string) error) execute {
	return execute(f)
}

type execute func(cmd *cobra.Command, args []string) error

func (e execute) Apply(cmd *cobra.Command) {
	cmd.RunE = (func(cmd *cobra.Command, args []string) error)(e)
}

func ExamplePrefixed(prefix string, examples string) example {
	return prefixedExample("  "+prefix+" ", examples)
}

func Example(value string) example {
	return prefixedExample("  ", value)
}

// ConfigureViper installs an [AfterAllHook] on the [cobra.Command]
// that rebind all your flags into viper with a new layout.
//
// Persistent flags on a root command can be accessed with `global-<flag>`
// Persistent flags on a sub-command can be accessed with `<cmd1>-<cmd2>-global-<flag>` where `<cmd1>-<cmd2>` is the command fully qualified path (see below for more details).
// Standard flags on a command can be accessed with `<cmd1>-<cmd2>-<flag>` where `<cmd1>-<cmd2>` is the command fully qualified path (see below for more details).
//
// For the following config:
//
//	Root("acme", "CLI sample application",
//		PersistentFlags(func(flags *pflag.FlagSet) { flags.String("auth", "", "Auth token") }),
//
//		Group("tools", "Tools for developers",
//			PersistentFlags(func(flags *pflag.FlagSet) { flags.Bool("dev", false, "Dev mode") }),
//
//			Command(toolsReadE, "read",
//				"Read command",
//				Flags(func(flags *pflag.FlagSet) { flags.Bool("skip-errors", false, "Skip read errors") }),
//			),
//
//			Command(toolsWriteE, "write",
//				"Write command",
//				Flags(func(flags *pflag.FlagSet) { flags.Bool("skip-errors", false, "Skip write errors") })),
//		),
//	)
//
// Which renders the follow CLI hierarchy of commands:
//
//	acme (--auth <auth>)
//	 tools (--dev)
//	   read (--skip-errors)
//	   write (--skip-errors)
//
// You can access the flags using [viper] sub-commands through this hierarchy:
//
//	viper.GetString("global-auth")             // acme --auth ""
//	viper.GetBool("tools-global-dev")          // acme tools --dev
//	viper.GetString("tools-read-skip-errors")  // acme tools read --skip-errors
//	viper.GetString("tools-write-skip-errors") // acme tools write --skip-errors
//
// And also with environment variables:
//
//	viper.GetString("global-auth")             // {PREFIX}_GLOBAL_AUTH=<auth> acme
//	viper.GetBool("tools-global-dev")          // {PREFIX}_TOOLS_GLOBAL_DEV=<dev> acme acme tools
//	viper.GetString("tools-read-skip-errors")  // {PREFIX}_TOOLS_READ_SKIP_ERRORS=<skip> acme tools read
//	viper.GetString("tools-write-skip-errors") // {PREFIX}_TOOLS_WRITE_SKIP_ERRORS=<skip> acme tools write
//
// Priority is:
// Flag definition via CLI overrides everyone else (--<flag>)
// Environment overrides values provided by config file or defaults ({PREFIX}_{ENV_KEY})
// Config file (if configure separately) overrides defaults values (if configure separately) (--<flag>)
// Defaults values defined on the flag definition directly
func ConfigureViper(envPrefix string) CommandOption {
	return AfterAllHook(func(cmd *cobra.Command) {
		ConfigureViperForCommand(cmd, envPrefix)
	})
}

// ConfigureVersion is an option that configures the `cobra.Command#Version` field
// automatically based fetching commit revision and date build from Golang available
// built-in build info.
//
// The version formatted output can take different forms depending on the state of
// 'vcs.revision' availability, 'vcs.time' availability and received 'version':
//
//	if vcs.revision == "" && vcs.time == "" return "{version}"
//	if vcs.revision != "" && vcs.time == "" return "{version} (Commit {vcs.revision[0:7]})"
//	if vcs.revision == "" && vcs.time == "" return "{version} (Built {vcs.date})"
//	if vcs.revision != "" && vcs.time != "" return "{version} (Commit {vcs.revision[0:7]}, Built {vcs.date})"
func ConfigureVersion(version string) CommandOption {
	return CommandOptionFunc(func(cmd *cobra.Command) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			panic("we should have been able to retrieve info from 'runtime/debug#ReadBuildInfo'")
		}

		commit := findSetting("vcs.revision", info.Settings)
		date := findSetting("vcs.time", info.Settings)

		var labels []string
		if len(commit) >= 7 {
			labels = append(labels, fmt.Sprintf("Commit %s", commit[0:7]))
		}

		if date != "" {
			labels = append(labels, fmt.Sprintf("Built %s", date))
		}

		if len(labels) == 0 {
			cmd.Version = version
		} else {
			cmd.Version = fmt.Sprintf("%s (%s)", version, strings.Join(labels, ", "))
		}
	})
}

func findSetting(key string, settings []debug.BuildSetting) (value string) {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}

	return ""
}

var isSpaceOnlyRegex = regexp.MustCompile(`^\s*$`)
var isCommentExampleLineRegex = regexp.MustCompile(`^\s*#`)

func prefixedExample(prefix string, value string) example {
	content := strings.TrimSpace(dedent.Dedent(value))
	lines := strings.Split(content, "\n")

	if len(lines) == 0 || len(lines) == 1 {
		// No line separator found, in all cases 'content' is our only "line"
		if isSpaceOnlyRegex.MatchString(content) {
			return example(content)
		}

		return example(prefix + content)
	}

	formatted := make([]string, len(lines))
	for i, line := range lines {
		switch {
		case isSpaceOnlyRegex.MatchString(line):
			formatted[i] = line
		case isCommentExampleLineRegex.MatchString(line):
			formatted[i] = "  " + line
		default:
			formatted[i] = prefix + line
		}
	}

	return example(strings.Join(formatted, "\n"))
}

type example string

func (e example) Apply(cmd *cobra.Command) {
	cmd.Example = string(e)
}

func Run(usage, short string, opts ...CommandOption) {
	cmd := Root(usage, short, opts...)

	visitAllCommands(cmd, func(cmd *cobra.Command) {
		if cmd.RunE != nil {
			cmd.RunE = silenceUsageOnError(cmd.RunE)
			cmd.SilenceUsage = false
		}
	})

	err := cmd.Execute()

	// FIXME: What is the right behavior on error from here?
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func visitAllCommands(cmd *cobra.Command, onCmd func(iterated *cobra.Command)) {
	onCmd(cmd)
	for _, subCommand := range cmd.Commands() {
		visitAllCommands(subCommand, onCmd)
	}
}

func Root(usage, short string, opts ...CommandOption) *cobra.Command {
	beforeAllHook := BeforeAllHook(func(cmd *cobra.Command) {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		if short != "" {
			cmd.Short = strings.TrimSpace(dedent.Dedent(short))
		}
	})

	return command(nil, usage, short, append([]CommandOption{beforeAllHook}, opts...)...)
}

func command(execute func(cmd *cobra.Command, args []string) error, usage, short string, opts ...CommandOption) *cobra.Command {
	command := &cobra.Command{}

	for _, opt := range opts {
		if _, ok := opt.(BeforeAllHook); ok {
			opt.Apply(command)
		}
	}

	command.Use = usage
	command.Short = short
	command.RunE = execute

	for _, opt := range opts {
		switch opt.(type) {
		case BeforeAllHook, AfterAllHook:
			continue
		default:
			opt.Apply(command)
		}
	}

	for _, opt := range opts {
		if _, ok := opt.(AfterAllHook); ok {
			opt.Apply(command)
		}
	}

	return command
}

func Dedent(input string) string {
	return strings.TrimSpace(dedent.Dedent(input))
}

// FlagDescription accepts a multi line indented description and transform it into a single line flag description.
// This method is used to make it easier to define long flag messages.
func FlagDescription(in string, args ...interface{}) string {
	return fmt.Sprintf(strings.Join(strings.Split(string(Description(in)), "\n"), " "), args...)
}

type CommandOnErrorHandler func(err error)

func OnError(onError func(err error)) CommandOption {
	return AfterAllHook(func(cmd *cobra.Command) {
		handler := CommandOnErrorHandler(onError)

		visitAllCommands(cmd, func(iterated *cobra.Command) {
			setCommandAnnotation(iterated, annotationOnError, handler)
		})
	})
}

// silenceUsageOnError performs a little trick so that error coming out of the actual executor
// does not trigger a rendering of the usage. By default, cobra prints the usage on flag/args error
// as well as on error coming form the executor (from the `cobra.Command#RunE` field).
//
// That is bad default behavior as in almost all cases, error coming from the executor are not
// usage error.
//
// The trick is to intercept the executor error, and if non-nil, before returning the actual
// error to cobra, we set `cmd.SilenceUsage = true` if the error is non-nil, which will
// properly avoid printing the usage.
func silenceUsageOnError(fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := fn(cmd, args)
		if err != nil {
			cmd.SilenceUsage = true

			onError, found := getCommandAnnotation(cmd, annotationOnError)
			if found {
				onError.(CommandOnErrorHandler)(err)
			}
		}

		return err
	}
}
