package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/outputs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IfTeamDirIsMissing_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook("", "pre-commit", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "Missing team dir argument")
}

func Test_IfHookNameIsMissing_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "Missing hook name argument")
}

func Test_IfGitHookIsCalledWithUnknownHookName_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "unknown-hook", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"unknown-hook\" does not exist.")
}

func Test_IfGitHookIsCalledWithGlobbingPatternInHookName_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "pre-commit*", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given git-hook \"pre-commit*\" does not exist.")
}

func Test_IfGitHookIsCalledWithPathThatDoesNotContainRepository_ShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	nonRepoPath := f.AbsolutePath("non-repo")

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "pre-commit", nonRepoPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "The given path \""+nonRepoPath+"\" does not contain a repository: failed reading branch short name: exit status 128")
}

func Test_IfTeamJsonIsMissing_HookShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	require.NoError(t, os.Remove(f.AbsolutePath("team.json")))

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "pre-commit", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining("Failed to load team configuration"))
}

func Test_IfTeamJsonIsCorrupted_HookShouldFail(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	require.NoError(t, os.WriteFile(f.AbsolutePath("team.json"), []byte("not valid json {{{"), 0o600))

	output, exitCode := f.MakeMrtCommandInTeamDir().
		GitHook(f.TeamDir, "pre-commit", f.ClonedRepositoryPath).
		Execute()

	require.NotEqual(t, 0, exitCode)
	output.AssertInOrder(t, outputs.HasLineContaining("Failed to load team configuration"))
}
