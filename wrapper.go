package vagrantexec

import (
	"fmt"
	"strings"

	"github.com/dominodatalab/vagrant-exec/command"
	log "github.com/sirupsen/logrus"
)

const binary = "vagrant"

// Wrapper is the default implementation of the Vagrant interface.
type Wrapper struct {
	executable string
	runner     command.Runner
	logger     log.FieldLogger
}

// New creates a new Vagrant CLI wrapper.
func New() Wrapper {
	logger := log.New()

	return Wrapper{
		executable: binary,
		logger:     logger,
		runner:     command.ShellRunner{},
	}
}

// Version returns the installed version of Vagrant. An error is returned if the command fails or the output is invalid.
func (w Wrapper) Version() (version string, err error) {
	out, err := w.exec("version", "--machine-readable")
	if err != nil {
		return
	}
	entries, err := parseMachineOutput(string(out))
	if err != nil {
		return
	}
	data, err := pluckEntryData(entries, "version-installed")
	if err != nil {
		return
	}

	return data[0], err
}

func (w Wrapper) exec(args ...string) ([]byte, error) {
	fullCmd := fmt.Sprintf("%s %s", w.executable, strings.Join(args, " "))

	w.logger.Debugf("Running command [%s]", fullCmd)
	bs, err := w.runner.Execute(w.executable, args...)
	w.logger.Debugf("Command output [%s]: %s", fullCmd, bs)

	return bs, err
}
