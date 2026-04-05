package git

import (
	"crypto/rand"
	"fmt"
)

const branchNameRandomBytes = 4

func UniqueBranchName() string {
	b := make([]byte, branchNameRandomBytes)
	_, _ = rand.Read(b)
	return fmt.Sprintf("branch-%x", b)
}
