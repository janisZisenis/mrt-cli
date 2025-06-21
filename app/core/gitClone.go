package core

import (
	"app/log"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"os/exec"
	"sync"
)

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
		log.Info("Failed to clone repository, skipping it.")
	}

	waitGroup.Wait()

	log.Success("Successfully cloned " + repositoryUrl)
}

func copyWithColor(dst io.Writer, src io.Reader) {
	purpleFatih := color.New(color.FgHiMagenta).SprintFunc()

	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			text := string(buf[:n])
			_, err = fmt.Fprintf(dst, "%s", purpleFatih(text))

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
