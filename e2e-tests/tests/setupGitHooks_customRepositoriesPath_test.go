package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfRepositoriesPathIsDot_CommittingOnBlockedBranch_ShouldBeBlocked(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoryName := "1_TestRepository"
	branchName := git.UniqueBranchName()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath("."),
		teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(repositoryName), f.AbsolutePath(repositoryName)).
		Execute()
	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()
	repositoryPath := f.AbsolutePath(repositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(repositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	require.Error(t, err)
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"commit\" not allowed on branch \""+branchName+"\"")
}

func Test_IfRepositoriesAreClonedToCustomPath_CommittingOnBlockedBranch_ShouldBeBlocked(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoriesPath := "some-path"
	repositoryName := "1_TestRepository"
	branchName := git.UniqueBranchName()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath(repositoriesPath),
		teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(repositoryName), f.AbsolutePath(repositoriesPath+"/"+repositoryName)).
		Execute()
	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()
	repositoryPath := f.AbsolutePath(repositoriesPath + "/" + repositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(repositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	require.Error(t, err)
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"commit\" not allowed on branch \""+branchName+"\"")
}

func Test_IfCustomRepositoriesPathDoesNotContainRepositories_InstallGitHooks_ShouldPrintNotFoundMessage(
	t *testing.T,
) {
	tests := []string{"some-path", "another-path"}

	for _, repositoriesPath := range tests {
		t.Run(repositoriesPath, func(t *testing.T) {
			testIfCustomRepositoriesPathDoesNotContainRepositoriesInstallGitHooksShouldPrintNotFoundMessage(
				t,
				repositoriesPath,
			)
		})
	}
}

func testIfCustomRepositoriesPathDoesNotContainRepositoriesInstallGitHooksShouldPrintNotFoundMessage(
	t *testing.T,
	repositoriesPath string,
) {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	repositoriesDir := f.AbsolutePath(repositoriesPath)
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath(repositoriesPath),
	)

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Installing git-hooks to repositories located in \""+repositoriesDir+"\""),
		outputs.HasLine("Did not find any repositories. Skip installing git-hooks."),
		outputs.HasLine("Done installing git-hooks."),
	)
}
