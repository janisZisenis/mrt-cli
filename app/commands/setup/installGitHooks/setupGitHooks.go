package installGitHooks

import (
	"app/core"
	"app/log"
	"path/filepath"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	log.Info("Installing git-hooks to repositories located in \"" + core.GetExecutablePath() + "/" + teamInfo.RepositoriesPath + "\"")

	repositories, _ := filepath.Glob(core.GetExecutablePath() + "/" + teamInfo.RepositoriesPath + "/*/.git")
	if len(repositories) == 0 {
		log.Info("Did not find any repositories. Skip installing git-hooks.")
	}

	for _, r := range repositories {
		log.Info("Installing git-hooks to \"" + r + "\"")
		writeHooks(r)
		log.Success("Done installing git-hooks to \"" + r + "\"")
	}

	log.Success("Done installing git-hooks.")
}

func writeHooks(repository string) {
	for _, hook := range core.GitHooks {
		writeGitHook(repository, hook)
	}
}
