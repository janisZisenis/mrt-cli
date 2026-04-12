package githook

import (
	"bufio"
	"io"
	"strings"
)

func getPushedBranchNames(reader io.Reader) []string {
	const refsHeadsPrefix = "refs/heads/"
	// git provides each pushed ref as: <local-ref> <local-sha1> <remote-ref> <remote-sha1>
	const prePushLineFieldCount = 4

	var branches []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < prePushLineFieldCount {
			continue
		}

		remoteRef := fields[2]
		if strings.HasPrefix(remoteRef, refsHeadsPrefix) {
			branches = append(branches, strings.TrimPrefix(remoteRef, refsHeadsPrefix))
		}
	}

	return branches
}
