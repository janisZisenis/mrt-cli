package githook

import (
	"app/core"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var commandName = "githook"
var repositoryPath = "repository-path"
var hookNameFlag = "hook-name"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Executes the specified git-hook for a specified repository",
		Run:   command,
	}

	command.Flags().String(hookNameFlag, "", "The name of the git-hook to be executed")
	command.Flags().String(repositoryPath, "", "The path to the repository")

	return command
}

func command(cmd *cobra.Command, args []string) {
	var teamInfo = core.LoadTeamConfiguration()
	hookName, _ := cmd.Flags().GetString(hookNameFlag)
	repositoryPath, _ := cmd.Flags().GetString(repositoryPath)

	currentBranchName := getCurrentBranchName(repositoryPath)

	switch hookName {
	case core.PreCommit:
		preCommitHook(teamInfo, currentBranchName)
	case core.PrePush:
		prePushHook(teamInfo, currentBranchName)
	default:
		commitMsgHook(teamInfo, currentBranchName, args)
	}
}

func getCurrentBranchName(repositoryPath string) string {
	repository, _ := git.PlainOpen(repositoryPath)
	currentBranch, _ := repository.Head()
	currentBranchName := currentBranch.Name().Short()
	return currentBranchName
}
