setup() {
	bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'

	one_cloned_repository_with_git_hooks_setup
}

teardown() {
	one_cloned_repository_with_git_hooks_teardown
}

assert_git_hook_fails_with_invalid_name() {
	local hook_name="$1"
	local repo_path="$2"

	run mrt_execute git-hook --hook-name "$hook_name" --repository-path "$repo_path"

	assert_failure
	assert_output --partial "The given git-hook \"$hook_name\" does not exist."
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
	assert_git_hook_fails_with_invalid_name "unknown-hook" "$(repository_dir)"
}

@test "passing hook name with globbing pattern should fail" {
	assert_git_hook_fails_with_invalid_name "pre-commit*" "$(repository_dir)"
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
	run mrt_execute git-hook --hook-name "pre-commit" --repository-path "$(test_env_dir)"

	assert_output --partial "The given path \"$(test_env_dir)\" does not contain a repository."
	assert_failure
}
