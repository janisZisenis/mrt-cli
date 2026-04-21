package installgithooks

import (
	"errors"
	"mrt-cli/app/core"
	"mrt-cli/app/log"
	"os"

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
		if errors.Is(err, core.ErrInvalidRepositoriesPath) {
			log.Errorf("repositoriesPath must be a relative path within the team repository")
			os.Exit(1)
		}
		log.Errorf("Failed to load team configuration")
	}
	setupGitHooks(teamInfo)
}
