package tests_test

import (
	"testing"

	"mrt-cli/go-e2e/fixtures"
	"mrt-cli/go-e2e/git"
	"mrt-cli/go-e2e/teamconfig"
)

const oneClonedRepoName = "1_TestRepository"

type oneClonedRepositoryWithGitHooksFixture struct {
	*fixtures.MrtFixture
	repositoryPath string
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
