package initialize

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMkdir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	fs := defaultFileSystem{}

	dir := fmt.Sprintf("%s/foobar", tmp)
	err = fs.Mkdir(dir, 0644)

	assert.NoError(t, err)
}
