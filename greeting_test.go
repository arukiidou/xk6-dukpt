package example

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_module_greeting(t *testing.T) {
	t.Parallel()

	mod := new(module)

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "without arg", arg: "", want: "Hello, World!"},
		{name: "with arg", arg: "Joe", want: "Hello, Joe!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, mod.greeting(tt.arg))
		})
	}
}
