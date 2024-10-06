package core

import (
	"os"
	"path/filepath"
)

var TeamDirectory string

func GetExecutionPath() string {
	if TeamDirectory != "" {
		return TeamDirectory
	}

	pwd, _ := os.Getwd()
	return pwd
}

func GetAbsoluteExecutionPath() string {
	absolute, _ := filepath.Abs(GetExecutionPath())
	return absolute
}

func GetExecutableName() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
