package fixtures

import (
	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type MrtFixture struct {
	t          *testing.T
	binaryPath string
	agent      *utils.Agent
	tempDir    string
}

func MakeMrtFixture(t *testing.T) *MrtFixture {
	t.Helper()

	agent, err := utils.StartSSHAgent()
	if err != nil {
		t.Fatalf("Could not start SSH agent.")
	}

	t.Cleanup(func() {
		if err := agent.Stop(); err != nil {
			t.Fatalf("%v", err)
		}
	})

	return &MrtFixture{
		t:          t,
		binaryPath: getBinaryPath(utils.GetRepoRootDir(), t),
		agent:      agent,
		tempDir:    t.TempDir(),
	}
}

func (f *MrtFixture) Parallel() *MrtFixture {
	f.t.Helper()
	f.t.Parallel()

	return f
}

func (f *MrtFixture) GitClone(repositoryName string, destination string) {
	utils.MakeGitCommand(f.agent.Env()).
		Clone(utils.MakeCloneUrlFrom(repositoryName), f.tempDir+"/"+destination).
		Execute()
}

func (f *MrtFixture) MakeMrtCommand() utils.MrtDirectedCommand {
	return utils.
		MakeMrtCommand(f.binaryPath, f.agent.Env()).
		RunInDirectory(f.tempDir)
}

func (f *MrtFixture) WriteTeamJson(withOptions ...utils.TeamConfigOption) {
	utils.WriteTeamJsonTo(f.tempDir, withOptions...)
}

func (f *MrtFixture) AssertRepositoryExists(repositoryName string, inFolder string) {
	f.t.Helper()
	utils.AssertDirectoryExists(f.t, f.tempDir+"/"+inFolder+"/"+repositoryName+"/.git")
}

func (f *MrtFixture) AssertFolderDoesNotExist(folder string) {
	f.t.Helper()
	utils.AssertDirectoryDoesNotExist(f.t, f.tempDir+"/"+folder)
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
