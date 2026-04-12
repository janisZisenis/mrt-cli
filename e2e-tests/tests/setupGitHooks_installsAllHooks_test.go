package tests_test

import (
	"mrt-cli/e2e-tests/assert"
	"mrt-cli/e2e-tests/fixtures"
	"path/filepath"
	"testing"
)

func Test_InstallGitHooks_ShouldInstallAllGitHooks(t *testing.T) {
	f := fixtures.MakeOneClonedRepositoryWithGitHooksFixture(t)
	hookNames := []string{
		"applypatch-msg",
		"pre-applypatch",
		"post-applypatch",
		"pre-commit",
		"pre-merge-commit",
		"prepare-commit-msg",
		"commit-msg",
		"post-commit",
		"pre-rebase",
		"post-checkout",
		"post-merge",
		"pre-push",
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
			hookFile := filepath.Join(f.ClonedRepositoryPath, ".git", "hooks", hookName)
			assert.FileContains(t, hookFile, "#!/bin/bash -e")
			assert.FileContains(t, hookFile, "git-hook")
			assert.FileContains(t, hookFile, "--hook-name \"$hook_name\"")
			assert.FileContains(t, hookFile, "--repository-path $PWD")
			assert.FileHasPermissions(t, hookFile, 0o700)
		})
	}
}
