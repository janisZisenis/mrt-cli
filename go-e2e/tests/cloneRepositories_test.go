package tests

import (
	"mrt-cli/go-e2e/assertions"
	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
	"testing"
)

func Test_IfTeamJsonDoesNotContainRepositoriesPath_Cloning_ShouldCloneRepositoryIntoDefaultFolder(t *testing.T) {
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

func Test_IfTeamJsonContainsARepositoryThatExistsOnTheRoot_Cloning_ShouldPrintOutSuccessMessage(t *testing.T) {
	t.Parallel()
	f := fixtures.MakeAuthenticatedFixture(t)
	repositoryURL := "git@github-testing:janisZisenisTesting/1_TestRepository.git"
	utils.WriteTeamJsonTo(f.TempDir,
		utils.WithRepositories([]string{repositoryURL}),
	)

	output := f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	output.AssertLineContains(t, 1, "Cloning "+repositoryURL)
	output.AssertLineMatchesRegex(t, 3, "Enumerating objects: [0-9]+, done.")
	output.Reversed().AssertLineContains(t, 1, "Successfully cloned "+repositoryURL)
}
