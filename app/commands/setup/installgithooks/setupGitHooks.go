package installgithooks

import (
	"path/filepath"

	"mrt-cli/app/core"
	"mrt-cli/app/log"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	reposDir := filepath.Join(core.GetExecutionPath(), teamInfo.RepositoriesPath)
	log.Infof("Installing git-hooks to repositories located in \"" + reposDir + "\"")

	pattern := filepath.Join(reposDir, "*", core.GitMetadataDir)
	repositories, _ := filepath.Glob(pattern)
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
	for _, hook := range core.GetGitHooks() {
		writeGitHook(repository, hook)
	}
}
