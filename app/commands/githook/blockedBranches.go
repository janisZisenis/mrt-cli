package githook

import (
	"app/core"
	"fmt"
	"os"
	"slices"
)

func failIfBranchIsBlocked(teamInfo core.TeamInfo, branch string, action string) {
	if slices.Contains(teamInfo.BlockedBranches, branch) {
		fmt.Println("Action \"" + action + "\" not allowed on branch \"" + branch + "\"")
		os.Exit(1)
	}
}
