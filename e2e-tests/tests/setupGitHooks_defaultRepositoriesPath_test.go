package tests_test

import (
	"os"
	"testing"

	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfRepositoriesPathContainsNonRepositoryFolder_InstallGitHooks_ShouldNotInstallGitHooks(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	folderPath := f.AbsolutePath(defaultRepositoriesPath + "/1_TestRepository")
	require.NoError(t, os.MkdirAll(folderPath, 0o750))

	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()

	f.AssertFolderDoesNotExist(defaultRepositoriesPath + "/1_TestRepository/.git/hooks")
}

func Test_IfRepositoriesPathContains2Repositories_CommittingOnBlockedBranchInSecondRepo_ShouldBeBlocked(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	branchName := git.UniqueBranchName()
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{
			git.MakeCloneURL(firstRepositoryName),
			git.MakeCloneURL(secondRepositoryName),
		}),
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(defaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(secondRepositoryName), f.AbsolutePath(defaultRepositoriesPath+"/"+secondRepositoryName)).
		Execute()
	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()
	secondRepositoryPath := f.AbsolutePath(defaultRepositoriesPath + "/" + secondRepositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(secondRepositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	require.Error(t, err)
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"commit\" not allowed on branch \""+branchName+"\"")
}

func Test_IfRepositoriesPathContains2Repositories_InstallGitHooks_ShouldPrintMessages(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	repositoriesDir := f.AbsolutePath(defaultRepositoriesPath)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(defaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(secondRepositoryName), f.AbsolutePath(defaultRepositoriesPath+"/"+secondRepositoryName)).
		Execute()

	output, _ := f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()

	output.AssertInOrder(t,
		outputs.HasLine("Installing git-hooks to repositories located in \""+repositoriesDir+"\""),
		outputs.HasLine("Installing git-hooks to \""+repositoriesDir+"/"+firstRepositoryName+"/.git\""),
		outputs.HasLine("Done installing git-hooks to \""+repositoriesDir+"/"+firstRepositoryName+"/.git\""),
		outputs.HasLine("Installing git-hooks to \""+repositoriesDir+"/"+secondRepositoryName+"/.git\""),
		outputs.HasLine("Done installing git-hooks to \""+repositoriesDir+"/"+secondRepositoryName+"/.git\""),
		outputs.HasLine("Done installing git-hooks."),
	)
}
