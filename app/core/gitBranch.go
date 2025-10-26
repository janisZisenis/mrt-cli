package core

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"app/log"
)

func GetCurrentBranchShortName(repoDir string) (string, error) {
	var stdout bytes.Buffer
	err := NewCommandBuilder().
		WithCommand("git").
		WithArgs("-C", repoDir, "rev-parse", "--abbrev-ref", "HEAD").
		WithStdout(&stdout).
		Run()
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
