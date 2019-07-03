package vagrantexec

import (
	"bufio"
	"errors"
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

// MachineState represents the state of a machine Vagrant is managing.
type MachineState int

func (s MachineState) String() string {
	return []string{
		"Unknown", "Aborted", "GuruMeditation", "Inaccessible", "NotCreated",
		"Paused", "PowerOff", "Stopping", "Running", "Saving", "Saved", "Stuck"}[s]
}

// ToMachineState converts a string into a MachineState. An Unknown state is returned if the string is invalid.
func ToMachineState(str string) MachineState {
	switch str {
	case "running":
		return Running
	case "not_created":
		return NotCreated
	case "saved":
		return Saved
	case "poweroff":
		return PowerOff
	case "aborted":
		return Aborted
	case "paused":
		return Paused
	case "stopping":
		return Stopping
	case "saving":
		return Saving
	case "stuck":
		return Stuck
	case "inaccessible":
		return Inaccessible
	case "gurumeditation":
		return GuruMeditation
	default:
		return Unknown
	}
}

// MachineStatus defines metadata that describes a single Vagrant machine.
type MachineStatus struct {
	Name     string
	Provider string
	State    MachineState
}

// machineOutputEntry defines all of the components in a single line of machine-readable output.
type machineOutputEntry struct {
	timestamp string
	target    string
	mType     string
	data      []string
}

// parseMachineReadable converts machine-readable output into a slice of machineOutputEntry.
func parseMachineReadable(machineOut string) (entries []machineOutputEntry, err error) {
	scanner := bufio.NewScanner(strings.NewReader(machineOut))
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")
		if len(row) < 4 {
			err = errors.New("invalid format")
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
