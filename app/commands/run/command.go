package run

import (
	"mrt-cli/app/commands/run/runscript"
	"mrt-cli/app/core"
	"path/filepath"

	"github.com/spf13/cobra"
)

const CommandName = "run"

func MakeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: "Executes a specified run command",
	}

	scriptPath := filepath.Join(".", runscript.GetScriptsPath())
	core.ForScriptInPathDo(scriptPath, func(scriptPath string, scriptName string) {
		command.AddCommand(runscript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
