package initialize

// Option represents a type that is able to override the default resources of
// the handler. These options are mainly used in a testing capacity.
type Option func(h *Handler)

// WithFileSystem is used to override the internal FileSystem of the handler.
func WithFileSystem(fs FileSystem) Option {
	return func(h *Handler) {
		h.fs = fs
	}
}

// WithStateManager is used to override the internal StateManager of the handler.
func WithStateManager(sm StateManager) Option {
	return func(h *Handler) {
		h.stateManager = sm
	}
}
