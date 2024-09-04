package setup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
)

var repositoryNotFoundError = "repository not found"
var notAuthenticatedError = "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain"

func clone(repository string, directory string) {
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
