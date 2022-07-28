package initialize

import (
	"testing"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

const (
	nameLine      = "What should we name this dir?"
	indexTypeLine = "Choose an index type"
)

func TestDefaultSurveyAsk(t *testing.T) {
	testCases := []struct {
		name  string
		input func(t *testing.T, c *expect.Console)
		wants Configuration
	}{
		{
			name: "default index option",
			input: func(t *testing.T, c *expect.Console) {
				c.ExpectString(nameLine)

				_, err := c.SendLine("foobar")

				assert.NoError(t, err)

				c.ExpectString(indexTypeLine)

				_, err = c.SendLine("")

				assert.NoError(t, err)

				c.ExpectEOF()
			},
			wants: Configuration{
				Name:      "foobar",
				IndexType: "timestamp",
			},
		},
		{
			name: "second index option",
			input: func(t *testing.T, c *expect.Console) {
				c.ExpectString(nameLine)

				_, err := c.SendLine("barfoo")
				assert.NoError(t, err)

				c.ExpectString(indexTypeLine)

				_, err = c.Send(string(terminal.KeyArrowDown))
				assert.NoError(t, err)

				_, err = c.SendLine("")
				assert.NoError(t, err)

				c.ExpectEOF()
			},
			wants: Configuration{
				Name:      "barfoo",
				IndexType: "sequential",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			pty, tty, err := pseudotty.Open()
			assert.NoError(t, err)

			term := vt10x.New(vt10x.WithWriter(tty))

			c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
			assert.NoError(t, err)

			defer c.Close()

			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				tt.input(t, c)
			}()

			stdio := terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}

			s := defaultSurvey{}

			res, err := s.Ask(survey.WithStdio(stdio.In, stdio.Out, stdio.Err))
			assert.NoError(t, err)

			assert.Equal(t, tt.wants, res)

			assert.NoError(t, c.Tty().Close())

			<-doneCh
		})
	}
}
