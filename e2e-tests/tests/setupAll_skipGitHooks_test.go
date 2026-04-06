package tests_test

import (
	"testing"

	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/teamconfig"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skipGitHooksFixture struct {
	*fixtures.MrtFixture
	blockedBranchName string
	repositoryPath    string
}

func setupRepoWithBlockedBranchButSkippedHooks(t *testing.T, extraOptions ...teamconfig.Option) skipGitHooksFixture {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()
	repositoryName := "1_TestRepository"
	blockedBranchName := git.UniqueBranchName()
	options := append(
		[]teamconfig.Option{
			teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
			teamconfig.WithBlockedBranches([]string{blockedBranchName}),
		},
		extraOptions...,
	)
	f.TeamConfigWriter().Write(options...)
	f.MakeMrtCommand().
		Setup().
		All("--skip-install-git-hooks").
		Execute()

	return skipGitHooksFixture{
		MrtFixture:        f,
		blockedBranchName: blockedBranchName,
		repositoryPath:    f.AbsolutePath(defaultRepositoriesPath + "/" + repositoryName),
	}
}

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingOnABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	f := setupRepoWithBlockedBranchButSkippedHooks(t)

	exitCode, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch(f.blockedBranchName, "some-message").
		Execute()

	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_PushingToABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	f := setupRepoWithBlockedBranchButSkippedHooks(t)
	t.Cleanup(func() {
		f.MakeGitCommand().
			InDirectory(f.repositoryPath).
			DeleteRemoteBranchIfExists(f.blockedBranchName).
			Execute()
	})
	f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch(f.blockedBranchName, "some-message").
		Execute()

	exitCode, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		Push(f.blockedBranchName).
		Execute()

	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingWithMissingPrefixInCommitMessage_ShouldNotBeRejected(t *testing.T) {
	f := setupRepoWithBlockedBranchButSkippedHooks(t, teamconfig.WithCommitPrefixRegex("Some-Prefix"))

	exitCode, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch(f.blockedBranchName, "some-message").
		Execute()

	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}
