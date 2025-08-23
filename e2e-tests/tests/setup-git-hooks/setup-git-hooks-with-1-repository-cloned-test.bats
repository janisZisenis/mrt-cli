setup() {
  bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'
  bats_load_library 'writeTeamFile.bash'

  one_cloned_repository_with_git_hooks_setup
}

teardown() {
  one_cloned_repository_with_git_hooks_teardown
}

@test "If team json contains blocked branch commiting on the blocked branches after setting up git-hooks should be blocked" {
	branchName="some-branch"
	writeBlockedBranches "$branchName"

	run commit_changes "$(repositoryDir)" $branchName

	assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "If team json contains blocked branch commiting on another blocked branch after setting up git-hooks should be allowed" {
	writeBlockedBranches "some-branch"

	run commit_changes "$(repositoryDir)" "another-branch"

	assert_success
}

@test "If team json contains 2 blocked branch commiting on second one after setting up git-hooks should be blocked" {
	branchName="some-branch"
	writeBlockedBranches "another-branch" "$branchName"

	run commit_changes "$(repositoryDir)" $branchName

	assert_output --partial "Action \"commit\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "If team json contains blocked branch pushing on the blocked after setting up git-hooks branch should be blocked" {
	branchName="$(unique_branch_name)"
	commit_changes "$(repositoryDir)" "$branchName"
	writeBlockedBranches "$branchName"

	push_changes "$(repositoryDir)" "$branchName"

	assert_output --partial "Action \"push\" not allowed on branch \"$branchName\""
	assert_failure
}

@test "if team json contains commitPrefixRegex 'commiting' with a message and a branch both without matching prefix is blocked" {
	commitPrefixRegex="Test-[0-9]+"
	writeCommitPrefixRegex "$commitPrefixRegex"

	run commit_changes "$(repositoryDir)" "no-prefix-branch" "no-prefix-message"

	assert_line --index 1 "The commit message needs a commit prefix, that matches the following regex $commitPrefixRegex."
	assert_line --index 2 "Either add the commit prefix to you commit message, or include it in the branch name."
	assert_line --index 3 "Use '--no-verify' to skip git-hooks."
	assert_failure
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has a matching prefix on a branch without prefix is not blocked" {
	matchingPrefix="Test-1"
	writeCommitPrefixRegex "Test-[0-9]+"

	run commit_changes "$(repositoryDir)" "no-prefix-branch" "$matchingPrefix: prefixed-message"

	assert_line --index 1 "The commit message contains an issue ID ($matchingPrefix). Good job!"
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix is not blocked" {
	commitPrefix=Asdf-99
	writeCommitPrefixRegex "Asdf-[0-9]+"

	run commit_changes "$(repositoryDir)" "feature/$commitPrefix/prefixed-branch" "not-prefix-message"

	assert_line --index 1 "Commit prefix '$commitPrefix' was found in current branch name, prepended to commit message."
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that starts with 'Merge branch' is not blocked" {
	test_merge_commit_messages_are_not_blocked "Merge branch"
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that starts with 'Merge remote-tracking branch' is not blocked" {
	test_merge_commit_messages_are_not_blocked "Merge remote-tracking branch"
}

test_merge_commit_messages_are_not_blocked() {
	commit_message=$1
	writeCommitPrefixRegex "Asdf-[0-9]+"

	run commit_changes "$(repositoryDir)" "no-prefix-branch" "$commit_message"

	assert_line --index 1 "Merge commit detected, skipping commit-msg hook."
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has a matching prefix leads to commit with given message" {
	commitMessage="Test-1: prefixed-message"
	writeCommitPrefixRegex "Test-[0-9]+"
	commit_changes "$(repositoryDir)" "no-prefix-branch" "$commitMessage"

	run get_commit_message_of_last_commit "$(repositoryDir)"

	assert_output "$commitMessage"
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix leads to commit with prefixed message" {
	matchingPrefix=Asdf-99
	commitMessage="not-prefixed-message"
	writeCommitPrefixRegex "Asdf-[0-9]+"
	commit_changes "$(repositoryDir)" "feature/$matchingPrefix/prefixed-branch" "$commitMessage"

	run get_commit_message_of_last_commit "$(repositoryDir)"
	assert_output "$matchingPrefix: $commitMessage"
}

@test "if team json does not contain commitPrefixRegex 'commiting' does not check for commit prefix" {
	run commit_changes "$(repositoryDir)" "not-prefixed-branch" "not-prefixed-message"

	refute_output --partial "JIRA-ID '' was found in current branch name, prepended to commit message."
}

@test "if team json does not contain commitPrefixRegex while 'commiting' a merge commit, it does not check for commit prefix" {
	run commit_changes "$(repositoryDir)" "not-prefixed-branch" "not-prefixed-message"

	refute_output --partial "JIRA-ID '' was found in current branch name, prepended to commit message."
}

@test "if team json does not contain commitPrefixRegex 'commiting' with a message that starts with 'Merge branch' does not check for prefix" {
	test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes "Merge branch"
}

@test "if team json does not contain commitPrefixRegex 'commiting' with a message that starts with 'Merge remote-tracking branch' does not check for prefix" {
	test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes "Merge remote-tracking branch"
}

test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes() {
	commit_message=$1

	run commit_changes "$(repositoryDir)" "no-prefix-branch" "$commit_message"

	refute_output --partial "Merge commit detected, skipping commit-msg hook."
	assert_success
}
