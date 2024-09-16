package runScript

import (
	"app/core"
	"github.com/spf13/cobra"
	"os"
)

const ScriptsPath = "/run/*/command"

func MakeCommand(scriptName string, scriptPath string) *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes run script: " + scriptName,
		Run: func(cmd *cobra.Command, args []string) {
			command(scriptPath, args)
		},
	}

	return command
}

func command(scriptPath string, args []string) {
	scriptArgs := append([]string{core.GetExecutablePath()}, args...)
	exitCode := core.ExecuteScript(scriptPath, scriptArgs)
	os.Exit(exitCode)
}
