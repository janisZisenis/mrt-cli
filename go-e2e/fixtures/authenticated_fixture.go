package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"testing"
)

func MakeAuthenticatedFixture(t *testing.T) *AgentFixture {
	t.Helper()
	f := MakeAgentFixture(t)

	keyPath := utils.GetRepoRootDir() + "/.ssh/private-key"

	if err := f.Agent.AddKey(keyPath); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("Added key %s to ssh-Agent", keyPath)
	}

	t.Cleanup(func() {
		if err := f.Agent.RemoveKey(keyPath); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return f
}
