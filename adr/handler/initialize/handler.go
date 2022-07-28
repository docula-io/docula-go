package initialize

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docula-io/docula/adr"
	"github.com/docula-io/docula/state"
)

var ErrAlreadyIntialized = errors.New("adr dir already initialized")

// Handler describes a type that is used to handle the initialize command.
type Handler struct {
	stateManager StateManager
	fs           FileSystem
	survey       Survey
}

// New acts as the default constructor for the Handler type. This method
// will initialize defaults for the internal resources, or will override them
// with any provided options. This method should be used instead of direct
// instantiation.
func New(opts ...Option) *Handler {
	h := &Handler{
		stateManager: state.NewManager(),
		fs:           &defaultFileSystem{},
		survey:       &defaultSurvey{},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Configuration represents a type that stores additional configuration for
// intializing an adr directory.
type Configuration struct {
	Name      string `survey:"name"`
	IndexType string `survey:"index"`
}

func (h *Handler) runSurvey() (Configuration, error) {
	answers, err := h.survey.Ask()
	if err != nil {
		return Configuration{}, err
	}

	return answers, nil
}

func (h *Handler) checkExistingADRs(s state.State, dir adr.Directory) error {
	// Check path is not already an ADR dir
	for _, adrDir := range s.ADR.Directories {
		if adrDir.Path == dir.Path {
			return ErrAlreadyIntialized
		}
	}

	return nil
}

func (h *Handler) createDir(path string) error {
	stateDir, err := h.stateManager.StateDir()
	if err != nil {
		return fmt.Errorf("obtain state path: %w", err)
	}

	absPath := fmt.Sprintf("%s%s", stateDir, path)

	// Create path if not exists
	if err = h.fs.Mkdir(absPath); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create adr dir: %w", err)
	}

	return nil
}

// Handle is the main Handler function. This function is used to initialize
// a new directory as an adr dir.
func (h *Handler) Handle(ctx context.Context, path string) error {
	path, err := h.stateManager.NormalizePath(path)
	if err != nil {
		return fmt.Errorf("normalize path: %w", err)
	}

	// Load state
	s, err := h.stateManager.Load()
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return fmt.Errorf("loading state: %w", err)
	}

	config, err := h.runSurvey()
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	if err = h.createDir(path); err != nil {
		return err
	}

	dir := adr.Directory{
		Path:  path,
		Name:  config.Name,
		Index: config.IndexType,
	}

	if err = h.checkExistingADRs(s, dir); err != nil {
		return err
	}

	// Update the ADR part
	s.ADR.Directories = append(s.ADR.Directories, dir)

	// Write the file
	if err = h.stateManager.Save(s); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	return nil
}
