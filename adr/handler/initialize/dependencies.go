//go:generate mockgen -source=dependencies.go -destination=./mocks.go -package=initialize -mock_names FileSystem=mockFileSystem,StateManager=mockStateManager,Survey=mockSurvey

package initialize

import (
	survey "github.com/AlecAivazis/survey/v2"

	"github.com/docula-io/docula/state"
)

// StateManager represents a type that is able to manage the docula state file.
type StateManager interface {
	Load() (state.State, error)
	NormalizePath(path string) (string, error)
	Save(state.State) error
	StateDir() (string, error)
}

// FileSystem represents a type that is able to manipulate the filesystem.
// This interface is typically a wrapper around the os package methods and
// is used to allow for improved testing.
type FileSystem interface {
	Mkdir(name string) error
}

// Survey represents a type that is able to get various inputs from stdin.
type Survey interface {
	Ask(opts ...survey.AskOpt) (Configuration, error)
}
