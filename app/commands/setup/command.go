package setup

import (
	"app/core"
	"github.com/spf13/cobra"
)

var commandName = "setup"
var skipHooksFlag = "skip-git-hooks"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Sets up you machine for development",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "The name of the git-hook to be executed")
	command.Flags().Lookup(skipHooksFlag).NoOptDefVal = "true"

	return command
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)

	teamInfo := core.LoadTeamConfiguration()
	setupRepositories(teamInfo)
	setupGitHooks(teamInfo, shouldSkipHooks)
}
