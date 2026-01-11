package core

import (
	"os"
	"path/filepath"
	"sync"
)

//nolint:gochecknoglobals // This is a global flag that we can set on every command to change the execution path
var (
	teamDirectory *string
	mu            sync.RWMutex
)

func SetTeamDirectory(directory *string) {
	mu.Lock()
	defer mu.Unlock()
	teamDirectory = directory
}

func GetExecutionPath() string {
	mu.RLock()
	defer mu.RUnlock()
	if teamDirectory != nil {
		return *teamDirectory
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
