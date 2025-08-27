setup() {
	bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'

	one_cloned_repository_with_git_hooks_setup
}

teardown() {
	one_cloned_repository_with_git_hooks_teardown
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
	local hook_name="unknown-hook"

	run mrt_execute git-hook --hook-name "$hook_name" --repository-path "$(repository_dir)"

	assert_output --partial "The given git-hook \"$hook_name\" does not exist."
	assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
	run mrt_execute git-hook --hook-name "pre-commit" --repository-path "$(test_env_dir)"

	assert_output --partial "The given path \"$(test_env_dir)\" does not contain a repository."
	assert_failure
}
