package runScript

import (
	"app/core"
	"os"

	"github.com/spf13/cobra"
)

var ScriptsPath = "/run/*/" + core.CommandFileName()

func MakeCommand(scriptName string, scriptPath string) *cobra.Command {
	var command = &cobra.Command{
		Use: scriptName,
		Run: func(cmd *cobra.Command, args []string) {
			command(scriptPath, args)
		},
	}

	return command
}

func command(scriptPath string, args []string) {
	scriptArgs := append([]string{core.GetAbsoluteExecutionPath()}, args...)
	exitCode := core.ExecuteScript(scriptPath, scriptArgs)
	os.Exit(exitCode)
}
