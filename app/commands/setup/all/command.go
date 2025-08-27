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
	setDefaultValueToTrue(command, skipHooksFlag)

	command.Flags().Bool(skipCloneFlag, false, "Skips cloning the repositories")
	setDefaultValueToTrue(command, skipCloneFlag)

	core.ForScriptInPathDo(teamDirectory+setupscript.GetScriptsPath(), func(_ string, scriptName string) {
		var skipFlag = skipFlagPrefix + scriptName
		command.Flags().Bool(skipFlag, false, "Skips setup command: "+scriptName)
		setDefaultValueToTrue(command, skipFlag)
	})

	return command
}

func setDefaultValueToTrue(command *cobra.Command, flag string) {
	command.Flags().Lookup(flag).NoOptDefVal = "true"
}

func command(cmd *cobra.Command, args []string) {
	shouldSkipHooks, _ := cmd.Flags().GetBool(skipHooksFlag)
	shouldSkipClone, _ := cmd.Flags().GetBool(skipCloneFlag)

	if !shouldSkipClone {
		clonerepositories.MakeCommand().Run(cmd, args)
	} else {
		log.Infof("Skipping clone-repositories step.")
	}

	if !shouldSkipHooks {
		installgithooks.MakeCommand().Run(cmd, args)
	} else {
		log.Infof("Skipping install-git-hooks step.")
	}

	executeAdditionalSetupScripts(cmd, args)
}

func executeAdditionalSetupScripts(cmd *cobra.Command, args []string) {
	log.Infof("Executing setup commands.")

	core.ForScriptInPathDo(core.GetExecutionPath()+setupscript.GetScriptsPath(), func(scriptPath string, scriptName string) {
		skipFlag, _ := cmd.Flags().GetBool(skipFlagPrefix + scriptName)
		if !skipFlag {
			setupscript.MakeCommand(scriptPath, scriptName).Run(cmd, args)
		} else {
			log.Infof("Skipping setup command: " + scriptName)
		}
	})

	log.Successf("Done executing setup commands.")
}
