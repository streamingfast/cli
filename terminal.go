package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// Prompt will ask the following "label" question to the user and will transform the received `string`
// into the generic type T via the `transformer` function. There is a few predefined transformer for
// common types like `PromptTypeString`, `PromptTypeInt`, `PromptTypeInt64`, `PromptTypeUint64`,
// `PromptTypeYesNo`. Using the right transformer will automatically infers the right T genric type.
//
//	userID := cli.Prompt("Please enter the user ID to issue the token to", cli.PromptTypeString)
//
// In this case, the `userID` variable will be of type `string`.
//
// If 'label' is a multiple line string, it will be split into separate lines
// and each line will be rendered with the appropriate style but once answered,
// only the final line changes, the previous written lines are not affected.
//
// This is why that rendering wise, favoring single line for 'label' will offer
// a better user experience.
func Prompt[T any](label string, transformer PromptTransformer[T], opts ...PromptOption) T {
	out, err := MaybePrompt(label, transformer, opts...)
	if err != nil {
		panic(fmt.Errorf("prompt failed: %w", err))
	}

	return out
}

// MaybePrompt is just like [Prompt] but may return an error
func MaybePrompt[T any](label string, transformer PromptTransformer[T], opts ...PromptOption) (T, error) {
	choice, err := PromptRaw(label, opts...)
	if err == nil {
		return transformer(choice)
	}

	var empty T
	return empty, err
}

// PromptConfirm is just like [Prompt] but enforce `IsConfirm` and returns a boolean which is either
// `true` for yes answer or `false` for a no answer.
func PromptConfirm(label string, opts ...PromptOption) (answer bool, wasAnswered bool) {
	answer, wasAnswered, err := MaybePromptConfirm(label, opts...)
	if err != nil {
		panic(fmt.Errorf("prompt confirm failed: %w", err))
	}

	return answer, wasAnswered
}

// MaybePromptConfirm is just like [PromptConfirm] but returns an error instead of panicking.
func MaybePromptConfirm(label string, opts ...PromptOption) (answer bool, wasAnswered bool, err error) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		wasAnswered = false
		return
	}

	opts = append([]PromptOption{WithPromptValidate("invalid", PrompValidateYesNo), WithPromptConfirm()}, opts...)

	answer, err = MaybePrompt(label, PromptTypeYesNo, opts...)
	return answer, err == nil, err
}

func PromptSelect[T any](label string, items []string, transformer PromptTransformer[T], opts ...PromptSelectOption) T {
	out, err := MaybePromptSelect(label, items, transformer, opts...)
	if err != nil {
		panic(fmt.Errorf("prompt select failed: %w", err))
	}

	return out
}

func MaybePromptSelect[T any](label string, items []string, transformer PromptTransformer[T], opts ...PromptSelectOption) (T, error) {
	options := promptSelectOptions{}
	for _, opt := range opts {
		opt.Apply(&options)
	}

	choice := promptui.Select{
		Label:     label,
		Items:     items,
		HideHelp:  true,
		Templates: options.selectTemplates,
	}

	_, selection, err := choice.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			// We received Ctrl-C, users wants to abort, nothing else to do, quit immediately
			Exit(1)
		}

		var empty T
		return empty, fmt.Errorf("running protocol prompt: %w", err)
	}

	return transformer(selection)
}

func splitPromptLines(label string, args ...any) (preLines []string, finalLine string) {
	content := fmt.Sprintf(label, args...)
	lines := slices.Collect(strings.Lines(content))

	sawNonEmptyLine := false
	var finalLines []string
	for i := len(lines) - 1; i >= 0; i-- {
		if !sawNonEmptyLine && strings.TrimSpace(lines[i]) == "" {
			continue
		}

		sawNonEmptyLine = true

		line := strings.Trim(strings.ReplaceAll(lines[i], "\t", "  "), "\n\r")
		finalLines = append(finalLines, line)
	}

	slices.Reverse(finalLines)

	if len(finalLines) == 0 {
		return nil, ""
	}

	return finalLines[:len(finalLines)-1], strings.TrimSpace(finalLines[len(finalLines)-1])
}

// AskConfirmation prompts the user for a yes/no confirmation.
//
// If 'label' is a multiple line string, it will be split into separate lines
// and each line will be rendered with the appropriate style but once answered,
// only the final line changes, the previous written lines are not affected.
//
// This is why that rendering wise, favoring single line for 'label' will offer
// a better user experience.
func AskConfirmation(label string, args ...any) (answeredYes bool, wasAnswered bool) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		wasAnswered = false
		return
	}

	preLines, finalLine := splitPromptLines(label, args...)
	if len(preLines) > 1 {
		for _, line := range preLines {
			fmt.Println(QuestionTextStyle.Render(line))
		}
	}

	prompt := promptui.Prompt{
		Label:       strings.TrimSuffix(finalLine, "?"),
		Default:     "N",
		AllowEdit:   true,
		IsConfirm:   true,
		HideEntered: false,
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

type PromptOption interface {
	Apply(opts *promptOptions)
}

type promptIsDefaultValue string

func (o promptIsDefaultValue) Apply(opts *promptOptions) {
	opts.defaultValue = string(o)
}

func WithPromptDefaultValue(in string) PromptOption {
	return promptIsDefaultValue(in)
}

type promptIsConfirmOption bool

func (o promptIsConfirmOption) Apply(opts *promptOptions) {
	opts.isConfirm = bool(o)
}

func WithPromptConfirm() PromptOption {
	return promptIsConfirmOption(true)
}

type validatePromptOption promptui.ValidateFunc

func (o validatePromptOption) Apply(opts *promptOptions) {
	opts.validate = promptui.ValidateFunc(o)
}

func WithPromptValidate(label string, fn promptui.ValidateFunc) PromptOption {
	return validatePromptOption(func(x string) error {
		err := fn(x)
		if err != nil {
			return fmt.Errorf(label+": %w", err)
		}

		return nil
	})
}

type promptTemplatesOption promptui.PromptTemplates

func (o *promptTemplatesOption) Apply(opts *promptOptions) {
	opts.promptTemplates = (*promptui.PromptTemplates)(o)
}

func WithPromptTemplates(templates *promptui.PromptTemplates) PromptOption {
	return (*promptTemplatesOption)(templates)
}

type promptOptions struct {
	validate        promptui.ValidateFunc
	isConfirm       bool
	promptTemplates *promptui.PromptTemplates
	defaultValue    string
}

func PromptRaw(label string, opts ...PromptOption) (answer string, err error) {
	options := promptOptions{}
	for _, opt := range opts {
		opt.Apply(&options)
	}

	templates := options.promptTemplates

	if templates == nil {
		templates = &promptui.PromptTemplates{
			Success: `{{ . | faint }}{{ ":" | faint}} `,
		}
	}

	if options.isConfirm {
		// We don't have no differences
		templates.Valid = `{{ "?" | blue}} {{ . | bold }} {{ "[y/N]" | faint}} `
		templates.Invalid = templates.Valid
	}

	preLines, finalLine := splitPromptLines(label)
	if len(preLines) > 1 {
		for _, line := range preLines {
			fmt.Println(HeaderStyle.Render(line))
		}
	}

	prompt := promptui.Prompt{
		Label:     finalLine,
		Templates: templates,
	}

	if options.validate != nil {
		prompt.Validate = options.validate
	}

	if prompt.Validate == nil && options.isConfirm {
		prompt.Validate = PrompValidateYesNo
	}

	if options.defaultValue != "" {
		prompt.Default = options.defaultValue
		prompt.AllowEdit = true
	}

	choice, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			Exit(1)
		}

		if prompt.IsConfirm && errors.Is(err, promptui.ErrAbort) {
			return "false", nil
		}

		return "", fmt.Errorf("running prompt: %w", err)
	}

	return choice, nil
}

// Various prompt transformer declared like you are using standard type and generics make the rest
//
// To be used like:
//
//	PromptT("Input your age", opts, PromptTypeUint64)
var (
	PromptTypeOptionalString = func(x string) (*string, error) {
		if x == "" {
			return nil, nil
		}
		return &x, nil
	}
	PromptTypeString = func(x string) (string, error) { return x, nil }
	PromptTypeInt    = strconv.Atoi
	PromptTypeInt64  = func(x string) (int64, error) { return strconv.ParseInt(x, 0, 64) }
	PromptTypeUint64 = func(x string) (uint64, error) { return strconv.ParseUint(x, 0, 64) }
	PromptTypeYesNo  = func(in string) (bool, error) { in = strings.ToLower(in); return in == "y" || in == "yes", nil }
)

// Various prompt validation declared like you are using standard validation.
//
// To be used like:
//
//	cli.PromptT("Input your age", cli.PromptTypeUint64, cli.WithPromptValidate("invalid age", cli.PrompValidateUint64))
var (
	PrompValidateString = func(x string) error { return nil }
	PrompValidateInt    = func(x string) error { _, err := strconv.Atoi(x); return err }
	PrompValidateInt64  = func(x string) error { _, err := strconv.ParseInt(x, 0, 64); return err }
	PrompValidateUint64 = func(x string) error { _, err := strconv.ParseUint(x, 0, 64); return err }

	PrompValidateYesNo = func(in string) error {
		if !confirmPromptRegex.MatchString(in) {
			return errors.New("answer with y/yes/Yes or n/no/No")
		}

		return nil
	}
)

var confirmPromptRegex = regexp.MustCompile("(y|Y|n|N|No|Yes|YES|NO)")

type PromptTransformer[T any] func(string) (T, error)

type promptSelectOptions struct {
	selectTemplates *promptui.SelectTemplates
}

type PromptSelectOption interface {
	Apply(opts *promptSelectOptions)
}

type promptSelectTemplatesOption promptui.SelectTemplates

func (o *promptSelectTemplatesOption) Apply(opts *promptSelectOptions) {
	opts.selectTemplates = (*promptui.SelectTemplates)(o)
}

func WithPromptSelectTemplates(templates *promptui.SelectTemplates) PromptSelectOption {
	return (*promptSelectTemplatesOption)(templates)
}
