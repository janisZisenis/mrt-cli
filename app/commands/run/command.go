package run

import (
	"path/filepath"

	"mrt-cli/app/commands/run/runscript"
	"mrt-cli/app/core"

	"github.com/spf13/cobra"
)

const CommandName = "run"

func MakeCommand(teamDirectory string) *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: "Executes a specified run command",
	}

	scriptPath := filepath.Join(teamDirectory, runscript.GetScriptsPath())
	core.ForScriptInPathDo(scriptPath, func(scriptPath string, scriptName string) {
		command.AddCommand(runscript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
