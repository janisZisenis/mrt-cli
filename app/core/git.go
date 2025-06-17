package core

import (
	"app/log"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

const notAuthenticatedError = "exit status 128"

func GetCurrentBranchShortName(repoDir string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "rev-parse", "--abbrev-ref", "HEAD")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		log.Error("The given path \"" + repoDir + "\" does not contain a repository.")
		os.Exit(1)
	}

	branchName := strings.TrimSpace(stdout.String())
	if branchName == "" {
		return "", errors.New("could not determine the current branch: output was empty")
	}

	return branchName, nil
}

func CloneRepository(repositoryUrl string, destination string) {
	cmd := exec.Command("git", "clone", "--progress", repositoryUrl, destination)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cloneError := cmd.Run()

	if cloneError != nil {
		if cloneError.Error() == notAuthenticatedError {
			log.Error("You have no access to " + repositoryUrl + ". Please make sure you have a valid ssh key in place.")
		}
	} else {
		log.Success("Successfully cloned " + repositoryUrl)
	}
}
