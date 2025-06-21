package setupScript

import (
	"app/core"
	"app/log"
	"strconv"

	"github.com/spf13/cobra"
)

var ScriptsPath = "/setup/*/" + core.CommandFileName()

func MakeCommand(scriptPath string, scriptName string) *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes setup command " + scriptName,
		Run: func(cmd *cobra.Command, args []string) {
			command(scriptName, scriptPath)
		},
	}

	return command
}

func command(scriptName string, filePath string) {
	log.Info("Execute setup command: " + scriptName)

	args := []string{core.GetAbsoluteExecutionPath()}
	exitCode := core.ExecuteScript(filePath, args)

	if exitCode == 0 {
		log.Success(scriptName + " executed successfully")
	} else {
		log.Error(scriptName + " failed with: exit status " + strconv.Itoa(exitCode))
	}
}
