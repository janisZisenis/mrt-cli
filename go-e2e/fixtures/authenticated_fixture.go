package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"testing"
)

func MakeAuthenticatedFixture(t *testing.T) *MrtFixture {
	t.Helper()
	f := MakeMrtFixture(t)

	keyPath := utils.GetRepoRootDir() + "/.ssh/private-key"

	if err := f.agent.AddKey(keyPath); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("Added key %s to ssh-Agent", keyPath)
	}

	t.Cleanup(func() {
		if err := f.agent.RemoveKey(keyPath); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return f
}
