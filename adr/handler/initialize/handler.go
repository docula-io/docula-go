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
}

// New acts as the default constructor for the Handler type. This method
// will initialize defaults for the internal resources, or will override them
// with any provided options. This method should be used instead of direct
// instantiation.
func New(opts ...Option) *Handler {
	h := &Handler{
		stateManager: state.NewManager(),
		fs:           &defaultFileSystem{},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
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

	// Check path is not already an ADR dir
	for _, adrDir := range s.ADR.Directories {
		if adrDir.Path == path {
			return ErrAlreadyIntialized
		}
	}

	stateDir, err := h.stateManager.StateDir()
	if err != nil {
		return fmt.Errorf("obtain state path: %w", err)
	}

	absPath := fmt.Sprintf("%s%s", stateDir, path)

	// Create path if not exists
	if err = h.fs.Mkdir(absPath, 0644); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create adr dir: %w", err)
	}

	// Update the ADR part
	s.ADR.Directories = append(s.ADR.Directories, adr.Directory{
		Path: path,
	})

	// Write the file
	if err = h.stateManager.Save(s); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	return nil
}
