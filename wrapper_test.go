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

func TestUp(t *testing.T) {
	w := New()
	mockUp := func(out []byte, err error) {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"up"}).Return(out, err)
		w.runner = runner
	}

	t.Run("success", func(t *testing.T) {
		mockUp([]byte("up output"), nil)
		assert.NoError(t, w.Up())
	})

	t.Run("error", func(t *testing.T) {
		mockUp(nil, errors.New("up failed"))
		assert.Error(t, w.Up())
	})
}

func TestHalt(t *testing.T) {
	w := New()
	mockHalt := func(out []byte, err error) {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"halt"}).Return(out, err)
		w.runner = runner
	}

	t.Run("success", func(t *testing.T) {
		mockHalt([]byte("halt output"), nil)
		assert.NoError(t, w.Halt())
	})

	t.Run("error", func(t *testing.T) {
		mockHalt(nil, errors.New("halt failed"))
		assert.Error(t, w.Halt())
	})
}

func TestDestroy(t *testing.T) {
	w := New()
	mockDestroy := func(out []byte, err error) {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"destroy", "--force"}).Return(out, err)
		w.runner = runner
	}

	t.Run("success", func(t *testing.T) {
		mockDestroy([]byte("destroy output"), nil)
		assert.NoError(t, w.Destroy())
	})

	t.Run("error", func(t *testing.T) {
		mockDestroy(nil, errors.New("destroy failed"))
		assert.Error(t, w.Destroy())
	})
}

func TestVersion(t *testing.T) {
	w := New()
	mockVersion := func(resp []byte, err error) {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"version", "--machine-readable"}).Return(resp, err)
		w.runner = runner
	}

	t.Run("success", func(t *testing.T) {
		mockVersion(ioutil.ReadFile("testdata/version"))

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
