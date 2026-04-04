package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"mrt-cli/go-e2e/internal"
)

type BaseCommand interface {
	InDirectory(path string) DirectedCommand
	Clone(repositoryURL string, destination string) ExecutableCommand
}

type DirectedCommand interface {
	MakeCommitOnNewBranch(branch string, message string) ExecutableCommand
	Push(branch string) ExecutableCommand
}

type ExecutableCommand interface {
	Execute() (int, error)
}

type Git struct {
	args   []string
	path   string
	sshEnv []string
}

func MakeCommand(sshEnv []string) BaseCommand {
	return &Git{sshEnv: sshEnv}
}

func (g *Git) InDirectory(path string) DirectedCommand {
	return &Git{
		args:   append(g.args, "-C", path),
		path:   path,
		sshEnv: g.sshEnv,
	}
}

func (g *Git) Clone(repositoryURL string, destination string) ExecutableCommand {
	return &Git{
		args:   append(g.args, "clone", repositoryURL, destination),
		sshEnv: g.sshEnv,
	}
}

func (g *Git) MakeCommitOnNewBranch(branch string, message string) ExecutableCommand {
	return &commitCommand{
		repositoryPath: g.path,
		branch:         branch,
		message:        message,
		sshEnv:         g.sshEnv,
	}
}

func (g *Git) Push(branch string) ExecutableCommand {
	return &Git{
		args:   append(g.args, "push", "--set-upstream", "origin", branch),
		sshEnv: g.sshEnv,
	}
}

func (g *Git) Execute() (int, error) {
	cmd := exec.CommandContext(context.Background(), "git", g.args...)
	cmd.Env = internal.MergeEnv(os.Environ(), g.sshEnv)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), fmt.Errorf("%s", string(outputBytes))
		}

		return -1, err
	}

	return 0, nil
}

type commitCommand struct {
	repositoryPath string
	branch         string
	message        string
	sshEnv         []string
}

func (c *commitCommand) Execute() (int, error) {
	steps := [][]string{
		{"git", "-C", c.repositoryPath, "checkout", "-b", c.branch},
		{"touch", c.repositoryPath + "/some_file"},
		{"git", "-C", c.repositoryPath, "add", "."},
		{"git", "-C", c.repositoryPath, "commit", "-m", c.message},
	}

	for _, args := range steps {
		exitCode, err := run(args, c.sshEnv)
		if exitCode != 0 || err != nil {
			return exitCode, err
		}
	}

	return 0, nil
}

func run(args []string, sshEnv []string) (int, error) {
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	cmd.Env = internal.MergeEnv(os.Environ(), sshEnv)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), fmt.Errorf("%s", string(outputBytes))
		}

		return -1, err
	}

	return 0, nil
}
