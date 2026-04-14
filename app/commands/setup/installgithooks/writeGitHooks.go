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

func getHookTemplate() string {
	return `
#!/bin/bash -e

hook_name=$(basename "$0")
` + getExecutableName() + " --team-dir " + getAbsoluteExecutionPath() + " " + githook.CommandName + ` --hook-name "$hook_name" --repository-path $PWD $@`
}

func writeGitHook(repositoryDirectory string, hookName string) {
	hooksPath := filepath.Join(repositoryDirectory, gitHooksDir)
	// #nosec G301 - githooks folder needs 0700 to be private (owner only)
	_ = os.MkdirAll(hooksPath, 0o700)
	hookFilePath := filepath.Join(hooksPath, hookName)
	// #nosec G306 - git hooks need to be executable by owner
	err := os.WriteFile(hookFilePath, []byte(getHookTemplate()), 0o700)
	if err != nil {
		log.Infof("unable to write file: " + err.Error())
	}
}
