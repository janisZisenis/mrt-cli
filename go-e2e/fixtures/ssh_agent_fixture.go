package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"testing"
)

type AgentFixture struct {
	t     *testing.T
	Agent *utils.Agent
}

func MakeAgentFixture(t *testing.T) *AgentFixture {
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

	return &AgentFixture{
		t:     t,
		Agent: agent,
	}
}
