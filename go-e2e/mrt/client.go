package mrt

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"mrt-cli/go-e2e/internal"
	"mrt-cli/go-e2e/outputs"
)

type ExecutableCommand interface {
	Execute() (*outputs.Output, int)
	WithStdin(input string) ExecutableCommand
}

type BaseCommand interface {
	RunInDirectory(directory string) DirectedCommand
	Setup() SetupCommand
	Execute() (*outputs.Output, int)
}

type DirectedCommand interface {
	Setup() SetupCommand
	Run(args ...string) ExecutableCommand
	Execute() (*outputs.Output, int)
}

type SetupCommand interface {
	Clone() ExecutableCommand
	Execute() (*outputs.Output, int)
}

type Mrt struct {
	binaryName string
	command    *exec.Cmd
	stdin      io.Reader
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

func (m *Mrt) Clone() ExecutableCommand {
	m.command.Args = append(m.command.Args, "clone-repositories")

	return m
}

func (m *Mrt) Run(args ...string) ExecutableCommand {
	m.command.Args = append(m.command.Args, "run")
	m.command.Args = append(m.command.Args, args...)

	return m
}

func (m *Mrt) WithStdin(input string) ExecutableCommand {
	m.stdin = bytes.NewBufferString(input)
	return m
}

func (m *Mrt) Execute() (*outputs.Output, int) {
	if m.stdin != nil {
		m.command.Stdin = m.stdin
	}

	byteOutput, err := m.command.CombinedOutput()
	out := string(byteOutput)

	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			panic("executing mrt command failed unexpectedly: " + err.Error())
		}
	}

	return outputs.Make(splitLines(out)), exitCode
}

func splitLines(out string) []string {
	if out == "" {
		return []string{}
	}

	return strings.Split(strings.TrimSpace(out), "\n")
}
