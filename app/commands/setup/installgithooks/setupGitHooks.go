package installgithooks

import (
	"mrt-cli/app/core"
	"mrt-cli/app/log"
	"path/filepath"
)

const (
	gitMetadataDir = ".git"
)

func setupGitHooks(teamInfo core.TeamInfo) {
	reposDir := filepath.Join(getAbsoluteExecutionPath(), teamInfo.RepositoriesPath)
	log.Infof("Installing git-hooks to repositories located in \"" + reposDir + "\"")

	pattern := filepath.Join(reposDir, "*", gitMetadataDir)
	repositories, err := filepath.Glob(pattern)
	if err != nil {
		log.Errorf("Failed to find repositories: %v", err)
		return
	}
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
