package core

import (
	"app/log"
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetCurrentBranchShortName(repoDir string) (string, error) {
	stdout, err := execute("git", "-C", repoDir, "rev-parse", "--abbrev-ref", "HEAD")

	if err != nil {
		log.Errorf("The given path \"" + repoDir + "\" does not contain a repository.")
		os.Exit(1)
	}

	branchName := strings.TrimSpace(stdout.String())
	if branchName == "" {
		return "", errors.New("could not determine the current branch: output was empty")
	}

	return branchName, nil
}

func execute(command string, args ...string) (bytes.Buffer, error) {
	const commandTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if ctxErr := ctx.Err(); ctxErr != nil {
		return stdout, ctxErr
	}

	if err != nil {
		return stdout, err
	}

	return stdout, nil
}
