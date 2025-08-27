package core

import (
	"app/log"
	"io"
	"os"
	"os/exec"
	"sync"
)

func CloneRepository(repositoryURL, destination string) {
	log.Infof("Cloning " + repositoryURL)

	cmd := exec.Command("git", "clone", "--progress", repositoryURL, destination)

	stdoutPipe, stdoutPipeErr := cmd.StdoutPipe()
	if stdoutPipeErr != nil {
		log.Errorf("Error getting StdoutPipe: %v\n", stdoutPipeErr)
		return
	}

	stderrPipe, stdErrPipeErr := cmd.StderrPipe()
	if stdErrPipeErr != nil {
		log.Errorf("Error getting StderrPipe: %v\n", stdErrPipeErr)
		return
	}

	if startErr := cmd.Start(); startErr != nil {
		log.Errorf("Error starting command: %v\n", startErr)
		return
	}

	var waitGroup sync.WaitGroup
	numberOfPipesToWaitFor := 2
	waitGroup.Add(numberOfPipesToWaitFor)

	go func() {
		defer waitGroup.Done()
		copyWithColor(&ColorWriter{Target: os.Stdout}, stdoutPipe)
	}()

	go func() {
		defer waitGroup.Done()
		copyWithColor(&ColorWriter{Target: os.Stderr}, stderrPipe)
	}()

	waitGroup.Wait()
	if err := cmd.Wait(); err != nil {
		log.Warningf("Failed to clone repository, skipping it.")
	}

	log.Successf("Successfully cloned " + repositoryURL)
}

func copyWithColor(dst io.Writer, src io.Reader) {
	numberOfBytes := 1024
	buf := make([]byte, numberOfBytes)
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			text := string(buf[:n])
			_, writeErr := dst.Write([]byte(text))

			if writeErr != nil {
				log.Errorf("Error writing to destination: %v\n", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			log.Errorf("Error reading from source: %v\n", readErr)
			break
		}
	}
}
