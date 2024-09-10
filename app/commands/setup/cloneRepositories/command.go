package cloneRepositories

import (
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
	"os"
)

const commandName = "clone-repositories"

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

	if len(teamInfo.Repositories) == 0 {
		log.Warning("Your team file does not contain any repositories")
		os.Exit(1)
	}

	SetupRepositories(teamInfo)
}
