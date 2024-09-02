package setup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:   "setup",
	Short: "Sets up you machine for development.",
	Long:  "Test",
	Run:   command,
}

var repositoryNotFoundError = "repository not found"
var notAuthenticatedError = "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain"

func command(cmd *cobra.Command, args []string) {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		teamInfo, err := LoadTeamConfiguration()
		if err != nil {
			fmt.Println("Could not read team file. Please make sure a \"team.json\" file exists next " +
				"to the executable and that it follows proper JSON syntax")
			os.Exit(1)
		}

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
			writePreCommitHook(hooksPath, teamInfo)
		}
	}
}

func getPreCommitHook(blockedBranches []string) string {
	return `
#!/bin/bash -e

branch="$(git rev-parse --abbrev-ref HEAD)"

blocked_branches=(` + printSlice(blockedBranches) + ` )
if [[ "${blocked_branches[@]}" =~ $branch ]]
then
	echo "Action \"commit\" not allowed on branch \"$branch\""
	exit 1
fi
`
}

func printSlice(s []string) string {
	if len(s) == 0 {
		return ""
	}

	return " " + fmt.Sprint(s[0]) + printSlice(s[1:])
}

func writePreCommitHook(hooksPath string, teamInfo *TeamInfo) {
	_ = os.MkdirAll(hooksPath, os.ModePerm)
	err := os.WriteFile(hooksPath+"pre-commit", []byte(getPreCommitHook(teamInfo.BlockedBranches)), 0755)
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
	return GetExecutablePath() + "/" + repositoriesPath + "/" + folderName
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
