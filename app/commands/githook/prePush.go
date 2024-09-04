package githook

import "app/core"

func prePushHook(teamInfo core.TeamInfo, branch string) {
	failIfBranchIsBlocked(teamInfo, branch, "push")
}
