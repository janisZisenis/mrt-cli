package installgithooks

import (
	"mrt-cli/app/core"
	"mrt-cli/app/log"

	"github.com/spf13/cobra"
)

const CommandName = "install-git-hooks"

func MakeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: "Installs git hooks to all repositories found in the " + core.RepositoriesPath + " from " + core.TeamFile,
		Run:   command,
	}

	return command
}

func command(_ *cobra.Command, _ []string) {
	teamInfo, err := core.LoadTeamConfiguration(".")
	if err != nil {
		log.Errorf("Failed to load team configuration")
	}
	setupGitHooks(teamInfo)
}
