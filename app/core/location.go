package core

import (
	"os"
	"path/filepath"
)

func GetExecutionPath() string {
    pwd, _ := os.Getwd()
    return pwd
}

func GetExecutable() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
