package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
)

func main() {
	var teamInfo = readTeamInfo()

	if len(teamInfo.Repositories) == 0 {
		fmt.Println("The " + teamFileName + " file does not contain any repositories")
		os.Exit(1)
	}

	for _, repository := range teamInfo.Repositories {
		url := "git@github.com:janisZisenis/" + repository + ".git"
		directory := getExecutablePath() + "/" + teamInfo.RepositoriesPath + "/" + repository

		// Clone the given repository to the given directory
		_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
			URL:               url,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		if cloneError != nil {
			fmt.Println(cloneError)
			os.Exit(1)
		}
	}
}
