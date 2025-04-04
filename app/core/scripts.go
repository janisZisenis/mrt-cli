package core

import (
	"app/log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func CommandFileName() string {
	return "command"
}

func ForScriptInPathDo(path string, do func(scriptPath string, scriptName string)) {
	scripts, err := filepath.Glob(path)
	if err != nil {
		log.Error("Error finding scripts in path: " + err.Error())
		return
	}

	for _, script := range scripts {
		dirPath := filepath.Dir(script)
		scriptName := filepath.Base(dirPath)

		do(script, scriptName)
	}
}

type ExitCode = int

func ExecuteScript(scriptPath string, args []string) ExitCode {
	script := exec.Command(scriptPath, args...)
	script.Stdout = os.Stdout
	script.Stdin = os.Stdin
	script.Stderr = os.Stderr
	err := script.Run()

	if err == nil {
		return 0
	}

	return ExitCode(extractExitCode(err, 1))
}

func extractExitCode(err error, defaultExitCode int) int {
	codeString := strings.TrimPrefix(err.Error(), "exit status ")
	extractedCode, conversionErr := strconv.Atoi(codeString)

	if conversionErr != nil {
		log.Error("Could not extract exit code from error: " + err.Error())
		return defaultExitCode
	}

	return extractedCode
}
