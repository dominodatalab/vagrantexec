// Package vagrantexec defines the types required to interface with the vagrant executable.
package vagrantexec

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dominodatalab/vagrant-exec/command"
	log "github.com/sirupsen/logrus"
)

const binary = "vagrant"

// Vagrant defines the interface for executing Vagrant commands.
type Vagrant interface {
	Up() error
	Halt() error
	Destroy() error
	Status() ([]MachineStatus, error)
	Version() (string, error)
	SSH(string, string) (string, error)

	PluginList() ([]Plugin, error)
	PluginInstall(Plugin) error
}

// Plugin encapsulates Vagrant plugin metadata.
type Plugin struct {
	Name     string
	Version  string
	Location string
}

// wrapper is the default implementation of the Vagrant Interface.
type wrapper struct {
	executable string
	runner     command.Runner
	logger     log.FieldLogger
}

// New creates a new Vagrant CLI wrapper targeting a directory where a Vagrantfile should exist.
func New(vagrantfileDir string, debug bool) Vagrant {
	if len(vagrantfileDir) == 0 {
		panic("vagrantfile dir cannot be empty")
	}
	runner := command.ShellRunner{
		Dir: vagrantfileDir,
	}

	logger := log.New()
	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	return wrapper{
		executable: binary,
		logger:     logger,
		runner:     runner,
	}
}

// Up creates and configures guest machines according to your Vagrantfile.
func (w wrapper) Up() error {
	w.logger.Info("Starting vagrant environment")
	return w.execLogOutput("up")
}

// Halt will gracefully shut down the guest operating system and power down the guest machine.
func (w wrapper) Halt() error {
	w.logger.Info("Stopping vagrant machines")
	return w.execLogOutput("halt")
}

// Destroy stops the running guest machines and destroys all of the resources created during the creation process.
func (w wrapper) Destroy() error {
	w.logger.Info("Deleting vagrant machines")
	return w.execLogOutput("destroy", "--force")
}

// Status reports the status of the machines Vagrant is managing.
func (w wrapper) Status() (statuses []MachineStatus, err error) {
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
func (w wrapper) Version() (version string, err error) {
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
// You can use an empty string as the nameOrID if you only have one VM defined in your Vagrantfile.
func (w wrapper) SSH(nameOrID, command string) (string, error) {
	cmdArgs := []string{"ssh", "--no-tty", "--command", command}
	if len(nameOrID) > 0 {
		cmdArgs = append(cmdArgs, nameOrID)
	}

	out, err := w.exec(cmdArgs...)
	return string(out), err
}

// PluginList returns a list of all installed plugins, their versions and install locations.
func (w wrapper) PluginList() (plugins []Plugin, err error) {
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
		if entry.mType == "ui" { // "ui" type may contain combined name/version data
			combinedData := entry.data[1]
			if strings.Contains(combinedData, "No plugins installed") {
				break
			}

			if ms := pluginMetadataExtractor.FindAllStringSubmatch(combinedData, -1); ms != nil {
				matches := ms[0][1:] // lens into captures
				plugins = append(plugins, Plugin{
					Name:     matches[0],
					Version:  matches[1],
					Location: matches[2],
				})
			}
		}
	}
	return
}

// PluginInstall installs a plugin with the given name or file path.
func (w wrapper) PluginInstall(plugin Plugin) error {
	if len(plugin.Name) == 0 {
		return errors.New("plugin must have a name")
	}
	cmdArgs := []string{"plugin", "install", plugin.Name}

	if len(plugin.Version) > 0 {
		cmdArgs = append(cmdArgs, "--plugin-version", plugin.Version)
	}
	if plugin.Location == "local" {
		cmdArgs = append(cmdArgs, "--local")
	}

	w.logger.Infof("Installing vagrant plugin: %s", plugin.Name)
	return w.execLogOutput(cmdArgs...)
}

// exec dispatches vagrant commands via the shell runner.
func (w wrapper) exec(args ...string) ([]byte, error) {
	fullCmd := fmt.Sprintf("%s %s", w.executable, strings.Join(args, " "))

	w.logger.Debugf("Running command [%s]", fullCmd)
	bs, err := w.runner.Execute(w.executable, args...)
	w.logger.Debugf("Command output [%s]: %s", fullCmd, bs)

	return bs, err
}

// execLogOutput logs the output of the command at an info level instead of returning it.
func (w wrapper) execLogOutput(args ...string) error {
	out, err := w.exec(args...)
	if len(out) > 0 {
		w.logger.Info(string(out))
	}
	return err
}
