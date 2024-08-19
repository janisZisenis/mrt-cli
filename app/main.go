package main

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"strings"
)

func main() {
	var teamInfo = readTeamInfo()

	if len(teamInfo.Repositories) == 0 {
		fmt.Println("The " + teamFileName + " file does not contain any repositories")
		os.Exit(1)
	}

	for _, repository := range teamInfo.Repositories {
		repositoryName := strings.Trim(repository[strings.LastIndex(repository, "/")+1:], ".git")
		directory := getExecutablePath() + "/" + teamInfo.RepositoriesPath + "/" + repositoryName

		_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
			URL:               repository,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		if cloneError != nil && errors.Is(cloneError, git.ErrRepositoryAlreadyExists) {
			fmt.Println("Repository " + repository + " already exists. Skipping it")
			continue
		}
	}
}
