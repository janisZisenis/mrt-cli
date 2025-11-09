package core

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"mrt-cli/app/log"
)

func CloneRepository(repositoryURL, destination string) {
	log.Infof("Cloning " + repositoryURL)

	var waitGroup sync.WaitGroup
	numberOfPipesToWaitFor := 2
	waitGroup.Add(numberOfPipesToWaitFor)

	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	cancel, wait, startErr := NewCommandBuilder().
		WithCommand("git").
		WithArgs("clone", "--progress", repositoryURL, destination).
		WithStdout(stdoutWriter).
		WithStderr(stderrWriter).
		Start()

	defer cancel()

	if startErr != nil {
		log.Errorf("Error starting command: %v\n", startErr)
		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
		return
	}

	go func() {
		if waitErr := wait(); waitErr != nil {
			if errors.Is(waitErr, context.DeadlineExceeded) {
				log.Warningf("Repository clone timed out after 5 minutes")
			} else {
				log.Warningf("Failed to clone repository, skipping it.")
			}
		}

		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stdout, stdoutReader)
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stderr, stderrReader)
	}()

	waitGroup.Wait()

	log.Successf("Successfully cloned " + repositoryURL)
}

func copyWithColor(dst io.Writer, src io.Reader) {
	colorWriter := ColorWriter{Target: dst}

	numberOfBytes := 1024
	buf := make([]byte, numberOfBytes)
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			text := string(buf[:n])
			_, writeErr := colorWriter.Write([]byte(text))

			if writeErr != nil {
				log.Errorf("Error writing to destination: %v\n", writeErr)
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			log.Errorf("Error reading from source: %v\n", readErr)
			break
		}
	}
}
