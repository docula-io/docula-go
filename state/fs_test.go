package state

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFileSystemCreate(t *testing.T) {
	tmp, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	fs := defaultFileSystem{}

	f, err := fs.Create(fmt.Sprintf("%s/foobar", tmp))

	assert.NotNil(t, f)
	assert.NoError(t, err)
}

func TestDefaultFileSystemReadFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.Remove(tmp.Name()))
	}()

	_, err = tmp.Write([]byte("foobar, foo bar!"))
	assert.NoError(t, err)

	assert.NoError(t, tmp.Close())

	fs := defaultFileSystem{}

	data, err := fs.ReadFile(tmp.Name())
	assert.NoError(t, err)

	assert.Equal(t, []byte("foobar, foo bar!"), data)
}

func TestDefaultFileSystemRemove(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	assert.NoError(t, err)

	_, err = tmp.Write([]byte("foobar, foo bar!"))
	assert.NoError(t, err)

	assert.NoError(t, tmp.Close())

	fs := defaultFileSystem{}

	err = fs.Remove(tmp.Name())
	assert.NoError(t, err)

	_, err = os.Stat(tmp.Name())
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultFileSystemRename(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	assert.NoError(t, err)

	tmpdir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	newPath := fmt.Sprintf("%s/foo", tmpdir)

	defer func() {
		assert.NoError(t, os.Remove(newPath))
	}()

	_, err = tmp.Write([]byte("foobar, foo bar!"))
	assert.NoError(t, err)

	assert.NoError(t, tmp.Close())

	fs := defaultFileSystem{}

	assert.NoError(t, fs.Rename(tmp.Name(), newPath))

	data, err := fs.ReadFile(newPath)
	assert.NoError(t, err)

	assert.Equal(t, []byte("foobar, foo bar!"), data)
}

func TestDefaultFileSystemGetwd(t *testing.T) {
	expected, err := os.Getwd()
	assert.NoError(t, err)

	fs := defaultFileSystem{}
	val, err := fs.Getwd()

	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestDefaultFileSystemStat(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.Remove(tmp.Name()))
	}()

	expected, err := os.Stat(tmp.Name())
	assert.NoError(t, err)

	fs := defaultFileSystem{}

	actual, err := fs.Stat(tmp.Name())
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
