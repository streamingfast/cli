package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prefixedExample(t *testing.T) {

	tests := []struct {
		name   string
		result example
		want   string
	}{
		{
			"empty",
			ExamplePrefixed("bin", ``),
			``,
		},
		{
			"non-empty but only spaces",
			ExamplePrefixed("bin", `       `),
			``,
		},
		{
			"single one",
			ExamplePrefixed("bin", `
				cmd1
			`),
			"  bin cmd1",
		},
		{
			"multi one",
			ExamplePrefixed("bin", `
				cmd1
				cmd2
			`),
			"  bin cmd1\n  bin cmd2",
		},
		{
			"comment not-prefixed",
			ExamplePrefixed("bin", `
				# Comment 1
				cmd1

				# Comment 2
				cmd2
			`),
			"  # Comment 1\n  bin cmd1\n\n  # Comment 2\n  bin cmd2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.result))
		})
	}
}
