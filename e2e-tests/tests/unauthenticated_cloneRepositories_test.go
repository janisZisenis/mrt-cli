package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	mrtclient "mrt-cli/e2e-tests/mrt"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"
)

func Test_IfTeamJsonContains2Repositories_Cloning_ShouldPrintDoneMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL("1_TestRepository"),
			git.MakeCloneURL("2_TestRepository"),
		}),
	)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.Reversed().AssertLineEquals(t, 0, "Cloning repositories done")
}

func Test_IfAuthenticationIsMissing_Cloning_ShouldPrintFailureMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
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
		outputs.HasLineContaining("Clone operation failed: "),
		outputs.HasLine(mrtclient.MsgFailedToCloneRepository),
	)
}

func Test_IfTeamJsonDoesNotExist_Cloning_ShouldPrintMessage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(
		t,
		0,
		"Could not read team file. To setup your repositories create a \"team.json\" file and add repositories to it.",
	)
}
