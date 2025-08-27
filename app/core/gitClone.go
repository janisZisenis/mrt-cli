package core

import (
	"app/log"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

func CloneRepository(repositoryURL, destination string) {
	log.Infof("Cloning " + repositoryURL)

	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "clone", "--progress", repositoryURL, destination)
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	var wg sync.WaitGroup
	wg.Add(2)

	if err := cmd.Start(); err != nil {
		log.Errorf("Error starting command: %v\n", err)
		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
		return
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Warningf("Failed to clone repository, skipping it.")
		}

		_ = stdoutWriter.Close()
		_ = stderrWriter.Close()
	}()

	go func() {
		defer wg.Done()
		copyWithColor(os.Stdout, stdoutReader)
	}()

	go func() {
		defer wg.Done()
		copyWithColor(os.Stderr, stderrReader)
	}()

	wg.Wait()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Warningf("Repository clone timed out after 5 minutes")
		return
	}

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
