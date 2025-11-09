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
	utils.WriteTeamJsonTo(f.TempDir,
		utils.WithRepositories([]string{"git@github-testing:janisZisenisTesting/" + repositoryName + ".git"}),
	)

	f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	assertions.TestDirectoryExists(t, f.TempDir+"/repositories/"+repositoryName+"/.git")
}
