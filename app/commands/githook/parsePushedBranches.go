package githook

import (
	"bufio"
	"io"
	"strings"
)

func getPushedBranchNames(reader io.Reader) []string {
	const refsHeadsPrefix = "refs/heads/"

	var branches []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		localRef := fields[0]
		if strings.HasPrefix(localRef, refsHeadsPrefix) {
			branches = append(branches, strings.TrimPrefix(localRef, refsHeadsPrefix))
		}
	}

	return branches
}
