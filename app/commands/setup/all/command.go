package all

import (
	"app/commands/setup/cloneRepositories"
	"app/commands/setup/installGitHooks"
	"app/commands/setup/setupScript"
	"app/core"
	"app/log"
	"github.com/spf13/cobra"
)

const scriptName = "all"
const skipFlagPrefix = "skip-"
const skipCloneFlag = skipFlagPrefix + cloneRepositories.CommandName
const skipHooksFlag = skipFlagPrefix + installGitHooks.CommandName

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes all setup commands",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "Skips setting the git-hooks")
	command.Flags().Lookup(skipHooksFlag).NoOptDefVal = "true"

	command.Flags().Bool(skipCloneFlag, false, "Skips cloning the repositories")
	command.Flags().Lookup(skipCloneFlag).NoOptDefVal = "true"

	core.ForScriptInPathDo(setupScript.ScriptsPath, func(filePath string, scriptName string) {
		var skipFlag = skipFlagPrefix + scriptName
		command.Flags().Bool(skipFlag, false, "Skips setup script: "+scriptName)
		command.Flags().Lookup(skipFlag).NoOptDefVal = "true"
	})

	return command
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)
	shouldSkipClone, _ := cmd.Flags().GetBool(skipCloneFlag)

	if !shouldSkipClone {
		cloneRepositories.MakeCommand().Run(cmd, args)
	} else {
		log.Info("Skipping clone-repositories step.")
	}

	if !shouldSkipHooks {
		installGitHooks.MakeCommand().Run(cmd, args)
	} else {
		log.Info("Skipping install-git-hooks step.")
	}

	executeAdditionalSetupScripts(cmd, args)
}

func executeAdditionalSetupScripts(cmd *cobra.Command, args []string) {
	log.Info("Executing additional setup-scripts.")

	core.ForScriptInPathDo(setupScript.ScriptsPath, func(scriptPath string, scriptName string) {
		skipFlag, _ := cmd.Flags().GetBool(skipFlagPrefix + scriptName)
		if !skipFlag {
			setupScript.MakeCommand(scriptPath, scriptName).Run(cmd, args)
		} else {
			log.Info("Skipping additional setup script: " + scriptName)
		}
	})

	log.Success("Done executing additional setup-scripts.")
}
