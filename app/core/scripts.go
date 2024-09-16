package core

import (
	"app/log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ForScriptInPathDo(path string, do func(scriptPath string, scriptName string)) {
	scripts, _ := filepath.Glob(GetExecutablePath() + path)

	for _, script := range scripts {
		segments := strings.Split(script, "/")
		scriptName := segments[len(segments)-2]

		do(script, scriptName)
	}
}

type ExitCode = int

func ExecuteScript(file string, args []string) ExitCode {
	script := exec.Command(file, args...)
	script.Stdout = os.Stdout
	script.Stdin = os.Stdin
	script.Stderr = os.Stderr
	err := script.Run()

	if err == nil {
		return 0
	}

	return extractExitCode(err, 1)
}

func extractExitCode(err error, defaultExitCode int) int {
	var codeString = strings.TrimPrefix(err.Error(), "exit status ")
	var extractedCode, conversionErr = strconv.Atoi(codeString)

	if conversionErr != nil {
		log.Error("Could not extract exit code from error: " + err.Error())
		return defaultExitCode
	}

	return extractedCode
}
