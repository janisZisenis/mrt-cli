package cloneRepositories

import (
	"app/core"
	"github.com/spf13/cobra"
)

var commandName = "clone-repositories"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Clones all repositories from " + core.TeamFile + " file",
		Run:   command,
	}

	return command
}

func command(cmd *cobra.Command, args []string) {
	teamInfo := core.LoadTeamConfiguration()

	SetupRepositories(teamInfo)
}
