package all

import (
	"path/filepath"

	"mrt-cli/app/commands/setup/clonerepositories"
	"mrt-cli/app/commands/setup/installgithooks"
	"mrt-cli/app/commands/setup/setupscript"
	"mrt-cli/app/core"
	"mrt-cli/app/log"

	"github.com/spf13/cobra"
)

const (
	scriptName     = "all"
	skipFlagPrefix = "skip-"
	skipCloneFlag  = skipFlagPrefix + clonerepositories.CommandName
	skipHooksFlag  = skipFlagPrefix + installgithooks.CommandName
)

func MakeCommand(teamDirectory string) *cobra.Command {
	command := &cobra.Command{
		Use:   scriptName,
		Short: "Executes all setup commands",
		Run:   command,
	}

	command.Flags().Bool(skipHooksFlag, false, "Skips setting the git-hooks")
	setDefaultValueToTrue(command, skipHooksFlag)

	command.Flags().Bool(skipCloneFlag, false, "Skips cloning the repositories")
	setDefaultValueToTrue(command, skipCloneFlag)

	scriptPath := filepath.Join(teamDirectory, setupscript.GetScriptsPath())
	core.ForScriptInPathDo(scriptPath, func(_ string, scriptName string) {
		skipFlag := skipFlagPrefix + scriptName
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

	scriptPath := filepath.Join(core.GetExecutionPath(), setupscript.GetScriptsPath())
	core.ForScriptInPathDo(
		scriptPath,
		func(scriptPath string, scriptName string) {
			skipFlag, _ := cmd.Flags().GetBool(skipFlagPrefix + scriptName)
			if !skipFlag {
				setupscript.MakeCommand(scriptPath, scriptName).Run(cmd, args)
			} else {
				log.Infof("Skipping setup command: " + scriptName)
			}
		})

	log.Successf("Done executing setup commands.")
}
