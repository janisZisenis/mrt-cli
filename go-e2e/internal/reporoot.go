package internal

import (
	"os"
	"path/filepath"
)

func GetRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	for {
		info, statErr := os.Stat(filepath.Join(dir, ".git"))

		if statErr == nil && info.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			panic("Reached root directory without finding .git folder")
		}

		dir = parent
	}
}
