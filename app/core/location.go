package core

import (
	"os"
	"path/filepath"
)

var (
	TeamDirectory string
)

func GetExecutionPath() string {
	if TeamDirectory != "" {
		return TeamDirectory
	}

	pwd, _ := os.Getwd()
	return pwd
}

func GetAbsoluteExecutionPath() string {
	absolutePath, _ := filepath.Abs(GetExecutionPath())
	return absolutePath
}

func GetExecutableName() string {
	executablePath, _ := os.Executable()
	return filepath.Base(executablePath)
}
