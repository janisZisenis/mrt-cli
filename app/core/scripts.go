package core

import (
	"os"
	"os/exec"
)

func ExecuteScript(file string, args []string) error {
	script := exec.Command(file, args...)
	script.Stdout = os.Stdout
	err := script.Run()

	return err
}
