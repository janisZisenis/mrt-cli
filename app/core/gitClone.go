package core

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"

	"mrt-cli/app/log"
)

const bufferSize = 64 * 1024

func CloneRepository(repositoryURL, destination string) error {
	var waitGroup sync.WaitGroup
	numberOfPipesToWaitFor := 3
	waitGroup.Add(numberOfPipesToWaitFor)

	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	bufferedStdout := bufio.NewReaderSize(stdoutReader, bufferSize)
	bufferedStderr := bufio.NewReaderSize(stderrReader, bufferSize)

	cancel, waitUntilCloneFinished, startErr := NewCommandBuilder().
		WithCommand("git").
		WithArgs("clone", "--progress", repositoryURL, destination).
		WithStdout(stdoutWriter).
		WithStderr(stderrWriter).
		Start()

	defer cancel()

	if startErr != nil {
		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
		return startErr
	}

	var cloneErr error
	go func() {
		defer waitGroup.Done()
		if waitErr := waitUntilCloneFinished(); waitErr != nil {
			cloneErr = waitErr
		}

		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stdout, bufferedStdout)
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(os.Stderr, bufferedStderr)
	}()

	waitGroup.Wait()

	if cloneErr != nil {
		return cloneErr
	}

	return nil
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
