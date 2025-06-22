package githook

import (
	"app/core"
	"app/log"
	"os"
	"regexp"
	"strings"
)

func prefixCommitMessage(teamInfo core.TeamInfo, branch string, args []string) {
	commitFile := args[0]
	data, _ := os.ReadFile(commitFile)
	commitMessage := string(data)

	if teamInfo.CommitPrefixRegex == "" {
		return
	}

	regex := regexp.MustCompile(teamInfo.CommitPrefixRegex)

	if strings.HasPrefix(commitMessage, "Merge branch") ||
		strings.HasPrefix(commitMessage, "Merge remote-tracking branch") {
		log.Infof("Merge commit detected, skipping commit-msg hook.")
		return
	}

	matchesFromMessage := regex.FindStringSubmatch(commitMessage)
	if len(matchesFromMessage) > 0 && strings.HasPrefix(commitMessage, matchesFromMessage[0]+": ") {
		log.Successf("The commit message contains an issue ID (" + matchesFromMessage[0] + "). Good job!")
		return
	}

	matchesFromBranch := regex.FindStringSubmatch(branch)
	if len(matchesFromBranch) > 0 {
		_ = os.WriteFile(commitFile, []byte(matchesFromBranch[0]+": "+commitMessage), 0640)
		log.Successf("Commit prefix '" + matchesFromBranch[0] + "' was found in current branch name, " +
			"prepended to commit message.")
		return
	}

	log.Errorf("The commit message needs a commit prefix, that matches the following regex " +
		teamInfo.CommitPrefixRegex + ".")
	log.Errorf("Either add the commit prefix to you commit message, or include it in the branch name.")
	log.Errorf("Use '--no-verify' to skip git-hooks.")
	os.Exit(1)
}
