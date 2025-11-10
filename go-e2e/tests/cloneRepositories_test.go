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

	assertions.AssertDirectoryExists(t, f.TempDir+"/repositories/"+repositoryName+"/.git")
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

func Test_IfTeamJsonContainsAlreadyClonedRepositories_Cloning_ClonesRemainingRepositoriesAndSkipsExistingOnes(t *testing.T) {
	t.Parallel()
	f := fixtures.MakeAuthenticatedFixture(t)
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	utils.WriteTeamJsonTo(f.TempDir,
		utils.WithRepositories([]string{
			"git@github-testing:janisZisenisTesting/" + firstRepositoryName + ".git",
			"git@github-testing:janisZisenisTesting/" + secondRepositoryName + ".git",
		}),
	)
	utils.MakeGitCommand(f.Agent.Env()).
		Clone("git@github-testing:janisZisenisTesting/"+firstRepositoryName+".git", f.TempDir+"/repositories").
		Execute()

	f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	assertions.AssertDirectoryExists(t, f.TempDir+"/repositories/"+firstRepositoryName+"/.git")
	assertions.AssertDirectoryExists(t, f.TempDir+"/repositories/"+secondRepositoryName+"/.git")
}

func Test_IfTeamJsonDoesNotContainAnyRepository_Cloning_Should_Not_Clone_Any_Repository(t *testing.T) {
	t.Parallel()
	f := fixtures.MakeAuthenticatedFixture(t)
	utils.WriteTeamJsonTo(f.TempDir,
		utils.WithRepositories([]string{}),
	)

	_ = f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	assertions.AssertDirectoryDoesNotExist(t, f.TempDir+"/repositories")
}

func Test_IfTeamJsonContainsNonExistingRepository_Cloning_ShouldPrintOutAFailureMessage(t *testing.T) {
	t.Parallel()
	f := fixtures.MakeAuthenticatedFixture(t)
	utils.WriteTeamJsonTo(f.TempDir,
		utils.WithRepositories([]string{"git@github-testing:janisZisenisTesting/nonExisting.git"}),
	)

	output := f.MakeMrtCommand().
		RunInDirectory(f.TempDir).
		Setup().
		Clone().
		Execute()

	output.AssertHasLine(t, "fatal: Could not read from remote repository.")
}
