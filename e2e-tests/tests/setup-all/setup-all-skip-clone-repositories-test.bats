setup() {
	bats_load_library 'repositories_path.bash'
	bats_load_library 'test_repositories.bash'
	bats_load_library 'fixtures/authenticated_fixture.bash'
	bats_load_library 'mrt/setup.bash'
	bats_load_library 'write_team_file.bash'

	authenticated_setup
}

teardown() {
	authenticated_teardown
}

@test "if setup is run with skipping the clone step it should not clone the repositories" {
	local repository="1_TestRepository"
	local repository_url
	repository_url="$(get_testing_repository_url "$repository")"
	write_repositories_urls "$repository_url"

	run mrt_setup_all --skip-clone-repositories

	assert_dir_not_exist "$(test_env_dir)/$(default_repositories_path)/$repository"
}

@test "if setup is run with skipping the clone step it should print a skip message" {
	local repository="1_TestRepository"
	local repository_url
	repository_url="$(get_testing_repository_url "$repository")"
	write_repositories_urls "$repository_url"

	run mrt_setup_all --skip-clone-repositories

	assert_line --index 0 "Skipping clone-repositories step."
}
