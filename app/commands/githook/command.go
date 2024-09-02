package githook

import (
	"app/core"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var branchFlag = "branch"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "githook",
		Short: "Executes the specified git-hook for a specified repository",
		Run:   command,
	}

	command.Flags().String(branchFlag, "", "The branch the commit hook was executed on")

	return command
}

func command(cmd *cobra.Command, args []string) {
	var teamInfo = core.LoadTeamConfiguration()
	branch, _ := cmd.Flags().GetString(branchFlag)

	if teamInfo.BlockedBranches.Contains(branch) {
		fmt.Println("Action \"commit\" not allowed on branch \"" + branch + "\"")
		os.Exit(1)
	}
}
