package vagrantexec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineStateString(t *testing.T) {
	testcases := []struct {
		state MachineState
		str   string
	}{
		{Unknown, "Unknown"},
		{Aborted, "Aborted"},
		{GuruMeditation, "GuruMeditation"},
		{Inaccessible, "Inaccessible"},
		{NotCreated, "NotCreated"},
		{Paused, "Paused"},
		{PowerOff, "PowerOff"},
		{Stopping, "Stopping"},
		{Running, "Running"},
		{Saving, "Saving"},
		{Saved, "Saved"},
		{Stuck, "Stuck"},
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.str, tc.state.String())
	}
}

func TestToMachineState(t *testing.T) {
	testcases := []struct {
		str   string
		state MachineState
	}{
		{"garbage", Unknown},
		{"aborted", Aborted},
		{"gurumeditation", GuruMeditation},
		{"inaccessible", Inaccessible},
		{"not_created", NotCreated},
		{"paused", Paused},
		{"poweroff", PowerOff},
		{"stopping", Stopping},
		{"running", Running},
		{"saving", Saving},
		{"saved", Saved},
		{"stuck", Stuck},
	}

	for _, tc := range testcases {
		state := ToMachineState(tc.str)
		assert.Equalf(t, tc.state, state, "expected %s, got %s", tc.state, state)
	}
}
