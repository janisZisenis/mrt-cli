package core

import (
	"app/log"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func executeGitCommand(repoDir string, args ...string) (string, string, error) {
	cmd := exec.Command("git", args...)

	if repoDir != "" {
		cmd.Dir = repoDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func GetCurrentBranchShortName(repoDir string) (string, error) {
	stdout, stderr, err := executeGitCommand(repoDir, "rev-parse", "--abbrev-ref", "HEAD")

	if err != nil {
		var errorMsg string

		if strings.Contains(stderr, "not a git repository") {
			errorMsg = "the specified directory is not a Git repository"
		} else {
			errorMsg = fmt.Sprintf("failed to retrieve the current branch: %s", stderr)
		}
		return "", errors.New(errorMsg)
	}

	branchName := strings.TrimSpace(stdout)
	if branchName == "" {
		return "", errors.New("could not determine the current branch: output was empty")
	}

	return branchName, nil
}

func CloneRepository(repoURL, destination string) {
	_, stderr, err := executeGitCommand("", "clone", repoURL, destination)

	if err != nil {
		if strings.Contains(stderr, "remote: Repository not found") {
			log.Error("repository not found: verify the URL and your access permissions")
		} else if strings.Contains(stderr, "fatal: destination path") && strings.Contains(stderr, "already exists") {
			log.Error("destination path already exists: choose a different directory")
		} else if strings.Contains(stderr, "not a valid URL") {
			log.Error("the provided URL is not a valid Git repository URL")
		} else {
			log.Error("failed to clone the repository: " + stderr)
		}
	}
}
