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

func GetExecutableName() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
