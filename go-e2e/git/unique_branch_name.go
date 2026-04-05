package git

import (
	"crypto/rand"
	"fmt"
)

func UniqueBranchName() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("branch-%x", b)
}
