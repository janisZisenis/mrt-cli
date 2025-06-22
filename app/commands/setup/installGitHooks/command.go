package installGitHooks

import (
	"app/core"
	"github.com/spf13/cobra"
)

const CommandName = "install-git-hooks"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CommandName,
		Short: "Installs git hooks to all repositories found in the " + core.RepositoriesPath + " from " + core.TeamFile,
		Run:   command,
	}

	return command
}

func command(_ *cobra.Command, _ []string) {
	teamInfo, _ := core.LoadTeamConfiguration()
	setupGitHooks(teamInfo)
}
