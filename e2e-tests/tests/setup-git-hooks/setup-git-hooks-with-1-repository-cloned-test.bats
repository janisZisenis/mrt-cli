setup() {
  bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'
  bats_load_library 'write_team_file.bash'

  one_cloned_repository_with_git_hooks_setup
}

teardown() {
  one_cloned_repository_with_git_hooks_teardown
}

@test "If team json contains blocked branch commiting on the blocked branches after setting up git-hooks should be blocked" {
	local branch_name="some-branch"
	write_blocked_branches "$branch_name"

	run commit_changes "$(repository_dir)" $branch_name

	assert_output --partial "Action \"commit\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "If team json contains blocked branch commiting on another blocked branch after setting up git-hooks should be allowed" {
	write_blocked_branches "some-branch"

	run commit_changes "$(repository_dir)" "another-branch"

	assert_success
}

@test "If team json contains 2 blocked branch commiting on second one after setting up git-hooks should be blocked" {
	local branch_name="some-branch"
	write_blocked_branches "another-branch" "$branch_name"

	run commit_changes "$(repository_dir)" $branch_name

	assert_output --partial "Action \"commit\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "If team json contains blocked branch pushing on the blocked after setting up git-hooks branch should be blocked" {
	local branch_name; branch_name="$(unique_branch_name)"
	commit_changes "$(repository_dir)" "$branch_name"
	write_blocked_branches "$branch_name"

	push_changes "$(repository_dir)" "$branch_name"

	assert_output --partial "Action \"push\" not allowed on branch \"$branch_name\""
	assert_failure
}

@test "if team json contains commitPrefixRegex 'commiting' with a message and a branch both without matching prefix is blocked" {
	local commit_prefix_regex="Test-[0-9]+"
	write_commit_prefix_regex "$commit_prefix_regex"

	run commit_changes "$(repository_dir)" "no-prefix-branch" "no-prefix-message"

	assert_line --index 1 "The commit message needs a commit prefix, that matches the following regex $commit_prefix_regex."
	assert_line --index 2 "Either add the commit prefix to you commit message, or include it in the branch name."
	assert_line --index 3 "Use '--no-verify' to skip git-hooks."
	assert_failure
}

@test "if team json contains commit_prefix_regex 'commiting' with a message that has a matching prefix on a branch without prefix is not blocked" {
	local matching_prefix="Test-1"
	write_commit_prefix_regex "Test-[0-9]+"

	run commit_changes "$(repository_dir)" "no-prefix-branch" "$matching_prefix: prefixed-message"

	assert_line --index 1 "The commit message contains an issue ID ($matching_prefix). Good job!"
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix is not blocked" {
	local commit_prefix="Asdf-99"
	write_commit_prefix_regex "Asdf-[0-9]+"

	run commit_changes "$(repository_dir)" "feature/$commit_prefix/prefixed-branch" "not-prefix-message"

	assert_line --index 1 "Commit prefix '$commit_prefix' was found in current branch name, prepended to commit message."
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that starts with 'Merge branch' is not blocked" {
	test_merge_commit_messages_are_not_blocked "Merge branch"
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that starts with 'Merge remote-tracking branch' is not blocked" {
	test_merge_commit_messages_are_not_blocked "Merge remote-tracking branch"
}

test_merge_commit_messages_are_not_blocked() {
	local commit_message="$1"
	write_commit_prefix_regex "Asdf-[0-9]+"

	run commit_changes "$(repository_dir)" "no-prefix-branch" "$commit_message"

	assert_line --index 1 "Merge commit detected, skipping commit-msg hook."
	assert_success
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has a matching prefix leads to commit with given message" {
	local commit_message="Test-1: prefixed-message"
	write_commit_prefix_regex "Test-[0-9]+"
	commit_changes "$(repository_dir)" "no-prefix-branch" "$commit_message"

	run get_commit_message_of_last_commit "$(repository_dir)"

	assert_output "$commit_message"
}

@test "if team json contains commitPrefixRegex 'commiting' with a message that has no matching prefix on a branch containing prefix leads to commit with prefixed message" {
	local matching_prefix="Asdf-99"
	local commit_message="not-prefixed-message"
	write_commit_prefix_regex "Asdf-[0-9]+"
	commit_changes "$(repository_dir)" "feature/$matching_prefix/prefixed-branch" "$commit_message"

	run get_commit_message_of_last_commit "$(repository_dir)"

	assert_output "$matching_prefix: $commit_message"
}

@test "if team json does not contain commitPrefixRegex 'commiting' does not check for commit prefix" {
	run commit_changes "$(repository_dir)" "not-prefixed-branch" "not-prefixed-message"

	refute_output --partial "JIRA-ID '' was found in current branch name, prepended to commit message."
}

@test "if team json does not contain commitPrefixRegex while 'commiting' a merge commit, it does not check for commit prefix" {
	run commit_changes "$(repository_dir)" "not-prefixed-branch" "not-prefixed-message"

	refute_output --partial "JIRA-ID '' was found in current branch name, prepended to commit message."
}

@test "if team json does not contain commitPrefixRegex 'commiting' with a message that starts with 'Merge branch' does not check for prefix" {
	test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes "Merge branch"
}

@test "if team json does not contain commitPrefixRegex 'commiting' with a message that starts with 'Merge remote-tracking branch' does not check for prefix" {
	test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes "Merge remote-tracking branch"
}

test_while_commiting_merge_commit_it_does_not_check_for_commit_prefixes() {
	local commit_message="$1"

	run commit_changes "$(repository_dir)" "no-prefix-branch" "$commit_message"

	refute_output --partial "Merge commit detected, skipping commit-msg hook."
	assert_success
}
