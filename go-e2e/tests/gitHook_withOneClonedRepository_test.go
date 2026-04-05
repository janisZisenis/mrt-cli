package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IfGitHookIsCalledWithUnknownHookName_ShouldFail(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)

	output, exitCode := f.MakeMrtCommand().
		GitHook("unknown-hook", f.repositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"unknown-hook\" does not exist.")
}

func Test_IfGitHookIsCalledWithGlobbingPatternInHookName_ShouldFail(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)

	output, exitCode := f.MakeMrtCommand().
		GitHook("pre-commit*", f.repositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"pre-commit*\" does not exist.")
}

func Test_IfGitHookIsCalledWithPathThatDoesNotContainRepository_ShouldFail(t *testing.T) {
	f := setupOneClonedRepositoryWithGitHooks(t)
	nonRepoPath := f.TeamDir()

	output, exitCode := f.MakeMrtCommand().
		GitHook("pre-commit", nonRepoPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given path \""+nonRepoPath+"\" does not contain a repository.")
}
