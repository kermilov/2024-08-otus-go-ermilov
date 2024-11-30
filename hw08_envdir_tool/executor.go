package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = env.ToEnv()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

func (env Environment) ToEnv() []string {
	envList := make([]string, 0, len(env))
	for name, value := range env {
		if value.NeedRemove {
			continue
		}
		envList = append(envList, name+"="+value.Value)
	}
	return envList
}
