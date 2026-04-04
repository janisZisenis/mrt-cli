package tests_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/teamconfig"
)

type skipGitHooksFixture struct {
	f                 *fixtures.MrtFixture
	blockedBranchName string
	repositoryPath    string
}

func setupRepoWithBlockedBranchButSkippedHooks(t *testing.T, extraOptions ...teamconfig.Option) skipGitHooksFixture {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	blockedBranchName := fmt.Sprintf("branch-%s", t.Name())
	options := append(
		[]teamconfig.Option{
			teamconfig.WithRepositories([]string{git.MakeCloneURL(repositoryName)}),
			teamconfig.WithBlockedBranches([]string{blockedBranchName}),
		},
		extraOptions...,
	)
	f.TeamConfigWriter().Write(options...)
	f.MakeMrtCommand().Setup().All("--skip-install-git-hooks").Execute()

	return skipGitHooksFixture{
		f:                 f,
		blockedBranchName: blockedBranchName,
		repositoryPath:    f.AbsolutePath(defaultRepositoriesPath + "/" + repositoryName),
	}
}

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingOnABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	fix := setupRepoWithBlockedBranchButSkippedHooks(t)

	exitCode, err := fix.f.MakeGitCommand().
		InDirectory(fix.repositoryPath).
		MakeCommitOnNewBranch(fix.blockedBranchName, "some-message").
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_PushingToABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	fix := setupRepoWithBlockedBranchButSkippedHooks(t)
	fix.f.MakeGitCommand().InDirectory(fix.repositoryPath).MakeCommitOnNewBranch(fix.blockedBranchName, "some-message").Execute()

	exitCode, err := fix.f.MakeGitCommand().
		InDirectory(fix.repositoryPath).
		Push(fix.blockedBranchName).
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingWithMissingPrefixInCommitMessage_ShouldNotBeRejected(t *testing.T) {
	fix := setupRepoWithBlockedBranchButSkippedHooks(t, teamconfig.WithCommitPrefixRegex("Some-Prefix"))

	exitCode, err := fix.f.MakeGitCommand().
		InDirectory(fix.repositoryPath).
		MakeCommitOnNewBranch(fix.blockedBranchName, "some-message").
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}
