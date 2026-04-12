package core

const (
	ApplypatchMsg       = "applypatch-msg"
	CommitMsg           = "commit-msg"
	FsmonitorWatchman   = "fsmonitor-watchman"
	P4Changelist        = "p4-changelist"
	P4PostChangelist    = "p4-post-changelist"
	P4PrepareChangelist = "p4-prepare-changelist"
	P4PreSubmit         = "p4-pre-submit"
	PostApplypatch      = "post-applypatch"
	PostCheckout        = "post-checkout"
	PostCommit          = "post-commit"
	PostIndexChange     = "post-index-change"
	PostMerge           = "post-merge"
	PostRewrite         = "post-rewrite"
	PreApplypatch       = "pre-applypatch"
	PreAutoGc           = "pre-auto-gc"
	PreCommit           = "pre-commit"
	PreMergeCommit      = "pre-merge-commit"
	PrePush             = "pre-push"
	PreRebase           = "pre-rebase"
	PrepareCommitMsg    = "prepare-commit-msg"
	SendemailValidate   = "sendemail-validate"
)

func IsGitHook(name string) bool {
	for _, hook := range GetGitHooks() {
		if hook == name {
			return true
		}
	}
	return false
}

func GetGitHooks() []string {
	return []string{
		ApplypatchMsg,
		CommitMsg,
		FsmonitorWatchman,
		P4Changelist,
		P4PostChangelist,
		P4PrepareChangelist,
		P4PreSubmit,
		PostApplypatch,
		PostCheckout,
		PostCommit,
		PostIndexChange,
		PostMerge,
		PostRewrite,
		PreApplypatch,
		PreAutoGc,
		PreCommit,
		PreMergeCommit,
		PrePush,
		PreRebase,
		PrepareCommitMsg,
		SendemailValidate,
	}
}
