package fixtures

import (
	"testing"

	"mrt-cli/go-e2e/utils"
)

var privateKeyPath = utils.GetRepoRootDir() + "/.ssh/private-key"

func MakeAuthenticatedFixture(t *testing.T) *MrtFixture {
	t.Helper()
	f := MakeMrtFixture(t)

	if err := f.agent.AddKey(privateKeyPath); err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		if err := f.agent.RemoveKey(privateKeyPath); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return f
}
