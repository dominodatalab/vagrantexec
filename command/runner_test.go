package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	sr := ShellRunner{}

	t.Run("success", func(t *testing.T) {
		out, err := sr.Execute("echo", "hello world")

		assert.NoError(t, err)
		assert.Equal(t, "hello world\n", string(out))
	})

	t.Run("error", func(t *testing.T) {
		_, err := sr.Execute("sh", "-c", "echo 'actual err msg' >&2 && exit 64")
		require.IsType(t, ExitError{}, err)

		ee := err.(ExitError)
		assert.Equal(t, 64, ee.ExitStatus())
		assert.Equal(t, "sh exited with status 64: actual err msg", ee.Error())
	})

	t.Run("not_executable", func(t *testing.T) {
		assert.Panics(t, func() { sr.Execute("garbage") })
	})
}
