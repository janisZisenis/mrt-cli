package core

import (
	"os"
	"path/filepath"
)

func GetAbsoluteExecutionPath() string {
	pwd, _ := os.Getwd()
	absolute, _ := filepath.Abs(pwd)
	return absolute
}

func GetExecutableName() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
