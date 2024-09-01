package main

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"strings"
)

var repositoryNotFoundError = "repository not found"
var notAuthenticatedError = "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		var teamInfo = readTeamInfo()

		if len(teamInfo.Repositories) == 0 {
			fmt.Println("The " + teamFileName + " file does not contain any repositories")
			os.Exit(1)
		}

		for _, repository := range teamInfo.Repositories {
			repositoryName := strings.Trim(repository[strings.LastIndex(repository, "/")+1:], ".git")

			var cloneDirectory = repositoryName
			for _, prefix := range teamInfo.RepositoriesPrefixes {
				if strings.HasPrefix(cloneDirectory, prefix) {
					cloneDirectory = strings.Replace(cloneDirectory, prefix, "", 1)
				}
			}

			directory := getExecutablePath() + "/" + teamInfo.RepositoriesPath + "/" + cloneDirectory

			_, cloneError := git.PlainClone(directory, false, &git.CloneOptions{
				URL:               repository,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if cloneError != nil {
				if errors.Is(cloneError, git.ErrRepositoryAlreadyExists) {
					fmt.Println("Repository " + repository + " already exists. Skipping it")
				}

				if cloneError.Error() == repositoryNotFoundError {
					fmt.Println("Repository " + repository + " was not found. Skipping it")
				}

				if cloneError.Error() == notAuthenticatedError {
					fmt.Println("You have no access to " + repository + ". Please make sure you have a valid ssh key in place.")
				}
			}
		}
	}
}
