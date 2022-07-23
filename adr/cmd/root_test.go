package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docula-io/docula/adr/cmd"
)

func TestRootCommand(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		wantsErr bool
	}{
		{
			name: "should have an init command",
			args: []string{
				"init", "--help",
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
			c := cmd.RootCmd()

			c.SetArgs(tt.args)

			err := c.Execute()

			if tt.wantsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
