package tests

import (
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
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

	utils.MrtNew(binaryName, agent.Env()).
		RunInDirectory(tempDir).
		Setup().
		Clone().
		Execute()

	assertions.TestDirectoryExists(t, tempDir+"/repositories/"+repositoryName+"/.git")
}
