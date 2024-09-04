package setup

import (
	"app/core"
	"github.com/spf13/cobra"
)

var commandName = "setup"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Sets up you machine for development",
		Run:   command,
	}
	return command
}

func command(cmd *cobra.Command, args []string) {
	teamInfo := core.LoadTeamConfiguration()
	setupRepositories(teamInfo)
}
