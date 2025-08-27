one_cloned_repository_with_git_hooks_setup() {
	bats_load_library 'fixtures/authenticated_fixture.bash'
	bats_load_library 'repositories_path.bash'
	bats_load_library 'mrt/setup.bash'
	bats_load_library 'git.bash'

	authenticated_setup

	clone_testing_repositories "$(test_env_dir)/$(default_repositories_path)" "$(repository_name)"
	mrt_setup_git_hooks
}

one_cloned_repository_with_git_hooks_teardown() {
	bats_load_library 'fixtures/authenticated_fixture.bash'

	authenticated_teardown
}

repository_dir() {
	bats_load_library 'repositories_path.bash'

	echo "$(test_env_dir)/$(default_repositories_path)/$(repository_name)"
}

repository_name() {
	echo "1_TestRepository"
}
