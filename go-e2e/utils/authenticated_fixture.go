package utils

import "testing"

func AuthenticatedFixture(t *testing.T, keyPath string) *Agent {
	t.Helper()
	agent := StartSSHAgentFixture(t)

	if err := agent.AddKey(keyPath); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("Added key %s to ssh-agent", keyPath)
	}

	t.Cleanup(func() {
		if err := agent.RemoveKey(keyPath); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return agent
}
