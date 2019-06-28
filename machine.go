package vagrantexec

import (
	"bufio"
	"fmt"
	"strings"
)

// machineOutputEntry defines all of the components in a single line of machine-readable output.
type machineOutputEntry struct {
	timestamp string
	target    string
	mType     string
	data      []string
}

// parseMachineOutput converts machine-readable output into a slice of entries.
func parseMachineOutput(machineOut string) (entries []machineOutputEntry, err error) {
	scanner := bufio.NewScanner(strings.NewReader(machineOut))
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")
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
