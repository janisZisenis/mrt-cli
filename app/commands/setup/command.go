package setup

import (
	"app/core"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "setup",
		Short: "Sets up you machine for development",
		Run:   command,
	}
	return command
}

var repositoryNotFoundError = "repository not found"
var notAuthenticatedError = "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain"

func command(cmd *cobra.Command, args []string) {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		teamInfo := core.LoadTeamConfiguration()

		if len(teamInfo.Repositories) == 0 {
			fmt.Println("Your team file does not contain any repositories")
			os.Exit(1)
		}

		for _, repositoryUrl := range teamInfo.Repositories {
			repositoryName := getRepositoryName(repositoryUrl)
			folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
			directory := getDirectory(teamInfo.RepositoriesPath, folderName)

			clone(repositoryUrl, directory)

			var hooksPath = directory + "/.git/hooks/"
			writePreCommitHook(hooksPath)
			writePrePushHook(hooksPath)
		}
	}
}

func getPreCommitHook() string {
	return `
#!/bin/bash -e

branch="$(git rev-parse --abbrev-ref HEAD)"
` + core.GetExecutable() + ` githook --branch $branch --hook-name pre-commit`
}

func writePreCommitHook(hooksPath string) {
	_ = os.MkdirAll(hooksPath, os.ModePerm)
	err := os.WriteFile(hooksPath+"pre-commit", []byte(getPreCommitHook()), 0755)
	if err != nil {
		fmt.Printf("unable to write file: %w", err)
	}
}

func getPrePushHook() string {
	return `
#!/bin/bash -e

branch="$(git rev-parse --abbrev-ref HEAD)"
` + core.GetExecutable() + ` githook --branch $branch --hook-name pre-push`
}

func writePrePushHook(hooksPath string) {
	_ = os.MkdirAll(hooksPath, os.ModePerm)
	err := os.WriteFile(hooksPath+"pre-push", []byte(getPrePushHook()), 0755)
	if err != nil {
		fmt.Printf("unable to write file: %w", err)
	}
}

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

func getRepositoryName(repositoryUrl string) string {
	return strings.Trim(repositoryUrl[strings.LastIndex(repositoryUrl, "/")+1:], ".git")
}

func getDirectory(repositoriesPath string, folderName string) string {
	return core.GetExecutablePath() + "/" + repositoriesPath + "/" + folderName
}

func getFolderName(repositoryName string, prefixes []string) string {
	var folderName = repositoryName
	for _, prefix := range prefixes {
		if strings.HasPrefix(folderName, prefix) {
			folderName = strings.Replace(folderName, prefix, "", 1)
		}
	}
	return folderName
}
