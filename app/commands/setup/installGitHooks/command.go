package installGitHooks

import (
	"app/core"
	"github.com/spf13/cobra"
)

var commandName = "install-git-hooks"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Installs git hooks to all repositories found in the " + core.RepositoriesPath + " from " + core.TeamFile,
		Run:   command,
	}

	return command
}

func command(cmd *cobra.Command, args []string) {
	teamInfo := core.LoadTeamConfiguration()
	setupGitHooks(teamInfo)
}
