package core

const (
	CommitMsg = "commit-msg"
	PreCommit = "pre-commit"
	PrePush   = "pre-push"
)

func GetGitHooks() []string {
	return []string{
		CommitMsg,
		PreCommit,
		PrePush,
	}
}
