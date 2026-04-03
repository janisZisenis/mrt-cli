package git

func MakeCloneURL(repositoryName string) string {
	return "git@github-testing:janisZisenisTesting/" + repositoryName + ".git"
}
