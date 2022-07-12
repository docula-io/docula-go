package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		wantsErr bool
	}{
		{
			name: "should have an adr command",
			args: []string{
				"adr", "help",
			},
		},
		{
			name: "should not have a foobar command",
			args: []string{
				"foobar",
			},
			wantsErr: true,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd()

			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.wantsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
