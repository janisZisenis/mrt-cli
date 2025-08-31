package setup

import (
	"github.com/spf13/cobra"

	"app/commands/setup/all"
	"app/commands/setup/clonerepositories"
	"app/commands/setup/installgithooks"
	"app/commands/setup/setupscript"
	"app/core"
)

const commandName = "setup"

func MakeCommand(teamDirectory string) *cobra.Command {
	command := &cobra.Command{
		Use:   commandName,
		Short: "Sets up you machine for development",
	}

	command.AddCommand(all.MakeCommand(teamDirectory))
	command.AddCommand(clonerepositories.MakeCommand())
	command.AddCommand(installgithooks.MakeCommand())

	core.ForScriptInPathDo(teamDirectory+setupscript.GetScriptsPath(), func(scriptPath string, scriptName string) {
		command.AddCommand(setupscript.MakeCommand(scriptPath, scriptName))
	})

	return command
}
