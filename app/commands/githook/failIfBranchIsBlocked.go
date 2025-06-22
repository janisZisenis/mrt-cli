package githook

import (
	"app/core"
	"app/log"
	"os"
	"slices"
)

func failIfBranchIsBlocked(teamInfo core.TeamInfo, branch string, action string) {
	if slices.Contains(teamInfo.BlockedBranches, branch) {
		log.Error("Action \"%s\" not allowed on branch \"%s\"", action, branch)
		os.Exit(1)
	}
}
