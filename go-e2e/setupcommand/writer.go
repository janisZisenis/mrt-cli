package setupcommand

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"mrt-cli/go-e2e/assert"
)

const (
	setupDir       = "setup"
	commandFile    = "command"
	spyFileSuffix  = "Executed"
	execPermission = 0o750
)

type Writer struct {
	teamDir string
}

func NewWriter(teamDir string) *Writer {
	return &Writer{teamDir: teamDir}
}

func (w *Writer) commandPath(commandName string) string {
	return filepath.Join(w.teamDir, setupDir, commandName, commandFile)
}

func (w *Writer) spyFilePath(commandName string) string {
	return w.commandPath(commandName) + spyFileSuffix
}

func (w *Writer) writeScript(commandName string, content string) {
	path := w.commandPath(commandName)

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		panic("setupcommand: failed to create command directory: " + err.Error())
	}

	if err := os.WriteFile(path, []byte(content), execPermission); err != nil {
		panic("setupcommand: failed to write command file: " + err.Error())
	}
}

// WriteSpyCommand writes a command script that records the arguments it receives.
func (w *Writer) WriteSpyCommand(commandName string) {
	script := fmt.Sprintf(`#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "$@" > "$SCRIPT_DIR/%s%s"
`, commandFile, spyFileSuffix)
	w.writeScript(commandName, script)
}

func (w *Writer) AssertSpyWasCalledWith(t *testing.T, commandName string, expectedArgs string) {
	t.Helper()
	assert.FileHasContent(t, w.spyFilePath(commandName), w.teamDir+" "+expectedArgs+"\n")
}

func (w *Writer) AssertSpyWasNotCalled(t *testing.T, commandName string) {
	t.Helper()
	assert.FileDoesNotExist(t, w.spyFilePath(commandName))
}
