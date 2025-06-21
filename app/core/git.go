package core

import (
	"app/log"
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

const (
	purple = "\033[35;1m"
	reset  = "\033[0m"
)

func CloneRepository(repositoryUrl, destination string) {
	log.Info("Cloning " + repositoryUrl + " into " + destination)

	cmd := exec.Command("git", "clone", "--progress", repositoryUrl, destination)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error getting StdoutPipe: %v\n", err)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error getting StderrPipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	go func() {
		copyWithColor(os.Stdout, stdoutPipe, purple, reset)
	}()

	go func() {
		copyWithColor(os.Stderr, stderrPipe, purple, reset)
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("%sFailed to clone repository: %v%s\n", reset, err, reset)
		return
	}

	log.Success("Successfully cloned " + repositoryUrl)
}

func copyWithColor(dst io.Writer, src io.Reader, colorStart, colorReset string) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			fmt.Fprintf(dst, "%s%s%s", colorStart, buf[:n], colorReset)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading from source: %v\n", err)
			break
		}
	}
}
