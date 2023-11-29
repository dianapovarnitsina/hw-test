package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
	if len(cmd) == 0 {
		return
	}
	var cm *exec.Cmd
	if len(cmd) == 1 {
		command := cmd[0] // this is done due to the linter err
		cm = exec.Command(command)
	} else {
		cm = exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	}

	cm.Stdout = os.Stdout
	cm.Stderr = os.Stderr
	if err := cm.Run(); err != nil {
		log.Printf("cannot execute command, error: %v\n", err)
		return 1
	}
	return
}
