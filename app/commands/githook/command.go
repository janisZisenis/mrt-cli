package githook

import (
	"mrt-cli/app/core"
	"mrt-cli/app/log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	CommandName    = "git-hook"
	repositoryPath = "repository-path"
	hookNameFlag   = "hook-name"
	teamDirFlag    = "team-dir"
	hookScriptsDir = "hook-scripts"
)

func MakeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   CommandName,
		Short: "Executes the specified git-hook for a specified repository",
		Run:   command,
	}

	command.Flags().String(hookNameFlag, "", "The name of the git-hook to be executed")
	command.Flags().String(repositoryPath, "", "The path to the repository")
	command.Flags().String(teamDirFlag, "", "The path to the team directory")

	return command
}

func command(cmd *cobra.Command, args []string) {
	teamDir, _ := cmd.Flags().GetString(teamDirFlag)

	if teamDir == "" {
		log.Errorf("Missing team dir argument")
		os.Exit(1)
	}

	teamInfo, err := core.LoadTeamConfiguration(teamDir)
	if err != nil {
		log.Errorf("Failed to load team configuration")
		os.Exit(1)
	}

	hookName, _ := cmd.Flags().GetString(hookNameFlag)
	repositoryPath, _ := cmd.Flags().GetString(repositoryPath)

	if repositoryPath == "" {
		log.Errorf("Missing repository path argument")
		os.Exit(1)
	}

	if hookName == "" {
		log.Errorf("Missing hook name argument")
		os.Exit(1)
	}

	currentBranchName := getCurrentBranchName(repositoryPath)

	if !core.IsGitHook(hookName) {
		log.Errorf("The given git-hook \"" + hookName + "\" does not exist.")
		os.Exit(1)
	}

	switch hookName {
	case core.PreCommit:
		failIfBranchIsBlocked(teamInfo, currentBranchName, "commit")
	case core.PrePush:
		for _, branch := range getPushedRemoteBranchNames(os.Stdin) {
			failIfBranchIsBlocked(teamInfo, branch, "push")
		}
	case core.CommitMsg:
		if commitMsgErr := prefixCommitMessage(teamInfo, currentBranchName, args); commitMsgErr != nil {
			log.Errorf(commitMsgErr.Error())
			os.Exit(1)
		}
	}

	executeAdditionalScripts(repositoryPath, hookName, args)
}

func executeAdditionalScripts(repositoryPath string, hookName string, args []string) {
	hookScriptsPath := filepath.Join(repositoryPath, hookScriptsDir, hookName, "*")
	files, err := filepath.Glob(hookScriptsPath)
	if err != nil {
		log.Errorf("Failed to find hook scripts: %v", err)
		return
	}
	for _, file := range files {
		exitCode := core.ExecuteScript(file, args)

		if exitCode != 0 {
			os.Exit(1)
		}
	}
}

func getCurrentBranchName(repositoryPath string) string {
	shortBranchName, err := core.GetCurrentBranchShortName(repositoryPath)
	if err != nil {
		log.Errorf("The given path \"%s\" does not contain a repository: %v", repositoryPath, err)
		os.Exit(1)
	}

	return shortBranchName
}
