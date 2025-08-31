package installgithooks

import (
	"os"

	"app/commands/githook"
	"app/core"
	"app/log"
)

func getHookTemplate() string {
	return `
#!/bin/bash -e

hook_name=$(basename "$0")
` + core.GetExecutableName() + " --team-dir " + core.GetAbsoluteExecutionPath() + " " + githook.CommandName + ` --hook-name "$hook_name" --repository-path $PWD $@`
}

func writeGitHook(repositoryDirectory string, hookName string) {
	hooksPath := repositoryDirectory + "/hooks/"
	// #nosec G301 - githooks folder needs 0755 to be executable
	_ = os.MkdirAll(hooksPath, 0o755)
	// #nosec G306 - githooks need 0755 to be executable
	err := os.WriteFile(hooksPath+hookName, []byte(getHookTemplate()), 0o755)
	if err != nil {
		log.Infof("unable to write file: " + err.Error())
	}
}
