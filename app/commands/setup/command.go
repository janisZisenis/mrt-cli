package setup

import (
	"app/commands/setup/all"
	"app/commands/setup/cloneRepositories"
	"app/commands/setup/installGitHooks"
	"app/commands/setup/setupScript"
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
	command.AddCommand(cloneRepositories.MakeCommand())
	command.AddCommand(installGitHooks.MakeCommand())

	core.ForScriptInPathDo(teamDirectory+setupScript.ScriptsPath, func(scriptPath string, scriptName string) {
		command.AddCommand(setupScript.MakeCommand(scriptPath, scriptName))
	})

	return command
}
