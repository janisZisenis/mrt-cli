package core

import (
	"os"
	"os/exec"
)

func ExecuteBash(file string, args []string) error {
	bashArgs := append([]string{file}, args...)
	script := exec.Command("/bin/bash", bashArgs...)

	//script.Stdin = os.Stdin
	script.Stdout = os.Stdout
	//script.Stderr = os.Stderr

	err := script.Run()

	return err
}
