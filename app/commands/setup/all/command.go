package all

import (
	"app/commands/setup/additionalScript"
	"app/commands/setup/cloneRepositories"
	"app/commands/setup/installGitHooks"
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
	"os"
)

var skipHooksFlag = "skip-git-hooks"
var scriptName = "all"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
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
		installGitHooks.MakeCommand().Run(cmd, args)
	}

	additionalScript.ForScriptInPathDo(additionalScript.ScriptsPath, func(scriptPath string) {
		additionalScript.MakeCommand(scriptPath).Run(cmd, args)
	})
}
