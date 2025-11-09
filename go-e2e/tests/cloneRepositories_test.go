package tests

import (
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
	"testing"
)

func TestCloneRepositoriesToCustomPath(t *testing.T) {
	t.Parallel()
	f := fixtures.MakeAuthenticatedFixture(t)

	repositoryName := "1_TestRepository"
	data := map[string]interface{}{
		"repositories": []string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"},
	}
	_ = utils.TeamConfigWriter(f.TempDir, data)

	f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	assertions.TestDirectoryExists(t, f.TempDir+"/repositories/"+repositoryName+"/.git")
}
