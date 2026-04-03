package git

import (
	"context"
	"os"
	"os/exec"

	"mrt-cli/go-e2e/internal"
)

type BaseCommand interface {
	Clone(repositoryURL string, destination string) CloneCommand
}

type CloneCommand interface {
	Execute()
}

type Git struct {
	command *exec.Cmd
}

func MakeCommand(sshEnv []string) BaseCommand {
	command := exec.CommandContext(context.Background(), "git")
	command.Env = internal.MergeEnv(os.Environ(), sshEnv)

	return &Git{
		command: command,
	}
}

func (g *Git) Clone(repositoryURL string, destination string) CloneCommand {
	g.command.Args = append(g.command.Args, "clone", repositoryURL, destination)

	return g
}

func (g *Git) Execute() {
	outputBytes, err := g.command.CombinedOutput()
	if err != nil {
		panic("executing git command failed: " + string(outputBytes))
	}
}
