package clonerepositories

import (
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
)

const CommandName = "clone-repositories"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CommandName,
		Short: "Clones all repositories from " + core.TeamFile + " file",
		Run:   command,
	}

	return command
}

func command(_ *cobra.Command, _ []string) {
	teamInfo, err := core.LoadTeamConfiguration()

	if err != nil {
		log.Info("Could not read team file. To setup your repositories create a \"" + core.TeamFile + "\" file and add repositories to it.")
		return
	}

	if len(teamInfo.Repositories) == 0 {
		log.Info("The team file does not contain any repositories, no repositories to clone.")
		return
	}

	CloneRepositories(teamInfo)
}
