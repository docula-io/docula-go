package state

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/docula-io/docula/adr"
)

// ErrNotFound describes an error in which the state file is not found in
// the current dir, or any of it's parents.
var ErrNotFound = errors.New("no state file found")

// State represents the docula state file which is associated with a project.
// This file is used to store the state of the docula changes.
type State struct {
	ADR adr.State `yaml:"adr"`
}

// Load attempts to load a state file from the project if it exists. It will
// look up the directory tree attempting to find a state file. If no
// state file exists, then the ErrNotFound error will be thrown. If any other
// errors occur, then a different error will be returned.
func Load() (State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return State{}, fmt.Errorf("get wd: %w", err)
	}

	return load(cwd)
}

func parentPath(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}

func load(path string) (State, error) {
	if len(path) == 0 {
		return State{}, ErrNotFound
	}

	filePath := fmt.Sprintf("%s/.docula", path)

	_, err := os.Stat(filePath)

	switch {
	case errors.Is(err, os.ErrNotExist):
		return load(parentPath(path))
	case err != nil:
		return State{}, fmt.Errorf("checking file: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return State{}, fmt.Errorf("reading .docula: %w", err)
	}

	var res State

	if err = yaml.Unmarshal(data, &res); err != nil {
		return State{}, fmt.Errorf("unmarshal docula state file: %w", err)
	}

	return res, nil
}
