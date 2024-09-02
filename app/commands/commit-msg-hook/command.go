package commit_msg_hook

import (
	"app/core"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var branchFlag = "branch"
var commitFileFlag = "commit-file"

func MakeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "commit-msg-hook",
		Short: "Executes the commit-msg git-hook",
		Run:   command,
	}

	command.Flags().String(branchFlag, "", "The branch the commit hook was executed on")
	command.Flags().String(commitFileFlag, "", "The commit file")

	return command
}

func command(cmd *cobra.Command, args []string) {
	var teamInfo = core.LoadTeamConfiguration()
	branch, _ := cmd.Flags().GetString(branchFlag)
	commitFile, _ := cmd.Flags().GetString(commitFileFlag)
	data, _ := os.ReadFile(commitFile)
	commitMessage := string(data)

	regex := regexp.MustCompile(teamInfo.JiraPrefixRegex)

	if strings.HasPrefix(commitMessage, "Merge branch") ||
		strings.HasPrefix(commitMessage, "Merge remote-tracking branch") {
		fmt.Println("Merge commit detected, skipping commit-msg hook.")
		os.Exit(0)
	}

	matchesFromMessage := regex.FindStringSubmatch(commitMessage)
	if len(matchesFromMessage) > 0 && strings.HasPrefix(commitMessage, matchesFromMessage[0]+": ") {
		fmt.Println("The commit message contains an issue ID (" + matchesFromMessage[0] + "). Good job!")
		os.Exit(0)
	}

	matchesFromBranch := regex.FindStringSubmatch(branch)
	if len(matchesFromBranch) > 0 {
		_ = os.WriteFile(commitFile, []byte(matchesFromBranch[0]+": "+commitMessage), 0640)
		fmt.Println("JIRA-ID '" + matchesFromBranch[0] + "' was found in current branch name, prepended to commit message.")
		os.Exit(0)
	}

	fmt.Println("The commit message needs a JIRA ID prefix.")
	fmt.Println("Either add the JIRA ID to you commit message, or include it in the branch name.")
	fmt.Println("Use '--no-verify' to skip git-hooks.")
	os.Exit(1)
}
