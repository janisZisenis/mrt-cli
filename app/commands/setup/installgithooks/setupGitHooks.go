package installgithooks

import (
	"app/core"
	"app/log"
	"path/filepath"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	log.Infof("Installing git-hooks to repositories located in \"" + core.GetExecutionPath() + "/" +
		teamInfo.RepositoriesPath + "\"")

	repositories, _ := filepath.Glob(core.GetExecutionPath() + "/" + teamInfo.RepositoriesPath + "/*/.git")
	if len(repositories) == 0 {
		log.Infof("Did not find any repositories. Skip installing git-hooks.")
	}

	for _, r := range repositories {
		log.Infof("Installing git-hooks to \"" + r + "\"")
		writeHooks(r)
		log.Successf("Done installing git-hooks to \"" + r + "\"")
	}

	log.Successf("Done installing git-hooks.")
}

func writeHooks(repository string) {
	for _, hook := range core.GitHooks {
		writeGitHook(repository, hook)
	}
}
