setup() {
  bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'
  bats_load_library 'scripts/script_writer.bash'

  one_cloned_repository_with_git_hooks_setup
}

teardown() {
  one_cloned_repository_with_git_hooks_teardown
}

@test "if pre-commit scripts exist 'committing' will execute them" {
	local scriptsPath; scriptsPath="$(repository_dir)/hook-scripts/pre-commit"
	local firstScript="$scriptsPath/script1"
	local secondScript="$scriptsPath/script2"
	write_spy_script "$firstScript"
	write_spy_script "$secondScript"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed "$firstScript"
	assert_script_was_executed "$secondScript"
}

@test "if commit-msg scripts exits with failure 'commiting' will also fail" {
	local scriptPath; scriptPath="$(repository_dir)/hook-scripts/commit-msg/script"
	write_stub_script "$scriptPath" "1" "some-output"

	run commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_failure
}

@test "if commit-msg scripts has output 'commiting' will contain the same output" {
	local scriptOutput="some-output"
	local scriptPath; scriptPath="$(repository_dir)/hook-scripts/commit-msg/script"
	write_stub_script "$scriptPath" "0" "$scriptOutput"

	run commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_output --partial "$scriptOutput"
}

@test "if pre-commit hook gets executed, it gets passed the git parameters" {
	local scriptPath; scriptPath="$(repository_dir)/hook-scripts/pre-commit/script"
	write_spy_script "$scriptPath"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed_with_parameters "$scriptPath" ""
}

@test "if pre-push hook gets executed, it gets passed the git parameters" {
	local scriptPath; scriptPath="$(repository_dir)/hook-scripts/pre-push/script"
	write_spy_script "$scriptPath"
	local branchName; branchName="$(unique_branch_name)"
	commit_changes "$(repository_dir)" "$branchName" "some-message"

	push_changes "$(repository_dir)" "$branchName"

	local originUrl; originUrl=$(git -C "$(repository_dir)" config --get remote.origin.url)
	local remoteName; remoteName=$(git remote)
	assert_script_was_executed_with_parameters "$scriptPath" "$remoteName $originUrl"
}

@test "if commit-msg hook gets executed, it gets passed the git parameters" {
	local scriptPath; scriptPath="$(repository_dir)/hook-scripts/commit-msg/script"
	write_spy_script "$scriptPath"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed_with_parameters "$scriptPath" ".git/COMMIT_EDITMSG"
}
