package githook

import "app/core"

func preCommitHook(teamInfo core.TeamInfo, branch string) {
	failIfBranchIsBlocked(teamInfo, branch, "commit")
}
