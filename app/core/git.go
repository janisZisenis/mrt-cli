package core

import (
	"app/log"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

const purple = "\033[35m"
const reset = "\033[0m"

func CloneRepository(repositoryUrl, destination string) error {
	log.Info("Cloning " + repositoryUrl + " into " + destination)

	cmd := exec.Command("git", "clone", repositoryUrl, destination)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		printColoredOutput(stderrBuf.String())
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	printColoredOutput(stdoutBuf.String())

	fmt.Println("Cloning repositories done")

	return nil
}

func printColoredOutput(output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("%s\t%s%s\n", purple, line, reset)
		}
	}
}
