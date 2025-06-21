bats_load_library 'test-case-with-1-cloned-repository-and-set-up-git-hooks'

@test "if subcommand 'git-hook' gets called with an unknown git hook it fails" {
	hookName="unknown-hook"

	run execute git-hook --hook-name "$hookName" --repository-path "$(repositoryDir)"

	assert_output --partial "The given git-hook \"$hookName\" does not exist."
	assert_failure
}

@test "if subcommand 'git-hook' gets called with a path that does not contain a repository it fails" {
	run execute git-hook --hook-name "pre-commit" --repository-path "$(testEnvDir)"

	assert_output --partial "The given path \"$(testEnvDir)\" does not contain a repository."
	assert_failure
}
