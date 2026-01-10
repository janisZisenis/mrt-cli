package githook

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"mrt-cli/app/core"
	"mrt-cli/app/log"
)

type MissingCommitMessageFileError struct{}

func (e *MissingCommitMessageFileError) Error() string {
	return "Missing commit message file argument"
}

type FailedToReadCommitMessageError struct {
	message error
}

func (e *FailedToReadCommitMessageError) Error() string {
	return fmt.Sprintf("Failed to read commit message file: %v", e.message)
}

type InvalidCommitPrefixRegexError struct {
	regex   string
	message error
}

func (e *InvalidCommitPrefixRegexError) Error() string {
	return fmt.Sprintf(`Invalid commit prefix regex in team.json:
CommitPrefixRegex: %s
Please fix the regex syntax in your team.json file
Details: %v`, e.regex, e.message)
}

type InvalidCommitMessageError struct {
	regex string
}

func (e *InvalidCommitMessageError) Error() string {
	return fmt.Sprintf(`The commit message needs a commit prefix that matches the following regex %s.
Either add the commit prefix to your commit message, or include it in the branch name.
Use '--no-verify' to skip git-hooks.`, e.regex)
}

func prefixCommitMessage(teamInfo core.TeamInfo, branch string, args []string) error {
	if len(args) == 0 {
		return &MissingCommitMessageFileError{}
	}
	commitFile := args[0]
	data, err := os.ReadFile(commitFile)
	if err != nil {
		return &FailedToReadCommitMessageError{message: err}
	}
	commitMessage := string(data)

	if teamInfo.CommitPrefixRegex == "" {
		return nil
	}

	regex, err := regexp.Compile(teamInfo.CommitPrefixRegex)
	if err != nil {
		return &InvalidCommitPrefixRegexError{regex: teamInfo.CommitPrefixRegex, message: err}
	}

	if strings.HasPrefix(commitMessage, "Merge branch") ||
		strings.HasPrefix(commitMessage, "Merge remote-tracking branch") {
		log.Infof("Merge commit detected, skipping commit-msg hook.")
		return nil
	}

	matchesFromMessage := regex.FindStringSubmatch(commitMessage)
	if len(matchesFromMessage) > 0 && strings.HasPrefix(commitMessage, matchesFromMessage[0]+": ") {
		log.Successf("The commit message contains an issue ID (" + matchesFromMessage[0] + "). Good job!")
		return nil
	}

	matchesFromBranch := regex.FindStringSubmatch(branch)
	if len(matchesFromBranch) > 0 {
		_ = os.WriteFile(commitFile, []byte(matchesFromBranch[0]+": "+commitMessage), 0o600)
		log.Successf("Commit prefix '" + matchesFromBranch[0] + "' was found in current branch name, " +
			"prepended to commit message.")
		return nil
	}

	return &InvalidCommitMessageError{regex: teamInfo.CommitPrefixRegex}
}
