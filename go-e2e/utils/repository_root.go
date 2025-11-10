package utils

import (
	"os"
	"path/filepath"
)

func GetRepoRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	for {
		info, err := os.Stat(filepath.Join(dir, ".git"))

		if err == nil && info.IsDir() {
			return dir
		} else {
			parent := filepath.Dir(dir)

			if parent == dir {
				panic("Reached root directory without finding .git folder")
			}

			dir = parent
		}
	}
}
