package installGitHooks

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
repositoryPath=$PWD
cd ` + core.GetExecutablePath() + `
` + core.GetExecutable() + " " + githook.CommandName + ` --hook-name "$hook_name" --repository-path $repositoryPath $@`
}

func getCommitMsgHookTemplate() string {
	return `
#!/bin/bash -e

hook_name=$(basename "$0")
repositoryPath=$PWD
commitFile=$(realpath $1)
cd ` + core.GetExecutablePath() + `
` + core.GetExecutable() + " " + githook.CommandName + ` --hook-name "$hook_name" --repository-path $repositoryPath $commitFile`
}

func writeGitHook(repositoryDirectory string, hookName string) {
    var template string
    if hookName == "commit-msg" {
        template = getCommitMsgHookTemplate()
    } else {
        template = getHookTemplate()
    }

	hooksPath := repositoryDirectory + "/hooks/"
	_ = os.MkdirAll(hooksPath, os.ModePerm)
	err := os.WriteFile(hooksPath+hookName, []byte(template), 0755)
	if err != nil {
		log.Info("unable to write file: " + err.Error())
	}
}
