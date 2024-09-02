package githook

import (
	"app/core"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"slices"
)

var branchFlag = "branch"
var hookNameFlag = "hook-name"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "githook",
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

	var action string
	if hookName == "pre-commit" {
		action = "commit"
	} else {
		action = "push"
	}

	if slices.Contains(teamInfo.BlockedBranches, branch) {
		fmt.Println("Action \"" + action + "\" not allowed on branch \"" + branch + "\"")
		os.Exit(1)
	}
}
