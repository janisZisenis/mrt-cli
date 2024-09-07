package setup

import (
	"app/core"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var commandName = "setup"
var skipHooksFlag = "skip-git-hooks"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Sets up you machine for development",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "The name of the git-hook to be executed")
	command.Flags().Lookup(skipHooksFlag).NoOptDefVal = "true"

	return command
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)

	teamInfo := core.LoadTeamConfiguration()

	if len(teamInfo.Repositories) == 0 {
		fmt.Println("Your team file does not contain any repositories")
		os.Exit(1)
	}

	setupRepositories(teamInfo)

	if !shouldSkipHooks {
		setupGitHooks(teamInfo)
	}

	files, _ := filepath.Glob(core.GetExecutablePath() + "/setup/*/command")
	for _, file := range files {
		segments := strings.Split(file, "/")
		commandName = segments[len(segments)-2]

		fmt.Println("Execute additional setup-script: " + commandName)

		args = []string{core.GetExecutablePath()}
		output, err := core.ExecuteBash(file, args)
		fmt.Print(output)

		if err != nil {
			fmt.Println(commandName + " failed with: " + err.Error())
		} else {
			fmt.Println(commandName + " executed successfully")
		}

		fmt.Println()
	}
}
