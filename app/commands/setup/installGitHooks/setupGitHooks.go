package installGitHooks

import (
	"app/core"
	"path/filepath"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	repositories, _ := filepath.Glob(core.GetExecutablePath() + "/" + teamInfo.RepositoriesPath + "/*/.git")
	for _, r := range repositories {
		writeHooks(r)
	}
}

func writeHooks(repository string) {
	for _, hook := range core.GitHooks {
		writeGitHook(repository, hook)
	}
}
