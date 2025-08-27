package core

const CommitMsg = "commit-msg"
const PreCommit = "pre-commit"
const PrePush = "pre-push"

func GetGitHooks() []string {
	return []string{
		CommitMsg,
		PreCommit,
		PrePush,
	}
}
