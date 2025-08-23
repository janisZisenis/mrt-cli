setup() {
  bats_load_library 'fixtures/one_cloned_repository_with_git_hooks_set_up_fixture.bash'

  one_cloned_repository_with_git_hooks_setup
}

teardown() {
  one_cloned_repository_with_git_hooks_teardown
}

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
	local hookName="unknown-hook"

	run mrtExecute git-hook --hook-name "$hookName" --repository-path "$(repositoryDir)"

	assert_output --partial "The given git-hook \"$hookName\" does not exist."
	assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
	run mrtExecute git-hook --hook-name "pre-commit" --repository-path "$(testEnvDir)"

	assert_output --partial "The given path \"$(testEnvDir)\" does not contain a repository."
	assert_failure
}
