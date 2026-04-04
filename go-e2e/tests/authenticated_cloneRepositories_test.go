package tests_test

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/outputs"
	"mrt-cli/go-e2e/teamconfig"
)

const defaultRepositoriesPath = "repositories"

func Test_IfTeamJsonDoesNotContainRepositoriesPath_Cloning_ShouldCloneRepositoryIntoDefaultFolder(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsARepository_Cloning_ShouldPrintOutSuccessMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryURL := git.MakeCloneURL("1_TestRepository")
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{repositoryURL}),
	)

	output, _ := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Cloning "+repositoryURL),
		outputs.HasLineContaining("Enumerating objects:"),
		outputs.HasLine("Successfully cloned "+repositoryURL),
	)
}

func Test_IfTeamJsonContainsAlreadyClonedRepositories_Cloning_ClonesRemainingRepositoriesAndSkipsExistingOnes(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(defaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()

	output, _ := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Failed to clone repository, skipping it."),
		outputs.HasLine("Cloning "+git.MakeCloneURL(secondRepositoryName)),
		outputs.HasLine("Successfully cloned "+git.MakeCloneURL(secondRepositoryName)),
	)
	f.AssertRepositoryExists(firstRepositoryName, defaultRepositoriesPath)
	f.AssertRepositoryExists(secondRepositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonDoesNotContainAnyRepository_Cloning_Should_Not_Clone_Any_Repository(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{}),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertFolderDoesNotExist(defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsNonExistingRepository_Cloning_ShouldPrintOutAFailureMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL("nonExistingRepository")}),
	)

	output, _ := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertHasLine(t, "fatal: Could not read from remote repository.")
}

func Test_IfTeamJsonContainsNonExistingAndExistingRepository_Cloning_ShouldCloneTheExistingOne(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL("nonExistingRepository"),
			git.MakeCloneURL(repositoryName),
		}),
	)

	output, _ := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Cloning "+git.MakeCloneURL("nonExistingRepository")),
		outputs.HasLine("Failed to clone repository, skipping it."),
		outputs.HasLine("Successfully cloned "+git.MakeCloneURL(repositoryName)),
	)
	f.AssertRepositoryExists(repositoryName, defaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixes_Cloning_ShouldTrimThePrefixesWhileCloningTheRepositories(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
		teamconfig.WithRepositoriesPrefixes([]string{"Prefix1_", "Prefix2_"}),
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
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
		teamconfig.WithRepositoriesPrefixes([]string{"FirstPrefix", "SecondPrefix"}),
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
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(repositoryName),
		}),
		teamconfig.WithRepositoriesPath(repositoriesPath),
	)

	f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, repositoriesPath)
}
