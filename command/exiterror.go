package command

import (
	"fmt"
	"strings"
)

// ExitError is created whenever a command exits with a non-zero status.
type ExitError struct {
	msg        string
	exitStatus int
}

func (e ExitError) Error() string {
	return e.msg
}

// ExitStatus returns the exit code of the exited process.
func (e ExitError) ExitStatus() int {
	return e.exitStatus
}

// newExistError creates a new ExitError with a descriptive message.
func newExitError(cmd string, exitStatus int, msg string) ExitError {
	return ExitError{
		msg:        fmt.Sprintf("%s exited with status %d: %s", cmd, exitStatus, strings.TrimSpace(msg)),
		exitStatus: exitStatus,
	}
}
