package githook

import (
	"app/core"
	"github.com/spf13/cobra"
)

var commandName = "githook"
var branchFlag = "branch"
var hookNameFlag = "hook-name"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Executes the specified git-hook for a specified repository",
		Run:   command,
	}

	command.Flags().String(branchFlag, "", "The branch the commit hook was executed on")
	command.Flags().String(hookNameFlag, "", "The name of the git-hook to be executed")

	return command
}

func command(cmd *cobra.Command, args []string) {
	var teamInfo = core.LoadTeamConfiguration()
	branch, _ := cmd.Flags().GetString(branchFlag)
	hookName, _ := cmd.Flags().GetString(hookNameFlag)

	switch hookName {
	case core.PreCommit:
		preCommitHook(teamInfo, branch)
	case core.PrePush:
		prePushHook(teamInfo, branch)
	default:
		commitMsgHook(branch, teamInfo, args)
	}
}
