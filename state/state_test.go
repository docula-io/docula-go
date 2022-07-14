package state_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docula-io/docula/adr"
	"github.com/docula-io/docula/state"
)

const exampleFile = `
adr:
  dirs:
  - path: foo/bar
  - path: docs/adr 
`

func TestLoadingState(t *testing.T) {
	type want struct {
		state state.State
		errIs error
		err   bool
	}

	testCases := []struct {
		name  string
		setup func(path string) error
		wants want
	}{
		{
			name:  "no state file in chain",
			setup: os.Chdir,
			wants: want{
				errIs: state.ErrNotFound,
			},
		},
		{
			name: "state file in current dir",
			setup: func(path string) error {
				f, err := os.Create(fmt.Sprintf("%s/.docula", path))
				if err != nil {
					return err
				}

				if _, err = f.Write([]byte(exampleFile)); err != nil {
					return err
				}

				f.Close()

				return nil
			},
			wants: want{
				state: state.State{
					ADR: adr.State{
						Directories: []adr.Directory{
							{
								Path: "foo/bar",
							},
							{
								Path: "docs/adr",
							},
						},
					},
				},
			},
		},
		{
			name: "state file in parent dir",
			setup: func(path string) error {
				p, err := os.MkdirTemp(path, "")
				if err != nil {
					return err
				}

				if err = os.Chdir(p); err != nil {
					return err
				}

				f, err := os.Create(fmt.Sprintf("%s/.docula", path))
				if err != nil {
					return err
				}

				f.Close()

				return nil
			},
		},
		{
			name: "non yaml state file in current dir",
			setup: func(path string) error {
				f, err := os.Create(fmt.Sprintf("%s/.docula", path))
				if err != nil {
					return err
				}

				contents := `foo:foo:[]`

				if _, err = f.Write([]byte(contents)); err != nil {
					return err
				}

				f.Close()

				return nil
			},
			wants: want{
				err: true,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			cwd, err := os.Getwd()
			assert.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Chdir(cwd))
			}()

			tmp, err := os.MkdirTemp("", "docula-state-test-*")
			assert.NoError(t, err)

			defer func() {
				assert.NoError(t, os.RemoveAll(tmp))
			}()

			assert.NoError(t, os.Chdir(tmp))

			assert.NoError(t, tt.setup(tmp))

			s, err := state.Load()

			switch {
			case tt.wants.errIs != nil:
				assert.ErrorIs(t, err, tt.wants.errIs)
			case tt.wants.err:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wants.state, s)
		})
	}
}
