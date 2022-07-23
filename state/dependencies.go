//go:generate mockgen -source=dependencies.go -destination=./mocks.go -package=state -mock_names FileSystem=mockFileSystem,File=mockFile

package state

import (
	"io"
	"os"
)

// FileSystem provides an interface that can interact with the file system.
// This interface is primarily used for testing. All of these methods are
// found in the `os` package.
type FileSystem interface {
	Create(name string) (File, error)
	ReadFile(name string) ([]byte, error)
	Remove(name string) error
	Rename(oldpath string, newpath string) error
	Getwd() (string, error)
	Stat(name string) (os.FileInfo, error)
}

// File provides an interface for an os.File to allow for testing without
// making modifications to the file system.
type File interface {
	Write([]byte) (int, error)
	io.Closer
}
