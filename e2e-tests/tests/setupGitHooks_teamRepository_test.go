package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	mrtclient "mrt-cli/e2e-tests/mrt"
	"mrt-cli/e2e-tests/outputs"
	"mrt-cli/e2e-tests/teamconfig"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfTeamDirIsAGitRepo_InstallGitHooks_ShouldPrintMessages(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)
	require.NoError(t, exec.Command("git", "init", f.TeamDir).Run())

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.Equal(t, 0, exitCode)
	output.AssertInOrder(t,
		outputs.HasLine(mrtclient.MsgInstallingGitHooksToTeamRepository),
		outputs.HasLine(mrtclient.MsgDoneInstallingGitHooksToTeamRepository),
	)
}

func Test_IfTeamDirIsNotAGitRepo_InstallGitHooks_ShouldNotMentionTeamRepository(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)

	output, exitCode := f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	require.Equal(t, 0, exitCode)
	output.AssertHasNoLineContaining(t, "team repository")
}

func Test_IfTeamDirIsAGitRepo_HooksSurviveMovingThePackage(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Authenticate()
	repositoryName := "1_TestRepository"
	branchName := git.UniqueBranchName()
	f.TeamConfigWriter().Write(
		teamconfig.WithBlockedBranches([]string{branchName}),
	)
	f.MakeGitCommand().
		Clone(git.MakeCloneURL(repositoryName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+repositoryName)).
		Execute()
	require.NoError(t, exec.Command("git", "init", f.TeamDir).Run())
	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()
	newTeamDir := filepath.Join(filepath.Dir(f.TeamDir), "moved-team-dir")
	require.NoError(t, os.Rename(f.TeamDir, newTeamDir))
	newRepoPath := filepath.Join(newTeamDir, mrtclient.DefaultRepositoriesPath, repositoryName)

	exitCode, err := f.MakeGitCommand().
		InDirectory(newRepoPath).
		MakeCommitOnNewBranch(branchName, "some-message").
		Execute()

	require.Error(t, err)
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, err.Error(), mrtclient.MsgActionNotAllowedOnBranch("commit", branchName))
}
