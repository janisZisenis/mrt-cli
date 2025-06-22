package setup

import (
	"app/commands/setup/all"
	"app/commands/setup/clonerepositories"
	"app/commands/setup/installgithooks"
	"app/commands/setup/setupsscript"
	"app/core"

	"github.com/spf13/cobra"
)

const commandName = "setup"

func MakeCommand(teamDirectory string) *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Sets up you machine for development",
	}

	command.AddCommand(all.MakeCommand(teamDirectory))
	command.AddCommand(clonerepositories.MakeCommand())
	command.AddCommand(installgithooks.MakeCommand())

	core.ForScriptInPathDo(teamDirectory+setupscript.ScriptsPath, func(scriptPath string, scriptName string) {
		command.AddCommand(setupscript.MakeCommand(scriptPath, scriptName))
	})

	return command
}
