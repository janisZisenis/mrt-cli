package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"mrt-cli/e2e-tests/git"
	"mrt-cli/e2e-tests/teamconfig"
	"testing"
)

const oneClonedRepoName = "1_TestRepository"

type oneClonedRepositoryWithGitHooksFixture struct {
	*fixtures.MrtFixture
	repositoryPath string
}

func (f oneClonedRepositoryWithGitHooksFixture) gitInRepo() git.DirectedCommand {
	return f.MakeGitCommand().InDirectory(f.repositoryPath)
}

func (f oneClonedRepositoryWithGitHooksFixture) configureBlockedBranches(branches []string) {
	f.TeamConfigWriter().Write(
		teamconfig.WithRepositories([]string{git.MakeCloneURL(oneClonedRepoName)}),
		teamconfig.WithBlockedBranches(branches),
	)
	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()
}

func setupOneClonedRepositoryWithGitHooks(
	t *testing.T,
	extraOptions ...teamconfig.Option,
) oneClonedRepositoryWithGitHooksFixture {
	t.Helper()
	f := fixtures.MakeMrtFixture(t).
		Authenticate().
		Parallel()

	options := append(
		[]teamconfig.Option{
			teamconfig.WithRepositories([]string{git.MakeCloneURL(oneClonedRepoName)}),
		},
		extraOptions...,
	)
	f.TeamConfigWriter().Write(options...)

	f.MakeGitCommand().
		Clone(git.MakeCloneURL(oneClonedRepoName), f.AbsolutePath(defaultRepositoriesPath+"/"+oneClonedRepoName)).
		Execute()

	f.MakeMrtCommand().
		Setup().
		InstallGitHooks().
		Execute()

	return oneClonedRepositoryWithGitHooksFixture{
		MrtFixture:     f,
		repositoryPath: f.AbsolutePath(defaultRepositoriesPath + "/" + oneClonedRepoName),
	}
}
