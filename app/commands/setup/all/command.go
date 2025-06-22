package all

import (
	"app/commands/setup/clonerepositories"
	"app/commands/setup/installgithooks"
	"app/commands/setup/setupsscript"
	"app/core"
	"app/log"

	"github.com/spf13/cobra"
)

const scriptName = "all"
const skipFlagPrefix = "skip-"
const skipCloneFlag = skipFlagPrefix + clonerepositories.CommandName
const skipHooksFlag = skipFlagPrefix + installgithooks.CommandName

func MakeCommand(teamDirectory string) *cobra.Command {
	var command = &cobra.Command{
		Use:   scriptName,
		Short: "Executes all setup commands",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "Skips setting the git-hooks")
	command.Flags().Lookup(skipHooksFlag).NoOptDefVal = "true"

	command.Flags().Bool(skipCloneFlag, false, "Skips cloning the repositories")
	command.Flags().Lookup(skipCloneFlag).NoOptDefVal = "true"

	core.ForScriptInPathDo(teamDirectory+setupscript.ScriptsPath, func(filePath string, scriptName string) {
		var skipFlag = skipFlagPrefix + scriptName
		command.Flags().Bool(skipFlag, false, "Skips setup command: "+scriptName)
		command.Flags().Lookup(skipFlag).NoOptDefVal = "true"
	})

	return command
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)
	shouldSkipClone, _ := cmd.Flags().GetBool(skipCloneFlag)

	if !shouldSkipClone {
		clonerepositories.MakeCommand().Run(cmd, args)
	} else {
		log.Info("Skipping clone-repositories step.")
	}

	if !shouldSkipHooks {
		installgithooks.MakeCommand().Run(cmd, args)
	} else {
		log.Info("Skipping install-git-hooks step.")
	}

	executeAdditionalSetupScripts(cmd, args)
}

func executeAdditionalSetupScripts(cmd *cobra.Command, args []string) {
	log.Info("Executing setup commands.")

	core.ForScriptInPathDo(core.GetExecutionPath()+setupscript.ScriptsPath, func(scriptPath string, scriptName string) {
		skipFlag, _ := cmd.Flags().GetBool(skipFlagPrefix + scriptName)
		if !skipFlag {
			setupscript.MakeCommand(scriptPath, scriptName).Run(cmd, args)
		} else {
			log.Info("Skipping setup command: " + scriptName)
		}
	})

	log.Success("Done executing setup commands.")
}
