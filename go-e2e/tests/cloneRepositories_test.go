package tests

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
)

var defaultRepositoriesPath = "repositories"

func Test_IfTeamJsonDoesNotContainRepositoriesPath_Cloning_ShouldCloneRepositoryIntoDefaultFolder(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	repositoryName := "1_TestRepository"
	f.WriteTeamJson(
		utils.WithRepositories([]string{utils.MakeCloneUrlFrom(repositoryName)}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsARepositoryThatExistsOnTheRoot_Cloning_ShouldPrintOutSuccessMessage(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	repositoryURL := utils.MakeCloneUrlFrom("1_TestRepository")
	f.WriteTeamJson(
		utils.WithRepositories([]string{repositoryURL}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineContains(t, 1, "Cloning "+repositoryURL)
	output.AssertLineMatchesRegex(t, 3, "Enumerating objects: [0-9]+, done.")
	output.Reversed().AssertLineContains(t, 1, "Successfully cloned "+repositoryURL)
}

func Test_IfTeamJsonContainsAlreadyClonedRepositories_Cloning_ClonesRemainingRepositoriesAndSkipsExistingOnes(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	f.WriteTeamJson(
		utils.WithRepositories([]string{
			utils.MakeCloneUrlFrom(firstRepositoryName),
			utils.MakeCloneUrlFrom(secondRepositoryName),
		}),
	)
	f.GitClone(firstRepositoryName, defaultRepositoriesPath)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(firstRepositoryName, defaultRepositoriesPath)
	f.AssertRepositoryExists(secondRepositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonDoesNotContainAnyRepository_Cloning_Should_Not_Clone_Any_Repository(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	f.WriteTeamJson(
		utils.WithRepositories([]string{}),
	)

	_ = f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertFolderDoesNotExist(defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsNonExistingRepository_Cloning_ShouldPrintOutAFailureMessage(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	f.WriteTeamJson(
		utils.WithRepositories([]string{utils.MakeCloneUrlFrom("nonExistingRepository")}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertHasLine(t, "fatal: Could not read from remote repository.")
}

func Test_IfTeamJsonContainsNonExistingAndExistingRepository_Cloning_ShouldCloneTheExistingOne(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	repositoryName := "1_TestRepository"
	f.WriteTeamJson(
		utils.WithRepositories([]string{
			utils.MakeCloneUrlFrom("nonExistingRepository"),
			utils.MakeCloneUrlFrom(repositoryName),
		}),
	)

	_ = f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixes_Cloning_ShouldTrimThePrefixesWhileCloningTheRepositories(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.WriteTeamJson(
		utils.WithRepositories([]string{
			utils.MakeCloneUrlFrom(firstRepositoryName),
			utils.MakeCloneUrlFrom(secondRepositoryName),
		}),
		utils.WithRepositoriesPrefixes([]string{"Prefix1_", "Prefix2_"}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists("TestRepository1", defaultRepositoriesPath)
	f.AssertRepositoryExists("TestRepository2", defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixesButUnprefixedRepositories_Cloning_ShouldNotTrim(t *testing.T) {
	f := fixtures.MakeAuthenticatedFixture(t).Parallel()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.WriteTeamJson(
		utils.WithRepositories([]string{
			utils.MakeCloneUrlFrom(firstRepositoryName),
			utils.MakeCloneUrlFrom(secondRepositoryName),
		}),
		utils.WithRepositoriesPrefixes([]string{"FirstPrefix", "SecondPrefix"}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(firstRepositoryName, defaultRepositoriesPath)
	f.AssertRepositoryExists(secondRepositoryName, defaultRepositoriesPath)
}
