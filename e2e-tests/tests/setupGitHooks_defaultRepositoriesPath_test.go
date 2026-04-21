package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	mrtclient "mrt-cli/e2e-tests/mrt"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfTeamJsonIsMissing_InstallGitHooks_ShouldFail(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.Equal(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining(mrtclient.MsgFailedToLoadTeamConfiguration))
}

func Test_IfTeamJsonIsCorrupted_InstallGitHooks_ShouldFail(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	require.NoError(t, os.WriteFile(f.AbsolutePath("team.json"), []byte("not valid json {{{"), 0o600))

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.Equal(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining(mrtclient.MsgFailedToLoadTeamConfiguration))
}

func Test_IfRepositoriesPathIsAbsolute_InstallGitHooks_ShouldExitNonZeroAndPrintError(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath("/absolute/path"),
	)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining(mrtclient.MsgInvalidRepositoriesPath))
}

func Test_IfRepositoriesPathEscapesTeamDir_InstallGitHooks_ShouldExitNonZeroAndPrintError(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath("../outside"),
	)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining(mrtclient.MsgInvalidRepositoriesPath))
}

func Test_IfRepositoriesPathIsValidRelativePath_InstallGitHooks_ShouldSucceed(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositoriesPath("some-relative-path"),
	)

	_, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.Equal(t, 0, exitCode)
}

func Test_IfRepositoriesPathContainsNonRepositoryFolder_InstallGitHooks_ShouldNotInstallGitHooks(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	folderPath := f.AbsolutePath(mrtclient.DefaultRepositoriesPath + "/1_TestRepository")
	require.NoError(t, os.MkdirAll(folderPath, 0o750))

	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	f.AssertFolderDoesNotExist(mrtclient.DefaultRepositoriesPath + "/1_TestRepository/.git/hooks")
}

func Test_IfRepositoriesPathContains2Repositories_CommittingOnBlockedBranchInSecondRepo_ShouldBeBlocked(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
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
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(secondRepositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+secondRepositoryName)).
		Execute()
	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()
	secondRepositoryPath := f.AbsolutePath(mrtclient.DefaultRepositoriesPath + "/" + secondRepositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(secondRepositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	require.Error(t, err)
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), mrtclient.MsgActionNotAllowedOnBranch("commit", branchName))
}

func Test_IfRepositoriesPathContains2Repositories_InstallGitHooks_ShouldPrintMessages(
	t *testing.T,
) {
	f := fixtures.MakeMrtFixture(t).
		Authenticate()
	firstRepositoryName := "1_TestRepository"
	secondRepositoryName := "2_TestRepository"
	repositoriesDir := f.AbsolutePath(mrtclient.DefaultRepositoriesPath)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(firstRepositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+firstRepositoryName)).
		Execute()
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(secondRepositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+secondRepositoryName)).
		Execute()

	output, _ := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	output.AssertInOrder(
		t,
		outputs.HasLine(mrtclient.MsgInstallingGitHooksToRepositoriesLocatedIn(repositoriesDir)),
		outputs.HasLine(
			"Installing git-hooks to \""+repositoriesDir+"/"+firstRepositoryName+"/.git\"",
		),
		outputs.HasLine(
			"Done installing git-hooks to \""+repositoriesDir+"/"+firstRepositoryName+"/.git\"",
		),
		outputs.HasLine(
			"Installing git-hooks to \""+repositoriesDir+"/"+secondRepositoryName+"/.git\"",
		),
		outputs.HasLine(
			"Done installing git-hooks to \""+repositoriesDir+"/"+secondRepositoryName+"/.git\"",
		),
		outputs.HasLine("Done installing git-hooks."),
	)
}
