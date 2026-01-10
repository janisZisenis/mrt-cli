package githook

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"mrt-cli/app/core"
	"mrt-cli/app/log"
)

const (
	CommandName    = "git-hook"
	repositoryPath = "repository-path"
	hookNameFlag   = "hook-name"
)

func MakeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: "Executes the specified git-hook for a specified repository",
		Run:   command,
	}

	command.Flags().String(hookNameFlag, "", "The name of the git-hook to be executed")
	command.Flags().String(repositoryPath, ".", "The path to the repository")

	return command
}

func command(cmd *cobra.Command, args []string) {
	teamInfo, _ := core.LoadTeamConfiguration()
	hookName, _ := cmd.Flags().GetString(hookNameFlag)
	repositoryPath, _ := cmd.Flags().GetString(repositoryPath)

	currentBranchName := getCurrentBranchName(repositoryPath)

	switch hookName {
	case core.PreCommit:
		failIfBranchIsBlocked(teamInfo, currentBranchName, "commit")
	case core.PrePush:
		failIfBranchIsBlocked(teamInfo, currentBranchName, "push")
	case core.CommitMsg:
		if err := prefixCommitMessage(teamInfo, currentBranchName, args); err != nil {
			log.Errorf(err.Error())
			os.Exit(1)
		}
	default:
		log.Errorf("The given git-hook \"" + hookName + "\" does not exist.")
		os.Exit(1)
	}

	executeAdditionalScripts(repositoryPath, hookName, args)
}

func executeAdditionalScripts(repositoryPath string, hookName string, args []string) {
	files, _ := filepath.Glob(repositoryPath + "/hook-scripts/" + hookName + "/*")
	for _, file := range files {
		exitCode := core.ExecuteScript(file, args)

		if exitCode != 0 {
			os.Exit(1)
		}
	}
}

func getCurrentBranchName(repositoryPath string) string {
	shortBranchName, err := core.GetCurrentBranchShortName(repositoryPath)
	if err != nil {
		log.Errorf("The given path \"" + repositoryPath + "\" does not contain a repository.")
		os.Exit(1)
	}

	return shortBranchName
}
