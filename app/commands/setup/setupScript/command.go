package setupScript

import (
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
	"strconv"
)

const ScriptsPath = "/setup/*/command"

func MakeCommand(scriptPath string, scriptName string) *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes setup script " + scriptName,
		Run: func(cmd *cobra.Command, args []string) {
			command(scriptName, scriptPath)
		},
	}

	return command
}

func command(scriptName string, filePath string) {
	log.Info("Execute setup-script: " + scriptName)

	args := []string{core.GetExecutablePath()}
	exitCode := core.ExecuteScript(filePath, args)

	if exitCode == 0 {
		log.Success(scriptName + " executed successfully")
	} else {
		log.Error(scriptName + " failed with: exit status " + strconv.Itoa(exitCode))
	}
}
