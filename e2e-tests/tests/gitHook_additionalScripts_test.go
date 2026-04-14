package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"testing"
)

func Test_IfAdditionalScriptExists_GitHookShouldExecuteIt(t *testing.T) {
	// pre-commit, pre-push, and commit-msg are excluded here because they have
	// built-in mrt logic and are covered via real git operations in
	// setupGitHooks_additionalScripts_test.go.
	hookNames := []string{
		"applypatch-msg",
		"pre-applypatch",
		"post-applypatch",
		"pre-merge-commit",
		"prepare-commit-msg",
		"post-commit",
		"pre-rebase",
		"post-checkout",
		"post-merge",
		"post-rewrite",
		"pre-auto-gc",
		"sendemail-validate",
		"fsmonitor-watchman",
		"post-index-change",
		"p4-changelist",
		"p4-prepare-changelist",
		"p4-post-changelist",
		"p4-pre-submit",
	}

	for _, hookName := range hookNames {
		t.Run(hookName, func(t *testing.T) {
			f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
			hooks := fixtures.NewHookScriptFixture(f.ClonedRepositoryPath)
			hooks.WriteSpyScript(hookName, "script")

			f.MakeMrtCommandInTeamDir().
				GitHook(f.TeamDir, hookName, f.ClonedRepositoryPath).
				Execute()

			hooks.AssertWasExecuted(t, hookName, "script")
		})
	}
}
