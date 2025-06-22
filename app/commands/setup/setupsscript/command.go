package setupscript

import (
	"strconv"

	"github.com/spf13/cobra"

	"app/core"
	"app/log"
)

var ScriptsPath = "/setup/*/" + core.CommandFileName()

func MakeCommand(scriptPath string, scriptName string) *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes setup command " + scriptName,
		Run: func(_ *cobra.Command, _ []string) {
			command(scriptName, scriptPath)
		},
	}

	return command
}

func command(scriptName string, filePath string) {
	log.Infof("Execute setup command: " + scriptName)

	args := []string{core.GetAbsoluteExecutionPath()}
	exitCode := core.ExecuteScript(filePath, args)

	if exitCode == 0 {
		log.Successf(scriptName + " executed successfully")
	} else {
		log.Errorf(scriptName + " failed with: exit status " + strconv.Itoa(exitCode))
	}
}
