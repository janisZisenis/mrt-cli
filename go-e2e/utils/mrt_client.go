package utils

import (
	"os"
	"os/exec"
	"strings"
)

type Mrt struct {
	binaryName string
	command    *exec.Cmd
}

func MakeMrtCommand(binaryPath string, env []string) *Mrt {
	command := exec.Command(binaryPath)
	command.Env = append(os.Environ(), env...)

	return &Mrt{
		binaryName: binaryPath,
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

func (m *Mrt) Execute() *Output {
	byteOutput, err := m.command.Output()
	output := string(byteOutput)

	if err != nil {
		panic("executing mrt command failed: " + output)
	}

	return MakeOutput(SplitLines(output))
}

func SplitLines(output string) []string {
	if output == "" {
		return []string{}
	}

	return strings.Split(strings.TrimSpace(output), "\n")
}

func (m *Mrt) makeMrtCommand() *exec.Cmd {
	cmd := exec.Command(m.binaryName)
	return cmd
}
