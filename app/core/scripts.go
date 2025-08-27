package core

import (
	"app/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CommandFileName() string {
	return "command"
}

func ForScriptInPathDo(path string, do func(scriptPath string, scriptName string)) {
	scripts, _ := filepath.Glob(path)

	for _, script := range scripts {
		dirPath := filepath.Dir(script)
		scriptName := filepath.Base(dirPath)

		do(script, scriptName)
	}
}

type ExitCode = int

func ExecuteScript(scriptPath string, args []string) ExitCode {
	err := NewCommandBuilder().
		WithCommand(scriptPath).
		WithArgs(args...).
		WithStdout(os.Stdout).
		WithStdin(os.Stdin).
		WithStderr(os.Stderr).
		Run()

	if err == nil {
		return 0
	}

	return extractExitCode(err, 1)
}

func extractExitCode(err error, defaultExitCode int) int {
	var codeString = strings.TrimPrefix(err.Error(), "exit status ")
	var extractedCode, conversionErr = strconv.Atoi(codeString)

	if conversionErr != nil {
		log.Errorf("Could not extract exit code from error: " + err.Error())
		return defaultExitCode
	}

	return extractedCode
}
