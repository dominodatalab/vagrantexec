package vagrantexec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dominodatalab/vagrant-exec/command"
	log "github.com/sirupsen/logrus"
)

const binary = "vagrant"

// Interface defines the supported Vagrant commands.
type Interface interface {
	Up() error
	Halt() error
	Destroy() error
	Status() ([]MachineStatus, error)
	Version() (string, error)
	SSH(string) (string, error)

	PluginList() ([]Plugin, error)
}

// Plugin encapsulates Vagrant plugin metadata.
type Plugin struct {
	Name     string
	Version  string
	Location string
}

// Wrapper is the default implementation of the Vagrant Interface.
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

// Status reports the status of the machines Vagrant is managing.
func (w Wrapper) Status() (statuses []MachineStatus, err error) {
	out, err := w.exec("status", "--machine-readable")
	if err != nil {
		return
	}
	machineInfo, err := parseMachineReadable(out)
	if err != nil {
		return
	}

	statusMap := map[string]*MachineStatus{}
	for _, entry := range machineInfo {
		if len(entry.target) == 0 {
			continue // skip when no target specified
		}

		var status *MachineStatus // fetch status or create when missing
		status, ok := statusMap[entry.target]
		if !ok {
			status = &MachineStatus{Name: entry.target}
			statusMap[entry.target] = status
		}

		switch entry.mType { // populate status fields
		case "provider-name":
			status.Provider = entry.data[0]
		case "state":
			status.State = ToMachineState(entry.data[0])
		}
	}

	for _, st := range statusMap {
		statuses = append(statuses, *st)
	}
	return statuses, nil
}

// Version displays the current version of Vagrant you have installed.
func (w Wrapper) Version() (version string, err error) {
	out, err := w.exec("version", "--machine-readable")
	if err != nil {
		return
	}
	vInfo, err := parseMachineReadable(out)
	if err != nil {
		return
	}
	data, err := pluckEntryData(vInfo, "version-installed")
	if err != nil {
		return
	}

	return data[0], err
}

// SSH executes a command on a Vagrant machine via SSH and returns the stdout/stderr output.
func (w Wrapper) SSH(command string) (string, error) {
	out, err := w.exec("ssh", "--no-tty", "--command", command)
	return string(out), err
}

// PluginList returns a list of all installed plugins, their versions and install locations.
func (w Wrapper) PluginList() (plugins []Plugin, err error) {
	out, err := w.exec("plugin", "list", "--machine-readable")
	if err != nil {
		return
	}
	pluginInfo, err := parseMachineReadable(out)
	if err != nil {
		return
	}
	pluginMetadataExtractor := regexp.MustCompile(`^([\w-]+)\s\((.*)%!\(VAGRANT_COMMA\)\s([a-z]+)\)$`)
	for _, entry := range pluginInfo {
		if entry.mType == "ui" {
			matches := pluginMetadataExtractor.FindAllStringSubmatch(entry.data[1], -1)[0][1:]
			plugins = append(plugins, Plugin{
				Name:     matches[0],
				Version:  matches[1],
				Location: matches[2],
			})
		}
	}
	return
}

// exec dispatches vagrant commands via the shell runner.
func (w Wrapper) exec(args ...string) ([]byte, error) {
	fullCmd := fmt.Sprintf("%s %s", w.executable, strings.Join(args, " "))

	w.logger.Debugf("Running command [%s]", fullCmd)
	bs, err := w.runner.Execute(w.executable, args...)
	w.logger.Debugf("Command output [%s]: %s", fullCmd, bs)

	return bs, err
}

// info will log non-empty input.
func (w Wrapper) info(out []byte) {
	if len(out) > 0 {
		w.logger.Info(string(out))
	}
}
