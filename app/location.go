package main

import (
	"os"
	"path"
)

func getExecutablePath() string {
	executable, _ := os.Executable()
	return path.Dir(executable)
}
