package githook

import (
	"os"
	"regexp"
	"strings"

	"mrt-cli/app/core"
	"mrt-cli/app/log"
)

func prefixCommitMessage(teamInfo core.TeamInfo, branch string, args []string) {
	if len(args) == 0 {
		log.Errorf("Missing commit message file argument")
		return
	}
	commitFile := args[0]
	data, err := os.ReadFile(commitFile)
	if err != nil {
		log.Errorf("Failed to read commit message file: %v", err)
		return
	}
	commitMessage := string(data)

	if teamInfo.CommitPrefixRegex == "" {
		return
	}

	regex, err := regexp.Compile(teamInfo.CommitPrefixRegex)
	if err != nil {
		log.Errorf("Invalid commit prefix regex in %v: %v", core.TeamFile, err)
		log.Errorf("CommitPrefixRegex: %s", teamInfo.CommitPrefixRegex)
		log.Errorf("Please fix the regex syntax in your %v file", core.TeamFile)
		return
	}

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
		_ = os.WriteFile(commitFile, []byte(matchesFromBranch[0]+": "+commitMessage), 0o600)
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
