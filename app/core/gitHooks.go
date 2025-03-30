package core

const (
	CommitMsg = "commit-msg"
	PreCommit = "pre-commit"
	PrePush   = "pre-push"
)

var GitHooks = []string{
	CommitMsg,
	PreCommit,
	PrePush,
}
