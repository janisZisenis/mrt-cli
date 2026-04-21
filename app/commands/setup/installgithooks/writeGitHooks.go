package installgithooks

import (
	"mrt-cli/app/commands/githook"
	"mrt-cli/app/log"
	"os"
	"path/filepath"
)

const (
	gitHooksDir = "hooks"
)

func getHookTemplate(relativePathToTeamDir string) string {
	return `
#!/bin/bash -e

hook_name=$(basename "$0")
` + getExecutableName() + ` --team-dir "$(cd "$(dirname "$0")/` + relativePathToTeamDir + `" && pwd)" ` + githook.CommandName + ` --hook-name "$hook_name" --repository-path $PWD $@`
}

func writeGitHook(repositoryDirectory string, hookName string, relativePathToTeamDir string) {
	hooksPath := filepath.Join(repositoryDirectory, gitHooksDir)
	// #nosec G301 - githooks folder needs 0700 to be private (owner only)
	if err := os.MkdirAll(hooksPath, 0o700); err != nil {
		log.Errorf("Failed to create hooks directory %q: %v", hooksPath, err)
		return
	}
	hookFilePath := filepath.Join(hooksPath, hookName)
	// #nosec G306 - git hooks need to be executable by owner
	if err := os.WriteFile(hookFilePath, []byte(getHookTemplate(relativePathToTeamDir)), 0o700); err != nil {
		log.Errorf("Failed to write hook file %q: %v", hookFilePath, err)
	}
}
