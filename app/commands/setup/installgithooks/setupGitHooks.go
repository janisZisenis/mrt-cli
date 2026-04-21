package installgithooks

import (
	"mrt-cli/app/core"
	"mrt-cli/app/log"
	"os"
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

	teamDir := getAbsoluteExecutionPath()

	for _, r := range repositories {
		log.Infof("Installing git-hooks to \"" + r + "\"")
		hookFileDir := filepath.Join(r, gitHooksDir)
		relPath, _ := filepath.Rel(hookFileDir, teamDir)
		writeHooks(r, relPath)
		log.Successf("Done installing git-hooks to \"" + r + "\"")
	}

	teamGitDir := filepath.Join(teamDir, gitMetadataDir)
	if _, statErr := os.Stat(teamGitDir); statErr == nil {
		log.Infof("Installing git-hooks to team repository")
		teamHookFileDir := filepath.Join(teamGitDir, gitHooksDir)
		relPath, _ := filepath.Rel(teamHookFileDir, teamDir)
		writeHooks(teamGitDir, relPath)
		log.Successf("Done installing git-hooks to team repository.")
	}

	log.Successf("Done installing git-hooks.")
}

func writeHooks(repository string, relativePathToTeamDir string) {
	for _, hook := range core.GetGitHooks() {
		writeGitHook(repository, hook, relativePathToTeamDir)
	}
}
