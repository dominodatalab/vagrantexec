package command

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sr := ShellRunner{}
		out, err := sr.Execute("echo", "hello world")

		require.NoError(t, err)
		assert.Equal(t, "hello world\n", string(out))
	})

	t.Run("in_dir", func(t *testing.T) {
		sr := ShellRunner{Dir: "/usr"}
		out, err := sr.Execute("pwd")

		require.NoError(t, err)
		assert.Equal(t, "/usr\n", string(out))
	})

	t.Run("exit_error", func(t *testing.T) {
		sr := ShellRunner{}
		_, err := sr.Execute("sh", "-c", "echo 'actual err msg' >&2 && exit 64")
		require.IsType(t, ExitError{}, err)

		ee := err.(ExitError)
		assert.Equal(t, 64, ee.ExitStatus())
		assert.Equal(t, "sh exited with status 64: actual err msg", ee.Error())
	})

	t.Run("not_executable", func(t *testing.T) {
		sr := ShellRunner{}
		_, err := sr.Execute("garbage")

		require.Error(t, err)
		assert.IsType(t, new(exec.Error), err)
	})
}
