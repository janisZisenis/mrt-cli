package tests_test

import (
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfTeamJsonContainsInvalidCommitPrefixRegex_Committing_ShouldFailGracefully(t *testing.T) {
	invalidRegex := "[invalid(regex"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex(invalidRegex))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("test-branch", "test-message").Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Invalid commit prefix regex in team.json:")
	assert.Contains(t, err.Error(), "CommitPrefixRegex: "+invalidRegex)
	assert.Contains(t, err.Error(), "Please fix the regex syntax in your team.json file")
}

func Test_IfCommitMessageFileCannotBeRead_HookShouldFailGracefully(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Test-[0-9]+"))

	output, exitCode := f.MakeMrtCommand().
		GitHook("commit-msg", f.repositoryPath, "/nonexistent/file.txt").
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining("Failed to read commit message file"))
}

func Test_IfCommitMessageFileArgumentIsMissing_HookShouldFailGracefully(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Test-[0-9]+"))

	output, exitCode := f.MakeMrtCommand().
		GitHook("commit-msg", f.repositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "Missing commit message file argument")
}

func Test_IfTeamJsonContainsBlockedBranch_CommittingOnBlockedBranch_ShouldBeBlocked(t *testing.T) {
	branchName := "some-branch"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithBlockedBranches([]string{branchName}))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch(branchName, "some-message").Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"commit\" not allowed on branch \""+branchName+"\"")
}

func Test_IfTeamJsonContainsBlockedBranch_CommittingOnAnotherBranch_ShouldBeAllowed(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithBlockedBranches([]string{"some-branch"}))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("another-branch", "some-message").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonContains2BlockedBranches_CommittingOnSecondOne_ShouldBeBlocked(t *testing.T) {
	branchName := "some-branch"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithBlockedBranches([]string{"another-branch", branchName}))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch(branchName, "some-message").Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"commit\" not allowed on branch \""+branchName+"\"")
}

func Test_IfTeamJsonContainsBlockedBranch_PushingOnBlockedBranch_ShouldBeBlocked(t *testing.T) {
	branchName := git.UniqueBranchName()
	f := setupOneClonedRepositoryWithGitHooks(t)
	_, _ = f.gitInRepo().MakeCommitOnNewBranch(branchName, "some-message").Execute()
	f.configureBlockedBranches([]string{branchName})

	exitCode, err := f.gitInRepo().Push(branchName).Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"push\" not allowed on branch \""+branchName+"\"")

	t.Cleanup(func() { _, _ = f.gitInRepo().DeleteRemoteBranchIfExists(branchName).Execute() })
}

func Test_IfTeamJsonContainsBlockedBranch_PushingNonBlockedBranchWhileOnBlockedBranch_ShouldBeAllowed(t *testing.T) {
	blockedBranchName := git.UniqueBranchName()
	featureBranchName := git.UniqueBranchName()
	f := setupOneClonedRepositoryWithGitHooks(t)
	_, _ = f.gitInRepo().MakeCommitOnNewBranch(featureBranchName, "some-message").Execute()
	_, _ = f.gitInRepo().CheckoutNewBranch(blockedBranchName).Execute()
	f.configureBlockedBranches([]string{blockedBranchName})

	exitCode, err := f.gitInRepo().Push(featureBranchName).Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	t.Cleanup(func() { _, _ = f.gitInRepo().DeleteRemoteBranchIfExists(featureBranchName).Execute() })
}

func Test_IfTeamJsonContainsBlockedBranch_PushingBlockedBranchWhileOnAnotherBranch_ShouldBeBlocked(t *testing.T) {
	blockedBranchName := git.UniqueBranchName()
	anotherBranchName := git.UniqueBranchName()
	f := setupOneClonedRepositoryWithGitHooks(t)
	_, _ = f.gitInRepo().MakeCommitOnNewBranch(blockedBranchName, "some-message").Execute()
	_, _ = f.gitInRepo().CheckoutNewBranch(anotherBranchName).Execute()
	f.configureBlockedBranches([]string{blockedBranchName})

	exitCode, err := f.gitInRepo().Push(blockedBranchName).Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), "Action \"push\" not allowed on branch \""+blockedBranchName+"\"")

	t.Cleanup(func() { _, _ = f.gitInRepo().DeleteRemoteBranchIfExists(blockedBranchName).Execute() })
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingWithNeitherMessageNorBranchMatchingPrefix_ShouldBeBlocked(
	t *testing.T,
) {
	commitPrefixRegex := "Test-[0-9]+"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex(commitPrefixRegex))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("no-prefix-branch", "no-prefix-message").Execute()

	require.Error(t, err)
	require.NotEqual(t, 0, exitCode)
	assert.Contains(
		t,
		err.Error(),
		"The commit message needs a commit prefix that matches the following regex "+commitPrefixRegex+".",
	)
	assert.Contains(
		t,
		err.Error(),
		"Either add the commit prefix to your commit message, or include it in the branch name.",
	)
	assert.Contains(t, err.Error(), "Use '--no-verify' to skip git-hooks.")
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingWithMatchingPrefixInMessage_ShouldNotBeBlocked(
	t *testing.T,
) {
	matchingPrefix := "Test-1"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Test-[0-9]+"))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("no-prefix-branch", matchingPrefix+": prefixed-message").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingOnBranchContainingPrefix_ShouldNotBeBlocked(
	t *testing.T,
) {
	commitPrefix := "Asdf-99"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Asdf-[0-9]+"))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("feature/"+commitPrefix+"/prefixed-branch", "not-prefix-message").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingWithMergeBranchMessage_ShouldNotBeBlocked(
	t *testing.T,
) {
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Asdf-[0-9]+"))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("no-prefix-branch", "Merge branch").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingWithMergeRemoteTrackingBranchMessage_ShouldNotBeBlocked(
	t *testing.T,
) {
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Asdf-[0-9]+"))

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("no-prefix-branch", "Merge remote-tracking branch").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingWithMatchingPrefixInMessage_CommitMessageShouldBePreserved(
	t *testing.T,
) {
	commitMessage := "Test-1: prefixed-message"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Test-[0-9]+"))
	_, _ = f.gitInRepo().MakeCommitOnNewBranch("no-prefix-branch", commitMessage).Execute()

	lastCommitMessage, err := f.gitInRepo().GetLastCommitMessage()

	require.NoError(t, err)
	assert.Equal(t, commitMessage, lastCommitMessage)
}

func Test_IfTeamJsonContainsCommitPrefixRegex_CommittingOnBranchContainingPrefix_CommitMessageShouldBePrefixed(
	t *testing.T,
) {
	matchingPrefix := "Asdf-99"
	commitMessage := "not-prefixed-message"
	f := setupOneClonedRepositoryWithGitHooks(t, teamconfig.WithCommitPrefixRegex("Asdf-[0-9]+"))
	_, _ = f.gitInRepo().MakeCommitOnNewBranch("feature/"+matchingPrefix+"/prefixed-branch", commitMessage).Execute()

	lastCommitMessage, err := f.gitInRepo().GetLastCommitMessage()

	require.NoError(t, err)
	assert.Equal(t, matchingPrefix+": "+commitMessage, lastCommitMessage)
}

func Test_IfTeamJsonDoesNotContainCommitPrefixRegex_Committing_ShouldNotCheckForPrefix(
	t *testing.T,
) {
	f := setupOneClonedRepositoryWithGitHooks(t)

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("not-prefixed-branch", "not-prefixed-message").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonDoesNotContainCommitPrefixRegex_CommittingWithMergeBranchMessage_ShouldNotCheckForPrefix(
	t *testing.T,
) {
	f := setupOneClonedRepositoryWithGitHooks(t)

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("not-prefixed-branch", "Merge branch").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func Test_IfTeamJsonDoesNotContainCommitPrefixRegex_CommittingWithMergeRemoteTrackingBranchMessage_ShouldNotCheckForPrefix(
	t *testing.T,
) {
	f := setupOneClonedRepositoryWithGitHooks(t)

	exitCode, err := f.gitInRepo().MakeCommitOnNewBranch("not-prefixed-branch", "Merge remote-tracking branch").Execute()

	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}
