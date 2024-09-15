package cloneRepositories

import (
	"app/core"
	"app/log"
	"strings"
)

func CloneRepositories(teamInfo core.TeamInfo) {
	log.Info("Start cloning repositories into \"" + teamInfo.RepositoriesPath + "\"")
	for _, repositoryUrl := range teamInfo.Repositories {
		repositoryName := getRepositoryName(repositoryUrl)
		folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
		repositoryDirectory := getRepositoryPath(teamInfo.RepositoriesPath, folderName)

		log.Info("Cloning " + repositoryUrl + " into " + "repositories" + "/" + folderName)
		clone(repositoryUrl, repositoryDirectory)
	}
	log.Success("Cloning repositories done")
}

func getRepositoryName(repositoryUrl string) string {
	return strings.TrimSuffix(repositoryUrl[strings.LastIndex(repositoryUrl, "/")+1:], ".git")
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

func getRepositoryPath(repositoriesPath string, folderName string) string {
	return core.GetExecutablePath() + "/" + repositoriesPath + "/" + folderName
}
