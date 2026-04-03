package mrt

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"mrt-cli/go-e2e/internal"
	"mrt-cli/go-e2e/output"
)

type BaseCommand interface {
	RunInDirectory(directory string) DirectedCommand
	Setup() SetupCommand
	Execute() *output.Output
}

type DirectedCommand interface {
	Setup() SetupCommand
	Execute() *output.Output
}

type SetupCommand interface {
	Clone() CloneCommand
	Execute() *output.Output
}

type CloneCommand interface {
	Execute() *output.Output
}

type Mrt struct {
	binaryName string
	command    *exec.Cmd
}

func MakeCommand(binaryPath string, sshEnv []string) BaseCommand {
	command := exec.CommandContext(context.Background(), binaryPath)
	command.Env = internal.MergeEnv(os.Environ(), sshEnv)

	return &Mrt{
		binaryName: binaryPath,
		command:    command,
	}
}

func (m *Mrt) RunInDirectory(directory string) DirectedCommand {
	m.command.Args = append(m.command.Args, "--team-dir", directory)

	return m
}

func (m *Mrt) Setup() SetupCommand {
	m.command.Args = append(m.command.Args, "setup")

	return m
}

func (m *Mrt) Clone() CloneCommand {
	m.command.Args = append(m.command.Args, "clone-repositories")

	return m
}

func (m *Mrt) Execute() *output.Output {
	byteOutput, err := m.command.CombinedOutput()
	out := string(byteOutput)

	if err != nil {
		panic("executing mrt command failed: " + out)
	}

	return output.Make(splitLines(out))
}

func splitLines(out string) []string {
	if out == "" {
		return []string{}
	}

	return strings.Split(strings.TrimSpace(out), "\n")
}
