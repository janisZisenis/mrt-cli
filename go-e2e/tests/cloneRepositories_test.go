package tests

import (
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
	"os"
	"os/exec"
	"testing"
)

func TestCloneRepositoriesToCustomPath(t *testing.T) {
	t.Parallel()
	repositoryDir := "/Users/jazi/Documents/Development/mrt/repositories/mrt-cli"
	agent := fixtures.AuthenticatedFixture(t, repositoryDir+"/.ssh/private-key")
	tempDir := t.TempDir()
	repositoryName := writeTeamConfig(tempDir)

	cloneRepositories(t, tempDir, agent)

	assertions.TestDirectoryExists(t, tempDir+"/repositories/"+repositoryName+"/.git")
}

func writeTeamConfig(tempDir string) string {
	repositoryName := "1_TestRepository"
	_ = utils.TeamConfigWriter(tempDir, map[string]interface{}{
		"repositories": []string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"},
	})
	return repositoryName
}

func cloneRepositories(t *testing.T, tempDir string, agent *utils.Agent) {
	cmd := exec.Command(binaryName, "--team-dir", tempDir, "setup", "clone-repositories")
	cmd.Env = append(os.Environ(), agent.Env()...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v", err)
	}
}
