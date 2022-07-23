package state

import "os"

type defaultFileSystem struct{}

func (d *defaultFileSystem) Create(name string) (File, error) {
	return os.Create(name)
}

func (d *defaultFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (d *defaultFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (d *defaultFileSystem) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (d *defaultFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

func (d *defaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
