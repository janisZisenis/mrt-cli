package installGitHooks

import (
	"app/core"
	"strings"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	repositoryDirectory := getRepositoryDir(teamInfo, teamInfo.Repositories[0])

	for _, hook := range core.GitHooks {
		writeGitHook(repositoryDirectory, hook)
	}
	//for _, repositoryUrl := range teamInfo.Repositories {
	//	repositoryDirectory := getRepositoryDir(teamInfo, repositoryUrl)
	//
	//	for _, hook := range core.GitHooks {
	//		writeGitHook(repositoryDirectory, hook)
	//	}
	//}
}

func getRepositoryDir(teamInfo core.TeamInfo, repositoryUrl string) string {
	repositoryName := getRepositoryName(repositoryUrl)
	folderName := getFolderName(repositoryName, teamInfo.RepositoriesPrefixes)
	return getDirectory(teamInfo.RepositoriesPath, folderName)
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

func getDirectory(repositoriesPath string, folderName string) string {
	return core.GetExecutablePath() + "/" + repositoriesPath + "/" + folderName
}
