package core

import (
	"app/log"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
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

func CloneRepository(repositoryUrl, destination string) {
	log.Info("Cloning " + repositoryUrl + " into " + destination)

	cmd := exec.Command("git", "clone", "--progress", repositoryUrl, destination)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%s\tFailed to create stdout pipe: %v%s\n", purple, err, reset)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("%s\tFailed to create stderr pipe: %v%s\n", purple, err, reset)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("%s\tFailed to start git clone command: %v%s\n", purple, err, reset)
		return
	}

	go processStream(stdoutPipe)
	go processStream(stderrPipe)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s\tFailed to clone repository: %v%s\n", purple, err, reset)
		return
	}

	fmt.Printf("Successfully cloned %s\n", repositoryUrl)
}

func processStream(stream io.ReadCloser) {
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			fmt.Printf("%s\t%s%s\n", purple, line, reset)
		}
	}
}
