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
	agent := fixtures.AuthenticatedFixture(t)
	tempDir := t.TempDir()
	repositoryName := "1_TestRepository"
	data := map[string]interface{}{
		"repositories": []string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"},
	}
	_ = utils.TeamConfigWriter(tempDir, data)

	cloneRepositories(t, tempDir, agent)

	assertions.TestDirectoryExists(t, tempDir+"/repositories/"+repositoryName+"/.git")
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
