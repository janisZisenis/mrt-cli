package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfPreCommitScriptsExist_Committing_ShouldExecuteThem(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteSpyScript("pre-commit", "script1")
	hooks.WriteSpyScript("pre-commit", "script2")

	_, err := f.GitInRepo().MakeCommitOnNewBranch("some-branch", "some-message").Execute()

	require.NoError(t, err)
	hooks.AssertWasExecuted(t, "pre-commit", "script1")
	hooks.AssertWasExecuted(t, "pre-commit", "script2")
}

func Test_IfCommitMsgScriptExitsWithFailure_Committing_ShouldAlsoFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteStubScript("commit-msg", "script", 1, "some-output")

	exitCode, err := f.GitInRepo().MakeCommitOnNewBranch("some-branch", "some-message").Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
}

func Test_IfCommitMsgScriptHasOutput_Committing_ShouldContainThatOutput(t *testing.T) {
	scriptOutput := "some-output"
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteStubScript("commit-msg", "script", 0, scriptOutput)

	output, exitCode, err := f.GitInRepo().MakeCommitOnNewBranch("some-branch", "some-message").ExecuteAndCaptureOutput()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	assert.Contains(t, output, scriptOutput)
}

func Test_IfPreCommitHookIsExecuted_ShouldReceiveEmptyParameters(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteSpyScript("pre-commit", "script")

	_, err := f.GitInRepo().MakeCommitOnNewBranch("some-branch", "some-message").Execute()

	require.NoError(t, err)
	hooks.AssertWasExecutedWith(t, "pre-commit", "script", "")
}

func Test_IfPrePushHookIsExecuted_ShouldReceiveRemoteNameAndURLAsParameters(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteSpyScript("pre-push", "script")
	branchName := git.UniqueBranchName()
	t.Cleanup(func() { _, _ = f.GitInRepo().DeleteRemoteBranchIfExists(branchName).Execute() })
	_, err := f.GitInRepo().MakeCommitOnNewBranch(branchName, "some-message").Execute()
	require.NoError(t, err)

	_, pushErr := f.GitInRepo().Push(branchName).Execute()

	require.NoError(t, pushErr)
	hooks.AssertWasExecutedWith(
		t,
		"pre-push",
		"script",
		"origin "+git.MakeCloneURL(fixtures.ClonedRepoName),
	)
}

func Test_IfCommitMsgHookIsExecuted_ShouldReceiveCommitMsgFilePathAsParameter(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
	hooks.WriteSpyScript("commit-msg", "script")

	_, err := f.GitInRepo().MakeCommitOnNewBranch("some-branch", "some-message").Execute()

	require.NoError(t, err)
	hooks.AssertWasExecutedWith(t, "commit-msg", "script", ".git/COMMIT_EDITMSG")
}
