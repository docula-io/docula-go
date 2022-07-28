package initialize_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/docula-io/docula/adr"
	"github.com/docula-io/docula/adr/handler/initialize"
	"github.com/docula-io/docula/state"
)

var defaultConfig = initialize.Configuration{
	Name:      "bar",
	IndexType: "timestamp",
}

func TestHandler(t *testing.T) {
	type setup struct {
		stateManager func(ctrl *gomock.Controller) initialize.StateManager
		fs           func(ctrl *gomock.Controller) initialize.FileSystem
		survey       func(ctrl *gomock.Controller) initialize.Survey
	}

	testCases := []struct {
		name  string
		setup setup
		input string
		wants error
	}{
		{
			name: "with no existing adr dir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("foo/bar").Return("foo/bar", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("/", nil)
					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "foo/bar",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/foo/bar").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "foo/bar",
			wants: nil,
		},
		{
			name: "with other adr dirs",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)

					s.EXPECT().NormalizePath("foo/bar").Return("foo/bar", nil)

					s.EXPECT().Load().Return(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "docs/adrs",
									Name:  "default",
									Index: "sequential",
								},
							},
						},
					}, nil)

					s.EXPECT().StateDir().Return("/", nil)

					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "docs/adrs",
									Name:  "default",
									Index: "sequential",
								},
								{
									Path:  "foo/bar",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/foo/bar").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "foo/bar",
			wants: nil,
		},
		{
			name: "with a complex state dir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("baz/foo").Return("baz/foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("/home/me/projects/", nil)
					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "baz/foo",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/home/me/projects/baz/foo").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "baz/foo",
			wants: nil,
		},
		{
			name: "without loading an existing state",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("baz/foo").Return("baz/foo", nil)
					s.EXPECT().Load().Return(state.State{}, state.ErrNotFound)
					s.EXPECT().StateDir().Return("/home/me/projects/", nil)
					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "baz/foo",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/home/me/projects/baz/foo").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "baz/foo",
			wants: nil,
		},
		{
			name: "with an existing adr dir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("foo/bar").Return("foo/bar", nil)
					s.EXPECT().StateDir().Return("/", nil)
					s.EXPECT().Load().Return(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "foo/bar",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}, nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/foo/bar").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "foo/bar",
			wants: initialize.ErrAlreadyIntialized,
		},
		{
			name: "failing to normalize path",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("", os.ErrInvalid)
					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					return initialize.NewmockFileSystem(ctrl)
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					return initialize.NewmockSurvey(ctrl)
				},
			},
			input: "hello/world",
			wants: os.ErrInvalid,
		},
		{
			name: "failing to load the survey",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(
						initialize.Configuration{}, os.ErrDeadlineExceeded,
					)
					return s
				},
			},
			input: "hello/world",
			wants: os.ErrDeadlineExceeded,
		},
		{
			name: "failing to load the state",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, os.ErrClosed)
					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					return s
				},
			},
			input: "hello/world",
			wants: os.ErrClosed,
		},
		{
			name: "failing to load the state dir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("", os.ErrProcessDone)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "hello/world",
			wants: os.ErrProcessDone,
		},
		{
			name: "failing to mkdir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("/home/user/", nil)
					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/home/user/foo").Return(os.ErrInvalid)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "hello/world",
			wants: os.ErrInvalid,
		},
		{
			name: "trying to mkdir on existing dir",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("/home/user/", nil)
					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "foo",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(nil)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/home/user/foo").Return(os.ErrExist)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "hello/world",
		},
		{
			name: "failing to save state",
			setup: setup{
				stateManager: func(ctrl *gomock.Controller) initialize.StateManager {
					s := initialize.NewmockStateManager(ctrl)
					s.EXPECT().NormalizePath("hello/world").Return("foo", nil)
					s.EXPECT().Load().Return(state.State{}, nil)
					s.EXPECT().StateDir().Return("/home/user/", nil)
					s.EXPECT().Save(state.State{
						ADR: adr.State{
							Directories: []adr.Directory{
								{
									Path:  "foo",
									Name:  "bar",
									Index: "timestamp",
								},
							},
						},
					}).Return(state.ErrInvalidPath)

					return s
				},
				fs: func(ctrl *gomock.Controller) initialize.FileSystem {
					fs := initialize.NewmockFileSystem(ctrl)
					fs.EXPECT().Mkdir("/home/user/foo").Return(nil)
					return fs
				},
				survey: func(ctrl *gomock.Controller) initialize.Survey {
					s := initialize.NewmockSurvey(ctrl)
					s.EXPECT().Ask().Return(defaultConfig, nil)
					return s
				},
			},
			input: "hello/world",
			wants: state.ErrInvalidPath,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := tt.setup.stateManager(ctrl)
			fs := tt.setup.fs(ctrl)
			survey := tt.setup.survey(ctrl)

			h := initialize.New(
				initialize.WithFileSystem(fs),
				initialize.WithStateManager(sm),
				initialize.WithSurvey(survey),
			)

			err := h.Handle(context.Background(), tt.input)

			if tt.wants != nil {
				assert.ErrorIs(t, err, tt.wants)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
