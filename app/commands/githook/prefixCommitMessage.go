package githook

import (
	"app/core"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func prefixCommitMessage(teamInfo core.TeamInfo, branch string, args []string) {
	commitFile := args[0]
	data, _ := os.ReadFile(commitFile)
	commitMessage := string(data)

	if teamInfo.JiraPrefixRegex == "" {
		return
	}

	regex := regexp.MustCompile(teamInfo.JiraPrefixRegex)

	if strings.HasPrefix(commitMessage, "Merge branch") ||
		strings.HasPrefix(commitMessage, "Merge remote-tracking branch") {
		fmt.Println("Merge commit detected, skipping commit-msg hook.")
		return
	}

	matchesFromMessage := regex.FindStringSubmatch(commitMessage)
	if len(matchesFromMessage) > 0 && strings.HasPrefix(commitMessage, matchesFromMessage[0]+": ") {
		fmt.Println("The commit message contains an issue ID (" + matchesFromMessage[0] + "). Good job!")
		return
	}

	matchesFromBranch := regex.FindStringSubmatch(branch)
	if len(matchesFromBranch) > 0 {
		_ = os.WriteFile(commitFile, []byte(matchesFromBranch[0]+": "+commitMessage), 0640)
		fmt.Println("JIRA-ID '" + matchesFromBranch[0] + "' was found in current branch name, prepended to commit message.")
		return
	}

	fmt.Println("The commit message needs a JIRA ID prefix.")
	fmt.Println("Either add the JIRA ID to you commit message, or include it in the branch name.")
	fmt.Println("Use '--no-verify' to skip git-hooks.")
	os.Exit(1)
}
