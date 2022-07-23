package state

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrNotFound describes an error in which the state file is not found in
// the current dir, or any of it's parents.
var ErrNotFound = errors.New("no state file found")

var ErrInvalidPath = errors.New("invalid path")

// Manager provides an interface that is able to load and save the state file.
type Manager struct {
	fs FileSystem
}

// NewManager acts as the default constructor for the manager instance.
// This method should be used over direct instantiation.
func NewManager(opts ...Option) *Manager {
	m := &Manager{
		fs: &defaultFileSystem{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Save will write the state struct to the file it was loaded from.
// If there is no directory associated with the state file, it will be
// saved in the current working directory.
func (m *Manager) Save(state State) error {
	path, err := m.obtainStatePath()
	if err != nil {
		return fmt.Errorf("obtaining state path: %w", err)
	}

	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	tmpPath := fmt.Sprintf("%s.tmp", path)

	tmp, err := m.fs.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("creating tmp buffer: %w", err)
	}

	if _, err = tmp.Write(data); err != nil {
		m.fs.Remove(tmpPath)
		return fmt.Errorf("writing to buffer: %w", err)
	}

	if err = tmp.Close(); err != nil {
		m.fs.Remove(tmpPath)
		return fmt.Errorf("closing tmp buffer: %w", err)
	}

	if err = m.fs.Rename(tmpPath, path); err != nil {
		m.fs.Remove(tmpPath)
		return fmt.Errorf("renaming tmp buffer: %w", err)
	}

	return nil
}

// Load attempts to load a state file from the project if it exists. It will
// look up the directory tree attempting to find a state file. If no
// state file exists, then the ErrNotFound error will be thrown. If any other
// errors occur, then a different error will be returned.
func (m *Manager) Load() (State, error) {
	path, err := m.findStatePath()
	if err != nil {
		return State{}, fmt.Errorf("find state path: %w", err)
	}

	data, err := m.fs.ReadFile(path)
	if err != nil {
		return State{}, fmt.Errorf("reading .docula: %w", err)
	}

	var res State

	if err = yaml.Unmarshal(data, &res); err != nil {
		return State{}, fmt.Errorf("unmarshal docula state file: %w", err)
	}

	return res, nil
}

func (m *Manager) findStatePath() (string, error) {
	cwd, err := m.fs.Getwd()
	if err != nil {
		return "", fmt.Errorf("get wd: %w", err)
	}

	return m.getPathFrom(cwd)
}

func (m *Manager) getPathFrom(path string) (string, error) {
	if len(path) == 0 {
		return "", ErrNotFound
	}

	filePath := fmt.Sprintf("%s/.docula", strings.TrimSuffix(path, "/"))

	_, err := m.fs.Stat(filePath)

	switch {
	case errors.Is(err, os.ErrNotExist):
		return m.getPathFrom(parentPath(path))
	case err != nil:
		return "", fmt.Errorf("checking file: %w", err)
	}

	return filePath, nil
}

func parentPath(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}

func (m *Manager) obtainStatePath() (string, error) {
	cwd, err := m.fs.Getwd()
	if err != nil {
		return "", fmt.Errorf("get wd: %w", err)
	}

	path, err := m.getPathFrom(cwd)
	if errors.Is(err, ErrNotFound) {
		return fmt.Sprintf("%s/.docula", strings.TrimSuffix(cwd, "/")), nil
	} else if err != nil {
		return "", fmt.Errorf("get path from cwd: %w", err)
	}

	return path, nil
}

// NormalizePath will convert a path into a relative path from the state
// files location.
func (m *Manager) NormalizePath(path string) (string, error) {
	filePath, err := m.obtainStatePath()
	if err != nil {
		return "", fmt.Errorf("obtain state path: %w", err)
	}

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("%s/%s", filePath, path)
	}

	comps := strings.Split(path, "/")

	resolved := make([]string, 0, len(comps))

	for _, x := range comps {
		switch x {
		case ".":
			continue
		case "..":
			resolved = resolved[:len(resolved)-1]
		default:
			resolved = append(resolved, x)
		}
	}

	resolvedPath := strings.Join(resolved, "/")

	if !strings.HasPrefix(resolvedPath, filePath) {
		return "", ErrInvalidPath
	}

	return strings.Replace(resolvedPath, filePath+"/", "", 1), nil
}

// StateDir returns the dir of the state file if it exists. If not state
// file does exist, then the current working directory will be returned.
func (m *Manager) StateDir() (string, error) {
	p, err := m.obtainStatePath()
	if err != nil {
		return "", fmt.Errorf("obtain state path: %w", err)
	}

	return strings.TrimSuffix(p, ".docula"), nil
}
