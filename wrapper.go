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

// Up creates and configures guest machines according to your Vagrantfile.
func (w Wrapper) Up() error {
	out, err := w.exec("up")
	if err == nil {
		w.info(out)
	}
	return err
}

// Halt will gracefully shut down the guest operating system and power down the guest machine.
func (w Wrapper) Halt() error {
	out, err := w.exec("halt")
	if err == nil {
		w.info(out)
	}
	return err
}

// Destroy stops the running guest machines and destroys all of the resources created during the creation process.
func (w Wrapper) Destroy() error {
	out, err := w.exec("destroy", "--force")
	if err == nil {
		w.info(out)
	}
	return err
}

// Version displays the current version of Vagrant you have installed.
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

func (w Wrapper) info(out []byte) {
	if len(out) > 0 {
		w.logger.Info(string(out))
	}
}
