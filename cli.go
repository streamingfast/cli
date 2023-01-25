package cli

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streamingfast/logging"
	"golang.org/x/crypto/ssh/terminal"
)

var zlog, _ = logging.PackageLogger("cli", "github.com/streamingfast/cli")

func CopyFile(inPath, outPath string) {
	inFile, err := os.Open(inPath)
	NoError(err, "Unable to open actual file %q", inPath)
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	NoError(err, "Unable to open expected file %q", outPath)
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	NoError(err, "Unable to copy file %q to %q", inPath, outPath)
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		// For this script, we don't care
		return false
	}

	return !stat.IsDir()
}

func DirectoryExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		// For this script, we don't care
		return false
	}

	return stat.IsDir()
}

func Ensure(condition bool, message string, args ...interface{}) {
	if !condition {
		Quit(message, args...)
	}
}

func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func Quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}

type CommandOption interface {
	apply(cmd *cobra.Command)
}

type CommandOptionFunc func(cmd *cobra.Command)

func (f CommandOptionFunc) apply(cmd *cobra.Command) {
	f(cmd)
}

type Flags func(flags *pflag.FlagSet)

func (f Flags) apply(cmd *cobra.Command) {
	f(cmd.Flags())
}

type PersistentFlags func(flags *pflag.FlagSet)

func (f PersistentFlags) apply(cmd *cobra.Command) {
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

func (a Args) apply(cmd *cobra.Command) {
	cmd.Args = cobra.PositionalArgs(a)
}

func Description(value string) description {
	return description(strings.TrimSpace(dedent.Dedent(value)))
}

type description string

func (d description) apply(cmd *cobra.Command) {
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

func (f BeforeAllHook) apply(cmd *cobra.Command) {
	f(cmd)
}

type AfterAllHook func(cmd *cobra.Command)

func (f AfterAllHook) apply(cmd *cobra.Command) {
	f(cmd)
}

func Execute(f func(cmd *cobra.Command, args []string) error) execute {
	return execute(f)
}

type execute func(cmd *cobra.Command, args []string) error

func (e execute) apply(cmd *cobra.Command) {
	cmd.RunE = (func(cmd *cobra.Command, args []string) error)(e)
}

func ExamplePrefixed(prefix string, examples string) example {
	return prefixedExample("  "+prefix+" ", examples)
}

func Example(value string) example {
	return prefixedExample("  ", value)
}

func ConfigureViper(prefix string) CommandOption {
	return AfterAllHook(func(cmd *cobra.Command) {
		ConfigureViperForCommand(cmd, prefix)
	})
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

func (e example) apply(cmd *cobra.Command) {
	cmd.Example = string(e)
}

func Run(usage, short string, opts ...CommandOption) {
	cmd := Root(usage, short, opts...)

	cmd.SilenceUsage = false
	cmd.RunE = silenceUsageOnError(cmd.RunE)

	err := cmd.Execute()

	// FIXME: What is the right behavior on error from here?
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// silenceUsageOnError performs a little trick so that error coming out of the actual executor
// does not trigger a rendering of the usage. By default, cobra prints the usage on flag/args error
// as well as on error coming form the executor.
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
		}

		return err
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
			opt.apply(command)
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
			opt.apply(command)
		}
	}

	for _, opt := range opts {
		if _, ok := opt.(AfterAllHook); ok {
			opt.apply(command)
		}
	}

	return command
}

func Dedent(input string) string {
	return strings.TrimSpace(dedent.Dedent(input))
}

func AskConfirmation(label string, args ...interface{}) (answeredYes bool, wasAnswered bool) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		wasAnswered = false
		return
	}

	prompt := promptui.Prompt{
		Label:       dedent.Dedent(fmt.Sprintf(label, args...)),
		Default:     "N",
		AllowEdit:   true,
		IsConfirm:   true,
		HideEntered: true,
	}

	_, err := prompt.Run()
	if err != nil {
		// zlog.Debug("unable to aks user to see diff right now, too bad", zap.Error(err))
		wasAnswered = false
		return
	}

	wasAnswered = true
	answeredYes = true

	return
}

// FlagDescription accepts a multi line indented description and transform it into a single line flag description.
// This method is used to make it easier to define long flag messages.
func FlagDescription(in string, args ...interface{}) string {
	return fmt.Sprintf(strings.Join(strings.Split(string(Description(in)), "\n"), " "), args...)
}
