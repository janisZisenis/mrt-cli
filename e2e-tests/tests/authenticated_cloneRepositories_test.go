package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	mrtclient "mrt-cli/e2e-tests/mrt"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"
)

func Test_IfTeamJsonDoesNotContainRepositoriesPath_Cloning_ShouldCloneRepositoryIntoDefaultFolder(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
	)

	f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonContainsARepository_Cloning_ShouldPrintOutSuccessMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoryURL := git.MakeCloneURL("1_TestRepository")
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{repositoryURL}),
	)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine(mrtclient.MsgCloning(repositoryURL)),
		outputs.HasLineContaining("Enumerating objects:"),
		outputs.HasLine(mrtclient.MsgSuccessfullyCloned(repositoryURL)),
	)
}

func Test_IfTeamJsonContainsAlreadyClonedRepositories_Cloning_ClonesRemainingRepositoriesAndSkipsExistingOnes(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine(mrtclient.MsgFailedToCloneRepository),
		outputs.HasLine(mrtclient.MsgCloning(git.MakeCloneURL(secondRepositoryName))),
		outputs.HasLine(mrtclient.MsgSuccessfullyCloned(git.MakeCloneURL(secondRepositoryName))),
	)
	f.AssertRepositoryExists(firstRepositoryName, mrtclient.DefaultRepositoriesPath)
	f.AssertRepositoryExists(secondRepositoryName, mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonDoesNotContainAnyRepository_Cloning_Should_Not_Clone_Any_Repository(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{}),
	)

	f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	f.AssertFolderDoesNotExist(mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonContainsNonExistingRepository_Cloning_ShouldPrintOutAFailureMessage(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL("nonExistingRepository")}),
	)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.AssertHasLine(t, "fatal: Could not read from remote repository.")
}

func Test_IfTeamJsonContainsNonExistingAndExistingRepository_Cloning_ShouldCloneTheExistingOne(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoryName := "1_TestRepository"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL("nonExistingRepository"),
			git.MakeCloneURL(repositoryName),
		}),
	)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Cloning "+git.MakeCloneURL("nonExistingRepository")),
		outputs.HasLine(mrtclient.MsgFailedToCloneRepository),
		outputs.HasLine("Successfully cloned "+git.MakeCloneURL(repositoryName)),
	)
	f.AssertRepositoryExists(repositoryName, mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixes_Cloning_ShouldTrimThePrefixesWhileCloningTheRepositories(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
		teamconfig.WithRepositoriesPrefixes([]string{"Prefix1_", "Prefix2_"}),
	)

	f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists("TestRepository1", mrtclient.DefaultRepositoriesPath)
	f.AssertRepositoryExists("TestRepository2", mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPrefixesButUnprefixedRepositories_Cloning_ShouldNotTrim(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	firstRepositoryName := "Prefix1_TestRepository1"
	secondRepositoryName := "Prefix2_TestRepository2"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
		teamconfig.WithRepositoriesPrefixes([]string{"FirstPrefix", "SecondPrefix"}),
	)

	f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(firstRepositoryName, mrtclient.DefaultRepositoriesPath)
	f.AssertRepositoryExists(secondRepositoryName, mrtclient.DefaultRepositoriesPath)
}

func Test_IfTeamJsonContainsRepositoriesPath_Cloning_ShouldCloneRepositoriesIntoGivenRepositoriesPath(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoryName := "1_TestRepository"
	repositoriesPath := "xyz"
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(repositoryName),
		}),
		teamconfig.WithRepositoriesPath(repositoriesPath),
	)

	f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	f.AssertRepositoryExists(repositoryName, repositoriesPath)
}
