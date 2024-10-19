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
	scripts, _ := filepath.Glob(path)

	for _, script := range scripts {
		dirPath := filepath.Dir(script)
		scriptName := filepath.Base(dirPath)

		do(script, scriptName)
	}
}

type ExitCode = int

func ExecuteScript(scriptPath string, args []string) ExitCode {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Error("The script " + scriptPath + " does not exist.")
	}

	script := exec.Command("C:/Program Files/Git/bin/bash.EXE", "./run/binary-location/command.sh")
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
