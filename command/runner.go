// Package command contains primitives for running external commands.
package command

import (
	"bytes"
	"os/exec"
)

// Runner provides an interface for running external commands.
type Runner interface {
	Execute(cmd string, args ...string) ([]byte, error)
}

// ShellRunner provides provides a simplified interface to exec.Command making it easier to process output and errors.
type ShellRunner struct {
	// Dir is the directory where the commands will be executed.
	Dir string
}

// Execute invokes a shell command with any number of arguments and returns standard output.
//
// If the command starts but does not complete successfully, an ExitError will be returned with output from standard
// error. Any other error will result in a panic.
func (r ShellRunner) Execute(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	c.Dir = r.Dir

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			err = newExitError(cmd, ee.ExitCode(), string(stderr.Bytes()))
		}
	}

	return stdout.Bytes(), err
}
