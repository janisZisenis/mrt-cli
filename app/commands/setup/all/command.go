package all

import (
	"app/commands/setup/cloneRepositories"
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var skipHooksFlag = "skip-git-hooks"
var commandName = "all"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   commandName,
		Short: "Executes all setup commands",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "Skips setting the git-hooks")
	command.Flags().Lookup(skipHooksFlag).NoOptDefVal = "true"

	return command
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)

	teamInfo := core.LoadTeamConfiguration()

	if len(teamInfo.Repositories) == 0 {
		log.Warning("Your team file does not contain any repositories")
		os.Exit(1)
	}

	cloneRepositories.MakeCommand().Run(cmd, args)

	if !shouldSkipHooks {
		teamInfo := core.LoadTeamConfiguration()
		setupGitHooks(teamInfo)
	}

	executeAdditionalSetupScripts()
}

func executeAdditionalSetupScripts() {
	files, _ := filepath.Glob(core.GetExecutablePath() + "/setup/*/command")
	for _, file := range files {
		segments := strings.Split(file, "/")
		commandName = segments[len(segments)-2]

		log.Info("Execute additional setup-script: " + commandName)

		args := []string{core.GetExecutablePath()}
		err := core.ExecuteScript(file, args)

		if err != nil {
			log.Error(commandName + " failed with: " + err.Error())
		} else {
			log.Success(commandName + " executed successfully")
		}

		log.Info("")
	}
}
