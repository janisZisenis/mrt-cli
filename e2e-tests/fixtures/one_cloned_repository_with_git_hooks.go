package fixtures

import (
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"
)

const (
	ClonedRepoName          = "1_TestRepository"
	defaultRepositoriesPath = "repositories"
)

type OneClonedRepositoryWithGitHooksFixture struct {
	*MrtFixture
	RepositoryPath string
}

func MakeOneClonedRepositoryWithGitHooksFixture(
	t *testing.T,
	extraOptions ...teamconfig.Option,
) *OneClonedRepositoryWithGitHooksFixture {
	t.Helper()
	f := MakeMrtFixture(t).
		Authenticate().
		Parallel()

	options := append(
		[]teamconfig.Option{
			teamconfig.WithRepositories([]string{git.MakeCloneURL(ClonedRepoName)}),
		},
		extraOptions...,
	)
	f.TeamConfigWriter().Write(options...)

	f.MakeGitCommand().
		Clone(git.MakeCloneURL(ClonedRepoName), f.AbsolutePath(defaultRepositoriesPath+"/"+ClonedRepoName)).
		Execute()

	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()

	return &OneClonedRepositoryWithGitHooksFixture{
		MrtFixture:     f,
		RepositoryPath: f.AbsolutePath(defaultRepositoriesPath + "/" + ClonedRepoName),
	}
}

func (f *OneClonedRepositoryWithGitHooksFixture) GitInRepo() git.DirectedCommand {
	return f.MakeGitCommand().InDirectory(f.RepositoryPath)
}

func (f *OneClonedRepositoryWithGitHooksFixture) ConfigureBlockedBranches(branches []string) {
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL(ClonedRepoName)}),
		teamconfig.WithBlockedBranches(branches),
	)
	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()
}
