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

const purple = "\033[35;1m"
const reset = "\033[0m"

func CloneRepository(repositoryUrl, destination string) {
	log.Info("Cloning " + repositoryUrl + " into " + destination)

	cmd := exec.Command("git", "clone", "--progress", repositoryUrl, destination)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create stdout pipe: %v\n", err)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to create stderr pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start git clone command: %v\n", err)
		return
	}

	go processStream(stdoutPipe)
	go processStream(stderrPipe)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("%sFailed to clone repository: %v%s\n", reset, err, reset)
		return
	}

	log.Success("Successfully cloned" + repositoryUrl)
}

func processStream(stream io.ReadCloser) {
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			subLines := strings.Split(line, "\r")

			for _, subLine := range subLines {
				if strings.TrimSpace(subLine) != "" {
					fmt.Printf("%s    %s%s\n", purple, strings.TrimSpace(subLine), reset)
				}
			}
		}
	}
}
