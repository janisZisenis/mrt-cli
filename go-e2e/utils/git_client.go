package utils

import (
	"os"
	"os/exec"
)

type Git struct {
	command *exec.Cmd
}

func MakeGitCommand(env []string) *Git {
	command := exec.Command("git")
	command.Env = append(os.Environ(), env...)

	return &Git{
		command: command,
	}
}

func (git *Git) Clone(repositoryUrl string, destination string) *Git {
	git.command.Args = append(git.command.Args, "clone", repositoryUrl, destination)

	return git
}

func (git *Git) Execute() {
	outputBytes, err := git.command.CombinedOutput()

	if err != nil {
		panic("executing git command failed: " + string(outputBytes))
	}
}
