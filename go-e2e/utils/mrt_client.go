package utils

import (
	"fmt"
	"os"
	"os/exec"
)

type Mrt struct {
	binaryName string
	command    *exec.Cmd
}

func MrtNew(binaryName string, env []string) *Mrt {
	command := exec.Command(binaryName)
	command.Env = append(os.Environ(), env...)

	return &Mrt{
		binaryName: binaryName,
		command:    command,
	}
}

func (m *Mrt) RunInDirectory(directory string) *Mrt {
	m.command.Args = append(m.command.Args, "--team-dir", directory)

	return m
}

func (m *Mrt) Setup() *Mrt {
	m.command.Args = append(m.command.Args, "setup")

	return m
}

func (m *Mrt) Clone() *Mrt {
	m.command.Args = append(m.command.Args, "clone-repositories")

	return m
}

type MrtExecutionError struct {
	TeamDir string
	Err     error
	Output  string
}

func (e *MrtExecutionError) Error() string {
	return fmt.Sprintf("failed to clone repositories in %s: %v\noutput: %s", e.TeamDir, e.Err, e.Output)
}

func (m *Mrt) Execute() {
	m.command.Stdout = os.Stdout
	m.command.Stderr = os.Stderr

	if err := m.command.Run(); err != nil {
		panic("executing mrt command failed: " + err.Error())
	}
}

func (m *Mrt) makeMrtCommand() *exec.Cmd {
	cmd := exec.Command(m.binaryName)
	return cmd
}
