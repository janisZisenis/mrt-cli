package tests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
)

func Test_IfPreCommitScriptsExist_Committing_ShouldExecuteThem(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteSpyScript("pre-commit", "script1")
	hooks.WriteSpyScript("pre-commit", "script2")

	_, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch("some-branch", "some-message").
		Execute()

	require.NoError(t, err)
	hooks.AssertWasExecuted(t, "pre-commit", "script1")
	hooks.AssertWasExecuted(t, "pre-commit", "script2")
}

func Test_IfCommitMsgScriptExitsWithFailure_Committing_ShouldAlsoFail(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteStubScript("commit-msg", "script", 1, "some-output")

	exitCode, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch("some-branch", "some-message").
		Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
}

func Test_IfCommitMsgScriptHasOutput_Committing_ShouldContainThatOutput(t *testing.T) {
	scriptOutput := "some-output"
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteStubScript("commit-msg", "script", 0, scriptOutput)

	output, exitCode, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch("some-branch", "some-message").
		ExecuteAndCaptureOutput()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	assert.Contains(t, output, scriptOutput)
}

func Test_IfPreCommitHookIsExecuted_ShouldReceiveEmptyParameters(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteSpyScript("pre-commit", "script")

	_, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch("some-branch", "some-message").
		Execute()

	require.NoError(t, err)
	hooks.AssertWasExecutedWith(t, "pre-commit", "script", "")
}

func Test_IfPrePushHookIsExecuted_ShouldReceiveRemoteNameAndURLAsParameters(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteSpyScript("pre-push", "script")
	branchName := git.UniqueBranchName()
	t.Cleanup(func() {
		f.MakeGitCommand().
			InDirectory(f.repositoryPath).
			DeleteRemoteBranchIfExists(branchName).
			Execute()
	})
	_, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()
	require.NoError(t, err)

	_, pushErr := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		Push(branchName).
		Execute()

	require.NoError(t, pushErr)
	hooks.AssertWasExecutedWith(t, "pre-push", "script", "origin "+git.MakeCloneURL(oneClonedRepoName))
}

func Test_IfCommitMsgHookIsExecuted_ShouldReceiveCommitMsgFilePathAsParameter(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	hooks := fixtures.NewHookScriptFixture(f.repositoryPath)
	hooks.WriteSpyScript("commit-msg", "script")

	_, err := f.MakeGitCommand().
		InDirectory(f.repositoryPath).
		MakeCommitOnNewBranch("some-branch", "some-message").
		Execute()

	require.NoError(t, err)
	hooks.AssertWasExecutedWith(t, "commit-msg", "script", ".git/COMMIT_EDITMSG")
}
