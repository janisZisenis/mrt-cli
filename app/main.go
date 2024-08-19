package main

import (
	"github.com/go-git/go-git/v5"
	"os"
	"path"
)

func getExecutablePath() string {
	executable, _ := os.Executable()
	return path.Dir(executable)
}

func main() {
	url := "git@github.com:janisZisenis/BoardGames.TDD-London-School.git"
	directory := getExecutablePath() + "/repositories/BoardGames.TDD-London-School"

	// Clone the given repository to the given directory
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		os.Exit(1)
	}
}
