package initialize

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMkdir(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "single depth dir",
			input: "foobar",
		},
		{
			name:  "multi depth dir",
			input: "foo/bar",
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tmp, err := os.MkdirTemp("", "")
			assert.NoError(t, err)

			fs := defaultFileSystem{}

			dir := fmt.Sprintf("%s/%s", tmp, tt.input)
			err = fs.Mkdir(dir)

			assert.NoError(t, err)
		})
	}
}
