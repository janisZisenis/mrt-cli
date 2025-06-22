package core

import (
	"os"
	"path/filepath"
)

var TeamDirectory *string

func SetTeamDirectory(directory *string) {
	TeamDirectory = directory
}

func GetExecutionPath() string {
	if TeamDirectory != nil {
		return *TeamDirectory
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
