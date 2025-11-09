package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"testing"
)

func StartSSHAgentFixture(t *testing.T) *utils.Agent {
	t.Helper()

	agent, err := utils.StartSSHAgent()
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("SSH Agent started with PID %d", agent.PID)
	}

	t.Cleanup(func() {
		if err := agent.Stop(); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return agent
}
