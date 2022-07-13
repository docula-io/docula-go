package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	errBoom := errors.New("boom")

	type want struct {
		err  error
		path string
	}

	testCases := []struct {
		name       string
		handlerRet error
		args       []string
		wants      want
	}{
		{
			name:       "happy path",
			handlerRet: nil,
			args:       []string{"./test"},
			wants: want{
				path: "./test",
			},
		},
		{
			name:       "bad args",
			handlerRet: nil,
			args:       []string{},
			wants: want{
				err: errors.New(""),
			},
		},
		{
			name:       "handler error",
			handlerRet: errBoom,
			args:       []string{"./mydir"},
			wants: want{
				err:  errBoom,
				path: "./mydir",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			h := func(ctx context.Context, path string) error {
				return tt.handlerRet
			}

			cmd := initCmd(h)

			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.wants.err != nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wants.err)
			}
		})
	}
}
