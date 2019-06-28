package vagrantexec

import (
	"fmt"
	"strings"

	"github.com/dominodatalab/vagrant-exec/command"
	log "github.com/sirupsen/logrus"
)

const binary = "vagrant"

type wrapper struct {
	executable string
	runner     command.Runner
	logger     log.FieldLogger
}

func New() *wrapper {
	logger := log.New()
	logger.Level = log.DebugLevel

	return &wrapper{
		executable: binary,
		logger:     logger,
		runner:     command.ShellRunner{},
	}
}

func (w wrapper) Version() (version string, err error) {
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

func (w wrapper) exec(args ...string) ([]byte, error) {
	fullCmd := fmt.Sprintf("%s %s", w.executable, strings.Join(args, " "))

	w.logger.Debugf("Running command [%s]", fullCmd)
	bs, err := w.runner.Execute(w.executable, args...)
	w.logger.Debugf("Command output [%s]: %s", fullCmd, bs)

	return bs, err
}
