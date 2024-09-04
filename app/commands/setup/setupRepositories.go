package setup

import (
	"app/core"
	"fmt"
	"os"
	"strings"
)

func setupRepositories(teamInfo core.TeamInfo) {
	if len(teamInfo.Repositories) == 0 {
		fmt.Println("Your team file does not contain any repositories")
		os.Exit(1)
	}

	for _, repositoryUrl := range teamInfo.Repositories {
		repositoryName := getRepositoryName(repositoryUrl)
		folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
		repositoryDirectory := getDirectory(teamInfo.RepositoriesPath, folderName)

		clone(repositoryUrl, repositoryDirectory)

		writeGitHook(repositoryDirectory, core.PreCommit)
		writeGitHook(repositoryDirectory, core.PrePush)
		writeGitHook(repositoryDirectory, core.CommitMsg)
	}
}

func getRepositoryName(repositoryUrl string) string {
	return strings.Trim(repositoryUrl[strings.LastIndex(repositoryUrl, "/")+1:], ".git")
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

func getDirectory(repositoriesPath string, folderName string) string {
	return core.GetExecutablePath() + "/" + repositoriesPath + "/" + folderName
}
