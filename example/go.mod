module github.com/dfuse-io/cli/example

go 1.15

replace github.com/dfuse-io/cli => ../

require (
	github.com/dfuse-io/cli v0.0.0-00010101000000-000000000000
	github.com/dfuse-io/logging v0.0.0-20210518215502-2d920b2ad1f2
	github.com/spf13/cobra v1.1.3
	go.uber.org/zap v1.16.0
)
