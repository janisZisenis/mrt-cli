package run

import (
	"app/commands/run/additionalRunScript"
	"github.com/spf13/cobra"
)

const CommandName = "run"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CommandName,
		Short: "Executes a specified run script",
	}

	additionalRunScript.ForScriptInPathDo(additionalRunScript.ScriptsPath, func(scriptPath string, scriptName string) {
		command.AddCommand(additionalRunScript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
