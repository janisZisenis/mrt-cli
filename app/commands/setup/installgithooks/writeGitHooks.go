package installgithooks

import (
	"app/commands/githook"
	"app/core"
	"app/log"
	"os"
)

func getHookTemplate() string {
	return `
#!/bin/bash -e

hook_name=$(basename "$0")
` + core.GetExecutableName() + " --team-dir " + core.GetAbsoluteExecutionPath() + " " + githook.CommandName + ` --hook-name "$hook_name" --repository-path $PWD $@`
}

func writeGitHook(repositoryDirectory string, hookName string) {
	hooksPath := repositoryDirectory + "/hooks/"
	_ = os.MkdirAll(hooksPath, os.ModePerm)
	err := os.WriteFile(hooksPath+hookName, []byte(getHookTemplate()), 0755)
	if err != nil {
		log.Info("unable to write file: " + err.Error())
	}
}
