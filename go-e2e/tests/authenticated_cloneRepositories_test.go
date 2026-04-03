package tests_test

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
)

const defaultRepositoriesPath = "repositories"

func Test_IfTeamJsonDoesNotContainRepositoriesPath_Cloning_ShouldCloneRepositoryIntoDefaultFolder(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{utils.MakeCloneURLFrom(repositoryName)}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsARepositoryThatExistsOnTheRoot_Cloning_ShouldPrintOutSuccessMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryURL := utils.MakeCloneURLFrom("1_TestRepository")
	f.WriteTeamJSON(
		utils.WithRepositories([]string{repositoryURL}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(t, 1, "Cloning "+repositoryURL)
	output.AssertLineMatchesRegex(t, 3, "Enumerating objects: [0-9]+, done.")
	output.Reversed().AssertLineEquals(t, 1, "Successfully cloned "+repositoryURL)
}

func Test_IfTeamJsonContainsAlreadyClonedRepositories_Cloning_ClonesRemainingRepositoriesAndSkipsExistingOnes(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{
			utils.MakeCloneURLFrom(firstRepositoryName),
			utils.MakeCloneURLFrom(secondRepositoryName),
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
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	f.WriteTeamJSON(
		utils.WithRepositories([]string{}),
	)

	_ = f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertFolderDoesNotExist(defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsNonExistingRepository_Cloning_ShouldPrintOutAFailureMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	f.WriteTeamJSON(
		utils.WithRepositories([]string{utils.MakeCloneURLFrom("nonExistingRepository")}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertHasLine(t, "fatal: Could not read from remote repository.")
}

func Test_IfTeamJsonContainsNonExistingAndExistingRepository_Cloning_ShouldCloneTheExistingOne(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{
			utils.MakeCloneURLFrom("nonExistingRepository"),
			utils.MakeCloneURLFrom(repositoryName),
		}),
	)

	_ = f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixes_Cloning_ShouldTrimThePrefixesWhileCloningTheRepositories(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{
			utils.MakeCloneURLFrom(firstRepositoryName),
			utils.MakeCloneURLFrom(secondRepositoryName),
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
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{
			utils.MakeCloneURLFrom(firstRepositoryName),
			utils.MakeCloneURLFrom(secondRepositoryName),
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

func Test_IfTeamJsonContainsRepositoriesPath_Cloning_ShouldCloneRepositoriesIntoGivenRepositoriesPath(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	repositoriesPath := "xyz"
	f.WriteTeamJSON(
		utils.WithRepositories([]string{
			utils.MakeCloneURLFrom(repositoryName),
		}),
		utils.WithRepositoriesPath(repositoriesPath),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, repositoriesPath)
}
