package run

import (
	"app/commands/run/runScript"
	"app/core"
	"github.com/spf13/cobra"
)

const CommandName = "run"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CommandName,
		Short: "Executes a specified run script",
	}

	core.ForScriptInPathDo(runScript.ScriptsPath, func(scriptPath string, scriptName string) {
		command.AddCommand(runScript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
