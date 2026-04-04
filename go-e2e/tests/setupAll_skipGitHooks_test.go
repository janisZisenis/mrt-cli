package tests_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/teamconfig"
)

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingOnABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	repositoryURL := git.MakeCloneURL(repositoryName)
	branchName := fmt.Sprintf("branch-%s", t.Name())
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{repositoryURL}),
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeMrtCommand().Setup().All("--skip-install-git-hooks").Execute()
	repositoryPath := f.AbsolutePath(defaultRepositoriesPath + "/" + repositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(repositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_PushingToABlockedBranch_ShouldNotBeRejected(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	repositoryURL := git.MakeCloneURL(repositoryName)
	branchName := fmt.Sprintf("branch-%s", t.Name())
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{repositoryURL}),
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeMrtCommand().Setup().All("--skip-install-git-hooks").Execute()
	repositoryPath := f.AbsolutePath(defaultRepositoriesPath + "/" + repositoryName)
	f.MakeGitCommand().InDirectory(repositoryPath).MakeCommitOnNewBranch(branchName, "some-message").Execute()

	exitCode, err := f.MakeGitCommand().
		InDirectory(repositoryPath).
		Push(branchName).
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func Test_IfSetupAllIsRunWithSkipGitHooks_CommittingWithMissingPrefixInCommitMessage_ShouldNotBeRejected(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate().Parallel()
	repositoryName := "1_TestRepository"
	repositoryURL := git.MakeCloneURL(repositoryName)
	branchName := fmt.Sprintf("branch-%s", t.Name())
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{repositoryURL}),
		teamconfig.WithBlockedBranches([]string{branchName}),
		teamconfig.WithCommitPrefixRegex("Some-Prefix"),
	)
	f.MakeMrtCommand().Setup().All("--skip-install-git-hooks").Execute()
	repositoryPath := f.AbsolutePath(defaultRepositoriesPath + "/" + repositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(repositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}
