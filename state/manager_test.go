package state_test

import (
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/docula-io/docula/adr"
	"github.com/docula-io/docula/state"
)

func setupFs(path string) func(ctrl *gomock.Controller) state.FileSystem {
	return func(ctrl *gomock.Controller) state.FileSystem {
		fs := state.NewmockFileSystem(ctrl)

		fs.EXPECT().Getwd().Return(path, nil)
		fs.EXPECT().Stat(path+"/.docula").Return(nil, nil)

		return fs
	}
}

func TestManagerSave(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(ctrl *gomock.Controller) state.FileSystem
		input state.State
		wants error
	}{
		{
			name: "happy path",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				f := state.NewmockFile(ctrl)
				fs.EXPECT().Create("/.docula.tmp").Return(f, nil)

				expects := "adr:\n    dirs: []\n"

				gomock.InOrder(
					f.EXPECT().Write([]byte(expects)).Return(0, nil),
					f.EXPECT().Close(),
					fs.EXPECT().Rename("/.docula.tmp", "/.docula").Return(nil),
				)
				return fs
			},
			input: state.State{},
			wants: nil,
		},
		{
			name: "with state",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				f := state.NewmockFile(ctrl)
				fs.EXPECT().Create("/.docula.tmp").Return(f, nil)

				expects := "adr:\n    dirs: []\n"

				gomock.InOrder(
					f.EXPECT().Write([]byte(expects)).Return(0, nil),
					f.EXPECT().Close(),
					fs.EXPECT().Rename("/.docula.tmp", "/.docula").Return(nil),
				)

				return fs
			},
			input: state.State{},
			wants: nil,
		},
		{
			name: "failing to create tmp buffer",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				fs.EXPECT().Create("/.docula.tmp").Return(nil, os.ErrInvalid)

				return fs
			},
			input: state.State{},
			wants: os.ErrInvalid,
		},
		{
			name: "failed to write to tmp buffer",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				f := state.NewmockFile(ctrl)
				fs.EXPECT().Create("/.docula.tmp").Return(f, nil)

				expects := "adr:\n    dirs: []\n"

				gomock.InOrder(
					f.EXPECT().Write([]byte(expects)).Return(0, os.ErrInvalid),
					fs.EXPECT().Remove("/.docula.tmp"),
				)

				return fs
			},
			input: state.State{},
			wants: os.ErrInvalid,
		},
		{
			name: "failed to close the tmp buffer",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				f := state.NewmockFile(ctrl)
				fs.EXPECT().Create("/.docula.tmp").Return(f, nil)

				expects := "adr:\n    dirs: []\n"

				gomock.InOrder(
					f.EXPECT().Write([]byte(expects)).Return(0, nil),
					f.EXPECT().Close().Return(os.ErrInvalid),
					fs.EXPECT().Remove("/.docula.tmp"),
				)

				return fs
			},
			input: state.State{},
			wants: os.ErrInvalid,
		},
		{
			name: "failed to rename the tmp buffer",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)

				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, nil)

				f := state.NewmockFile(ctrl)
				fs.EXPECT().Create("/.docula.tmp").Return(f, nil)

				expects := "adr:\n    dirs: []\n"

				gomock.InOrder(
					f.EXPECT().Write([]byte(expects)).Return(0, nil),
					f.EXPECT().Close().Return(nil),
					fs.EXPECT().Rename("/.docula.tmp", "/.docula").Return(os.ErrInvalid),
					fs.EXPECT().Remove("/.docula.tmp"),
				)

				return fs
			},
			input: state.State{},
			wants: os.ErrInvalid,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fs := tt.setup(ctrl)
			manager := state.NewManager(state.WithFileSystem(fs))

			err := manager.Save(tt.input)
			if tt.wants != nil {
				assert.ErrorIs(t, err, tt.wants)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManagerLoad(t *testing.T) {
	const stateFile = `adr:
  dirs:
    - path: foo
    - path: bar
`

	type want struct {
		err    error
		anyErr bool
		state  state.State
	}

	testCases := []struct {
		name  string
		setup func(ctrl *gomock.Controller) state.FileSystem
		wants want
	}{
		{
			name: "happy path",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/foo/bar/boo", nil)
				fs.EXPECT().Stat("/foo/bar/boo/.docula").Return(nil, nil)
				fs.EXPECT().ReadFile("/foo/bar/boo/.docula").Return([]byte(stateFile), nil)

				return fs
			},
			wants: want{
				state: state.State{
					ADR: adr.State{
						Directories: []adr.Directory{
							{
								Path: "foo",
							},
							{
								Path: "bar",
							},
						},
					},
				},
			},
		},
		{
			name: "state file in parent dir",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/foo/bar/boo", nil)

				gomock.InOrder(
					fs.EXPECT().Stat("/foo/bar/boo/.docula").Return(nil, os.ErrNotExist),
					fs.EXPECT().Stat("/foo/bar/.docula").Return(nil, os.ErrNotExist),
					fs.EXPECT().Stat("/foo/.docula").Return(nil, nil),
				)
				fs.EXPECT().ReadFile("/foo/.docula").Return([]byte(stateFile), nil)

				return fs
			},
			wants: want{
				state: state.State{
					ADR: adr.State{
						Directories: []adr.Directory{
							{
								Path: "foo",
							},
							{
								Path: "bar",
							},
						},
					},
				},
			},
		},
		{
			name: "fail to get current working dir",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("", os.ErrInvalid)

				return fs
			},
			wants: want{
				err: os.ErrInvalid,
			},
		},
		{
			name: "fail to get read file",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/foo/bar/boo", nil)
				fs.EXPECT().Stat("/foo/bar/boo/.docula").Return(nil, nil)
				fs.EXPECT().ReadFile("/foo/bar/boo/.docula").Return(nil, os.ErrInvalid)

				return fs
			},
			wants: want{
				err: os.ErrInvalid,
			},
		},
		{
			name: "bad yaml",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/foo/bar/boo", nil)
				fs.EXPECT().Stat("/foo/bar/boo/.docula").Return(nil, nil)
				fs.EXPECT().ReadFile("/foo/bar/boo/.docula").Return([]byte("adr: foo"), nil)

				return fs
			},
			wants: want{
				anyErr: true,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fs := tt.setup(ctrl)

			manager := state.NewManager(state.WithFileSystem(fs))

			s, err := manager.Load()

			if tt.wants.anyErr {
				assert.Error(t, err)
			} else if tt.wants.err != nil {
				assert.ErrorIs(t, err, tt.wants.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wants.state, s)
		})
	}
}

func TestManagerNormalizePath(t *testing.T) {
	type want struct {
		path string
		err  bool
	}

	testCases := []struct {
		name  string
		input string
		setup func(ctrl *gomock.Controller) state.FileSystem
		wants want
	}{
		{
			name:  "redudnant parent",
			input: "./foo/../foo",
			setup: setupFs("/home/foo/bar"),
			wants: want{
				path: "foo",
			},
		},
		{
			name:  "parent at the start",
			input: "../foo/../foo",
			setup: setupFs("/home/bar"),
			wants: want{
				err: true,
			},
		},
		{
			name:  "current dir",
			input: "./",
			setup: setupFs("/home/docula"),
			wants: want{
				path: "",
			},
		},
		{
			name:  "standard subdir",
			input: "foo/bar/baz",
			setup: setupFs("/home/docula"),
			wants: want{
				path: "foo/bar/baz",
			},
		},
		{
			name:  "missing .docula",
			input: "foo/bar/baz",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, os.ErrNotExist)
				return fs
			},
			wants: want{
				path: "foo/bar/baz",
			},
		},
		{
			name:  "err loading file stat",
			input: "foo/bar/baz",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("/", nil)
				fs.EXPECT().Stat("/.docula").Return(nil, os.ErrInvalid)
				return fs
			},
			wants: want{
				err: true,
			},
		},
		{
			name:  "err getting cwd",
			input: "foo/bar/baz",
			setup: func(ctrl *gomock.Controller) state.FileSystem {
				fs := state.NewmockFileSystem(ctrl)
				fs.EXPECT().Getwd().Return("", os.ErrInvalid)
				return fs
			},
			wants: want{
				err: true,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fs := tt.setup(ctrl)
			manager := state.NewManager(state.WithFileSystem(fs))

			res, err := manager.NormalizePath(tt.input)
			if tt.wants.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wants.path, res)
		})
	}
}
