setup() {
  bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'
  bats_load_library 'scripts/script_writer.bash'

  one_cloned_repository_with_git_hooks_setup
}

teardown() {
  one_cloned_repository_with_git_hooks_teardown
}

@test "if pre-commit scripts exist 'committing' will execute them" {
	local scripts_path; scripts_path="$(repository_dir)/hook-scripts/pre-commit"
	local first_script="$scripts_path/script1"
	local second_script="$scripts_path/script2"
	write_spy_script "$first_script"
	write_spy_script "$second_script"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed "$first_script"
	assert_script_was_executed "$second_script"
}

@test "if commit-msg scripts exits with failure 'commiting' will also fail" {
	local script_path; script_path="$(repository_dir)/hook-scripts/commit-msg/script"
	write_stub_script "$script_path" "1" "some-output"

	run commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_failure
}

@test "if commit-msg scripts has output 'commiting' will contain the same output" {
	local script_output="some-output"
	local script_path; script_path="$(repository_dir)/hook-scripts/commit-msg/script"
	write_stub_script "$script_path" "0" "$script_output"

	run commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_output --partial "$script_output"
}

@test "if pre-commit hook gets executed, it gets passed the git parameters" {
	local script_path; script_path="$(repository_dir)/hook-scripts/pre-commit/script"
	write_spy_script "$script_path"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed_with_parameters "$script_path" ""
}

@test "if pre-push hook gets executed, it gets passed the git parameters" {
	local script_path; script_path="$(repository_dir)/hook-scripts/pre-push/script"
	write_spy_script "$script_path"
	local branch_name; branch_name="$(unique_branch_name)"
	commit_changes "$(repository_dir)" "$branch_name" "some-message"

	push_changes "$(repository_dir)" "$branch_name"

	local origin_url; origin_url=$(git -C "$(repository_dir)" config --get remote.origin.url)
	local remote_name; remote_name=$(git remote)
	assert_script_was_executed_with_parameters "$script_path" "$remote_name $origin_url"
}

@test "if commit-msg hook gets executed, it gets passed the git parameters" {
	local script_path; script_path="$(repository_dir)/hook-scripts/commit-msg/script"
	write_spy_script "$script_path"

	commit_changes "$(repository_dir)" "some-branch" "some-message"

	assert_script_was_executed_with_parameters "$script_path" ".git/COMMIT_EDITMSG"
}
