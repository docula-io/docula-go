package initialize

import "os"

type defaultFileSystem struct{}

func (f *defaultFileSystem) Mkdir(name string) error {
	const dirPerms = os.FileMode(0755)
	return os.MkdirAll(name, dirPerms)
}
