package tests_test

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/outputs"
	"mrt-cli/go-e2e/teamconfig"
)

func Test_IfTeamJsonContainsRepositoriesPath_Cloning_ShouldPrintMessageAboutCloningIntoThatPath(t *testing.T) {
	for _, repositoryPath := range []string{"some-path", "another-path"} {
		t.Run(repositoryPath, func(t *testing.T) {
			testCloningIntoRepositoriesPath(t, repositoryPath)
		})
	}
}

func testCloningIntoRepositoriesPath(t *testing.T, repositoryPath string) {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).Parallel()
	f.WriteTeamJSON(
		teamconfig.WithRepositoriesPath(repositoryPath),
		teamconfig.WithRepositories([]string{git.MakeCloneURL("1_TestRepository")}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(t, 0, "Start cloning repositories into \""+repositoryPath+"\"")
	output.Reversed().AssertLineEquals(t, 0, "Cloning repositories done")
}

func Test_IfTeamJsonContains2Repositories_Cloning_ShouldPrintDoneMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	f.WriteTeamJSON(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL("1_TestRepository"),
			git.MakeCloneURL("2_TestRepository"),
		}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.Reversed().AssertLineEquals(t, 0, "Cloning repositories done")
}

func Test_IfAuthenticationIsMissing_Cloning_ShouldPrintFailureMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	repositoryURL := git.MakeCloneURL("1_TestRepository")
	f.WriteTeamJSON(
		teamconfig.WithRepositories([]string{repositoryURL}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Cloning "+repositoryURL),
		outputs.HasLineContaining("Clone operation failed: "),
		outputs.HasLine("Failed to clone repository, skipping it."),
	)
}

func Test_IfTeamJsonDoesNotContainAnyRepositories_Cloning_ShouldPrintMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	f.WriteTeamJSON(
		teamconfig.WithRepositories([]string{}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(t, 0, "The team file does not contain any repositories, no repositories to clone.")
}

func Test_IfTeamJsonDoesNotExist_Cloning_ShouldPrintMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(
		t,
		0,
		"Could not read team file. To setup your repositories create a \"team.json\" file and add repositories to it.",
	)
}
