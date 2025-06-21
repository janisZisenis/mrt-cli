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
	"sync"
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

func CloneRepository(repositoryUrl, destination string) {
	log.Info("Cloning " + repositoryUrl)

	cmd := exec.Command("git", "clone", "--progress", repositoryUrl, destination)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error getting StdoutPipe: %v\n", err)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Error("Error getting StderrPipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Error("Error starting command: %v\n", err)
		return
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stdout, stdoutPipe)
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stderr, stderrPipe)
	}()

	if err := cmd.Wait(); err != nil {
		log.Warning("Failed to clone repository, skipping it.")
	}

	waitGroup.Wait()

	log.Success("Successfully cloned " + repositoryUrl)
}

func copyWithColor(dst io.Writer, src io.Reader) {
	const (
		purple = "\033[35;1m"
		reset  = "\033[0m"
	)

	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err = fmt.Fprintf(dst, "%s%s%s", purple, buf[:n], reset)

			if err != nil {
				log.Error("Error writing to destination: %v\n", err)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("Error reading from source: %v\n", err)
			break
		}
	}
}
