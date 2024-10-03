package core

import (
	"os"
	"path/filepath"
)

func GetExecutionPath() string {
	pwd, _ := os.Getwd()
	return pwd
}

func GetExecutableName() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
