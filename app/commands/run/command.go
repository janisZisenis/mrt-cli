package run

import (
	"app/commands/run/runscript"
	"app/core"

	"github.com/spf13/cobra"
)

const CommandName = "run"

func MakeCommand(teamDirectory string) *cobra.Command {
	var command = &cobra.Command{
		Use:   CommandName,
		Short: "Executes a specified run command",
	}

	core.ForScriptInPathDo(teamDirectory+runscript.ScriptsPath, func(scriptPath string, scriptName string) {
		command.AddCommand(runscript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
