package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MrtFixture struct {
	binaryPath string
	agent      *utils.Agent
	TempDir    string
}

func MakeMrtFixture(t *testing.T) *MrtFixture {
	t.Helper()

	agent, err := utils.StartSSHAgent()
	if err != nil {
		t.Fatalf("Could not start SSH Agent.")
	}

	t.Cleanup(func() {
		if err := agent.Stop(); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return &MrtFixture{
		binaryPath: getBinaryPath(utils.GetRepoRootDir(), t),
		agent:      agent,
		TempDir:    t.TempDir(),
	}
}

func (m *MrtFixture) MakeMrtCommand() *utils.Mrt {
	return utils.MakeMrtCommand(m.binaryPath, m.agent.Env())
}

func getBinaryPath(repositoryDir string, t *testing.T) string {
	cmd := exec.Command("mrt", "--team-dir", repositoryDir, "run", "binary-location")
	binaryPathBytes, err := cmd.Output()

	output := string(binaryPathBytes)

	if err != nil {
		t.Fatalf("failed to run mrt command: %s", output)
	}

	binaryPath := stringTrimNewline(output)

	if binaryPath == "" {
		t.Fatalf("command returned empty directory")
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("binary not found at: %s", binaryPath)
	}

	return binaryPath
}

func stringTrimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return strings.TrimSpace(s)
}
