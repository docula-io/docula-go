package cmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docula-io/docula/cmd"
)

func TestExecute(t *testing.T) {
	assert.NoError(t, cmd.Execute(context.Background()))
}
