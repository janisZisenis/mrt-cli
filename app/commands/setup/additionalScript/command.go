package additionalScript

import (
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

const ScriptsPath = "/setup/*/command"

func ForScriptInPathDo(scriptsPath string, do func(filePath string)) {
	files, _ := filepath.Glob(core.GetExecutablePath() + scriptsPath)
	for _, file := range files {
		do(file)
	}
}

func MakeCommand(filePath string) *cobra.Command {
	segments := strings.Split(filePath, "/")
	scriptName := segments[len(segments)-2]

	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes additional setup script " + scriptName,
		Run: func(cmd *cobra.Command, args []string) {
			command(scriptName, filePath)
		},
	}

	return command
}

func command(scriptName string, filePath string) {
	log.Info("Execute additional setup-script: " + scriptName)

	args := []string{core.GetExecutablePath()}
	err := core.ExecuteScript(filePath, args)

	if err != nil {
		log.Error(scriptName + " failed with: " + err.Error())
	} else {
		log.Success(scriptName + " executed successfully")
	}
}
