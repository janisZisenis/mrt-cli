package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IfGitHookIsCalledWithUnknownHookName_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommand().
		GitHook("unknown-hook", f.RepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"unknown-hook\" does not exist.")
}

func Test_IfGitHookIsCalledWithGlobbingPatternInHookName_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommand().
		GitHook("pre-commit*", f.RepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"pre-commit*\" does not exist.")
}

func Test_IfGitHookIsCalledWithPathThatDoesNotContainRepository_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	nonRepoPath := f.AbsolutePath("non-repo")

	output, exitCode := f.MakeMrtCommand().
		GitHook("pre-commit", nonRepoPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given path \""+nonRepoPath+"\" does not contain a repository.")
}
