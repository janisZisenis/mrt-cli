package core

const CommitMsg = "commit-msg"
const PreCommit = "pre-commit"
const PrePush = "pre-push"

var GitHooks = []string{
	CommitMsg,
	PreCommit,
	PrePush,
}
