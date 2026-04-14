package installgithooks

import (
	"os"
	"path/filepath"
)

func getAbsoluteExecutionPath() string {
	pwd, _ := os.Getwd()
	absolute, _ := filepath.Abs(pwd)
	return absolute
}

func getExecutableName() string {
	executable, _ := os.Executable()
	return filepath.Base(executable)
}
