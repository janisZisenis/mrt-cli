package mrt

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"mrt-cli/e2e-tests/internal"
	"mrt-cli/e2e-tests/outputs"
)

type ExecutableCommand interface {
	Execute() (*outputs.Output, int)
	ExecuteWithInput(input string) (*outputs.Output, int)
}

type RunCommand interface {
	SubCommand(name string, args ...string) ExecutableCommand
	Execute() (*outputs.Output, int)
}

type BaseCommand interface {
	RunInDirectory(directory string) DirectedCommand
	Setup() SetupCommand
	Execute() (*outputs.Output, int)
}

type DirectedCommand interface {
	Setup() SetupCommand
	Run() RunCommand
	GitHook(hookName string, repositoryPath string, args ...string) ExecutableCommand
	Execute() (*outputs.Output, int)
}

type SetupCommand interface {
	Clone() ExecutableCommand
	InstallGitHooks() ExecutableCommand
	All(args ...string) ExecutableCommand
	SubCommand(name string, args ...string) ExecutableCommand
	Execute() (*outputs.Output, int)
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

func (m *Mrt) Clone() ExecutableCommand {
	m.command.Args = append(m.command.Args, "clone-repositories")

	return m
}

func (m *Mrt) InstallGitHooks() ExecutableCommand {
	m.command.Args = append(m.command.Args, "install-git-hooks")

	return m
}

func (m *Mrt) All(args ...string) ExecutableCommand {
	m.command.Args = append(m.command.Args, "all")
	m.command.Args = append(m.command.Args, args...)

	return m
}

func (m *Mrt) SubCommand(name string, args ...string) ExecutableCommand {
	m.command.Args = append(m.command.Args, name)
	if len(args) > 0 {
		m.command.Args = append(m.command.Args, "--")
		m.command.Args = append(m.command.Args, args...)
	}

	return m
}

func (m *Mrt) GitHook(hookName string, repositoryPath string, args ...string) ExecutableCommand {
	m.command.Args = append(m.command.Args, "git-hook", "--hook-name", hookName, "--repository-path", repositoryPath)
	m.command.Args = append(m.command.Args, args...)

	return m
}

func (m *Mrt) Run() RunCommand {
	m.command.Args = append(m.command.Args, "run")

	return m
}

func (m *Mrt) ExecuteWithInput(input string) (*outputs.Output, int) {
	m.command.Stdin = bytes.NewBufferString(input)
	return m.execute()
}

func (m *Mrt) Execute() (*outputs.Output, int) {
	return m.execute()
}

func (m *Mrt) execute() (*outputs.Output, int) {
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
