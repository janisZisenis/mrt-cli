package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"testing"
)

var privateKeyPath = utils.GetRepoRootDir() + "/.ssh/private-key"

func MakeAuthenticatedFixture(t *testing.T) *MrtFixture {
	t.Helper()
	f := MakeMrtFixture(t)

	if err := f.Agent.AddKey(privateKeyPath); err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		if err := f.Agent.RemoveKey(privateKeyPath); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return f
}
