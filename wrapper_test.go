package vagrantexec

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRunner struct {
	mock.Mock
}

func (m *mockRunner) Execute(cmd string, cmdargs ...string) ([]byte, error) {
	args := m.Called(cmd, cmdargs)
	if output, ok := args.Get(0).([]byte); ok {
		return output, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestVersion(t *testing.T) {
	w := New()
	mockVersion := func(resp []byte, err error) {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"version", "--machine-readable"}).Return(resp, err)
		w.runner = runner
	}

	t.Run("success", func(t *testing.T) {
		mockVersion(ioutil.ReadFile("testdata/version.out"))

		version, err := w.Version()
		assert.NoError(t, err)
		assert.Equal(t, "2.2.5", version)
	})

	t.Run("bad_output", func(t *testing.T) {
		mockVersion([]byte("bad output"), nil)

		_, err := w.Version()
		require.Error(t, err)
		assert.Equal(t, "invalid format", err.Error())
	})

	t.Run("error", func(t *testing.T) {
		mockVersion(nil, errors.New("something went wrong"))

		_, err := w.Version()
		assert.Error(t, err)
	})
}
