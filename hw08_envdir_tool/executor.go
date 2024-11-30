package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Env = env.ToEnv()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()

	if err != nil {
		if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
			return err.(*exec.ExitError).ExitCode()
		} else {
			return 1
		}
	}
	return 0
}
