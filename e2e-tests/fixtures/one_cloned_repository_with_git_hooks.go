package fixtures

import (
	"mrt-cli/e2e-tests/git"
	mrtclient "mrt-cli/e2e-tests/mrt"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"
)

const ClonedRepoName = "1_TestRepository"

type OneClonedRepositoryWithGitHooksFixture struct {
	*MrtFixture
	ClonedRepositoryPath string
}

func MakeOneClonedRepositoryWithGitHooksFixture(
	t *testing.T,
	extraOptions ...teamconfig.Option,
) *OneClonedRepositoryWithGitHooksFixture {
	t.Helper()
	f := MakeMrtFixture(t).
		Authenticate()

	options := append(
		[]teamconfig.Option{
			teamconfig.WithRepositories([]string{git.MakeCloneURL(ClonedRepoName)}),
		},
		extraOptions...,
	)
	f.TeamConfigWriter().Write(options...)

	f.MakeGitCommand().
		Clone(git.MakeCloneURL(ClonedRepoName), f.AbsolutePath(mrtclient.DefaultRepositoriesPath+"/"+ClonedRepoName)).
		Execute()

	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()

	return &OneClonedRepositoryWithGitHooksFixture{
		MrtFixture:           f,
		ClonedRepositoryPath: f.AbsolutePath(mrtclient.DefaultRepositoriesPath + "/" + ClonedRepoName),
	}
}

func (f *OneClonedRepositoryWithGitHooksFixture) GitInClonedRepository() git.DirectedCommand {
	return f.MakeGitCommand().InDirectory(f.ClonedRepositoryPath)
}

func (f *OneClonedRepositoryWithGitHooksFixture) ConfigureBlockedBranches(branches []string) {
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL(ClonedRepoName)}),
		teamconfig.WithBlockedBranches(branches),
	)
	f.MakeMrtCommandInTeamDir().
		Setup().
		InstallGitHooks().
		Execute()
}
