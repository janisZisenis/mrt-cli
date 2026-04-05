package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"mrt-cli/go-e2e/internal"
)

type BaseCommand interface {
	InDirectory(path string) DirectedCommand
	Clone(repositoryURL string, destination string) ExecutableCommand
}

type DirectedCommand interface {
	MakeCommitOnNewBranch(branch string, message string) ExecutableCommand
	Push(branch string) ExecutableCommand
	DeleteRemoteBranchIfExists(branch string) ExecutableCommand
	GetLastCommitMessage() (string, error)
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

func (g *Git) DeleteRemoteBranchIfExists(branch string) ExecutableCommand {
	return &deleteRemoteBranchIfExistsCommand{
		git: &Git{args: append(g.args, "push", "origin", "--delete", branch), sshEnv: g.sshEnv},
	}
}

func (g *Git) GetLastCommitMessage() (string, error) {
	args := append(g.args, "log", "-1", "--pretty=%B") //nolint:gocritic // intentional append to copy slice
	cmd := exec.CommandContext(context.Background(), "git", args...)
	cmd.Env = internal.MergeEnv(os.Environ(), g.sshEnv)

	outputBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w", err)
	}

	return strings.TrimSpace(string(outputBytes)), nil
}

type deleteRemoteBranchIfExistsCommand struct {
	git *Git
}

func (d *deleteRemoteBranchIfExistsCommand) Execute() (int, error) {
	exitCode, err := d.git.Execute()
	if err != nil && strings.Contains(err.Error(), "remote ref does not exist") {
		return 0, nil
	}

	return exitCode, err
}

func (g *Git) Execute() (int, error) {
	//nolint:gosec // args are controlled by internal callers
	cmd := exec.CommandContext(context.Background(), "git", g.args...)
	cmd.Env = internal.MergeEnv(os.Environ(), g.sshEnv)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
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
	//nolint:gosec // args are controlled by internal callers
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	cmd.Env = internal.MergeEnv(os.Environ(), sshEnv)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), fmt.Errorf("%s", string(outputBytes))
		}

		return -1, err
	}

	return 0, nil
}
