package utils

import (
	"context"
	"os"
	"os/exec"
)

type GitBaseCommand interface {
	Clone(repositoryURL string, destination string) GitCloneCommand
}

type GitCloneCommand interface {
	Execute()
}

type Git struct {
	command *exec.Cmd
}

func MakeGitCommand(env []string) GitBaseCommand {
	command := exec.CommandContext(context.Background(), "git")
	command.Env = mergeEnv(os.Environ(), env)

	return &Git{
		command: command,
	}
}

func (git *Git) Clone(repositoryURL string, destination string) GitCloneCommand {
	git.command.Args = append(git.command.Args, "clone", repositoryURL, destination)

	return git
}

func (git *Git) Execute() {
	outputBytes, err := git.command.CombinedOutput()
	if err != nil {
		panic("executing git command failed: " + string(outputBytes))
	}
}
