package core

import (
	"os"
	"path"
)

func GetExecutablePath() string {
	return path.Dir(GetExecutable())
}

func GetExecutable() string {
	executable, _ := os.Executable()
	return executable
}
