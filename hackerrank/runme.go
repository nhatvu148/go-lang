package main

import (
	"fmt"
	"os/exec"
)

func main() {
	args := []string{
		"cmd",
		"/C",
		"C:/Users/nhatv/AppData/Local/Apps/BETA_CAE_Systems/ansa_v21.0.1/ansa64.bat",
	}

	CmdExec(args...)
}

func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	fmt.Printf("Exec: %v", args)

	cmd := exec.Command(baseCmd, cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
