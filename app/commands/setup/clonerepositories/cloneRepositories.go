package clonerepositories

import (
	"strings"

	"app/core"
	"app/log"
)

func CloneRepositories(teamInfo core.TeamInfo) {
	log.Infof("Start cloning repositories into \"" + teamInfo.RepositoriesPath + "\"")
	for _, repositoryURL := range teamInfo.Repositories {
		repositoryName := getRepositoryName(repositoryURL)
		folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
		repositoryDirectory := getRepositoryPath(teamInfo.RepositoriesPath, folderName)

		core.CloneRepository(repositoryURL, repositoryDirectory)
	}
	log.Successf("Cloning repositories done")
}

func getRepositoryName(repositoryURL string) string {
	return strings.TrimSuffix(repositoryURL[strings.LastIndex(repositoryURL, "/")+1:], ".git")
}

func getFolderName(repositoryName string, prefixes []string) string {
	folderName := repositoryName
	for _, prefix := range prefixes {
		if strings.HasPrefix(folderName, prefix) {
			folderName = strings.Replace(folderName, prefix, "", 1)
		}
	}
	return folderName
}

func getRepositoryPath(repositoriesPath string, folderName string) string {
	return core.GetExecutionPath() + "/" + repositoriesPath + "/" + folderName
}
