package initialize

import "os"

type defaultFileSystem struct{}

func (f *defaultFileSystem) Mkdir(name string, perms os.FileMode) error {
	return os.Mkdir(name, perms)
}
