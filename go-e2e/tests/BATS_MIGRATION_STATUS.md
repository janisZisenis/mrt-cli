# Bats Migration Status

## Files migrated to Go

### `e2e-tests/tests/git-hooks/git-hook-with-1-cloned-repository-test.bats`
**Status: safe to delete**

All 3 tests covered in `gitHook_withOneClonedRepository_test.go`.

---

### `e2e-tests/tests/setup-git-hooks/setup-git-hooks-additional-scripts-test.bats`
**Status: one weakened test — resolve before deleting**

Covered in `setupGitHooks_additionalScripts_test.go`.

| Bats test | Go test | Status |
|---|---|---|
| pre-commit scripts are executed | `Test_IfPreCommitScriptsExist_Committing_ShouldExecuteThem` | ✅ |
| commit-msg script failure causes commit to fail | `Test_IfCommitMsgScriptExitsWithFailure_Committing_ShouldAlsoFail` | ✅ |
| commit-msg script output appears in commit output | `Test_IfCommitMsgScriptHasOutput_Committing_ShouldSucceedAndScriptOutputShouldPassThrough` | ⚠️ Weakened — bats asserts the script's stdout appears in the commit output; Go only asserts exit code 0. Fix: capture stdout from `commitCommand` on success so the output can be asserted. |
| pre-commit hook receives empty params | `Test_IfPreCommitHookIsExecuted_ShouldReceiveEmptyParameters` | ✅ |
| pre-push hook receives remote name + URL | `Test_IfPrePushHookIsExecuted_ShouldReceiveRemoteNameAndURLAsParameters` | ✅ |
| commit-msg hook receives `.git/COMMIT_EDITMSG` | `Test_IfCommitMsgHookIsExecuted_ShouldReceiveCommitMsgFilePathAsParameter` | ✅ |

---

### `e2e-tests/tests/setup-git-hooks/setup-git-hooks-with-1-repository-cloned-test.bats`
**Status: safe to delete**

All 16 tests covered in `setupGitHooks_withOneClonedRepository_test.go`.
