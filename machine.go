package vagrantexec

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	// Unknown represents any unhandled state.
	Unknown MachineState = iota
	// Aborted means the VM was abruptly stopped without properly closing the session.
	Aborted
	// GuruMeditation means that an internal error in VirtualBox caused the VM to fail.
	GuruMeditation
	// Inaccessible means that VirtualBox can't find your VM configuration.
	Inaccessible
	// NotCreated means the VM has not been created.
	NotCreated
	// Paused means the VM may have been paused by VirtualBox.
	Paused
	// PowerOff means the VM has been turned off.
	PowerOff
	// Stopping means the VM is in the process of stopping.
	Stopping
	// Running means the VM is up and running.
	Running
	// Saving means the VM is currently saving its state.
	Saving
	// Saved means the VM has been suspended.
	Saved
	// Stuck means that VirtualBox is unable to recover the current state of the VM.
	Stuck
)

var (
	// stateStrList contains a list of state string representations based on the ordinal values of the constants.
	stateStrList = []string{
		"Unknown", "Aborted", "GuruMeditation", "Inaccessible", "NotCreated",
		"Paused", "PowerOff", "Stopping", "Running", "Saving", "Saved", "Stuck",
	}

	// strStateMap maps vagrant state output to their corresponding constants.
	strStateMap = map[string]MachineState{
		"running":        Running,
		"not_created":    NotCreated,
		"saved":          Saved,
		"poweroff":       PowerOff,
		"aborted":        Aborted,
		"paused":         Paused,
		"stopping":       Stopping,
		"saving":         Saving,
		"stuck":          Stuck,
		"inaccessible":   Inaccessible,
		"gurumeditation": GuruMeditation,
	}
)

// MachineState denotes the state of a machine Vagrant is managing.
//
// A list of all available states can be found here: https://github.com/hashicorp/vagrant/blob/4ce8d84f7e6709e4478612a9f0810dc686076ee0/templates/locales/en.yml#L2056
type MachineState int

// ToMachineState converts a string into a MachineState. An Unknown state is returned if the string is invalid.
func ToMachineState(str string) MachineState {
	return strStateMap[str]
}

// String returns a string representation of the MachineState.
func (s MachineState) String() string {
	return stateStrList[s]
}

// MachineStatus encompasses the machine metadata provided by Vagrant.
type MachineStatus struct {
	Name     string
	Provider string
	State    MachineState
}

// IsRunning returns true if the virtual machine is in a running state.
func (m MachineStatus) IsRunning() bool {
	return m.State == Running
}

// machineOutputEntry defines all of the components in a single line of machine-readable output.
//
// See https://www.vagrantup.com/docs/cli/machine-readable.html#format for more details.
type machineOutputEntry struct {
	timestamp string
	target    string
	mType     string
	data      []string
}

// parseMachineReadable converts machine-readable output into a slice of machineOutputEntry.
func parseMachineReadable(machineOut []byte) (entries []machineOutputEntry, err error) {
	scanner := bufio.NewScanner(strings.NewReader(string(machineOut)))
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")
		if len(row) < 4 {
			err = fmt.Errorf("invalid machine-readable format: %s", row)
			return
		}

		entries = append(entries, machineOutputEntry{
			timestamp: row[0],
			target:    row[1],
			mType:     row[2],
			data:      row[3:],
		})
	}
	err = scanner.Err()
	return
}

// pluckEntryData extracts a single data field from a collection of entries.
func pluckEntryData(entries []machineOutputEntry, messageType string) ([]string, error) {
	for _, e := range entries {
		if e.mType == messageType {
			return e.data, nil
		}
	}
	return nil, fmt.Errorf("cannot pluck data for message type: %s", messageType)
}
