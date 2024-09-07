package core

import "os/exec"

func ExecuteBash(file string, args []string) (string, error) {
	bashArgs := append([]string{file}, args...)
	script := exec.Command("/bin/bash", bashArgs...)
	output, err := script.Output()

	return string(output), err
}
