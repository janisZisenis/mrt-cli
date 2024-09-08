package cloneRepositories

import (
	"app/log"
	"errors"
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
			log.Warning("Repository " + repository + " already exists. Skipping it")
		}

		if cloneError.Error() == repositoryNotFoundError {
			log.Error("Repository " + repository + " was not found. Skipping it")
		}

		if cloneError.Error() == notAuthenticatedError {
			log.Error("You have no access to " + repository + ". Please make sure you have a valid ssh key in place.")
		}
	} else {
		log.Success("Successfully cloned " + repository)
	}
}
