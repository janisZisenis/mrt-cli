package clonerepositories

import (
	"strings"

	"mrt-cli/app/core"
	"mrt-cli/app/log"
)

func CloneRepositories(teamInfo core.TeamInfo) {
	log.Infof("Start cloning repositories into \"" + teamInfo.RepositoriesPath + "\"")
	for _, repositoryURL := range teamInfo.Repositories {
		repositoryName := getRepositoryName(repositoryURL)
		folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
		repositoryDirectory := getRepositoryPath(teamInfo.RepositoriesPath, folderName)

	    log.Infof("Cloning " + repositoryURL)
		if err := core.CloneRepository(repositoryURL, repositoryDirectory); err != nil {
            log.Warningf("Failed to clone repository, skipping it.")
			continue
		}
		log.Successf("Successfully cloned " + repositoryURL)
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
			folderName = strings.TrimPrefix(folderName, prefix)
			break
		}
	}
	return folderName
}

func getRepositoryPath(repositoriesPath string, folderName string) string {
	return core.GetExecutionPath() + "/" + repositoriesPath + "/" + folderName
}
