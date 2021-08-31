module github.com/streamingfast/cli/example

go 1.15

replace github.com/streamingfast/cli => ../

require (
	github.com/spf13/cobra v1.1.3
	github.com/streamingfast/cli v0.0.0-00010101000000-000000000000
	github.com/streamingfast/logging v0.0.0-20210811175431-f3b44b61606a // indirect
	go.uber.org/zap v1.16.0
)
