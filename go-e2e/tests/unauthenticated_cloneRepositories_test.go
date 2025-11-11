package tests

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/utils"
)

func Test_IfTeamJsonContainsRepositoryAndSomeRepositoryPath_Cloning_ShouldPrintOutMessageAboutCloning(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	repositoryName := "1_TestRepository"
	f.WriteTeamJson(
		utils.WithRepositoriesPath("some-path"),
		utils.WithRepositories([]string{utils.MakeCloneUrlFrom(repositoryName)}),
	)

	output := f.MakeMrtCommand().
		Setup().
		Clone().
		Execute()

	output.AssertLineEquals(t, 0, "Start cloning repositories into \"some-path\"")
	output.AssertLineEquals(t, 9, "Cloning repositories done")
}
