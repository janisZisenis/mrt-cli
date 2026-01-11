package setup

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"mrt-cli/app/commands/setup/all"
	"mrt-cli/app/commands/setup/clonerepositories"
	"mrt-cli/app/commands/setup/installgithooks"
	"mrt-cli/app/commands/setup/setupscript"
	"mrt-cli/app/core"
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

	scriptPath := filepath.Join(teamDirectory, setupscript.GetScriptsPath())
	core.ForScriptInPathDo(scriptPath, func(scriptPath string, scriptName string) {
		command.AddCommand(setupscript.MakeCommand(scriptPath, scriptName))
	})

	return command
}
