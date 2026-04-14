package run

import (
	"mrt-cli/app/commands/run/runscript"
	"mrt-cli/app/core"

	"github.com/spf13/cobra"
)

const (
	CommandName = "run"
	ShortHelp   = "Executes a specified run command"
	LongHelp    = ShortHelp + ".\n\n" +
		"No run commands found.\n\n" +
		"To add a run command, create an executable script at:\n\n" +
		"  run/<command-name>/command\n\n" +
		"Example:\n\n" +
		"  run/build/command\n" +
		"  run/test/command\n" +
		"  run/lint/command"
)

func MakeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: ShortHelp,
		Long:  LongHelp,
	}

	core.ForScriptInPathDo(runscript.GetScriptsPath(), func(scriptPath string, scriptName string) {
		command.AddCommand(runscript.MakeCommand(scriptName, scriptPath))
	})

	return command
}
