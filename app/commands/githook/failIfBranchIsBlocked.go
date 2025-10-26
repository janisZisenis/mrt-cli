package githook

import (
	"os"
	"slices"

	"app/core"
	"app/log"
)

func failIfBranchIsBlocked(teamInfo core.TeamInfo, branch string, action string) {
	if slices.Contains(teamInfo.BlockedBranches, branch) {
		log.Errorf("Action \"%s\" not allowed on branch \"%s\"", action, branch)
		os.Exit(1)
	}
}
