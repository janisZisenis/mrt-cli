package fixtures

import (
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MrtFixture struct {
	t          *testing.T
	binaryPath string
	Agent      *utils.Agent
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
		t:          t,
		binaryPath: getBinaryPath(utils.GetRepoRootDir(), t),
		Agent:      agent,
		TempDir:    t.TempDir(),
	}
}

func (m *MrtFixture) GitClone(repositoryName string, destination string) {
	utils.MakeGitCommand(m.Agent.Env()).
		Clone(utils.MakeCloneUrlFrom(repositoryName), m.TempDir+"/"+destination).
		Execute()
}

func (m *MrtFixture) MakeMrtCommand() *utils.Mrt {
	return utils.
		MakeMrtCommand(m.binaryPath, m.Agent.Env()).
		RunInDirectory(m.TempDir)
}

func (m *MrtFixture) WriteTeamJson(withOptions ...utils.TeamConfigOption) {
	utils.WriteTeamJsonTo(m.TempDir, withOptions...)
}

func (m *MrtFixture) AssertRepositoryExists(repositoryName string, inFolder string) {
	assertions.AssertDirectoryExists(m.t, m.TempDir+"/"+inFolder+"/"+repositoryName+"/.git")
}

func (m *MrtFixture) AssertFolderDoesNotExist(folder string) {
	assertions.AssertDirectoryDoesNotExist(m.t, m.TempDir+"/"+folder)
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
