package vagrantexec

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/dominodatalab/vagrant-exec/command"
	"github.com/sirupsen/logrus"
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

func newTestWrapper(runner command.Runner) wrapper {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	return wrapper{
		executable: binary,
		logger:     logger,
		runner:     runner,
	}
}

func TestUp(t *testing.T) {
	mockUp := func(out []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"up"}).Return(out, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockUp([]byte("up output"), nil)
		assert.NoError(t, w.Up())
	})

	t.Run("error", func(t *testing.T) {
		w := mockUp(nil, errors.New("up failed"))
		assert.Error(t, w.Up())
	})
}

func TestHalt(t *testing.T) {
	mockHalt := func(out []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"halt"}).Return(out, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockHalt([]byte("halt output"), nil)
		assert.NoError(t, w.Halt())
	})

	t.Run("error", func(t *testing.T) {
		w := mockHalt(nil, errors.New("halt failed"))
		assert.Error(t, w.Halt())
	})
}

func TestDestroy(t *testing.T) {
	mockDestroy := func(out []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"destroy", "--force"}).Return(out, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockDestroy([]byte("destroy output"), nil)
		assert.NoError(t, w.Destroy())
	})

	t.Run("error", func(t *testing.T) {
		w := mockDestroy(nil, errors.New("destroy failed"))
		assert.Error(t, w.Destroy())
	})
}

func TestStatus(t *testing.T) {
	mockStatus := func(out []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"status", "--machine-readable"}).Return(out, err)
		return newTestWrapper(runner)
	}

	t.Run("one_machine", func(t *testing.T) {
		w := mockStatus(ioutil.ReadFile("testdata/status-single"))

		statuses, err := w.Status()
		require.NoError(t, err)

		expected := []MachineStatus{
			{
				Name:     "srv-1",
				Provider: "virtualbox",
				State:    NotCreated,
			},
		}
		assert.EqualValues(t, expected, statuses)
	})

	t.Run("multi_machine", func(t *testing.T) {
		w := mockStatus(ioutil.ReadFile("testdata/status-multiple"))

		statuses, err := w.Status()
		require.NoError(t, err)

		expected := []MachineStatus{
			{
				Name:     "srv-1",
				Provider: "virtualbox",
				State:    Running,
			},
			{
				Name:     "srv-2",
				Provider: "virtualbox",
				State:    PowerOff,
			},
		}
		assert.ElementsMatch(t, expected, statuses)
	})

	t.Run("error", func(t *testing.T) {
		w := mockStatus(nil, errors.New("runner error"))

		_, err := w.Status()
		assert.Error(t, err)
	})
}

func TestVersion(t *testing.T) {
	mockVersion := func(resp []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"version", "--machine-readable"}).Return(resp, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockVersion(ioutil.ReadFile("testdata/version"))

		version, err := w.Version()
		assert.NoError(t, err)
		assert.Equal(t, "2.2.5", version)
	})

	t.Run("bad_output", func(t *testing.T) {
		w := mockVersion([]byte("bad output"), nil)

		_, err := w.Version()
		require.Error(t, err)
		assert.Equal(t, "invalid format", err.Error())
	})

	t.Run("error", func(t *testing.T) {
		w := mockVersion(nil, errors.New("something went wrong"))

		_, err := w.Version()
		assert.Error(t, err)
	})
}

func TestSSH(t *testing.T) {
	sshCmd := "my-command 1 2 3"

	mockSSH := func(resp []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"ssh", "--no-tty", "--command", sshCmd}).Return(resp, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockSSH([]byte("command output"), nil)

		output, err := w.SSH(sshCmd)
		assert.NoError(t, err)
		assert.Equal(t, "command output", output)
	})

	t.Run("error", func(t *testing.T) {
		w := mockSSH(nil, errors.New("runner error"))

		_, err := w.SSH(sshCmd)
		assert.Error(t, err)
	})
}

func TestPluginList(t *testing.T) {
	mockPluginList := func(resp []byte, err error) wrapper {
		runner := new(mockRunner)
		runner.On("Execute", "vagrant", []string{"plugin", "list", "--machine-readable"}).Return(resp, err)
		return newTestWrapper(runner)
	}

	t.Run("success", func(t *testing.T) {
		w := mockPluginList(ioutil.ReadFile("testdata/plugin-list"))

		actual, err := w.PluginList()
		require.NoError(t, err)

		expected := []Plugin{
			{
				Name:     "vagrant-disksize",
				Version:  "0.1.3",
				Location: "global",
			},
			{
				Name:     "vagrant-ip-show",
				Version:  "0.0.4",
				Location: "global",
			},
		}
		assert.EqualValues(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		w := mockPluginList(nil, errors.New("runner error"))

		_, err := w.PluginList()
		assert.Error(t, err)
	})
}

func TestPluginInstall(t *testing.T) {
	//mockPluginList := func(resp []byte, err error) {}

	// no name
	// version
	// local
	// error
}
