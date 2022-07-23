package state

// Option provides a function that is able to override the internals of a
// Manager instance. These options should rarely be used for anything other
// than testing, as the mananger will initialize sane defaults when using the
// NewManager method.
type Option func(*Manager)

// WithFileSystem provides an option to override the internal FileSystem
// interface that is found in the state.Manager.
func WithFileSystem(fs FileSystem) Option {
	return func(m *Manager) {
		m.fs = fs
	}
}
